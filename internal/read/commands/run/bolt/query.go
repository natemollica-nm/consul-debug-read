package bolt

import (
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"consul-debug-read/internal/read/commands/run"
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

	return 0
}

func runBoltQuery(name string) (string, error) {
	values, err := run.RunTimeBundle.GetMetricValuesMemDB(name, true, false, false)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve memDB values error=%v", err)
	}
	return values, nil
}

const synopsis = "Runs consul-debug-read boltDB query when running as long-running app"
const help = `
Usage: consul-debug-read run bolt-query [options]

  Queries boltDB backend when running consul-debug-read as background process.
`
