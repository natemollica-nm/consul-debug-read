package cmd

import (
	funcs "consul-debug-read/lib"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// setDebugPathCmd represents the setDebugPath command
var (
	debugFile       string
	err             error
	setDebugPathCmd = &cobra.Command{
		Use:   "set-debug-path",
		Short: "Set the consul debug extracted bundle path for processing",
		Long: `consul-debug-read set-debug-path [options]

To change which bundle you're focusing on for analysis run this command to set
the path to the desired extracted debug bundle directory. Simply hard-sets the config.yaml
to point the debugPath: entry to the desired location. This is so the command doesn't require
passing the --debug-path flag upon every run.

Example:
	$ consul-debug-read set-debug-path --path ./bundles/consul-debug-2023-10-04T18-29-47Z
      debug path set to: bundles/consul-debug-2023-10-04T18-29-47Z

	$ consul-debug-read agent members

Extraction:
	$ consul-debug-read set-debug-path --file ./bundles/124722consul-debug-2023-10-04T18-29-47Z.tar.gz
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if debugFile, err = cmd.Flags().GetString("file"); err != nil {
				return fmt.Errorf("failed to retrieve --file flag %v\n", err)
			}
			if debugPath, err = cmd.Flags().GetString("path"); err != nil {
				return fmt.Errorf("failed to retrieve --path flag %v\n", err)
			}
			if debugFile != "" {
				if isValidBundle := strings.HasSuffix(debugFile, ".tar.gz"); isValidBundle {
					log.Printf("file passed in for extraction: %s\n", debugFile)
					if debugPath, err = funcs.SelectAndExtractTarGzFilesInDir(debugFile); err != nil {
						return fmt.Errorf("set-debug-path: failed to extract and select debug bundle %v\n", err)
					}
				} else {
					return fmt.Errorf("invalid debug file format passed in with --file flag - must be .tar.gz")
				}
			} else if debugPath != "" {
				if isFile := strings.HasSuffix(debugPath, ".tar.gz"); isFile {
					return fmt.Errorf("--path used with .tar.gz file, provide path to extracted bundle or use --file to extract bundle and set path")
				}
			}

			debugPath = strings.TrimSuffix(debugPath, "/")
			if _, err := os.Stat(debugPath); os.IsNotExist(err) {
				return fmt.Errorf("directory does not exists: %s - %v\n", debugPath, err)
			} else {
				viper.Set("debugPath", debugPath)
			}

			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to write the configuration file: %v\n", err)
			} else {
				log.Printf("config.yaml consul-debug-read debug-path has been set => %s\n", debugPath)
				log.Printf("[WARN] CONSUL_DEBUG_PATH env var will take precedence over this if set. Unset this variable or override it's value using 'unset' or 'export' as necessary.")
			}
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("config.yaml not found or error reading config: %v\n", err)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(setDebugPathCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setDebugPathCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setDebugPathCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setDebugPathCmd.Flags().StringVarP(&debugPath, "path", "d", "", "File path to directory containing consul-debug.tar.gz bundle(s).")
	setDebugPathCmd.Flags().StringVarP(&debugFile, "file", "f", "", "File path to single consul-debug.tar.gz bundle.")
	setDebugPathCmd.MarkFlagsMutuallyExclusive("path", "file")
}
