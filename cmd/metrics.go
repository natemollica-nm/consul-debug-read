package cmd

import (
	funcs "consul-debug-read/lib"
	"consul-debug-read/lib/types"
	"consul-debug-read/metrics"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

// metricsCmd represents the metrics command
var (
	keyMetricNames = map[string][]string{
		"Transaction Timing":               transactionTimingMetrics,
		"Leadership Stability":             leaderShipMetrics,
		"Certificate Authority Expiration": certAuthority,
		"Autopilot":                        autoPilot,
		"Memory Utilization":               memory,
		"Network Activity":                 network,
		"Raft Thread Saturation":           raftThreadSaturation,
		"Raft Replication Capacity":        replicationCapacity,
		"BoltDB Performance":               boldDBPerformance,
	}
	transactionTimingMetrics = []string{
		"consul.kvs.apply",
		"consul.txn.apply",
		"consul.raft.apply",
		"consul.raft.commitTime",
	}
	leaderShipMetrics = []string{
		"consul.raft.leader.lastContact",
		"consul.raft.state.candidate",
		"consul.raft.state.leader",
		"consul.server.isLeader",
	}
	certAuthority = []string{
		"consul.mesh.active-root-ca.expiry",
		"consul.mesh.active-signing-ca.expiry",
		"consul.agent.tls.cert.expiry",
	}
	autoPilot = []string{
		"consul.autopilot.healthy",
		"consul.autopilot.failure_tolerance",
	}
	memory = []string{
		"consul.runtime.alloc_bytes",
		"consul.runtime.sys_bytes",
		"consul.runtime.total_gc_pause_ns",
	}
	network = []string{
		"consul.client.rpc",
		"consul.client.rpc.exceeded",
		"consul.client.rpc.failed",
	}
	raftThreadSaturation = []string{
		"consul.raft.thread.main.saturation",
		"consul.raft.thread.fsm.saturation",
	}
	replicationCapacity = []string{
		"consul.raft.fsm.lastRestoreDuration",
		"consul.raft.leader.oldestLogAge",
		"consul.raft.rpc.installSnapshot",
	}
	boldDBPerformance = []string{
		"consul.raft.boltdb.freelistBytes",
		"consul.raft.boltdb.logsPerBatch",
		"consul.raft.boltdb.storeLogs",
		"consul.raft.boltdb.writeCapacity",
	}
	rateLimitingMetrics = []string{
		"consul.client.rpc",
		"consul.client.rpc.failed",
		"consul.client.rpc.exceeded",
		"consul.rpc.rate_limit.exceeded",
		"consul.rpc.rate_limit.log_dropped",
	}
	metricsCmd = &cobra.Command{
		Use:   "metrics",
		Short: "Ingest metrics.json from consul debug bundle",
		Long: `Read metrics information from specified bundle and return timestamped values.

Example usage:
	Display summary of bundle capture	
		$ consul-debug-read metrics

	Display full list of queryable metric names
		$ consul-debug-read metrics --list 
	
	Retrieve all timestamped captures of metric
		$ consul-debug-read metrics --name <name_of_metric>
	
	Sort metric capture by value (highest to lowest)
		$ consul-debug-read metrics --name <name_of_metric> --sort-by-value
	
	Skip hashidoc metric name validation:
		$ consul-debug-read metrics --name <valid_name_but_not_in_docs> --validate-metric-name=false
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
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
				l, _ := cmd.Flags().GetBool("list")
				listTransactionTiming, _ := cmd.Flags().GetBool("list-transaction-timing")
				h, _ := cmd.Flags().GetBool("host")
				if !l && !listTransactionTiming {
					log.Printf("debug-path:  '%s'\n", debugPath)
					// don't read in metrics.json if we don't have to
					if h {
						if err := debugBundle.DecodeJSON(debugPath, "host"); err != nil {
							return fmt.Errorf("failed to decode bundle: %v", err)

						}
						log.Printf("successfully read-in host.json from: '%s'\n", debugPath)
					} else {
						if err := debugBundle.DecodeJSON(debugPath, "agent"); err != nil {
							log.Fatalf("failed to decode 'agent.json' - %v", err)
						}
						if err := debugBundle.DecodeJSON(debugPath, "host"); err != nil {
							log.Fatalf("failed to decode 'host.json' - %v", err)

						}
						if err := debugBundle.DecodeJSON(debugPath, "index"); err != nil {
							log.Fatalf("failed to decode 'index.json' - %v", err)
						}
						if err := debugBundle.DecodeJSON(debugPath, "metrics"); err != nil {
							log.Fatalf("failed to decode 'metrics.json' - %v", err)
						}
						log.Printf("successfully read-in agent, host, metrics, and index from:  '%s'\n", debugPath)
					}
				}
			} else {
				return fmt.Errorf("debug-path is null")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// no metrics.json ingestion required
			l, _ := cmd.Flags().GetBool("list")
			listTransactionTiming, _ := cmd.Flags().GetBool("list-transaction-timing")
			keyMetrics, _ := cmd.Flags().GetBool("get-key-metrics")
			rateLimiting, _ := cmd.Flags().GetBool("rate-limiting")
			transactionTiming, _ := cmd.Flags().GetBool("transaction-timing")
			leadershipChanges, _ := cmd.Flags().GetBool("leadership-changes")
			h, _ := cmd.Flags().GetBool("host")

			// requires metrics.json
			n, _ := cmd.Flags().GetString("name")
			summary, _ := cmd.Flags().GetBool("summary")
			validateName, _ := cmd.Flags().GetBool("verify")
			byValue, _ := cmd.Flags().GetBool("sort-by-value")
			short, _ := cmd.Flags().GetBool("short")
			telegraf, _ := cmd.Flags().GetBool("telegraf")

			showSummary := func() error {
				m := debugBundle.Metrics
				index := debugBundle.Index
				host := debugBundle.Host
				metricsFile := fmt.Sprintf(debugPath + "/metrics.json")
				fmt.Printf("\nMetrics Bundle Summary: %s\n", metricsFile)
				fmt.Println("----------------------------------------------")
				fmt.Println("Datacenter:", debugBundle.Agent.Config.Datacenter)
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
			hostMetrics := func() {
				host := debugBundle.Host
				conv := types.ByteConverter{}
				hostFile := fmt.Sprintf(debugPath + "/host.json")
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
				if err := showSummary(); err != nil {
					return err
				}
			case l:
				if err := metrics.ListMetrics(true, false); err != nil {
					return err
				}
			case listTransactionTiming:
				if err := metrics.ListMetrics(false, true); err != nil {
					return err
				}
			case n != "":
				values, err := debugBundle.GetMetricValues(n, validateName, byValue, short)
				if err != nil {
					return err
				}
				fmt.Println(values)
			case keyMetrics:
				funcs.ClearScreen()
				for keyMetricTitle, row := range keyMetricNames {
					fmt.Printf("\n[Next Key Metric] => %s: press [ENTER] to retrieve values", keyMetricTitle)
					fmt.Scanln()
					funcs.ClearScreen()
					for _, name := range row {
						doneCh := make(chan bool)
						go funcs.Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
						values, err := debugBundle.GetMetricValues(name, false, byValue, short)
						if err != nil {
							return err
						}
						doneCh <- true // Stop the dot printing goroutine
						close(doneCh)
						fmt.Printf("\n%s\n", values)
					}
				}
			case rateLimiting:
				funcs.ClearScreen()
				for _, name := range rateLimitingMetrics {
					fmt.Printf("\n[Rate Limiting Metrics] => %s: press [ENTER] to retrieve values", name)
					fmt.Scanln()
					doneCh := make(chan bool)
					go funcs.Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
					values, err := debugBundle.GetMetricValues(name, false, byValue, short)
					if err != nil {
						return err
					}
					doneCh <- true // Stop the dot printing goroutine
					close(doneCh)
					fmt.Printf("%s\n", values)
				}
			case transactionTiming:
				funcs.ClearScreen()
				for _, name := range transactionTimingMetrics {
					fmt.Printf("\n[Transaction Timing Metrics] => %s: press [ENTER] to retrieve values", name)
					fmt.Scanln()
					doneCh := make(chan bool)
					go funcs.Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
					values, err := debugBundle.GetMetricValues(name, false, byValue, short)
					if err != nil {
						return err
					}
					doneCh <- true // Stop the dot printing goroutine
					close(doneCh)
					fmt.Printf("%s\n", values)
				}
			case leadershipChanges:
				funcs.ClearScreen()
				for _, name := range leaderShipMetrics {
					fmt.Printf("\n[Transaction Timing Metrics] => %s: press [ENTER] to retrieve values", name)
					fmt.Scanln()
					doneCh := make(chan bool)
					go funcs.Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
					values, err := debugBundle.GetMetricValues(name, false, byValue, short)
					if err != nil {
						return err
					}
					doneCh <- true // Stop the dot printing goroutine
					close(doneCh)
					fmt.Printf("%s\n", values)
				}
			case h:
				hostMetrics()
			case telegraf:
				if err := debugBundle.GenerateTelegrafMetrics(); err != nil {
					return err
				}
			default:
				if err := showSummary(); err != nil {
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(metricsCmd)
	metricsCmd.Flags().Bool("summary", false, "Retrieve metrics summary info from bundle.")
	metricsCmd.Flags().Bool("host", false, "Retrieve Host specific metrics.")
	metricsCmd.Flags().BoolP("list", "l", false, "List available metric names to parse with by name.")
	metricsCmd.Flags().Bool("list-transaction-timing", false, "List key metrics for Consul txn and kv transaction timing.")
	metricsCmd.Flags().Bool("get-key-metrics", false, "Retrieve key metric values for Consul from debug bundle.")
	metricsCmd.Flags().Bool("rate-limiting", false, "Retrieve key rate limit metric values for Consul from debug bundle.")
	metricsCmd.Flags().Bool("transaction-timing", false, "Retrieve key transaction timing metric values for Consul from debug bundle.")
	metricsCmd.Flags().Bool("leadership-changes", false, "Retrieve key raft leadership stability metric values for Consul from debug bundle.")
	metricsCmd.Flags().StringP("name", "n", "", "Retrieve specific metric timestamped values by name.")
	metricsCmd.Flags().Bool("sort-by-value", false, "Parse metric value by name and sort results by value vice timestamp order.")
	metricsCmd.Flags().Bool("short", false, "Only print timestamp, value, and labels columns.")
	metricsCmd.Flags().Bool("verify", true, "Performs metric name validation with hashicorp docs.")
	metricsCmd.Flags().Bool("telegraf", false, "Generate telegraf compatible metrics file for ingesting offline metrics.")
}
