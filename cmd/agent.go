package cmd

import (
	"consul-debug-read/lib"
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

and much more...
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, _ := cmd.Flags().GetBool("summary")
		m, _ := cmd.Flags().GetBool("members")
		l, _ := cmd.Flags().GetBool("local")
		w, _ := cmd.Flags().GetBool("remote")
		r, _ := cmd.Flags().GetBool("raft")

		// Get Metrics object
		var agent = debugBundle.Agent
		var members = debugBundle.Members
		var agentFile = fmt.Sprintf(debugPath + "/agent.json")
		var membersFile = fmt.Sprintf(debugPath + "/members.json")
		raftConfig := lib.ConvertToValidJSON(agent.Stats.Raft.LatestConfiguration)
		raftConfigFormatted, err := lib.ExecuteJQ(raftConfig, ".")
		if err != nil {
			return err
		}

		switch {
		case summary:
			fmt.Printf("Agent Configuration Summary: %s\n", agentFile)
			fmt.Println("----------------------")
			fmt.Println("Datacenter:", agent.Config.Datacenter)
			fmt.Println("Primary DC:", agent.Config.PrimaryDatacenter)
			fmt.Println("Version:", agent.Config.Version)
			fmt.Println("Server:", agent.Config.Server)
			fmt.Println("NodeName:", agent.Config.NodeName)
			fmt.Println("Serf Member Count (WAN + LAN):", len(members))
			fmt.Printf("Latest Raft Configuration: \n%s", raftConfigFormatted)
		case m:
			fmt.Printf("Membership Summary \nfile: %s\n", membersFile)
			fmt.Println("----------------------")
			fmt.Println("Primary DC:", agent.Config.PrimaryDatacenter)
			fmt.Println("Datacenter:", agent.Config.Datacenter)
			fmt.Println("Member Count:", len(members))
			fmt.Println("----------------------")
			for i, member := range members {
				var status string
				switch {
				case member.Status == 0:
					status = "None"
				case member.Status == 1:
					status = "Alive"
				case member.Status == 2:
					status = "Leaving"
				case member.Status == 3:
					status = "Left"
				case member.Status == 4:
					status = "Failed"
				}
				if l {
					if member.Tags.Dc == agent.Config.Datacenter {
						fmt.Printf("Member Name: %s [%d]\n", member.Name, i+1)
						fmt.Printf("Member Version: %s\n", member.Tags.Build)
						fmt.Printf("Member Datacenter: %s\n", member.Tags.Dc)
						fmt.Printf("Member Address: %s:%d\n", member.Addr, member.Port)
						fmt.Printf("Member Status: %s\n", status)
						fmt.Println("----------------------")
					}
				} else if w {
					if member.Tags.Dc != agent.Config.Datacenter {
						fmt.Printf("Member Name: %s [%d]\n", member.Name, i+1)
						fmt.Printf("Member Version: %s\n", member.Tags.Build)
						fmt.Printf("Member Datacenter: %s\n", member.Tags.Dc)
						fmt.Printf("Member Address: %s:%d\n", member.Addr, member.Port)
						fmt.Printf("Member Status: %s\n", status)
						fmt.Println("----------------------")
					}
				} else {
					fmt.Printf("Member Name: %s [%d]\n", member.Name, i+1)
					fmt.Printf("Member Version: %s\n", member.Tags.Build)
					fmt.Printf("Member Datacenter: %s\n", member.Tags.Dc)
					fmt.Printf("Member Address: %s:%d\n", member.Addr, member.Port)
					fmt.Printf("Member Status: %s\n", status)
					fmt.Println("----------------------")
				}
			}
		case r:
			fmt.Printf("Agent Raft Configuration: %s\n", agentFile)
			fmt.Println("----------------------")
			fmt.Printf("Latest Raft Configuration: \n%s", raftConfigFormatted)

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.Flags().BoolP("summary", "s", false, "Retrieve agent configuration summary.")
	agentCmd.Flags().BoolP("members", "m", false, "Retrieve agent members summary.")
	agentCmd.Flags().Bool("local", false, "Retrieve agent members summary for agent's local datacenter.")
	agentCmd.Flags().Bool("remote", false, "Retrieve agent members summary for agent's local datacenter.")
	agentCmd.Flags().BoolP("raft", "r", false, "Retrieve agent raft configuration summary.")
}
