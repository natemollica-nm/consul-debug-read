package types

import (
	"consul-debug-read/metrics"
	telemetry "consul-debug-read/metrics"
	"fmt"
	"github.com/ryanuber/columnize"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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

type ByValue []string

func (m ByValue) Len() int      { return len(m) }
func (m ByValue) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByValue) Less(i, j int) bool {
	columns_i := strings.Split(m[i], "\x1f")
	columns_j := strings.Split(m[j], "\x1f")
	var value_i, value_j float64
	if len(columns_i) >= 2 && len(columns_i) <= 4 {
		value_i, _ = strconv.ParseFloat(strings.TrimRight(columns_i[1], "%"), 64)
		value_j, _ = strconv.ParseFloat(strings.TrimRight(columns_j[1], "%"), 64)
	} else {
		value_i, _ = strconv.ParseFloat(strings.TrimRight(columns_i[4], "%"), 64)
		value_j, _ = strconv.ParseFloat(strings.TrimRight(columns_j[4], "%"), 64)
	}

	// using '>' vice '<' to sort from highest -> lowest
	return value_i > value_j
}

// GetMetricValues
// 1. Retrieves value of metric by passed in name string
// 2. (if applicable) Sorts metric dataset by value (highest-to-lowest) vice the default timestamp order

// GetMetricValues
// 1. if no --skip-name-validation flag passed, validate metric name with telemetry hashidoc
// 2. retrieve metric unit and type from telemetry page
// 3. retrieve the metric all values by name
// 4. perform conversion to readable format (time/bytes)
// 5. columnize the results mapping timestamp to values
func (b *Debug) GetMetricValues(name string, validate, byValue, short bool) (string, error) {
	result := []string{fmt.Sprintf("\x1f%s\x1f", name)}
	underline := fmt.Sprintf(strings.Repeat("-", len(name)))
	result = append(result, fmt.Sprintf("\x1f%s\x1f", underline))
	if short {
		result = append(result, "Timestamp\x1fValue\x1f")

	} else {
		result = append(result, "Timestamp\x1fMetric\x1fType\x1fUnit\x1fValue\x1f")
	}
	timeReg := regexp.MustCompile("^ns$|^ms$|^seconds$|^hours$")
	bytesReg := regexp.MustCompile("bytes")
	percentageReg := regexp.MustCompile("percentage")

	stringInfo, telemetryInfo, err := telemetry.GetTelemetryMetrics()
	if err != nil {
		return "", err
	}
	validateName := func(n string, info string) error {
		// This metric name is dynamic and can be anything that the customer uses for service names
		reg := regexp.MustCompile(`^consul\.proxy\..+$`)
		if reg.MatchString(n) {
			fmt.Printf("built-in mesh proxy prefix used: %s\n", name)
			return nil
		}
		// list of metrics contains the name somewhere, return with no error
		if strings.Contains(info, n) {
			// log.Printf("[metrics-name-validation]: validated metric name %s\n", name)
			return nil
		}
		return fmt.Errorf(fmt.Sprintf("[metrics-name-validation] '%s' not a valid telemetry metric name\n  visit: %s for full list of consul telemetry metrics", name, telemetry.TelemetryURL))
	}
	if validate {
		if err = validateName(name, stringInfo); err != nil {
			return "", err
		}
	}
	// log.Printf("[metric-value-by-name]: retrieving metric unit and type")
	unit, metricType := GetUnitAndType(name, telemetryInfo)
	var dataMaps [][]map[string]interface{}
	var label []string
	// log.Printf("[metric-value-by-name]: performing readable conversion for value and mapping labels")
	for _, metric := range b.Metrics.Metrics {
		data := metric.ExtractMetricValueByName(name)
		dataMaps = append(dataMaps, data)
		for _, info := range data {
			mName := info["name"].(string)
			mValue := info["value"]
			mLabels := info["labels"].(map[string]interface{})
			mTimestamp := info["timestamp"]
			label = []string{}
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
					vFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", mValue), 64)
					percent := vFloat * 100.00
					v = fmt.Sprintf("%.2f%%", percent)

				} else {
					v = fmt.Sprintf("%v", mValue)
				}
				if short {
					if len(label) > 0 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f",
							mTimestamp, v, label))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f",
							mTimestamp, v))
					}
				} else {
					if len(label) > 0 {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							mTimestamp, mName, metricType, unit, v, label))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							mTimestamp, mName, metricType, unit, v))
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
	if len(result) == 3 {
		result = []string{fmt.Sprintf("*\x1f%s\x1f=>\x1fnil\x1fvalue(s)\x1freturned\x1f", name)}
	}
	if byValue {
		sort.Sort(ByValue(result[3:]))
	}
	// log.Printf("[metric-value-by-name]: finishing processing metric by name. columnizing....")
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	// log.Printf("[metric-value-by-name]: returning values")
	return output, nil
}

// MetricValueExtractor is an interface for extracting metric values by name
type MetricValueExtractor interface {
	ExtractMetricValueByName(metricName string) interface{}
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

// ExtractMetricValueByName ExtractMetricValueByName: Interface implementation for MetricValueExtractor
func (m Metric) ExtractMetricValueByName(metricName string) []map[string]interface{} {
	var matches []map[string]interface{}

	regex := regexp.MustCompile(".*" + metricName)

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

// nonNegativeDifference calculates the non-negative difference between two float64 values.
func nonNegativeDifference(a, b float64) float64 {
	diff := a - b
	if diff >= 0 {
		return diff
	}
	return 0 // Return the absolute value of the difference if < 0
}

// CalculateGCRate calculates the rate of Garbage Collection (GC) in nanoseconds per minute.
func CalculateGCRate(currentValueData, previousValueData []map[string]interface{}) (string, error) {
	var rate string

	currentValue, ok := currentValueData[0]["value"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid 'value' field in data")
	}
	previousValue, ok := previousValueData[0]["value"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid 'value' field in data")
	}
	// Calculate the non-negative difference in GC pause times
	diff := nonNegativeDifference(currentValue, previousValue)
	timeCurrent, err := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", currentValueData[0]["timestamp"]))
	if err != nil {
		return "", err
	}
	timePrevious, err := time.Parse("2006-01-02 15:04:05 -0700 MST", fmt.Sprintf("%s", previousValueData[0]["timestamp"]))
	if err != nil {
		return "", err
	}
	// consul debug caputures default to 5m/30s capture intervals (>= v1.16.x)
	//
	timeDiff := timeCurrent.Sub(timePrevious).Seconds()
	if diff >= 0 && timeDiff > 0 {
		rate, err = ConvertToReadableTime(diff/(timeDiff/60), "ns") // convert to ns/min to most-readable-time/minute
		if err != nil {
			return "", err
		}
		rate = fmt.Sprintf("%s/min", rate)
	}
	if rate == "" {
		rate = "-"
	}
	return rate, nil
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
			return fmt.Sprintf("%.2fms", v/float64(nsInMs)), nil
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
