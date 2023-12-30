package load

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"consul-debug-read/internal/read/commands/run"
	set "consul-debug-read/internal/read/commands/set-debug-path"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"strings"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	path string
	file string

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
	c.flags.StringVar(&c.path, "path", "", "File path to set for debug bundle reading analysis")
	c.flags.StringVar(&c.file, "file", "", "File path to .tar.gz set for debug bundle reading analysis")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

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

	var extractedPath string
	var err error
	var ok bool
	if c.path != "" {
		hclog.L().Debug("attempting to set with -path filepath", "path", c.path)
		extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.path)
		if err != nil {
			hclog.L().Error("failed to extract bundle from path", "path", c.path, "err", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		if ok, err := set.ValidateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		if ok, err := set.UpdateDebugReadConfig(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		c.ui.Output("consul-debug-path set successfully")
	} else if c.file != "" {
		hclog.L().Debug("attempting to set with -file filepath", "file", c.file)
		if ok = strings.HasSuffix(c.path, ".tar.gz"); ok {
			extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.path)
			if err != nil {
				hclog.L().Error("failed to extract bundle from path", "path", c.path, "err", err)
				c.ui.Error("failed to set consul-debug-read path")
				return 1
			}
		}
		if ok, err = set.ValidateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		if ok, err = set.UpdateDebugReadConfig(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		c.ui.Output("consul-debug-path set successfully")
	}
	path, err := run.GetPath()
	if err != nil {
		c.ui.Error("failed to retrieve updated path from configuration file")
		return 1
	}
	run.RunTimeBundle.Backend = read.NewBackend()
	c.ui.Output("successfully initialized backend boltDB")
	c.ui.Output(fmt.Sprintf("consul-debug-read running | v%s", read.Version))

	c.ui.Output("starting bundle serialization to boltDB")
	c.ui.Output("reading in agent.json")
	if err = run.RunTimeBundle.DecodeJSON(path, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	c.ui.Output("reading in host.json")
	if err = run.RunTimeBundle.DecodeJSON(path, "host"); err != nil {
		hclog.L().Error("failed to decode host.json", "error", err)
		return 1
	}
	c.ui.Output("reading in index.json")
	if err = run.RunTimeBundle.DecodeJSON(path, "index"); err != nil {
		hclog.L().Error("failed to decode index.json", "error", err)
		return 1
	}
	c.ui.Output("reading in metrics.json")
	if err = run.RunTimeBundle.DecodeJSON(path, "metrics"); err != nil {
		hclog.L().Error("failed to decode metrics.json", "error", err)
		return 1
	}
	c.ui.Output("successfully read in bundle contents to boltDB")

	return 0
}

const synopsis = "Loads consul debug bundle into runtime configuration"
const help = `
Usage: consul-debug-read load [options]

  Triggers bundle serialization and upload to the runtime boltDB in-memory
  database.
`
