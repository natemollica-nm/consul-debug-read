package get

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

func (c *cmd) Help() string { return commands.Usage(showDebugPathHelp, c.flags) }

func (c *cmd) Synopsis() string {
	return synopsis
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
	hclog.L().Debug("checking CONSUL_DEBUG_PATH env variable")
	if ok, path := c.CheckPathEnvVarSet(); ok {
		c.ui.Output(path)
		return 0
	}
	hclog.L().Debug("CONSUL_DEBUG_PATH unset, rendering config path setting")
	if ok, path := c.RenderPathFromConfig(); ok {
		c.ui.Output(path)
	}
	return 0
}

func (c *cmd) CheckPathEnvVarSet() (bool, string) {
	if path := os.Getenv(read.DebugReadEnvVar); path != "" {
		return true, path
	}
	return false, ""
}

func (c *cmd) RenderPathFromConfig() (bool, string) {
	cmdYamlCfg, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		hclog.L().Error("error reading consul-debug-read user config file", "filepath", read.DebugReadConfigFullPath, "error", err)
		return false, ""
	}
	var currentPathSetting read.ReaderConfig
	err = yaml.Unmarshal(cmdYamlCfg, &currentPathSetting)
	if err != nil {
		hclog.L().Error("error deserializing YAML contents", "filepath", read.DebugReadConfigFullPath, "error", err)
		return false, ""
	}
	if currentPathSetting.DebugDirectoryPath == "" {
		hclog.L().Warn("empty or null consul-debug-path set", "warn", read.DebugReadConfigFullPath)
	}
	return true, currentPathSetting.DebugDirectoryPath
}

const synopsis = `Show the current debug bundle path under analysis`
const showDebugPathHelp = `
Shows currently set consul-debug-read command debug path as set in
config.yaml viper configuration file. 

To change file-path use consul-debug-read set --path <path_to_debug_bundle> to alter.

Example:
	$ consul-debug-read get
	bundles/consul-debug-2023-10-04T18-29-47Z
`
