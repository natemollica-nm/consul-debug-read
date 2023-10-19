package types

import "regexp"

type Gauge struct {
	Name   string   `json:"Name"`
	Value  float64  `json:"Value"`
	Labels struct{} `json:"Labels"`
}

type Points struct {
	Name   string   `json:"Name"`
	Points float64  `json:"Points"`
	Labels struct{} `json:"Labels"`
}

type Counters struct {
	Name   string   `json:"Name"`
	Count  int      `json:"Count"`
	Rate   float64  `json:"Rate"`
	Sum    float64  `json:"Sum"`
	Min    float64  `json:"Min"`
	Max    float64  `json:"Max"`
	Mean   float64  `json:"Mean"`
	Stddev float64  `json:"Stddev"`
	Labels struct{} `json:"Labels"`
}

type Samples struct {
	Name   string   `json:"Name"`
	Count  int      `json:"Count"`
	Rate   float64  `json:"Rate"`
	Sum    float64  `json:"Sum"`
	Min    float64  `json:"Min"`
	Max    float64  `json:"Max"`
	Mean   float64  `json:"Mean"`
	Stddev float64  `json:"Stddev"`
	Labels struct{} `json:"Labels"`
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

// ExtractMetricValueByName: Interface implementation for MetricValueExtractor
func (m Metric) ExtractMetricValueByName(metricName string) interface{} {
	regex := regexp.MustCompile(".*" + metricName)
	for _, gauge := range m.Gauges {
		if regex.MatchString(gauge.Name) {
			return gauge.Value
		}
	}
	for _, point := range m.Points {
		if regex.MatchString(point.Name) {
			return point.Points
		}
	}
	for _, counter := range m.Counters {
		if regex.MatchString(counter.Name) {
			return counter.Count
		}
	}
	for _, sample := range m.Samples {
		if regex.MatchString(sample.Name) {
			return sample.Count
		}
	}
	// Return nil or an appropriate value if the metric is not found
	return nil
}
