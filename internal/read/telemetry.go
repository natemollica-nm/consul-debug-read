package read

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ryanuber/columnize"
	"log"
	"net/http"
	"strings"
)

type AgentTelemetryMetric struct {
	Name string
	Unit string
	Type string
}

//type ConsulTelemetryEndpoints struct {
//	KeyMetricNames
//}

// TODO: Make the URL Agent Version adaptable (i.e., alter URL string to corresponding version)
const (
	TelemetryURL            = "https://developer.hashicorp.com/consul/docs/agent/telemetry"
	telegrafMetricsFilePath = "metrics/telegraf"
)

func GetTelemetryMetrics() (string, []AgentTelemetryMetric, error) {
	// Define a data structure to store metric endpoints.
	telemetryMetrics := []string{"Metric\x1fUnit\x1fType"}

	// Send an HTTP GET request to the Consul telemetry metrics reference page.
	response, err := http.Get(TelemetryURL)
	if err != nil {
		return "", []AgentTelemetryMetric{}, err
	}
	defer response.Body.Close()
	cleanup := func(err error) error {
		_ = response.Body.Close()
		return err
	}
	// Parse the HTML content.
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", []AgentTelemetryMetric{}, err
	}

	// Extract metric endpoints from the HTML.
	doc.Find("table tbody tr").Each(func(index int, rowHtml *goquery.Selection) {
		// Parse and extract metric endpoint.
		endpoint := rowHtml.Find("td:nth-child(1)").Text()
		metricUnit := rowHtml.Find("td:nth-child(3)").Text()
		metricType := rowHtml.Find("td:nth-child(4)").Text()
		if strings.HasPrefix(endpoint, "consul") {
			telemetryMetrics = append(telemetryMetrics, fmt.Sprintf("%s\x1f%s\x1f%s\x1f",
				endpoint, metricUnit, metricType))
		}
	})
	var telemetryInfo []AgentTelemetryMetric
	for i, line := range telemetryMetrics {
		infoSections := strings.Split(line, string([]byte{0x1f}))
		if len(infoSections) < 3 || i == 0 {
			continue
		}
		info := AgentTelemetryMetric{
			Name: infoSections[0],
			Unit: infoSections[1],
			Type: infoSections[2],
		}
		telemetryInfo = append(telemetryInfo, info)
	}
	if err := response.Body.Close(); err != nil {
		return "", []AgentTelemetryMetric{}, cleanup(err)
	}
	// Build output string in columnized format for readability
	output := columnize.Format(telemetryMetrics, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, telemetryInfo, nil
}

func ListMetrics() (string, error) {
	var latestMetrics string
	var err error
	if latestMetrics, _, err = GetTelemetryMetrics(); err != nil {
		return "", err
	}
	fmt.Printf("\nConsul Telemetry Metric Names (pulled from: %s)\n\n", TelemetryURL)
	return latestMetrics, nil
}

func (b *Debug) GenerateTelegrafMetrics() error {
	metrics := b.Metrics.Metrics
	log.Printf("converting metrics timestamps to RFC3339")
	for i := range metrics {
		telegrafMetrics := metrics[i]
		ts := metrics[i].Timestamp
		timestampRFC, err := ToRFC3339(ts)
		if err != nil {
			return err
		}
		telegrafMetrics.Timestamp = timestampRFC

		data, err := json.MarshalIndent(telegrafMetrics, "", "  ")
		if err != nil {
			return err
		}
		// Write out the resultant metrics.json file.
		// Must be 0644 because this is written by the consul-k8s user but needs
		// to be readable by the consul user
		metricsFile := fmt.Sprintf("%s/metrics-%d.json", telegrafMetricsFilePath, i)
		if err = WriteFileWithPerms(metricsFile, string(data), 0755); err != nil {
			return fmt.Errorf("error writing RFC3339 formatted metrics to %s: %v", telegrafMetricsFilePath, err)
		}
	}
	fmt.Printf("telegraf metrics generated successfully to %s", telegrafMetricsFilePath)
	return nil
}
