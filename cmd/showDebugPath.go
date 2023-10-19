package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// showDebugPathCmd represents the showDebugPath command
var showDebugPathCmd = &cobra.Command{
	Use:   "show-debug-path",
	Short: "Show currently configured extracted debug bundle filepath",
	Long: `Shows currently set consul-debug-read command debug path as set in
config.yaml viper configuration file. 

To change file-path use consul-debug-read set-debug-path --path <path_to_debug_bundle> to alter.

Example:
	$ consul-debug-read show-debug-path
	  
	2023/10/19 08:50:22 show-debug-path: debug-path => 'bundles/consul-debug-2023-10-04T18-29-47Z'
`,
	RunE: func(cmd *cobra.Command, args []string) error {
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
		log.Printf("debug-path => '%s'\n", debugPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showDebugPathCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showDebugPathCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showDebugPathCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
