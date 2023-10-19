package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// membersCmd represents the members command
var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Parses members.json and formats to typical 'consul members -wan' output",
	Long: `Templates the 'standardOutput()' function from the 'consul members' command' and 
ingests and parses <debug_path>/members.json for useful output"'. 

For example:
	consul-debug-read agent members -d bundles/consul-debug-2023-10-04T18-29-47Z
.`,
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
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("compiled wan membership list (source node/dc: %s/%s)\n", debugBundle.Agent.Config.NodeName, debugBundle.Agent.Config.Datacenter)
		membersOutput := debugBundle.MembersStandard()
		fmt.Print(membersOutput)
	},
}

func init() {
	agentCmd.AddCommand(membersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	membersCmd.Flags().Bool("wan", false, "Retrieve agent members summary for agent's wan fed members.")
}
