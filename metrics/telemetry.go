package metrics

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

const (
	telemetryURL = "https://developer.hashicorp.com/consul/docs/agent/telemetry"
)

func getTelemetryMetrics() ([]string, error) {
	// Define a data structure to store metric endpoints.
	var telemetryMetrics []string

	// Send an HTTP GET request to the Consul telemetry metrics reference page.
	response, err := http.Get(telemetryURL)
	if err != nil {
		return telemetryMetrics, err
	}
	defer response.Body.Close()

	// Parse the HTML content.
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return telemetryMetrics, err
	}

	// Extract metric endpoints from the HTML.
	doc.Find("table tbody tr").Each(func(index int, rowHtml *goquery.Selection) {
		// Parse and extract metric endpoint.
		endpoint := rowHtml.Find("td:nth-child(1)").Text()
		if strings.HasPrefix(endpoint, "consul") {
			telemetryMetrics = append(telemetryMetrics, endpoint)
		}
	})

	return telemetryMetrics, nil
}

func ListMetrics() error {
	var latestMetrics []string
	var err error
	if latestMetrics, err = getTelemetryMetrics(); err != nil {
		return err
	}
	for _, metricName := range latestMetrics {
		fmt.Println(metricName)
	}
	return nil
}

func ValidateMetricName(name string) error {
	var latestMetrics []string
	var err error
	if latestMetrics, err = getTelemetryMetrics(); err != nil {
		return err
	}
	for _, metricName := range latestMetrics {
		reg := regexp.MustCompile(`^consul\.proxy\..+$`)
		if reg.MatchString(name) {
			fmt.Printf("built-in mesh proxy prefix used: %s\n", name)
			return nil
		}
		if strings.Contains(metricName, name) {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("%s not a valid telemetry metric name\n  visit: %s for full list of consul telemetry metrics", name, telemetryURL))
}
