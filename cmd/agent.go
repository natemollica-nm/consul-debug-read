package cmd

import (
	"consul-debug-read/cmd/config"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

// agentCmd represents the agent command
var (
	agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "Debug bundle agent.json information parsing.",
		Long: `The agent flag will ingest the agent.json and parse for additional information pertaining to the agent.
This includes:
  - Consul Versioning
  - Server Agent Status
  - Client Agent Status
  - Known Serf Members
  - Current Raft Configuration
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := os.LookupEnv(envDebugPath); ok {
				envPath := os.Getenv(envDebugPath)
				envPath = strings.TrimSuffix(envPath, "/")
				if _, err := os.Stat(envPath); os.IsNotExist(err) {
					return fmt.Errorf("directory does not exists: %s - %v\n", envPath, err)
				} else {
					debugPath = envPath
					if config.Verbose {
						log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", debugPath)
					}
				}
			} else {
				debugPath = viper.GetString("debugPath")
				if config.Verbose {
					log.Printf("using config.yaml debug path setting - %s\n", debugPath)
				}
			}
			if debugPath != "" {
				if config.Verbose {
					log.Printf("debug-path:  '%s'\n", debugPath)
				}
				if err := debugBundle.DecodeJSON(debugPath, "agent"); err != nil {
					return fmt.Errorf("failed to decode agent.json %v", err)
				}
				if err := debugBundle.DecodeJSON(debugPath, "members"); err != nil {
					return fmt.Errorf("failed to decode members.json %v", err)
				}
				if config.Verbose {
					log.Printf("successfully read-in bundle from:  '%s'\n", debugPath)
				}
			} else {
				return fmt.Errorf("debug-path is null")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			summary, _ := cmd.Flags().GetBool("summary")
			c, _ := cmd.Flags().GetBool("config")
			// Get Metrics object
			var agent = debugBundle.Agent
			// var members = debugBundle.Members
			var agentFile = fmt.Sprintf(debugPath + "/agent.json")

			switch {
			case summary:
				if config.Verbose {
					log.Printf("agent summary: configuration rendered from: %s\n", agentFile)
				}
				agent.AgentSummary()
			case c:
				if config.Verbose {
					log.Printf("agent hcl config: configuration rendered from: %s\n", agentFile)
				}
				cfg, err := agent.AgentConfigFull()
				if err != nil {
					return err
				}
				fmt.Println(cfg)
			default:
				//fmt.Printf("Agent Configuration Summary:\n")
				//fmt.Println("----------------------")
				//fmt.Println("Serf Member Count (wan members):", len(members))
				//agent.AgentSummary()
				//fmt.Printf("debug file: %s\n", agentFile)
				if err := cmd.Usage(); err != nil {
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().BoolP("summary", "s", false, "Retrieve agent configuration summary.")
	agentCmd.Flags().Bool("config", false, "Retrieve agent configuration in HCL format")
}
