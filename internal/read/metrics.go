package read

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"regexp"
	"sort"
	"strings"
)

type Gauge struct {
	Name   string                 `json:"Name"`
	Value  float64                `json:"Value"`
	Labels map[string]interface{} `json:"Labels"`
}

type Points struct {
	Name   string                 `json:"Name"`
	Points float64                `json:"Points"`
	Labels map[string]interface{} `json:"Labels"`
}

type Counters struct {
	Name   string                 `json:"Name"`
	Count  int                    `json:"Count"`
	Rate   float64                `json:"Rate"`
	Sum    float64                `json:"Sum"`
	Min    float64                `json:"Min"`
	Max    float64                `json:"Max"`
	Mean   float64                `json:"Mean"`
	Stddev float64                `json:"Stddev"`
	Labels map[string]interface{} `json:"Labels"`
}

type Samples struct {
	Name   string                 `json:"Name"`
	Count  int                    `json:"Count"`
	Rate   float64                `json:"Rate"`
	Sum    float64                `json:"Sum"`
	Min    float64                `json:"Min"`
	Max    float64                `json:"Max"`
	Mean   float64                `json:"Mean"`
	Stddev float64                `json:"Stddev"`
	Labels map[string]interface{} `json:"Labels"`
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
	metricsMap map[string][]map[string]interface{}
}

type Index struct {
	Version      int      `json:"Version"`
	AgentVersion string   `json:"AgentVersion"`
	Interval     string   `json:"Interval"`
	Duration     string   `json:"Duration"`
	Targets      []string `json:"Targets"`
}

// GetMetricValues
// 1. unless --validate set to false (i.e., --validate=false), validate metric name with telemetry hashidoc
// 2. retrieve metric unit and type from telemetry page
// 3. retrieve the metric all values by name
// 4. perform conversion to readable format (time/bytes)
// 5. columnize the results mapping timestamp to values
func (b *Debug) GetMetricValues(name string, validate, byValue, short bool) (string, error) {
	stringInfo, telemetryInfo, err := GetTelemetryMetrics()
	if err != nil {
		return "", err
	}
	if validate {
		if ok := validateName(name, stringInfo); !ok {
			errString := fmt.Sprintf("[metrics-name-validation] '%s' not a valid telemetry metric name\n  visit: %s for full list of consul telemetry metrics", name, TelemetryURL)
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
	metricData, found := b.Metrics.metricsMap[name]
	if !found {
		// Metric not found in the metrics map ==> nil return
		result = []string{fmt.Sprintf("*\x1f%s\x1f=>\x1fnil\x1fvalue(s)\x1freturned\x1f", name)}
		output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
		return output, nil
	}

	// Prepare dataMaps for each metric type
	var dataMaps [][]map[string]interface{}
	var label []string
	// Iterate over metric data and construct result
	for _, data := range metricData {
		timestamp := data["timestamp"].(string)
		mValue := data["value"]
		mLabels := data["labels"].(map[string]interface{})
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
	if len(label) > 0 {
		result[2] += "Labels\x1f"
	}
	if name == "consul.runtime.total_gc_pause_ns" {
		result[2] += "gc/min\x1f"
		// Calculate the GC rate and add it to each line
		for i := 0; i < len(dataMaps); i++ {
			if i == 0 {
				// no previous time stamped gc gauge value to perform non-neg diff calc
				result[i+3] = fmt.Sprintf("%s%s\x1f", result[i+3], "-")
			} else {
				// Calculate the rate using your CalculateGCRate function or another method
				rate, err := CalculateGCRate(dataMaps[i], dataMaps[i-1])
				if err != nil {
					return "", fmt.Errorf("error calculating rate: %v", err)
				}
				// Append the calculated rate to the line
				result[i+3] = fmt.Sprintf("%s%s\x1f", result[i+3], rate)
			}
		}
	}
	if byValue {
		sort.Sort(ByValue(result[3:]))
	}
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, nil
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

// MetricValueExtractor is an interface for extracting metric values.
type MetricValueExtractor interface {
	ExtractMetricValueByName(metricName string) []map[string]interface{}
}

func (m *Metrics) BuildMetricsMap() {
	m.metricsMap = make(map[string][]map[string]interface{})

	for _, metric := range m.Metrics {
		timestamp := metric.Timestamp

		// Iterate over Gauges and add to the map
		for _, gauge := range metric.Gauges {
			metricData := map[string]interface{}{
				"name":      gauge.Name,
				"timestamp": timestamp,
				"value":     gauge.Value,
				"labels":    gauge.Labels,
			}
			m.metricsMap[gauge.Name] = append(m.metricsMap[gauge.Name], metricData)
		}

		// Iterate over Points and add to the map
		for _, point := range metric.Points {
			metricData := map[string]interface{}{
				"name":      point.Name,
				"timestamp": timestamp,
				"value":     point.Points,
				"labels":    point.Labels,
			}
			m.metricsMap[point.Name] = append(m.metricsMap[point.Name], metricData)
		}

		// Iterate over Counters and add to the map
		for _, counter := range metric.Counters {
			metricData := map[string]interface{}{
				"name":      counter.Name,
				"timestamp": timestamp,
				"value":     counter.Count,
				"labels":    counter.Labels,
			}
			m.metricsMap[counter.Name] = append(m.metricsMap[counter.Name], metricData)
		}

		// Iterate over Samples and add to the map
		for _, sample := range metric.Samples {
			metricData := map[string]interface{}{
				"name":      sample.Name,
				"timestamp": timestamp,
				"value":     sample.Mean,
				"labels":    sample.Labels,
			}
			m.metricsMap[sample.Name] = append(m.metricsMap[sample.Name], metricData)
		}
	}
}

// ExtractMetricValueByName extracts metric values by name.
func (m Metric) ExtractMetricValueByName(metricName string) []map[string]interface{} {
	var matches []map[string]interface{}
	regex := regexp.MustCompile(".*" + metricName)

	// Loop through Gauges and extract matching metrics
	for _, gauge := range m.Gauges {
		if regex.MatchString(gauge.Name) {
			match := map[string]interface{}{
				"name":      gauge.Name,
				"value":     gauge.Value,
				"labels":    gauge.Labels,
				"timestamp": m.Timestamp,
			}
			matches = append(matches, match)
		}
	}

	// Loop through Points and extract matching metrics
	for _, point := range m.Points {
		if regex.MatchString(point.Name) {
			match := map[string]interface{}{
				"name":      point.Name,
				"value":     point.Points,
				"labels":    point.Labels,
				"timestamp": m.Timestamp,
			}
			matches = append(matches, match)
		}
	}

	// Loop through Counters and extract matching metrics
	for _, counter := range m.Counters {
		if regex.MatchString(counter.Name) {
			match := map[string]interface{}{
				"name":      counter.Name,
				"value":     counter.Count,
				"labels":    counter.Labels,
				"timestamp": m.Timestamp,
			}
			matches = append(matches, match)
		}
	}

	// Loop through Samples and extract matching metrics
	for _, sample := range m.Samples {
		if regex.MatchString(sample.Name) {
			match := map[string]interface{}{
				"name":      sample.Name,
				"value":     sample.Mean,
				"labels":    sample.Labels,
				"timestamp": m.Timestamp,
			}
			matches = append(matches, match)
		}
	}

	return matches
}

func (b *Debug) Summary() string {
	title := "Metrics Bundle Summary"
	ul := strings.Repeat("-", len(title))
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
		len(b.Metrics.Metrics),
		b.Metrics.Metrics[0].Timestamp,
		b.Metrics.Metrics[len(b.Metrics.Metrics)-1].Timestamp)
}
