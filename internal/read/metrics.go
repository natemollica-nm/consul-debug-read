package read

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"regexp"
	"sort"
	"strings"
)

type Gauge struct {
	Name   string            `json:"Name"`
	Value  float64           `json:"Value"`
	Labels map[string]string `json:"Labels"`
}

type Points struct {
	Name   string            `json:"Name"`
	Points float64           `json:"Points"`
	Labels map[string]string `json:"Labels"`
}

type Counters struct {
	Name   string            `json:"Name"`
	Count  int               `json:"Count"`
	Rate   float64           `json:"Rate"`
	Sum    float64           `json:"Sum"`
	Min    float64           `json:"Min"`
	Max    float64           `json:"Max"`
	Mean   float64           `json:"Mean"`
	Stddev float64           `json:"Stddev"`
	Labels map[string]string `json:"Labels"`
}

type Samples struct {
	Name   string            `json:"Name"`
	Count  int               `json:"Count"`
	Rate   float64           `json:"Rate"`
	Sum    float64           `json:"Sum"`
	Min    float64           `json:"Min"`
	Max    float64           `json:"Max"`
	Mean   float64           `json:"Mean"`
	Stddev float64           `json:"Stddev"`
	Labels map[string]string `json:"Labels"`
}

type Metric struct {
	Timestamp string     `json:"Timestamp"`
	Gauges    []Gauge    `json:"Gauges"`
	Points    []Points   `json:"Points"`
	Counters  []Counters `json:"Counters"`
	Samples   []Samples  `json:"Samples"`
}

type Metrics struct {
	Metrics    []Metric
	MetricsMap map[string][]map[string]interface{}
}

type Index struct {
	Version      int      `json:"Version"`
	AgentVersion string   `json:"AgentVersion"`
	Interval     string   `json:"Interval"`
	Duration     string   `json:"Duration"`
	Targets      []string `json:"Targets"`
}

// getUnitAndType returns the Unit and Type for a given Name.
func getUnitAndType(name string, telemetry []AgentTelemetryMetric) (string, string) {
	for _, metric := range telemetry {
		if metric.Name == name {
			return metric.Unit, metric.Type
		} else if name == "*" {
			return metric.Unit, metric.Type
		}
	}
	return "-", "-"
}

// BuildMetricsIndex
// Builds metrics map from the ingested metrics.json,
// extracts metric name, value, labels, and timestamp
// for retrieval via query, and uploads to boltDB.
func (b *Debug) BuildMetricsIndex() {
	b.Metrics.MetricsMap = make(map[string][]map[string]interface{})

	for _, metric := range b.Metrics.Metrics {
		timestamp := metric.Timestamp

		// Iterate over Gauges and add to the map
		for _, gauge := range metric.Gauges {
			metricData := map[string]interface{}{
				"name":      gauge.Name,
				"timestamp": timestamp,
				"value":     gauge.Value,
				"labels":    gauge.Labels,
			}
			b.Metrics.MetricsMap[gauge.Name] = append(b.Metrics.MetricsMap[gauge.Name], metricData)
		}

		// Iterate over Points and add to the map
		for _, point := range metric.Points {
			metricData := map[string]interface{}{
				"name":      point.Name,
				"timestamp": timestamp,
				"value":     point.Points,
				"labels":    point.Labels,
			}
			b.Metrics.MetricsMap[point.Name] = append(b.Metrics.MetricsMap[point.Name], metricData)
		}

		// Iterate over Counters and add to the map
		for _, counter := range metric.Counters {
			metricData := map[string]interface{}{
				"name":      counter.Name,
				"timestamp": timestamp,
				"value":     counter.Count,
				"labels":    counter.Labels,
			}
			b.Metrics.MetricsMap[counter.Name] = append(b.Metrics.MetricsMap[counter.Name], metricData)
		}

		// Iterate over Samples and add to the map
		for _, sample := range metric.Samples {
			metricData := map[string]interface{}{
				"name":      sample.Name,
				"timestamp": timestamp,
				"value":     sample.Mean,
				"labels":    sample.Labels,
			}
			b.Metrics.MetricsMap[sample.Name] = append(b.Metrics.MetricsMap[sample.Name], metricData)
		}
	}
}

func formatMetricValue(mValue interface{}, unit string) (string, error) {
	switch value := mValue.(type) {
	case float64:
		if timeReg.MatchString(unit) {
			return ConvertToReadableTime(value, unit)
		} else if percentageReg.MatchString(unit) {
			return fmt.Sprintf("%.2f%%", value*100.00), nil
		} else {
			return fmt.Sprintf("%v", value), nil
		}
	case int:
		if timeReg.MatchString(unit) {
			return ConvertToReadableTime(float64(value), unit)
		} else if percentageReg.MatchString(unit) {
			return fmt.Sprintf("%.2f%%", float64(value)*100.00), nil
		} else {
			return fmt.Sprintf("%v", value), nil
		}
	case nil:
		return "nil", nil
	default:
		return fmt.Sprintf("%v", mValue), nil
	}
}

