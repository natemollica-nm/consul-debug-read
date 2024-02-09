package parse

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"consul-debug-read/internal/read/log"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"os"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	sort    bool
	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.BoolVar(&c.sort, "sort", false, "Parse metric value by name and sort results by value vice timestamp order")
	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

func (c *cmd) Synopsis() string { return synopsis }

func (c *cmd) Run(args []string) int {
	var result []log.LogEntry
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
	cmdYamlCfg, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		hclog.L().Error("error reading consul-debug-read user config file", "filepath", read.DebugReadConfigFullPath, "error", err)
		return 1
	}
	var cfg read.ReaderConfig
	err = yaml.Unmarshal(cmdYamlCfg, &cfg)
	if err != nil {
		hclog.L().Error("error deserializing YAML contents", "filepath", read.DebugReadConfigFullPath, "error", err)
		return 1
	}
	if cfg.DebugDirectoryPath == "" {
		hclog.L().Error("empty or null consul-debug-path setting", "error", read.DebugReadConfigFullPath)
		return 1
	}
	logFile := cfg.DebugDirectoryPath + "/consul.log"
	result, err = log.ParseLogFile(logFile, "")
	if err != nil {
		hclog.L().Error("error parsing log file", "file", logFile, "error", err)
		return 1
	}
	counts := log.AggregateEntries(result)
	out := log.RPCCounts(counts)
	c.ui.Output(out)
	return 0
}

const synopsis = `Log parser for debug bundle captured consul.log file`
const help = `
Requires:
    - Valid consul monitor or consul log file with '.log' extension
    - TRACE level capture enabled on agent's log or monitor
      - agent cmd: -log-level=trace
      - agent conf: log_level="trace"
      - monitor cmd: consul monitor -log-level=trace

Description:
    Parses consul trace logs for all (default) or specified ([method]) rpc method calls and provides
        => Rate-per-minute count of rpc call(s) sorted from highest to lowest
        => Total log capture count of rpc call(s) sorted from highest to lowest

Usage: 
    consul-debug-read log [options]


Options:
  [method]: Specify an RPC method to filter results (e.g., 'Catalog.NodeServiceList') - Optional
`
