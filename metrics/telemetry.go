package metrics

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ryanuber/columnize"
	"net/http"
	"regexp"
	"strings"
)

const (
	telemetryURL = "https://developer.hashicorp.com/consul/docs/agent/telemetry"
)

func getTelemetryMetrics() (string, error) {
	// Define a data structure to store metric endpoints.
	telemetryMetrics := []string{"Metric\x1fUnit\x1fType"}
	telemetryMetrics = append(telemetryMetrics, fmt.Sprintf("%s\x1f%s\x1f%s\x1f",
		"----------------", "----------------", "----------------"))

	// Send an HTTP GET request to the Consul telemetry metrics reference page.
	response, err := http.Get(telemetryURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Parse the HTML content.
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
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
	output, err := columnize.Format(telemetryMetrics, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	if err != nil {
		return "", nil
	}
	return output, nil
}

func ListMetrics() error {
	var latestMetrics string
	var err error
	if latestMetrics, err = getTelemetryMetrics(); err != nil {
		return err
	}
	fmt.Printf("\nConsul Telemetry Metric Names (pulled from: %s)\n\n", telemetryURL)
	fmt.Println(latestMetrics)
	return nil
}

func ValidateMetricName(name string) error {
	var latestMetrics string
	var err error
	if latestMetrics, err = getTelemetryMetrics(); err != nil {
		return err
	}
	reg := regexp.MustCompile(`^consul\.proxy\..+$`)
	if reg.MatchString(name) {
		fmt.Printf("built-in mesh proxy prefix used: %s\n", name)
		return nil
	}
	if strings.Contains(latestMetrics, name) {
		return nil
	}
	return errors.New(fmt.Sprintf("%s not a valid telemetry metric name\n  visit: %s for full list of consul telemetry metrics", name, telemetryURL))
}
