package types

import (
	"consul-debug-read/metrics"
	"fmt"
	"reflect"
	"regexp"
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

type Metrics struct{ Metrics []Metric }

type MetricsIndex struct {
	Version      int      `json:"Version"`
	AgentVersion string   `json:"AgentVersion"`
	Interval     string   `json:"Interval"`
	Duration     string   `json:"Duration"`
	Targets      []string `json:"Targets"`
}

// MetricValueExtractor is an interface for extracting metric values by name
type MetricValueExtractor interface {
	ExtractMetricValueByName(metricName string) interface{}
}

// ExtractMetricValueByName ExtractMetricValueByName: Interface implementation for MetricValueExtractor
func (m Metric) ExtractMetricValueByName(metricName string) []map[string]interface{} {
	var matches []map[string]interface{}

	regex := regexp.MustCompile(".*" + metricName)

	for _, gauge := range m.Gauges {
		if regex.MatchString(gauge.Name) {
			match := map[string]interface{}{
				"name":   gauge.Name,
				"value":  gauge.Value,
				"labels": gauge.Labels,
			}
			matches = append(matches, match)
		}
	}
	for _, point := range m.Points {
		if regex.MatchString(point.Name) {
			match := map[string]interface{}{
				"name":   point.Name,
				"value":  point.Points,
				"labels": point.Labels,
			}
			matches = append(matches, match)
		}
	}
	for _, counter := range m.Counters {
		if regex.MatchString(counter.Name) {
			match := map[string]interface{}{
				"name":   counter.Name,
				"value":  counter.Count,
				"labels": counter.Labels,
			}
			matches = append(matches, match)
		}
	}
	for _, sample := range m.Samples {
		if regex.MatchString(sample.Name) {
			match := map[string]interface{}{
				"name":   sample.Name,
				"value":  sample.Mean,
				"labels": sample.Labels,
			}
			matches = append(matches, match)
		}
	}
	if len(matches) > 0 {
		return matches
	} else {
		return nil
	}
}

// GetUnitAndType returns the Unit and Type for a given Name.
func GetUnitAndType(name string, telemetry []metrics.AgentTelemetryMetric) (string, string) {
	for _, metric := range telemetry {
		if metric.Name == name {
			return metric.Unit, metric.Type
		} else if name == "*" {
			return metric.Unit, metric.Type
		}
	}
	return "-", "-"
}

// ByteConverter
// Struct used to implement the ConvertToReadableBytes interface function for int and float64
// byte conversion.
type ByteConverter struct{}

func (bc ByteConverter) ConvertToReadableBytes(value interface{}) string {
	switch v := value.(type) {
	case int:
		return ConvertIntBytes(v)
	case float64:
		return ConvertFloatBytes(v)
	default:
		return "Unsupported type"
	}
}

func ConvertIntBytes(bytes int) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
		tb = 1024 * gb
	)

	switch {
	case bytes >= tb:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(tb))
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func ConvertFloatBytes(bytes float64) string {
	const (
		kb = 1024
		mb = 1024 * kb
		gb = 1024 * mb
		tb = 1024 * gb
	)

	switch {
	case bytes >= tb:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(tb))
	case bytes >= gb:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%.4f bytes", bytes)
	}
}

// TimeConverter is the interface for converting time units.
type TimeConverter interface {
	Convert(timeValue interface{}) (string, error)
}

func ConvertToReadableTime(value interface{}, units string) (string, error) {
	var converter TimeConverter

	switch units {
	case "ns":
		converter = NanosecondsConverter{}
	case "ms":
		converter = MillisecondsConverter{}
	case "seconds":
		converter = SecondsConverter{}
	case "hours":
		converter = HoursConverter{}
	}
	v, err := converter.Convert(value)
	if err != nil {
		return "", err
	}
	return v, nil
}

// NanosecondsConverter implements TimeConverter for nanoseconds.
type NanosecondsConverter struct{}

func (n NanosecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		nsInMs     = 1e6
		nsInSecond = 1e9
		nsInHour   = 3.6e12
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case v >= nsInHour:
			return fmt.Sprintf("%.2fh", float64(v)/float64(nsInHour)), nil
		case v >= nsInSecond:
			return fmt.Sprintf("%.2fs", float64(v)/float64(nsInSecond)), nil
		case v >= nsInMs:
			return fmt.Sprintf("%.2fms", float64(v)/float64(nsInMs)), nil
		default:
			return fmt.Sprintf("%dns", v), nil
		}
	case float64:
		switch {
		case v >= nsInHour:
			return fmt.Sprintf("%.2fh", v/float64(nsInHour)), nil
		case v >= nsInSecond:
			return fmt.Sprintf("%.2fs", v/float64(nsInSecond)), nil
		case v >= nsInMs:
			return fmt.Sprintf("%.4fms", v/float64(nsInMs)), nil
		default:
			return fmt.Sprintf("%.4fns", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// MillisecondsConverter implements TimeConverter for milliseconds.
type MillisecondsConverter struct{}

func (m MillisecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		msInSecond = 1e3
		msInHour   = 3.6e6
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case v >= msInHour:
			return fmt.Sprintf("%.2fh", float64(v)/float64(msInHour)), nil
		case v >= msInSecond:
			return fmt.Sprintf("%.2fs", float64(v)/float64(msInSecond)), nil
		default:
			return fmt.Sprintf("%.4fms", float64(v)), nil
		}
	case float64:
		switch {
		case v >= msInHour:
			return fmt.Sprintf("%.2fh", v/float64(msInHour)), nil
		case v >= msInSecond:
			return fmt.Sprintf("%.2fs", v/float64(msInSecond)), nil
		default:
			return fmt.Sprintf("%.4fms", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// SecondsConverter implements TimeConverter for seconds.
type SecondsConverter struct{}

func (s SecondsConverter) Convert(timeValue interface{}) (string, error) {
	const (
		secondsInHour = 3600
	)

	switch v := timeValue.(type) {
	case int:
		switch {
		case v >= secondsInHour:
			return fmt.Sprintf("%.2fh", float64(v)/float64(secondsInHour)), nil

		default:
			return fmt.Sprintf("%.2fs", float64(v)), nil
		}
	case float64:
		switch {
		case v >= secondsInHour:
			return fmt.Sprintf("%.2fh", v/float64(secondsInHour)), nil
		default:
			return fmt.Sprintf("%.2fs", v), nil
		}
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}

// HoursConverter implements TimeConverter for hours.
type HoursConverter struct{}

func (h HoursConverter) Convert(timeValue interface{}) (string, error) {
	switch v := timeValue.(type) {
	case int:
		return fmt.Sprintf("%.2fh", float64(v)), nil
	case float64:
		return fmt.Sprintf("%.2fh", v), nil
	default:
		return "", fmt.Errorf("unsupported type: %v", reflect.TypeOf(timeValue))
	}
}
