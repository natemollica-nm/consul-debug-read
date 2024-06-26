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

type Cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &Cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *Cmd) Help() string { return commands.Usage(help, c.flags) }

func (c *Cmd) Synopsis() string {
	return synopsis
}

func (c *Cmd) Run(args []string) int {
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
	hclog.L().Debug("rendering debug path setting from config.yaml")
	if path, ok := RenderPathFromConfig(); ok {
		c.ui.Output(path)
	}
	return 0
}

func RenderPathFromConfig() (string, bool) {
	var path string
	var config read.ReaderConfig

	currentData, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		hclog.L().Error("error reading consul-debug-read user config file", "filepath", read.DebugReadConfigFullPath, "error", err)
		return "", false
	}

	err = yaml.Unmarshal(currentData, &config)
	if err != nil {
		hclog.L().Error("error deserializing YAML contents", "filepath", read.DebugReadConfigFullPath, "error", err)
		return "", false
	}

	// Logic to ensure that if the terminal session sets the CONSUL_DEBUG_PATH
	// environment variable, we default to using the env var over the configured
	// path from ~/.consul-debug-read/config.yaml
	if path = os.Getenv(read.DebugReadEnvVar); path != "" {
		var extractedPath string
		hclog.L().Debug("configuring path from rendered CONSUL_DEBUG_PATH setting", read.DebugReadEnvVar, path)
		if extractedPath, err = read.SelectAndExtractTarGzFilesInDir(path); err != nil {
			hclog.L().Error("failed to extract bundle from path", "path", path, "err", err)
			return "", false
		}
		return extractedPath, true
	} else if config.DebugDirectoryPath == "" {
		hclog.L().Warn("empty or null consul-debug-path set", "warn", read.DebugReadConfigFullPath)
		return config.DebugDirectoryPath, true
	} else {
		hclog.L().Debug("configuring path from rendered $HOME/.consul-debug-read/config.yaml setting", "DebugPath", config.DebugDirectoryPath)
		return config.DebugDirectoryPath, true
	}
}

const synopsis = `Show the current debug bundle path under analysis`
const help = `
Shows current debug path setting as set in $HOME/.consul-debug-read/config.yaml. 

Example:
	$ consul-debug-read config current-path
`
