package summary

import (
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/config/get"
	"consul-debug-read/internal/read/commands/flags"
	"consul-debug-read/internal/read/log"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"time"
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

	logFile := path + "/consul.log"
	var entries []log.LogEntry
	var out string

	loggingSummary := map[string]string{
		log.ErrorLevel: "",
		log.WarnLevel:  "",
		log.DebugLevel: "",
		log.TraceLevel: "",
	}
	for k, _ := range loggingSummary {
		entries, err = log.ParseLog(logFile, k, "", time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		counts := log.AggregateLogEntries(entries, k, log.MessageSelect)
		out = log.FormatCounts(counts, "message")
		loggingSummary[k] = out
	}
	c.ui.Output(loggingSummary[log.ErrorLevel])
	c.ui.Output(loggingSummary[log.WarnLevel])
	c.ui.Output(loggingSummary[log.DebugLevel])
	c.ui.Output(loggingSummary[log.TraceLevel])
	return 0
}

const synopsis = `Returns log-specific data points of interest overview`
const help = `
Usage: 
    consul-debug-read log summary [options]
`
