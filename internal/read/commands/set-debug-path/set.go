package set_debug_path

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
	"os/exec"
	"strings"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	path    string
	file    string
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

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(setDebugPathHelp, c.flags) }

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

	if _, err := os.Stat(read.DebugReadConfigDirPath); os.IsNotExist(err) {
		hclog.L().Info("default configuration filepath not found, attempting to create and populate", "file", read.DebugReadConfigFullPath)
		err := os.MkdirAll(read.DebugReadConfigDirPath, 0755)
		if err != nil {
			hclog.L().Error("failed to create directory", "error", err)
			return 1
		}
	}

	if _, err := os.Stat(read.DebugReadConfigFullPath); err != nil {
		if os.IsNotExist(err) {
			hclog.L().Debug("configuring default debug path to current directory", "dir", read.CurrentDir)
			// Create Default Configuration File
			config := &read.ReaderConfig{
				DebugDirectoryPath: read.CurrentDir,
			}
			defaultCfgBytes, err := yaml.Marshal(&config)
			if err != nil {
				hclog.L().Error("failed to create default configuration file", "error", err)
				return 1
			}
			err = os.WriteFile(read.DebugReadConfigFullPath, defaultCfgBytes, 0755)
			if err != nil {
				hclog.L().Error("failed to create write to configuration file", "error", err)
				return 1
			}
		}
	}

	var extractedPath string
	var err error
	hclog.L().Debug("reading env var for configuration file update", "env", read.DebugReadEnvVar)
	if path := os.Getenv(read.DebugReadEnvVar); path != "" {
		if ok, err := validateDebugPath(path); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err)
		}
		hclog.L().Debug("env variable set, updating config file", "CONSUL_DEBUG_PATH", path)
		if ok, err := updateDebugReadConfig(path); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			return 1
		}
		c.ui.Output("consul-debug-path set successfully")
	} else if c.path != "" {
		extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.path)
		if err != nil {
			hclog.L().Error("failed to extract bundle from path", "path", c.path, "err", err)
			return 1
		}
		if ok, err := validateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err)
		}
		if ok, err := updateDebugReadConfig(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			return 1
		}
		c.ui.Output("consul-debug-path set successfully")
	} else if c.file != "" {
		if ok := strings.HasSuffix(c.path, ".tar.gz"); ok {
			extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.path)
			if err != nil {
				hclog.L().Error("failed to extract bundle from path", "path", c.path, "err", err)
				return 1
			}
		}
		if ok, err := validateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err)
		}
		if ok, err := updateDebugReadConfig(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			return 1
		}
		c.ui.Output("consul-debug-path set successfully")
	}
	return 0
}

func updateDebugReadConfig(updatePath string) (bool, error) {
	config := &read.ReaderConfig{
		DebugDirectoryPath: updatePath,
	}
	configBytes, err := yaml.Marshal(&config)
	if err != nil {
		hclog.L().Error("failed to create default configuration file", "error", err)
		return false, err
	}
	err = os.WriteFile(read.DebugReadConfigFullPath, configBytes, 0755)
	if err != nil {
		hclog.L().Error("failed to create write to configuration file", "error", err)
		return false, err
	}
	return true, nil
}

const synopsis = `Changes which bundle you're focusing on for analysis`
const setDebugPathHelp = `consul-debug-read set-debug-path [options]

Validates the path contents or extracts a valid .tar.gz bundle and points to this valid directory path for processing.

-path can be either:
  * consul-debug extracted contents (valid agent.json, metrics.json, host.json, and index.json) or
  * path to multiple bundles available for extraction and path setting

-file can be either:
  * path to valid consul debug .tar.gz archive or
  * path to multiple bundles available for extraction and path setting

Example (-path):
	$ consul-debug-read set-debug-path --path bundles/consul-debug-2023-10-04T18-29-47Z

Example (-path) for dir containing multiple .tar.gz bundles:
	$ consul-debug-read set-debug-path --path bundles

	select a .tar.gz file to extract:
	1: 124722consul-debug-2023-10-04T18-29-47Z.tar.gz
	2: 124722consul-debug-2023-10-11T17-33-55Z.tar.gz
	3: 124722consul-debug-2023-10-11T17-43-15Z.tar.gz
	4: 124722consul-debug-eu-01-stag.tar.gz
	5: 124722consul-debug-eu-133-stag-default.tar.gz
	6: 124722consul-debug-us-135-stag-default.tar.gz
	7: 124722consul-debug-us-east-stag.tar.gz
	enter the number of the file to extract: 

Example (-file) for extraction:
	$ consul-debug-read set-debug-path --file bundles/124722consul-debug-2023-10-11T17-43-15Z.tar.gz
`

func validateDebugPath(path string) (bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	var metricsJson, agentJson, membersJson, hostJson, indexJson bool
	for _, file := range files {
		switch file.Name() {
		case "metrics.json":
			metricsJson = true
		case "agent.json":
			agentJson = true
		case "host.json":
			hostJson = true
		case "index.json":
			indexJson = true
		case "members.json":
			membersJson = true
		case "cluster.json":
			clusterJsonPath := path + "/" + file.Name()
			membersJsonPath := path + "/members.json"
			if err := os.Rename(clusterJsonPath, membersJsonPath); err != nil {
				return false, err
			}
			membersJson = true
		}
	}
	if agentJson && membersJson && hostJson && indexJson {
		// older debug bundles separated v1/agent/metrics captures into each interval
		// if so, try and merge the metrics.json files into one large metrics.json
		// for ingestion.
		if !metricsJson {
			// "metrics.json" not found in the current directory
			// Run the "merge-metrics.sh" script with debugPath as an argument
			scriptPath := "scripts/merge-metrics.sh"
			cmd := exec.Command(scriptPath, path)
			if _, err := cmd.CombinedOutput(); err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
