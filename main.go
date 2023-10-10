package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	// Define a flag for the input JSON file
	inputFile := flag.String("file", "", "Input JSON file")
	flag.Parse()
	// Check if the input file flag is provided
	if *inputFile == "" {
		log.Fatalf("Please provide an input JSON file using the -file flag.")
	}

	jsonData, err := os.Open(*inputFile)
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

	for _, metric := range metrics {
		fmt.Println("Timestamp:", metric.Timestamp)
		fmt.Println("Number of Gauages:", len(metric.Gauges))
		fmt.Println("Number of Points:", len(metric.Points))
		fmt.Println("Number of Counters:", len(metric.Counters))
		fmt.Println("Number of Samples:", len(metric.Samples))
	}

}
