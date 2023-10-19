package cmd

import (
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
				log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", debugPath)
			}
		} else {
			debugPath = viper.GetString("debugPath")
			log.Printf("using config.yaml debug path setting - %s\n", debugPath)
		}
		if ok := debugPath != ""; ok {
			log.Printf("debug-path:  '%s'\n", debugPath)
			err := debugBundle.DecodeJSON(debugPath)
			if err != nil {
				return fmt.Errorf("failed to decode bundle: %v", err)
			}
			log.Printf("Successfully read-in bundle from:  '%s'\n", debugPath)
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
		log.Printf("compiled latest raft configuration (source node/dc: %s/%s)\n", debugBundle.Agent.Config.NodeName, debugBundle.Agent.Config.Datacenter)
		fmt.Println(raftConfiguration)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(raftConfigurationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// raftConfigurationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// raftConfigurationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
