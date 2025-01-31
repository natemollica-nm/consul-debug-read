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

// GetMetricValues / extracts all timestamped occurrences of metric values by name
func (b *Debug) GetMetricValues(name string, validate, byValue, short bool) (string, error) {
	var err error

	stringInfo, telemetryInfo, _ := GetTelemetryMetrics()
	if validate {
		if ok := validateName(name, stringInfo); !ok {
			errString := fmt.Sprintf("'%s' not a valid telemetry metric name\n  visit: %s for a full list of consul telemetry metrics", name, TelemetryURL)
			return "", fmt.Errorf(errString)
		}
	}
	unit, metricType := getUnitAndType(name, telemetryInfo)

	// Build Metrics Information Title
	result := []string{fmt.Sprintf("\x1f%s\x1f", name)}
	//nolint:staticcheck
	ul := fmt.Sprintf(strings.Repeat("-", len(name)))
	result = append(result, fmt.Sprintf("\x1f%s\x1f", ul))
	if short {
		result = append(result, "Timestamp\x1fValue\x1f")

	} else {
		result = append(result, "Timestamp\x1fMetric\x1fType\x1fUnit\x1fValue\x1f")
	}

	// Retrieve metric data from the metrics map
	metricData, found := b.Metrics.extractMetricValueByName(name)
	if !found {
		// Metric not found in the metrics map ==> nil return
		result = []string{fmt.Sprintf("*\x1f%s\x1f=>\x1fnil\x1fvalue(s)\x1freturned\x1f", name)}
		output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
		return output, nil
	}

	var label []string
	// Iterate over metric data and construct result
	for _, data := range metricData {
		for _, scrape := range data {
			timestamp := scrape["timestamp"].(string)
			mValue := scrape["value"]
			mLabels := scrape["labels"].(map[string]string)
			for k, v := range mLabels {
				label = append(label, fmt.Sprintf("%s=%v", k, v))
			}
			if mValue != nil {
				var v string
				if timeReg.MatchString(unit) {
					v, err = ConvertToReadableTime(mValue, unit)
					if err != nil {
						return "", err
					}
				} else if bytesReg.MatchString(unit) {
					conv := ByteConverter{}
					v = conv.ConvertToReadableBytes(mValue)
				} else if percentageReg.MatchString(unit) {
					vFloat := mValue.(float64)
					percent := vFloat * 100.00
					v = fmt.Sprintf("%.2f%%", percent)
				} else {
					v = fmt.Sprintf("%v", mValue)
				}
				totalLabels := len(label)
				if short {
					if totalLabels > 0 && totalLabels <= 3 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f",
							timestamp, v, label))
					} else if totalLabels > 3 {
						label = label[:6-3]
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s + %d more ...",
							timestamp, v, label, totalLabels))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f",
							timestamp, v))
					}
				} else {
					if totalLabels > 0 && totalLabels <= 3 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							timestamp, name, metricType, unit, v, label))
					} else if totalLabels > 3 {
						label = label[:6-3]
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s + %d more ...",
							timestamp, name, metricType, unit, v, label, totalLabels))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							timestamp, name, metricType, unit, v))
					}
				}
			}
		}

	}
	if len(label) > 0 {
		result[2] += "Labels\x1f"
	}
	if name == "consul.runtime.total_gc_pause_ns" {
		result[2] += "gc/min\x1f"
		// Calculate the GC rate and add it to each line
		for _, values := range metricData {
			var rate string
			for i := 0; i < len(values); i++ {
				if i == 0 {
					// no previous time stamped gc gauge value to perform non-neg diff calc
					result[i+3] = fmt.Sprintf("%s%s\x1f", result[i+3], "-")
				} else {
					gcPause := values[i]
					prevGcPause := values[i-1]
					// Calculate non-negative difference in gc_pause rate
					rate, err = CalculateGCRate(gcPause, prevGcPause)
					if err != nil {
						return "", fmt.Errorf("error calculating rate: %v", err)
					}

					// Append the calculated rate to the line
					result[i+3] = fmt.Sprintf("%s%s\x1f", result[i+3], rate)
				}
			}
		}
	}
	if byValue {
		sort.Sort(ByValue(result[3:]))
	}
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, nil
}

// matchMetricsByRegex matches metric names using a given regex and returns the matching data.
func matchMetricsByRegex(metricsMap map[string][]map[string]interface{}, pattern string) ([][]map[string]interface{}, bool) {
	regex := regexp.MustCompile(pattern)
	var matches [][]map[string]interface{}
	found := false
	for name, data := range metricsMap {
		if regex.MatchString(name) {
			matches = append(matches, data)
			found = true
		}
	}
	return matches, found
}

// extractMetricValueByName uses regex to pull the matching metrics data from the metrics map.
// It returns a slice of maps containing the matched metrics and a boolean indicating if the metric was found.
func (m Metrics) extractMetricValueByName(metricName string) ([][]map[string]interface{}, bool) {
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
