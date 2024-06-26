package raft

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/config/get"
	"consul-debug-read/internal/read/commands/flags"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(raftCommandHelp, c.flags) }

func (c *cmd) Synopsis() string { return synopsis }

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		c.ui.Error(fmt.Sprintf("Failed to parse flags: %v", err))
		return 1
	}
	if c.verbose && c.silent {
		c.ui.Error(fmt.Sprintf("Cannot specify both -silent and -verbose"))
		return 1
	}

	level := hclog.Info
	if c.verbose {
		level = hclog.Debug
	} else if c.silent {
		level = hclog.Off
	}
	var result string

	commands.InitLogging(c.ui, level)

	var ok bool
	var err error
	var path string
	if path, ok = get.RenderPathFromConfig(); !ok {
		hclog.L().Error("error rendering debug filepath", "filepath", path, "error", err)
		return 1
	}

	var data read.Debug
	if err := data.DecodeJSON(path, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	if err := data.DecodeJSON(path, "members"); err != nil {
		hclog.L().Error("failed to decode members.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in agent cmd information from bundle")
	hclog.L().Debug("compiling raft configuration from agent.json and members.json")
	result, err = data.RaftListPeers()
	if err != nil {
		hclog.L().Error("failed to retrieve raft list peers from debug bundle", "error", err)
		return 1
	}
	c.ui.Output(result)
	return 0
}

const synopsis = `Retrieve agent's latest raft configuration summary'`
const raftCommandHelp = `Parses latest raft-configuration from Agent.json capture within bundle.

Example usage and output:
$ consul-debug-read raft-configuration
Node                      ID                                   Address           State    Voter AppliedIndex CommitIndex
consul-i-0aa97949095868769 666e152f-7316-81aa-848b-3f4719564404 10.2.101.211:8300 follower true  -            -
consul-i-0ba0dff4180ec2dc7 4f36f7ab-240a-61a6-c5e1-b78ce62813a2 10.2.4.230:8300   follower true  -            -
consul-i-08e67d882fe525809 1baa8d56-a9ae-adf7-5309-b12460c3e6c5 10.2.64.253:8300  follower true  -            -
consul-i-05a474f75fea384bb 263fd5e5-fbd7-90b1-a904-4ab3c53b74f7 10.2.17.109:8300  leader   true  2801009780   2801009780
consul-i-06033dd57876bf1a7 eca79896-dad9-1713-94a2-c2b35a37d7df 10.2.4.89:8300    follower true  -            -`
