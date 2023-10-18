package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Debug bundle agent.json information parsing.",
	Long: `The agent flag will ingest the agent.json and parse for additional information pertaining to the agent.
This includes:
  - Consul Versioning
  - Server Agent Status
  - Client Agent Status
  - Known Serf Members
  - Current Raft Configuration
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, _ := cmd.Flags().GetBool("summary")

		// Get Metrics object
		var agent = debugBundle.Agent
		var members = debugBundle.Members
		var agentFile = fmt.Sprintf(debugPath + "/agent.json")

		switch {
		case summary:
			fmt.Printf("Agent Configuration Summary: %s\n", agentFile)
			fmt.Println("----------------------")
			agent.AgentSummary()
		default:
			fmt.Printf("Agent Configuration Summary: %s\n", agentFile)
			fmt.Println("----------------------")
			fmt.Println("Serf Member Count (WAN + LAN):", len(members))
			agent.AgentSummary()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().BoolP("summary", "s", false, "Retrieve agent configuration summary.")
}
