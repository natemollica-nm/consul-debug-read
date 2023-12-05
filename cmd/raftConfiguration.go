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

// raftConfigurationCmd represents the raftConfiguration command
var raftConfigurationCmd = &cobra.Command{
	Use:   "raft-configuration",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		if ok := debugPath != ""; ok {
			if config.Verbose {
				log.Printf("debug-path:  '%s'\n", debugPath)
			}
			if err := debugBundle.DecodeJSON(debugPath, "agent"); err != nil {
				return fmt.Errorf("failed to decode bundle: %v", err)
			}
			if err := debugBundle.DecodeJSON(debugPath, "members"); err != nil {
				return fmt.Errorf("failed to decode bundle: %v", err)
			}
			if config.Verbose {
				log.Printf("Successfully read-in agent and members from:  '%s'\n", debugPath)
			}
		} else {
			return fmt.Errorf("debug-path is null")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		raftConfiguration, err := debugBundle.RaftListPeers()
		if err != nil {
			return err
		}
		if config.Verbose {
			log.Printf("compiled latest raft configuration (source node/dc: %s/%s)\n", debugBundle.Agent.Config.NodeName, debugBundle.Agent.Config.Datacenter)
		}
		fmt.Println(raftConfiguration)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(raftConfigurationCmd)
}
