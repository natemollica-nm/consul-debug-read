package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// membersCmd represents the members command
var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "Parses members.json and formats to typical 'consul members -wan' output",
	Long: `Templates the 'standardOutput()' function from the 'consul members' command' and 
ingests and parses <debug_path>/members.json for useful output"'. 

For example:
	consul-debug-read agent members -d bundles/consul-debug-2023-10-04T18-29-47Z
.`,
	Run: func(cmd *cobra.Command, args []string) {
		wan, _ := cmd.Flags().GetBool("wan")
		if wan {
			fmt.Println("wan called")
		} else {
			membersOutput := debugBundle.MembersStandard()
			fmt.Print(membersOutput)
		}

	},
}

func init() {
	agentCmd.AddCommand(membersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	membersCmd.Flags().Bool("wan", false, "Retrieve agent members summary for agent's wan fed members.")
}