// GetMetricValues / extracts all timestamped occurrences of metric values by name
func (b *Debug) GetMetricValues(name string, validate, byValue, short bool) (string, error) {
	var err error

	// Get telemetry metrics
	stringInfo, telemetryInfo, _ := GetTelemetryMetrics()
	if validate {
		if ok := validateName(name, stringInfo); !ok {
			errString := fmt.Sprintf("'%s' not a valid telemetry metric name\n  visit: %s for a full list of consul telemetry metrics", name, TelemetryURL)
			return "", fmt.Errorf(errString)
		}
	}

	// Retrieve metric data and matching metric names
	metricData, matchedNames, found := b.Metrics.extractMetricValueByName(name)
	if !found {
		// No metrics found matching the given name
		result := []string{fmt.Sprintf("*\x1f%s\x1f=>\x1fnil\x1fvalue(s)\x1freturned\x1f", name)}
		output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
		return output, nil
	}

	var result []string

	// Iterate through matched metric names and process data
	for _, matchedName := range matchedNames {
		unit, metricType := getUnitAndType(matchedName, telemetryInfo)

		// Build header for each matched metric name
		result = append(result, fmt.Sprintf("\x1f%s\x1f", matchedName))
		result = append(result, fmt.Sprintf("\x1f%s\x1f", strings.Repeat("-", len(matchedName))))
		if short {
			result = append(result, "Timestamp\x1fValue\x1f")
		} else {
			result = append(result, "Timestamp\x1fMetric\x1fType\x1fUnit\x1fValue\x1f")
		}

		// Process metric data for the current matched name
		for _, data := range metricData {
			for _, scrape := range data {
				timestamp := scrape["timestamp"].(string)
				mValue := scrape["value"]
				mLabels := scrape["labels"].(map[string]string)

				// Construct labels
				var labels []string
				for k, v := range mLabels {
					labels = append(labels, fmt.Sprintf("%s=%v", k, v))
				}

				// Process metric value
				var formattedValue string
				formattedValue, err = formatMetricValue(mValue, unit)
				if err != nil {
					return "", err
				}

				// Add metric record to the result
				totalLabels := len(labels)
				if short {
					if totalLabels > 0 && totalLabels <= 3 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f", timestamp, formattedValue, labels))
					} else if totalLabels > 3 {
						labels = labels[:6-3]
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s + %d more ...", timestamp,
							formattedValue, labels, totalLabels))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f", timestamp, formattedValue))
					}
				} else {
					if totalLabels > 0 && totalLabels <= 3 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							timestamp, matchedName, metricType, unit, formattedValue, labels))
					} else if totalLabels > 3 {
						labels = labels[:6-3]
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s + %d more ...",
							timestamp, matchedName, metricType, unit, formattedValue, labels, totalLabels))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							timestamp, matchedName, metricType, unit, formattedValue))
					}
				}
			}
		}
	}

	// Add label information if applicable
	if len(result) > 2 && strings.Contains(result[2], "Labels") {
		result[2] += "Labels\x1f"
	}
	if name == "consul.runtime.total_gc_pause_ns" {
		result[2] += "gc/min\x1f"
		// Calculate GC rates and update the result
		for _, values := range metricData {
			var rate string
			for i := range values {
				if i == 0 {
					result[i+3] = fmt.Sprintf("%s-\x1f", result[i+3])
				} else {
					gcPause := values[i]
					prevGCPause := values[i-1]
					rate, err = CalculateGCRate(gcPause, prevGCPause)
					if err != nil {
						return "", fmt.Errorf("error calculating rate: %v", err)
					}
					result[i+3] = fmt.Sprintf("%s%s\x1f", result[i+3], rate)
				}
			}
		}
	}

	// Sort results by value if requested
	if byValue {
		sort.Sort(ByValue(result[3:]))
	}

	// Format the result into a columnized string and return
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, nil
}

// matchMetricsByRegex matches metric names using a given regex and returns the matching data and metric names.
func matchMetricsByRegex(metricsMap map[string][]map[string]interface{}, pattern string) ([][]map[string]interface{}, []string, bool) {
	regex := regexp.MustCompile(pattern)
	var matches [][]map[string]interface{}
	var matchedNames []string
	found := false

	for name, data := range metricsMap {
		if regex.MatchString(name) {
			matches = append(matches, data)
			matchedNames = append(matchedNames, name)
			found = true
		}
	}

	return matches, matchedNames, found
}

// extractMetricValueByName uses regex to pull the matching metrics data and metric names from the metrics map.
// It returns a slice of matched data, a slice of matched names, and a boolean indicating if the metric was found.
func (m Metrics) extractMetricValueByName(metricName string) ([][]map[string]interface{}, []string, bool) {
	// Replace * with regex wildcard .* if present
	if strings.Contains(metricName, "*") {
		pattern := strings.ReplaceAll(regexp.QuoteMeta(metricName), `\*`, ".*")
		return matchMetricsByRegex(m.MetricsMap, pattern)
	}

	// No wildcard case
	return matchMetricsByRegex(m.MetricsMap, `.*`+regexp.QuoteMeta(metricName))
}

func (b *Debug) Summary() string {
	title := "Metrics Bundle Summary"
	ul := strings.Repeat("-", len(title))
	captures, _ := b.numberOfCaptures()
	return fmt.Sprintf("%s\n%s\nDatacenter: %v\nHostname: %s\nAgent Version: %s\nRaft State: %s\nInterval: %s\nDuration: %s\nCapture Targets: %v\nTotal Captures: %d\nCapture Time Start: %s\nCapture Time Stop: %s\n",
		title,
		ul,
		b.Agent.Config.Datacenter,
		b.Host.Host.Hostname,
		b.Index.AgentVersion,
		b.Agent.Stats.Raft.State,
		b.Index.Interval,
		b.Index.Duration,
		b.Index.Targets,
		captures,
		b.Metrics.Metrics[0].Timestamp,
		b.Metrics.Metrics[len(b.Metrics.Metrics)-1].Timestamp)
}
