package cmd

import (
	funcs "consul-debug-read/lib"
	bundle "consul-debug-read/lib/types"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var (
	debugPath   string
	debugFile   string
	debugBundle bundle.Debug
	rootCmd     = &cobra.Command{
		Use:   "consul-debug-read",
		Short: "A simple CLI tool for parsing a Consul agent debug bundle",
		Long: `consul-debug-read cli tool

The tool is designed to aid in quickly parsing key metrics,
agent, and consul host information from a 'consul debug' cmd bundle capture.
`, PersistentPreRun: func(cmd *cobra.Command, args []string) {
			extract, _ := cmd.Flags().GetBool("extract")
			if extract {
				debugPath, _ = funcs.SelectAndExtractTarGzFilesInDir(debugPath)
			} else {
				debugPath = strings.TrimSuffix(debugPath, "/")
			}
			if debugPath != "" {
				err := debugBundle.DecodeJSON(debugPath)
				if err != nil {
					fmt.Printf("failed to decode bundle: %v", err)
					os.Exit(1)
				}
				log.Printf("Successfully read-in bundle from:  '%s'\n\n", debugPath)
			} else {
				err := cmd.Help()
				if err != nil {
					return
				}
			}
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&debugPath, "debug-path", "d", "", "File path to directory containing consul-debug.tar.gz bundle(s).")
	rootCmd.PersistentFlags().StringVarP(&debugFile, "debug-file", "f", "", "File path to single consul-debug.tar.gz bundle.")
	rootCmd.MarkFlagsMutuallyExclusive("debug-path", "debug-file")
	rootCmd.PersistentFlags().BoolP("extract", "x", false, "Flag indicating bundle requires extraction.")
}
