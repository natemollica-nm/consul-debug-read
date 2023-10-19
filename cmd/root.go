package cmd

import (
	bundle "consul-debug-read/lib/types"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	envDebugPath = "CONSUL_DEBUG_PATH"
)

// rootCmd represents the base command when called without any subcommands
var (
	debugPath   string
	debugBundle bundle.Debug
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
				log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", debugPath)
			} else {
				debugPath = viper.GetString("debugPath")
				log.Printf("using config.yaml debug path setting - %s\n", debugPath)
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
}

func initConfig() {
	// Load configuration from a file (config.yaml)
	viper.SetConfigName("config") // Set the name of the configuration file (config.yaml, config.json, etc.)
	viper.AddConfigPath(".")      // Search for the configuration file in the current directory
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Config file (./config.yaml) not found or error reading config: %s\n", err)
	}
	viper.SetEnvPrefix("CONSUL_DEBUG")
	viper.AutomaticEnv()
}
