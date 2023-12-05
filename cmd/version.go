package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "0.0.3"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of consul-debug-read",
		Long:  `All software has versions. This is consul-debug-read's'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("consul-debug-read: v%s", version)
			return nil
		},
	}
)
