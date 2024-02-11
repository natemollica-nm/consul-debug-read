package info

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
	c.flags.StringVar(&c.source, "source", "", "Capture INFO messages from specific sources (e.g., \"agent.http\", \"agent.server\")")
	c.flags.BoolVar(&c.messageCount, "message-count", false, "Parse log for INFO messages and return timestamp sorted list of messages received")
	c.flags.BoolVar(&c.sourceCount, "source-count", false, "Parse log for INFO messages and return count sorted (descending order) list of messages received from specific sources")

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
	var out string

	switch {
	case c.source != "" && !c.sourceCount:
		hclog.L().Debug("parsing info bundle log file [INFO] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.InfoLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		out = log.FormatLog(entries)
	case c.sourceCount:
		hclog.L().Debug("parsing info bundle log file [INFO] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.InfoLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		hclog.L().Debug("aggregating [INFO] messages by logged entry source type", "log-file", logFile)
		counts := log.AggregateLogEntries(entries, log.InfoLevel, log.SourceSelect)
		out = log.FormatCounts(counts, "source")
	case c.messageCount:
		hclog.L().Debug("parsing info bundle log file [INFO] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.InfoLevel, c.source, time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		hclog.L().Debug("aggregating [INFO] messages by message string", "log-file", logFile)
		counts := log.AggregateLogEntries(entries, log.InfoLevel, log.MessageSelect)
		out = log.FormatCounts(counts, "message")
	default:
		hclog.L().Debug("parsing info bundle log file [INFO] messages", "log-file", logFile)
		entries, err = log.ParseLog(logFile, log.InfoLevel, "", time.Time{}, time.Time{})
		if err != nil {
			hclog.L().Error("error parsing log file", "file", logFile, "error", err)
			return 1
		}
		out = log.FormatLog(entries)
	}

	c.ui.Output(out)
	return 0
}

const synopsis = `Parses debug bundle log for [INFO] messages`
const help = `
Usage: 
    consul-debug-read log parse-warn [options]

Parses consul debug bundle logs for processing [INFO] messages
`
