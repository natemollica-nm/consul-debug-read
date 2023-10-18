package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// raftConfigurationCmd represents the raftConfiguration command
var raftConfigurationCmd = &cobra.Command{
	Use:   "raft-configuration",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		raftConfiguration, err := debugBundle.RaftListPeers()
		if err != nil {
			return err
		}
		fmt.Println(raftConfiguration)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(raftConfigurationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// raftConfigurationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// raftConfigurationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
