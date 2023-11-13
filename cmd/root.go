package cmd

import (
	bundle "consul-debug-read/lib/types"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const (
	envDebugPath = "CONSUL_DEBUG_PATH"
)

// rootCmd represents the base command when called without any subcommands
var (
	debugPath   string
	debugBundle bundle.Debug
	verbose     bool
	rootCmd     = &cobra.Command{
		Use:   "consul-debug-read",
		Short: "A simple CLI tool for parsing a Consul agent debug bundle",
		Long: `consul-debug-read cli tool

The tool is designed to aid in quickly parsing key metrics,
agent, and consul host information from a 'consul debug' cmd bundle capture.
`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := os.LookupEnv(envDebugPath); ok {
				debugPath = os.Getenv(envDebugPath)
				if verbose {
					log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", debugPath)
				}

			} else {
				debugPath = viper.GetString("debugPath")
				if verbose {
					log.Printf("using config.yaml debug path setting - %s\n", debugPath)
				}
			}
			if err := cmd.Help(); err != nil {
				return err
			}
			fmt.Printf("\n  ==> current debug-path: '%s'\n", debugPath)
			return nil
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
	initConfig()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	// Load configuration from a file (config.yaml)
	viper.SetConfigName(".consul-debug-read") // Set the name of the configuration file (config.yaml, config.json, etc.)
	viper.AddConfigPath(home)                 // Search for the configuration file in the home directory
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("No config file found. A new one will be created.")
		err := viper.SafeWriteConfigAs(filepath.Join(home, ".consul-debug-read.yaml"))
		if err != nil {
			return
		}
	}
	viper.SetEnvPrefix("CONSUL_DEBUG")
	viper.AutomaticEnv()
}
