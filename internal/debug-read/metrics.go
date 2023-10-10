package debug_read

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Metric struct {
	Timestamp string `json:"Timestamp"`
	Gauges    []struct {
		Name   string   `json:"Name"`
		Value  float64  `json:"Value"`
		Labels struct{} `json:"Labels"`
	} `json:"Gauges"`
	Points []struct {
		Name   string
		Points []float32
	} `json:"Points"`
	Counters []struct {
		Name   string   `json:"Name"`
		Count  int      `json:"Count"`
		Rate   float64  `json:"Rate"`
		Sum    float64  `json:"Sum"`
		Min    float64  `json:"Min"`
		Max    float64  `json:"Max"`
		Mean   float64  `json:"Mean"`
		Stddev float64  `json:"Stddev"`
		Labels struct{} `json:"Labels"`
	} `json:"Counters"`
	Samples []struct {
		Name   string   `json:"Name"`
		Count  int      `json:"Count"`
		Rate   float64  `json:"Rate"`
		Sum    float64  `json:"Sum"`
		Min    float64  `json:"Min"`
		Max    float64  `json:"Max"`
		Mean   float64  `json:"Mean"`
		Stddev float64  `json:"Stddev"`
		Labels struct{} `json:"Labels"`
	} `json:"Samples"`
}

func GetMetrics(inputFile string) ([]Metric, int) {
	// Check if the input file flag is provided
	if inputFile == "" {
		log.Fatalf("Please provide an input JSON file using the -file flag.")
	}

	jsonData, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}
	// defer closing input file until main() is complete.
	defer func(jsonData *os.File) {
		err := jsonData.Close()
		if err != nil {

		}
	}(jsonData)

	// create multi-metric object for metrics.json ingestion
	var metrics []Metric

	// open json decode for serialized json data
	decoder := json.NewDecoder(jsonData)
	for {
		// Parse the JSON data into a variable Metrics struct
		var metric Metric

		err := decoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}
		// add metric objects to overall metrics object
		metrics = append(metrics, metric)
	}
	return metrics, 0
}

func (c *) Run() {
	err := GetMetrics(runConfig *Config) ()
}
