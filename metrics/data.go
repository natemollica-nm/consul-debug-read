package metrics

import "fmt"

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
		return fmt.Sprintf("%d bytes", bytes)
	}
}
