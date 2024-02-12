package summary

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
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

	var data read.Debug
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "members"); err != nil {
		hclog.L().Error("failed to decode members.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in agent information from bundle")

	var result string
	result = agentSummary(data)
	c.ui.Output(result)
	return 0
}

func agentSummary(bundle read.Debug) string {
	return bundle.Agent.Summary()
}

const synopsis = `Returns agent-specific information in summarized format`
const help = `
Usage: 
    consul-debug-read agent summary [options]
`
