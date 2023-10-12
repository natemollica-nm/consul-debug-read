package cmd

import (
	mFuncs "consul-debug-read/metrics"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Ingest metrics.json from consul debug bundle",
	Long: `Read metrics information from specified bundle and return timestamped values.

Example usage:
	
	consul-debug-read metrics -d consul-debug-2023-10-04T18-29-47Z.tar.gz `,
	Run: func(cmd *cobra.Command, args []string) {
		var isExtractPath = false
		a, _ := cmd.Flags().GetBool("all")
		l, _ := cmd.Flags().GetBool("list")
		g, _ := cmd.Flags().GetBool("gauges")
		p, _ := cmd.Flags().GetBool("points")
		c, _ := cmd.Flags().GetBool("counters")
		s, _ := cmd.Flags().GetBool("samples")

		if l {
			if err := mFuncs.ListMetrics(); err != nil {
				fmt.Printf("failed to retrieve list of metrics: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// Read-in debug path
		debugPath, _ := cmd.Flags().GetString("debug-file-path")
		extractPath, _ := cmd.Flags().GetString("use-extract-bundle")

		if debugPath == "" && extractPath != "" {
			isExtractPath = true
			debugPath = extractPath
			debugPath = strings.TrimSuffix(debugPath, "/")
		} else {
			debugPath = strings.TrimSuffix(debugPath, "/")
			err := mFuncs.SelectAndExtractTarGzFilesInDir(debugPath, debugPath)
			if err != nil {
				fmt.Printf("failed to extract bundle: %v", err)
				os.Exit(1)
			}
		}
		// Ingest metrics.json
		m, err := mFuncs.ImportMetrics(debugPath, isExtractPath)
		if err != nil {
			fmt.Printf("failed to parse metrics: %v", err)
		}
		// Interpret metrics specific flags
		n, _ := cmd.Flags().GetString("name")

		if n != "" {
			if err := mFuncs.ValidateMetricName(n); err != nil {
				fmt.Printf("metric name validation failed: %v\n", err)
				os.Exit(1)
			}
			for _, metric := range m.Metrics {
				conv := mFuncs.ByteConverter{}
				value := metric.ExtractMetricValueByName(n)
				if value != nil {
					fmt.Printf("%s '%s': %v\n", metric.Timestamp, n, conv.ConvertToReadableBytes(value))
				} else {
					fmt.Printf("%s '%s': nil value returned\n", metric.Timestamp, n)
				}
			}
			return
		}

		switch {
		case a:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Gauges:", len(metric.Gauges))
				fmt.Println("Number of Points:", len(metric.Points))
				fmt.Println("Number of Counters:", len(metric.Counters))
				fmt.Println("Number of Samples:", len(metric.Samples))
			}
		case g:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Gauges:", len(metric.Gauges))
			}
		case p:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Points:", len(metric.Points))
			}
		case c:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Counters:", len(metric.Counters))
			}
		case s:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Samples:", len(metric.Samples))
			}
		default:
			for _, metric := range m.Metrics {
				fmt.Println("Timestamp:", metric.Timestamp)
				fmt.Println("Number of Gauges:", len(metric.Gauges))
				fmt.Println("Number of Points:", len(metric.Points))
				fmt.Println("Number of Counters:", len(metric.Counters))
				fmt.Println("Number of Samples:", len(metric.Samples))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metricsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	metricsCmd.Flags().BoolP("all", "a", false, "Retrieve metrics summary info from bundle.")
	metricsCmd.Flags().BoolP("gauges", "g", false, "Retrieve Gauges metrics summary info only.")
	metricsCmd.Flags().BoolP("points", "p", false, "Retrieve Points metrics summary info only.")
	metricsCmd.Flags().BoolP("counters", "c", false, "Retrieve Counters metrics summary info only.")
	metricsCmd.Flags().BoolP("samples", "s", false, "Retrieve Samples metrics summary info only.")
	metricsCmd.Flags().BoolP("list", "l", false, "List available metric names to parse with by name.")
	metricsCmd.Flags().StringP("name", "n", "", "Retrieve specific metric timestamped values by name.")
}
