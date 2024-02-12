package show

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
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

	hclog.L().Debug("rendering consul-debug-read config.yaml from user's home directory")
	if config, ok := c.renderConfigurationSettings(); ok {
		c.ui.Output(config)
	}
	return 0
}

func (c *cmd) CheckPathEnvVarSet() (string, bool) {
	if read.EnvVarPathSetting != "" {
		path, _ := filepath.Abs(read.EnvVarPathSetting)
		return path, true
	}
	return "", false
}

func (c *cmd) renderConfigurationSettings() (string, bool) {
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

	// Build configuration title output
	title := "consul-debug-read configuration settings"
	ul := fmt.Sprintf(strings.Repeat("-", len(title)))
	menu := []string{fmt.Sprintf("\x1f%s\x1f", title)}
	menu = append(menu, fmt.Sprintf("\x1f%s\x1f", ul))

	var env string
	var set bool
	if env, set = c.CheckPathEnvVarSet(); !set {
		env = "<UNSET>"
	}

	menu = append(menu, fmt.Sprintf("Setting\x1fValue\x1f"))
	menu = append(menu, fmt.Sprintf("-------\x1f-----\x1f"))
	menu = append(menu, fmt.Sprintf("Configuration File Location\x1f%s", config.ConfigFile))
	menu = append(menu, fmt.Sprintf("Debug Bundle Path\x1f%s (%s)", config.DebugDirectoryPath, config.PathRenderedFrom))
	menu = append(menu, fmt.Sprintf("CONSUL_DEBUG_PATH\x1f%s", env))
	output := columnize.Format(menu, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output, true
}

const synopsis = `Show configuration details for the consul-debug-read tool`
const help = `
Shows currently settings currently used from $HOME/.consul-debug-read/config.yaml

Example:
	$ consul-debug-read config show
`
