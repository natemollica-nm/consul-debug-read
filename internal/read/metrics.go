package read

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ryanuber/columnize"
	bolt "go.etcd.io/bbolt"
	"sort"
	"strings"
)

const (
	GaugeType    = "Gauges"
	PointsType   = "Points"
	CountersType = "Counters"
	SamplesType  = "Samples"
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

// GetMetricValues / retrieve metrics captures for specified metrics by name
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
	metricData, found := b.Metrics.MetricsMap[name]
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
		mLabels := data["labels"].(map[string]string)
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

// BoltDBUpload /
// Builds metrics map from the ingested metrics.json,
// extracts metric name, value, labels, and timestamp
// for retrieval via query, and uploads to boltDB.
func (b *Debug) BoltDBUpload() error {
	b.Metrics.MetricsMap = make(map[string][]map[string]interface{})
	var err error
	var bucket, gaugeBucket, pointsBucket, countersBucket, samplesBucket *bolt.Bucket
	var encodedMetric []byte
	err = b.Backend.DB.Update(func(tx *bolt.Tx) error {
		bucket, err = tx.CreateBucketIfNotExists([]byte("Metrics"))
		if err != nil {
			return err
		}
		for i, metric := range b.Metrics.Metrics {
			timestamp := metric.Timestamp
			err = bucket.Put([]byte(fmt.Sprintf("%d", i)), []byte(timestamp))
			if err != nil {
				return err
			}
			gaugeBucket, err = bucket.CreateBucketIfNotExists([]byte("Gauges"))
			if err != nil {
				return err
			}
			pointsBucket, err = bucket.CreateBucketIfNotExists([]byte("Points"))
			if err != nil {
				return err
			}
			countersBucket, err = bucket.CreateBucketIfNotExists([]byte("Counters"))
			if err != nil {
				return err
			}
			samplesBucket, err = bucket.CreateBucketIfNotExists([]byte("Samples"))
			if err != nil {
				return err
			}

			// Iterate over Gauges and add to the map
			for _, gauge := range metric.Gauges {
				metricData := map[string]interface{}{
					"name":      gauge.Name,
					"timestamp": timestamp,
					"value":     gauge.Value,
					"labels":    gauge.Labels,
					"type":      GaugeType,
				}
				b.Metrics.MetricsMap[gauge.Name] = append(b.Metrics.MetricsMap[gauge.Name], metricData)
				// Encode metric for storage
				encodedMetric, err = json.Marshal(metricData)
				if err != nil {
					return err
				}
				key := fmt.Sprintf("%s:%s", gauge.Name, metric.Timestamp) // Adjust as needed
				err = gaugeBucket.Put([]byte(key), encodedMetric)
				if err != nil {
					return err
				}
			}
			// Iterate over Points and add to the map
			for _, point := range metric.Points {
				metricData := map[string]interface{}{
					"name":      point.Name,
					"timestamp": timestamp,
					"value":     point.Points,
					"labels":    point.Labels,
					"type":      PointsType,
				}
				b.Metrics.MetricsMap[point.Name] = append(b.Metrics.MetricsMap[point.Name], metricData)
				// Encode metric for storage
				encodedMetric, err = json.Marshal(metricData)
				if err != nil {
					return err
				}
				key := fmt.Sprintf("%s:%s", point.Name, metric.Timestamp) // Adjust as needed
				err = pointsBucket.Put([]byte(key), encodedMetric)
				if err != nil {
					return err
				}
			}
			// Iterate over Counters and add to the map
			for _, counter := range metric.Counters {
				metricData := map[string]interface{}{
					"name":      counter.Name,
					"timestamp": timestamp,
					"value":     counter.Count,
					"labels":    counter.Labels,
					"type":      CountersType,
				}
				b.Metrics.MetricsMap[counter.Name] = append(b.Metrics.MetricsMap[counter.Name], metricData)
				encodedMetric, err = json.Marshal(metricData)
				if err != nil {
					return err
				}
				key := fmt.Sprintf("%s:%s", counter.Name, metric.Timestamp) // Adjust as needed
				err = countersBucket.Put([]byte(key), encodedMetric)
				if err != nil {
					return err
				}
			}
			// Iterate over Samples and add to the map
			for _, sample := range metric.Samples {
				metricData := map[string]interface{}{
					"name":      sample.Name,
					"timestamp": timestamp,
					"value":     sample.Mean,
					"labels":    sample.Labels,
					"type":      SamplesType,
				}
				b.Metrics.MetricsMap[sample.Name] = append(b.Metrics.MetricsMap[sample.Name], metricData)
				encodedMetric, err = json.Marshal(metricData)
				if err != nil {
					return err
				}
				key := fmt.Sprintf("%s:%s", sample.Name, metric.Timestamp) // Adjust as needed
				err = samplesBucket.Put([]byte(key), encodedMetric)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// GetMetricValuesMemDB / extracts all timestamped occurrences of metric values by name
func (b *Debug) GetMetricValuesMemDB(name string, validate, byValue, short bool) (string, error) {
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
	// Retrieve metric data from BoltDB instead of the metrics map
	info, found := b.Metrics.MetricsMap[name]
	if !found {
		// Metric not found in the metrics map ==> nil return
		result = []string{fmt.Sprintf("*\x1f%s\x1f=>\x1fnil\x1fvalue(s)\x1freturned\x1f", name)}
		output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
		return output, nil
	}
	var metricData []map[string]interface{}
	err := b.Backend.DB.View(func(tx *bolt.Tx) error {
		rootBucket := tx.Bucket([]byte("Metrics"))
		if rootBucket == nil {
			return fmt.Errorf("metrics bucket does not exist")
		}

		// Determine the appropriate bucket for the metric type based on the name
		// This assumes you have a way to map 'name' to the correct metric type
		metricTypeBucketName := info[0]["type"]
		metricTypeBucket := rootBucket.Bucket([]byte(metricTypeBucketName.(string)))
		if metricTypeBucket == nil {
			return fmt.Errorf("%s bucket does not exist", metricTypeBucketName)
		}

		c := metricTypeBucket.Cursor()
		prefix := []byte(name + ":")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var metric map[string]interface{}
			if err := json.Unmarshal(v, &metric); err != nil {
				return err
			}
			metricData = append(metricData, metric)
		}
		return nil
	})
	if err != nil {
		return "", err
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
