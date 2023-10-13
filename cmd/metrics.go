package cmd

import (
	funcs "consul-debug-read/lib"
	mFuncs "consul-debug-read/metrics"
	"fmt"
	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Ingest metrics.json from consul debug bundle",
	Long: `Read metrics information from specified bundle and return timestamped values.

Example usage:
	
	consul-debug-read metrics -x -d consul-debug-2023-10-04T18-29-47Z.tar.gz `,
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, _ := cmd.Flags().GetBool("summary")
		l, _ := cmd.Flags().GetBool("list")
		g, _ := cmd.Flags().GetBool("gauges")
		p, _ := cmd.Flags().GetBool("points")
		c, _ := cmd.Flags().GetBool("counters")
		s, _ := cmd.Flags().GetBool("samples")

		// If list called just get and list available metrics and return
		if l {
			if err := mFuncs.ListMetrics(); err != nil {
				return err
			}
			return nil
		}

		// Get Metrics object
		m := debugBundle.Metrics
		metricsFile := fmt.Sprintf(debugPath + "/metrics.json")

		// Interpret metrics specific flags
		n, _ := cmd.Flags().GetString("name")
		if n != "" {
			if err := mFuncs.ValidateMetricName(n); err != nil {
				return err
			}
			for _, metric := range m.Metrics {
				conv := funcs.ByteConverter{}
				value := metric.ExtractMetricValueByName(n)
				if value != nil {
					fmt.Printf("%s '%s': %v\n", metric.Timestamp, n, conv.ConvertToReadableBytes(value))
				} else {
					fmt.Printf("%s '%s': nil value returned\n", metric.Timestamp, n)
				}
			}
			return nil
		}

		switch {
		case summary:
			start := m.Metrics[0].Timestamp
			stop := m.Metrics[len(m.Metrics)-1].Timestamp

			fmt.Printf("\nMetrics Bundle Summary: %s\n", metricsFile)
			fmt.Println("----------------------")
			fmt.Println("Start:", start)
			fmt.Println("Stop:", stop)
			fmt.Printf("Duration: %v\n", funcs.TimeStampDuration(start, stop))
			fmt.Println("Number of Captures:", len(m.Metrics))
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
		return nil
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
	metricsCmd.Flags().Bool("summary", false, "Retrieve metrics summary info from bundle.")
	metricsCmd.Flags().BoolP("gauges", "g", false, "Retrieve Gauges metrics summary info only.")
	metricsCmd.Flags().BoolP("points", "p", false, "Retrieve Points metrics summary info only.")
	metricsCmd.Flags().BoolP("counters", "c", false, "Retrieve Counters metrics summary info only.")
	metricsCmd.Flags().BoolP("samples", "s", false, "Retrieve Samples metrics summary info only.")
	metricsCmd.Flags().BoolP("list", "l", false, "List available metric names to parse with by name.")
	metricsCmd.Flags().StringP("name", "n", "", "Retrieve specific metric timestamped values by name.")
}
