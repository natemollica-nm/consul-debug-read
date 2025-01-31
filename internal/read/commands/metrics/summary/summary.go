package summary

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

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

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
	hclog.L().Debug("reading in agent.json", "filepath", path)
	if err = data.DecodeJSON(path, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in host.json", "filepath", path)
	if err = data.DecodeJSON(path, "host"); err != nil {
		hclog.L().Error("failed to decode host.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in index.json", "filepath", path)
	if err = data.DecodeJSON(path, "index"); err != nil {
		hclog.L().Error("failed to decode index.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in metrics.json", "filepath", path)
	if err = data.DecodeJSON(path, "metrics"); err != nil {
		hclog.L().Error("failed to decode metrics.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in bundle contents")

	var result string
	result = data.Summary()
	c.ui.Output(result)
	return 0
}

const synopsis = `Returns metrics-specific information in summarized format`
const help = `
Usage: 
    consul-debug-read metrics summary [options]
`
