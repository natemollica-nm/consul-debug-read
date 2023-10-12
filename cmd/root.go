package cmd

/*
Copyright © 2023 NAME HERE nathan.mollica@hashicorp.com
*/

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "consul-debug-read",
	Short: "A simple CLI tool for parsing a Consul agent debug bundle",
	Long: `consul-debug-read cli tool

The tool is designed to aid in quickly parsing key metrics,
agent, and consul host information from a 'consul debug' cmd bundle capture.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
	},
}

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

	rootCmd.PersistentFlags().StringP("debug-file-path", "d", "", "File path to directory containing consul-debug tar.gz bundle.")
	rootCmd.PersistentFlags().StringP("use-extract-bundle", "f", "", "Path to already extracted debug bundle's root directory.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringP("debug-file-path", "p", "", "Path to extracted debug bundle.")
}
