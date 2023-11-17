package cmd

import (
	"fmt"
	"github.com/spf13/viper"
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
	  
	bundles/consul-debug-2023-10-04T18-29-47Z
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, ok := os.LookupEnv(envDebugPath); ok {
			envPath := os.Getenv(envDebugPath)
			envPath = strings.TrimSuffix(envPath, "/")
			if _, err := os.Stat(envPath); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exists: %s - %v\n", envPath, err)
			} else {
				debugPath = envPath
			}
		} else {
			debugPath = viper.GetString("debugPath")
		}
		fmt.Println(debugPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showDebugPathCmd)
}
