package cmd

import (
	"consul-debug-read/cmd/config"
	funcs "consul-debug-read/lib"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// setDebugPathCmd represents the setDebugPath command
var (
	debugFile       string
	err             error
	setDebugPathCmd = &cobra.Command{
		Use:   "set-debug-path",
		Short: "Changes which bundle you're focusing on for analysis",
		Long: `consul-debug-read set-debug-path [options]

Validates the path contents or extracts a valid .tar.gz bundle and points to this valid directory path for processing.

--path can be either:
  * consul-debug extracted contents (valid agent.json, metrics.json, host.json, and index.json) or
  * path to multiple bundles available for extraction and path setting

--file can be either:
  * path to valid consul debug .tar.gz archive or
  * path to multiple bundles available for extraction and path setting

Example (--path):
	$ consul-debug-read set-debug-path --path bundles/consul-debug-2023-10-04T18-29-47Z

Example (--path) for dir containing multiple .tar.gz bundles:
	$ consul-debug-read set-debug-path --path bundles

	select a .tar.gz file to extract:
	1: 124722consul-debug-2023-10-04T18-29-47Z.tar.gz
	2: 124722consul-debug-2023-10-11T17-33-55Z.tar.gz
	3: 124722consul-debug-2023-10-11T17-43-15Z.tar.gz
	4: 124722consul-debug-eu-01-stag.tar.gz
	5: 124722consul-debug-eu-133-stag-default.tar.gz
	6: 124722consul-debug-us-135-stag-default.tar.gz
	7: 124722consul-debug-us-east-stag.tar.gz
	enter the number of the file to extract: 

Example (--file) for extraction:
	$ consul-debug-read set-debug-path --file bundles/124722consul-debug-2023-10-11T17-43-15Z.tar.gz
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if debugFile, err = cmd.Flags().GetString("file"); err != nil {
				return fmt.Errorf("[set-debug-path] failed to retrieve --file flag %v\n", err)
			}
			if debugPath, err = cmd.Flags().GetString("path"); err != nil {
				return fmt.Errorf("[set-debug-path] failed to retrieve --path flag %v\n", err)
			}

			if debugFile != "" {
				if isValidBundle := strings.HasSuffix(debugFile, ".tar.gz"); isValidBundle {
					if config.Verbose {
						log.Printf("[set-debug-path] file passed in for extraction: %s\n", debugFile)
					}
					if debugPath, err = funcs.SelectAndExtractTarGzFilesInDir(debugFile); err != nil {
						return fmt.Errorf("[set-debug-path] failed to extract and select debug bundle %v\n", err)
					}
				} else if fileInfo, err := os.Stat(debugFile); err == nil {
					if fileInfo.IsDir() {
						if debugPath, err = funcs.SelectAndExtractTarGzFilesInDir(debugFile); err != nil {
							return fmt.Errorf("[set-debug-path] failed to extract and select debug bundle %v\n", err)
						}
					}
				} else {
					return fmt.Errorf("[set-debug-path] invalid debug file format passed in with --file flag - must be .tar.gz")
				}
			}
			if debugPath != "" {
				if isFile := strings.HasSuffix(debugPath, ".tar.gz"); isFile {
					return fmt.Errorf("[set-debug-path] --path used with .tar.gz file, provide path to extracted bundle or use --file to extract bundle and set path")
				}
				// Verify contents to path are contain valid debug contents
				consulDebugFiles, err := os.ReadDir(debugPath)
				if err != nil {
					return fmt.Errorf("[set-debug-path] failed to list files in directory path %s - %v\n", debugPath, err)
				}
				var metricsJson, agentJson, membersJson, hostJson, indexJson bool
				for _, file := range consulDebugFiles {
					switch file.Name() {
					case "metrics.json":
						metricsJson = true
					case "agent.json":
						agentJson = true
					case "host.json":
						hostJson = true
					case "index.json":
						indexJson = true
					case "members.json":
						membersJson = true
					case "cluster.json":
						clusterJsonPath := debugPath + "/" + file.Name()
						membersJsonPath := debugPath + "/members.json"
						if err := os.Rename(clusterJsonPath, membersJsonPath); err != nil {
							return fmt.Errorf("[set-debug-path] error renaming file %s to members.json: %v\n", clusterJsonPath, err)
						} else {
							if config.Verbose {
								log.Printf("[set-debug-path] renamed %s to members.json\n", clusterJsonPath)
							}
						}
						membersJson = true
					}
				}
				if agentJson && membersJson && hostJson && indexJson {
					// older debug bundles separated v1/agent/metrics captures into each interval
					// if so, try and merge the metrics.json files into one large metrics.json
					// for ingestion.
					if !metricsJson {
						// "metrics.json" not found in the current directory
						// Run the "merge-metrics.sh" script with debugPath as an argument
						scriptPath := "scripts/merge-metrics.sh"
						cmd := exec.Command(scriptPath, debugPath)
						if config.Verbose {
							log.Printf("attempting to merge sub-directory metrics.json files (older debug capture bundle)...")
						}
						if _, err := cmd.CombinedOutput(); err != nil {
							return fmt.Errorf("error running %s: %v\n", scriptPath, err)
						} else {
							if config.Verbose {
								log.Printf("ran %s to merge %s sub-directory metrics.json files\n", scriptPath, debugPath)
							}
						}
					}
					if config.Verbose {
						log.Printf("[set-debug-path] path contents validated!\n")
					}
				} else if fileInfo, err := os.Stat(debugPath); err == nil {
					// if path is directory and no index.json, then it's probably a path to a dir containing multiple .tar.gz files
					if fileInfo.IsDir() && !indexJson {
						if debugPath, err = funcs.SelectAndExtractTarGzFilesInDir(debugPath); err != nil {
							return fmt.Errorf("[set-debug-path] failed to extract and select debug bundle - %v\n", err)
						}
					}
				} else {
					return fmt.Errorf("[set-debug-path] invalid path request - %v\n", err)
				}
			}

			debugPath = strings.TrimSuffix(debugPath, "/")
			if _, err := os.Stat(debugPath); os.IsNotExist(err) {
				return fmt.Errorf("[set-debug-path] directory does not exists: %s - %v\n", debugPath, err)
			} else {
				viper.Set("debugPath", debugPath)
			}
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("[set-debug-path] failed to write the configuration file -- %v\n", err)
			} else {
				if config.Verbose {
					log.Printf("[set-debug-path] config.yaml consul-debug-read debug-path has been set => %s\n", debugPath)
					log.Printf("[set-debug-path][WARN] 'CONSUL_DEBUG_PATH' env var will take precedence over %s."+
						"Unset this variable or override it's value using 'unset' or 'export' as necessary.", viper.ConfigFileUsed())
				}
			}
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("[set-debug-path] config.yaml not found or error reading config: %v\n", err)
			}
			fmt.Println("set consul-debug-read path successful =>", debugPath)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(setDebugPathCmd)
	setDebugPathCmd.Flags().StringVarP(&debugPath, "path", "d", "", "File path to directory containing consul-debug.tar.gz bundle(s).")
	setDebugPathCmd.Flags().StringVarP(&debugFile, "file", "f", "", "File path to single consul-debug.tar.gz bundle.")
	setDebugPathCmd.MarkFlagsMutuallyExclusive("path", "file")
}
