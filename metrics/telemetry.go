package metrics

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ryanuber/columnize"
	"net/http"
	"strings"
)

// TODO: Make the URL Agent Version adaptable (i.e., alter URL string to corresponding version)
const (
	TelemetryURL = "https://developer.hashicorp.com/consul/docs/agent/telemetry"
)

type AgentTelemetryMetric struct {
	Name string
	Unit string
	Type string
}

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

func GetTransactionTimingMetrics() (string, []AgentTelemetryMetric, error) {
	// Define a data structure to store metric endpoints.
	telemetryMetrics := []string{"\x1fTransaction Timing\x1f"}
	telemetryMetrics = append(telemetryMetrics, "Metric\x1fUnit\x1fType")

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
		if endpoint == "consul.kvs.apply" || endpoint == "consul.txn.apply" || endpoint == "consul.raft.apply" || endpoint == "consul.raft.commitTime" {
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

func ListMetrics(all, transactionTiming bool) error {
	var latestMetrics string
	var err error
	if all {
		if latestMetrics, _, err = GetTelemetryMetrics(); err != nil {
			return err
		}
	}
	if transactionTiming {
		if latestMetrics, _, err = GetTransactionTimingMetrics(); err != nil {
			return err
		}
	}

	fmt.Printf("\nConsul Telemetry Metric Names (pulled from: %s)\n\n", TelemetryURL)
	fmt.Println(latestMetrics)
	return nil
}
