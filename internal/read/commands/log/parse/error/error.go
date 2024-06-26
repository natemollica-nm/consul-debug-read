package error

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

	source string

	messageCount bool
	sourceCount  bool

	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.StringVar(&c.source, "source", "", "Capture error messages from specific sources (e.g., \"agent.http\", \"agent.server\")")
	c.flags.BoolVar(&c.messageCount, "message-count", false, "Parse log for error messages and return timestamp sorted list of messages received")
	c.flags.BoolVar(&c.sourceCount, "source-count", false, "Parse log for error messages and return count sorted (descending order) list of messages received from specific sources")

	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

func (c *cmd) Synopsis() string { return synopsis }

func (c *cmd) Run(args []string) int {
	var entries []log.LogEntry
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
	case c.source != "" && !c.sourceCount:
		hclog.L().Debug("parsing debug bundle log file [ERROR] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.ErrorLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		hclog.L().Debug("running general parse collection on debug log for all [ERROR] messages")
		out = log.FormatLog(entries)
	case c.sourceCount:
		hclog.L().Debug("parsing debug bundle log file [ERROR] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.ErrorLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		hclog.L().Debug("aggregating [ERROR] messages by logged entry source type")
		counts := log.AggregateLogEntries(entries, log.ErrorLevel, log.SourceSelect)
		out = log.FormatCounts(counts, "source")
	case c.messageCount:
		hclog.L().Debug("parsing debug bundle log file [ERROR] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.ErrorLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		hclog.L().Debug("aggregating [ERROR] messages by message string")
		counts := log.AggregateLogEntries(entries, log.ErrorLevel, log.MessageSelect)
		out = log.FormatCounts(counts, "message")
	default:
		hclog.L().Debug("parsing debug bundle log file [ERROR] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.ErrorLevel, "", time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		out = log.FormatLog(entries)
	}

	c.ui.Output(out)
	return 0
}

const synopsis = `Parses debug bundle log for [ERROR] messages`
const help = `
Usage: 
    consul-debug-read log parse-error [options]

Parses consul debug bundle logs for processing [ERROR] messages
`
