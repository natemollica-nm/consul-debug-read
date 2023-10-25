package cmd

import (
	funcs "consul-debug-read/lib"
	"consul-debug-read/lib/types"
	"consul-debug-read/metrics"
	"fmt"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Ingest metrics.json from consul debug bundle",
	Long: `Read metrics information from specified bundle and return timestamped values.

Example usage:
	$ consul-debug-read metrics

	$ consul-debug-read metrics --name <name_of_metric>

	$ consul-debug-read metrics --list 

	$ consul-debug-read metrics --gauges
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		l, _ := cmd.Flags().GetBool("list")
		if _, ok := os.LookupEnv(envDebugPath); ok {
			envPath := os.Getenv(envDebugPath)
			envPath = strings.TrimSuffix(envPath, "/")
			if _, err := os.Stat(envPath); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exists: %s - %v\n", envPath, err)
			} else {
				debugPath = envPath
				log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", debugPath)
			}
		} else {
			debugPath = viper.GetString("debugPath")
			log.Printf("using config.yaml debug path setting - %s\n", debugPath)
		}
		if debugPath != "" {
			if !l {
				log.Printf("debug-path:  '%s'\n", debugPath)
				err := debugBundle.DecodeJSON(debugPath)
				if err != nil {
					return fmt.Errorf("failed to decode bundle: %v", err)
				}
				log.Printf("successfully read-in bundle from:  '%s'\n", debugPath)
			}
		} else {
			return fmt.Errorf("debug-path is null")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, _ := cmd.Flags().GetBool("summary")
		l, _ := cmd.Flags().GetBool("list")
		n, _ := cmd.Flags().GetString("name")
		h, _ := cmd.Flags().GetBool("host")
		telegraf, _ := cmd.Flags().GetBool("telegraf")

		buildMetricsData := func() (types.Metrics, types.MetricsIndex, types.Host, types.ByteConverter, string, string) {
			// Get Metrics object
			m := debugBundle.Metrics
			index := debugBundle.Index
			host := debugBundle.Host
			conv := types.ByteConverter{}
			metricsFile := fmt.Sprintf(debugPath + "/metrics.json")
			hostFile := fmt.Sprintf(debugPath + "/host.json")
			return m, index, host, conv, metricsFile, hostFile
		}

		showSummary := func() error {
			m, index, host, _, metricsFile, _ := buildMetricsData()
			fmt.Printf("\nMetrics Bundle Summary: %s\n", metricsFile)
			fmt.Println("----------------------------------------------")
			fmt.Println("Hostname:", host.Host.Hostname)
			fmt.Println("Agent Version:", index.AgentVersion)
			fmt.Println("Raft State:", debugBundle.Agent.Stats.Raft.State)
			fmt.Println("Interval:", index.Interval)
			fmt.Println("Duration:", index.Duration)
			fmt.Println("Capture Targets:", index.Targets)
			fmt.Println("Total Captures:", len(m.Metrics))
			fmt.Printf("Capture Time Start: %s\n", m.Metrics[0].Timestamp)
			fmt.Printf("Capture Time Stop: %s\n", m.Metrics[len(m.Metrics)-1].Timestamp)
			return nil
		}

		// metrics --name <metric_name>
		// 1. if no --skip-name-validation flag passed, validate metric name with telemetry hashidoc
		// 2. retrieve metric unit and type from telemetry page
		// 3. retrieve the metric value by name, and aggregate the results
		// 4. perform conversion to readable format (time/bytes)
		// 5. columnize the results mapping timestamp to values
		metricValueByName := func(name string) error {
			m, _, _, conv, _, _ := buildMetricsData()
			if skip, _ := cmd.Flags().GetBool("skip-name-validation"); skip {
				log.Printf("=> skipping metric name validation with hashicorp docs")
			} else {
				if err := metrics.ValidateMetricName(name); err != nil {
					return err
				}
			}

			var telemetryInfo []metrics.AgentTelemetryMetric
			_, telemetryInfo, err = metrics.GetTelemetryMetrics()
			if err != nil {
				return err
			}
			unit, metricType := types.GetUnitAndType(name, telemetryInfo)
			timeReg := regexp.MustCompile("ns|ms|seconds|hours")
			bytesReg := regexp.MustCompile("bytes")

			result := []string{"Timestamp\x1fMetric\x1fType\x1fUnit\x1fValue\x1fLabels\x1f"}
			for _, metric := range m.Metrics {
				data := metric.ExtractMetricValueByName(name)
				for _, info := range data {
					m_name := info["name"].(string)
					m_value := info["value"]
					m_labels := info["labels"].(map[string]interface{})
					var label []string
					for k, v := range m_labels {
						label = append(label, fmt.Sprintf("{%s: %v}", k, v))
					}
					if m_value != nil {
						var v string
						if timeReg.MatchString(unit) {
							v, err = types.ConvertToReadableTime(m_value, unit)
							if err != nil {
								return err
							}
						} else if bytesReg.MatchString(unit) {
							v = conv.ConvertToReadableBytes(m_value)
						} else {
							v = fmt.Sprintf("%v", m_value)
						}
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							metric.Timestamp, m_name, metricType, unit, v, label))
					} else {
						result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f%s\x1f",
							metric.Timestamp, m_name, metricType, unit, "<nil>", "-"))
					}
				}
			}
			output, err := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
			if err != nil {
				return err
			}
			fmt.Printf("\n%s\n", output)
			return nil
		}

		hostMetrics := func() {
			_, _, host, conv, _, hostFile := buildMetricsData()
			bootTimeStamp := time.Unix(int64(host.Host.BootTime), 0)
			bootTime := bootTimeStamp.Format("2006-01-02 15:04:05 MST")
			upTime := funcs.ConvertSecondsReadable(host.Host.Uptime)
			fmt.Printf("\nHost Metrics Summary: %s\n", hostFile)
			fmt.Println("----------------------")
			fmt.Println("OS:", host.Host.Os)
			fmt.Println("Host Name", host.Host.Hostname)
			fmt.Println("Architecture:", host.Host.KernelArch)
			fmt.Println("Number of Cores:", len(host.CPU))
			fmt.Println("CPU Vendor ID:", host.CPU[0].VendorID)
			fmt.Println("CPU Model Name:", host.CPU[0].ModelName)
			fmt.Printf("Platform: %s | %s\n", host.Host.Platform, host.Host.PlatformVersion)
			fmt.Println("Running Since:", bootTime)
			fmt.Println("Uptime at Capture:", upTime)
			fmt.Printf("\nHost Memory Metrics Summary:\n")
			fmt.Println("----------------------")
			fmt.Printf("Used: %s  (%.2f%%)\n", conv.ConvertToReadableBytes(host.Memory.Used), host.Memory.UsedPercent)
			fmt.Println("Total Available:", conv.ConvertToReadableBytes(host.Memory.Available))
			fmt.Println("Total:", conv.ConvertToReadableBytes(host.Memory.Total))
			fmt.Printf("\nHost Disk Metrics Summary:\n")
			fmt.Println("----------------------")
			fmt.Printf("Used: %s  (%.2f%%)\n", conv.ConvertToReadableBytes(host.Disk.Used), host.Disk.UsedPercent)
			fmt.Println("Free:", conv.ConvertToReadableBytes(host.Disk.Free))
			fmt.Println("Total:", conv.ConvertToReadableBytes(host.Disk.Total))
		}

		switch {
		case summary:
			err := showSummary()
			if err != nil {
				return err
			}
		case l:
			if err := metrics.ListMetrics(); err != nil {
				return err
			}
		case n != "":
			err := metricValueByName(n)
			if err != nil {
				return err
			}
		case h:
			hostMetrics()
		case telegraf:
			err := debugBundle.GenerateTelegrafMetrics()
			if err != nil {
				return err
			}
		default:
			err := showSummary()
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)
	metricsCmd.Flags().Bool("summary", false, "Retrieve metrics summary info from bundle.")
	metricsCmd.Flags().Bool("host", false, "Retrieve Host specific metrics.")
	metricsCmd.Flags().BoolP("list", "l", false, "List available metric names to parse with by name.")
	metricsCmd.Flags().StringP("name", "n", "", "Retrieve specific metric timestamped values by name.")
	metricsCmd.Flags().Bool("skip-name-validation", false, "Skip metric name validation with hashicorp docs.")
	metricsCmd.Flags().Bool("telegraf", false, "Generate telegraf compatible metrics file for ingesting offline metrics.")
}
