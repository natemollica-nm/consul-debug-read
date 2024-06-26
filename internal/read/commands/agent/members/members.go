package members

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

	silent  bool
	verbose bool
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

func (c *cmd) Help() string { return commands.Usage(cmdHelp, c.flags) }

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

	commands.InitLogging(c.ui, level)

	var ok bool
	var err error
	var path string
	if path, ok = get.RenderPathFromConfig(); !ok {
		hclog.L().Error("error rendering debug filepath", "filepath", path, "error", err)
		return 1
	}

	var data read.Debug
	if err = data.DecodeJSON(path, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	if err = data.DecodeJSON(path, "members"); err != nil {
		hclog.L().Error("failed to decode members.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in agent cmd information from bundle")

	result := agentMembers(data.Agent)
	c.ui.Output(result)
	return 0
}

func agentMembers(agent read.Agent) string {
	return agent.MembersStandard()
}

const synopsis = "Parses members.json and formats to typical 'consul members -wan' output"
const cmdHelp = `Templates the 'standardOutput()' function from the 'consul members' command' and 
ingests and parses <debug_path>/members.json for useful output"'. 

For example:
	consul-debug-read agent members -d bundles/consul-debug-2023-10-04T18-29-47Z`
