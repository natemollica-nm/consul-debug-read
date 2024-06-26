package rpccounts

import (
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/config/get"
	"consul-debug-read/internal/read/commands/flags"
	"consul-debug-read/internal/read/log"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	method string

	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.StringVar(&c.method, "method", "", "Specify a specific RPC method for filtering results (i.e., 'Catalog.NodeServiceList')")
	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

func (c *cmd) Synopsis() string { return synopsis }

func (c *cmd) Run(args []string) int {
	var entries []log.Entry
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

	logFile := path + "/consul.log"
	var out string

	switch {
	case c.method != "":
		entries, err = log.ParseRPCMethods(logFile, c.method)
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		counts := log.AggregateRPCEntries(entries)
		out = log.RPCCounts(counts)
	default:
		entries, err = log.ParseRPCMethods(logFile, "")
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		counts := log.AggregateRPCEntries(entries)
		out = log.RPCCounts(counts)
	}

	c.ui.Output(out)
	return 0
}

const synopsis = `Parses debug bundle log for [TRACE] messages pertaining to RPC Rate Limiting`
const help = `
Usage: 
    consul-debug-read log parse-rpc-counts [options]

Parses consul trace logs for all (default) or specified ([method]) rpc method calls and provides
	=> Rate-per-minute count of rpc call(s) sorted from highest to lowest
	=> Total log capture count of rpc call(s) sorted from highest to lowest

Requires:
    - Valid consul monitor or consul log file with '.log' extension
    - TRACE level capture enabled on agent's log or monitor
      - agent cmd:   '-log-level=trace'
      - agent conf:  'log_level=trace'
      - monitor cmd: 'consul monitor -log-level=trace'`
