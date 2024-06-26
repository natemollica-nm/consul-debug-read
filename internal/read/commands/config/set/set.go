package set

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
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
	c.flags.StringVar(&c.file, "file", "", "File path to .tar.gz set for debug bundle reading analysis")

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

	var path, extractedPath string
	var err error
	var ok, usePath, useFile, useEnvVar bool

	hclog.L().Debug("checking CONSUL_DEBUG_PATH env var (if set)", "env", read.DebugReadEnvVar)

	if path, err = read.ExtractEnvironmentPath(); err != nil {
		c.ui.Error("failed to retrieve absolute path of CONSUL_DEBUG_PATH environment variable setting")
		return 1
	}

	hclog.L().Debug("CONSUL_DEBUG_PATH env var setting", "path", path)

	if path != "" && c.path == "" && c.file == "" {
		useEnvVar = true
	} else if c.file != "" {
		useFile = true
	} else if c.path != "" {
		usePath = true
	} else {
		c.ui.Error("[FAILED] No path settings were passed in! set your path using the `-path` or `-file` flags or by setting the CONSUL_DEBUG_PATH environment variable")
		return 1
	}

	if useEnvVar {
		hclog.L().Debug("attempting to set with CONSUL_DEBUG_PATH env variable", "path", path)
		extractedPath, err = read.SelectAndExtractTarGzFilesInDir(path)
		if err != nil {
			hclog.L().Error("failed to extract bundle from path", "path", path, "err", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		hclog.L().Debug("successfully extracted bundle from path", "path", path, "extractedPath", extractedPath)
		if ok, err = ValidateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err, "path", extractedPath)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		hclog.L().Debug("env variable set, updating config file", "CONSUL_DEBUG_PATH", path)
		if ok, err = UpdateCurrentPath(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		hclog.L().Debug("using env var setting", read.DebugReadEnvVar, extractedPath)
		c.ui.Output(fmt.Sprintf("\nconsul-debug-path set successfully using CONSUL_DEBUG_PATH env var => %s\n", extractedPath))
	} else if usePath {
		hclog.L().Debug("attempting to set with -path filepath", "path", c.path)
		extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.path)
		if err != nil {
			hclog.L().Error("failed to extract bundle from path", "path", c.path, "err", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		hclog.L().Debug("successfully extracted bundle from path", "path", c.path, "extractedPath", extractedPath)
		if ok, err = ValidateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err, "path", extractedPath)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		if ok, err = UpdateCurrentPath(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			c.ui.Error("failed to set consul-debug-read path using -path")
			return 1
		}
		c.ui.Output(fmt.Sprintf("\nconsul-debug-path set successfully => %s\n", extractedPath))
	} else if useFile {
		hclog.L().Debug("attempting to set with -file filepath", "file", c.file)
		if ok = strings.HasSuffix(c.file, ".tar.gz"); ok {
			extractedPath, err = read.SelectAndExtractTarGzFilesInDir(c.file)
			if err != nil {
				hclog.L().Error("failed to extract bundle from file", "file", c.file, "err", err)
				c.ui.Error("failed to set consul-debug-read path using -file")
				return 1
			}
		}
		if ok, err = ValidateDebugPath(extractedPath); !ok {
			hclog.L().Error("extracted bundle is invalid and does not contain all required debug bundle file extracts", "error", err, "path", extractedPath)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		if ok, err = UpdateCurrentPath(extractedPath); !ok {
			hclog.L().Error("failed update debug-read configuration file", "error", err)
			c.ui.Error("failed to set consul-debug-read path")
			return 1
		}
		c.ui.Output(fmt.Sprintf("\nconsul-debug-path set successfully => %s\n", extractedPath))
	}
	return 0
}

func UpdateCurrentPath(updatePath string) (bool, error) {
	var config read.ReaderConfig

	// Retrieve configuration file location and open file
	currentData, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		hclog.L().Error("error reading consul-debug-read user config file", "filepath", read.DebugReadConfigFullPath, "error", err)
		return false, err
	}
	// render current contents into read.ReaderConfig struct
	err = yaml.Unmarshal(currentData, &config)
	if err != nil {
		hclog.L().Error("error deserializing YAML contents", "filepath", read.DebugReadConfigFullPath, "error", err)
		return false, err
	}

	// Update the configuration path setting in the struct
	if path := os.Getenv(read.DebugReadEnvVar); path != "" {
		config.DebugDirectoryPath = updatePath
		config.DebugEnvVarSetting = updatePath
		config.PathRenderedFrom = "env:CONSUL_DEBUG_PATH"
	} else {
		config.DebugDirectoryPath = updatePath
		config.DebugEnvVarSetting = "<UNSET>"
		config.PathRenderedFrom = "cli:-path|-file"
	}

	// Marshal the struct to bytes and write to file
	updatedConfig, err := yaml.Marshal(&config)
	if err != nil {
		hclog.L().Error("failed to create default configuration file", "error", err)
		return false, err
	}

	// Write the changes to file
	err = os.WriteFile(read.DebugReadConfigFullPath, updatedConfig, 0755)
	if err != nil {
		hclog.L().Error("failed to create write to configuration file", "error", err)
		return false, err
	}
	return true, nil
}

const synopsis = `Changes which bundle you're focusing on for analysis`
const setDebugPathHelp = `consul-debug-read config set [options]

Validates the path contents or extracts a valid .tar.gz bundle and points to this valid directory path for processing.

-path can be either:
  * consul-debug extracted contents (valid agent.json, metrics.json, host.json, and index.json) or
  * path to multiple bundles available for extraction and path setting

-file can be either:
  * path to valid consul debug .tar.gz archive or
  * path to multiple bundles available for extraction and path setting

Example (-path):
	$ consul-debug-read config set -path bundles/consul-debug-2023-10-04T18-29-47Z

Example (-path) for dir containing multiple .tar.gz bundles:
	$ consul-debug-read config set -path bundles/

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
	$ consul-debug-read config set-path -file bundles/124722consul-debug-2023-10-11T17-43-15Z.tar.gz
`

func ValidateDebugPath(path string) (bool, error) {
	dir, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer dir.Close()

	entries, err := dir.ReadDir(0)
	if err != nil {
		return false, err
	}

	var metricsJson, agentJson, membersJson, hostJson, indexJson, consulLog bool
	for _, file := range entries {
		if file.IsDir() {
			continue
		}
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
		case "consul.log":
			consulLog = true
		case "cluster.json":
			clusterJsonPath := path + "/" + file.Name()
			membersJsonPath := path + "/members.json"
			if err = os.Rename(clusterJsonPath, membersJsonPath); err != nil {
				return false, err
			}
			membersJson = true
		}
	}
	if agentJson && membersJson && hostJson && indexJson {
		if !metricsJson {
			err = ConcatenateMetrics(path)
			if err != nil {
				return false, err
			}
		}
		if !consulLog {
			err = RetrieveFirstConsulLog(path)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	return false, fmt.Errorf("invalid path setting passed in | file-check: metrics=%v, agent=%v, host=%v, index=%v, members=%v, log=%v", metricsJson, agentJson, hostJson, indexJson, membersJson, consulLog)
}

// ConcatenateMetrics reads all metrics.json files in the subdirectories of bundle
// and appends their contents to a single metrics.json file in root bundle dir.
//
// Older debug bundles separated v1/agent/metrics captures into each interval
// so in this case we try and merge the metrics.json files into one large metrics.json
// for ingestion.
func ConcatenateMetrics(path string) error {
	var err error
	// Open or create the final metrics.json file for appending
	outputFile, err := os.OpenFile(filepath.Join(path, "metrics.json"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("unable to open or create metrics.json file: %w", err)
	}
	defer outputFile.Close()
	cleanup := func(err error) error {
		_ = outputFile.Close()
		return err
	}

	// Read directories in the debugPath
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			captureDir := filepath.Join(path, entry.Name())
			metricsPath := filepath.Join(captureDir, "metrics.json")

			// Check if metrics.json exists in the directory
			if _, err = os.Stat(metricsPath); err == nil {
				var metricsData []byte
				metricsData, err = os.ReadFile(metricsPath)
				if err != nil {
					return fmt.Errorf("error reading %v: %w", metricsPath, err)
				}

				// Append the metrics data to the final metrics.json file
				_, err = outputFile.Write(metricsData)
				if err != nil {
					return fmt.Errorf("error writing to metrics.json: %w", err)
				}

				// Optionally, append a newline between entries
				_, err = outputFile.WriteString("\n")
				if err != nil {
					return fmt.Errorf("error writing newline to metrics.json: %w", err)
				}
			}
		}
	}
	if err = outputFile.Close(); err != nil {
		return cleanup(err)
	}

	return nil
}

// RetrieveFirstConsulLog searches subdirectories of debug path for the first "consul.log"
// and copies it to path root. If no "consul.log" is found, it returns an error.
func RetrieveFirstConsulLog(path string) error {
	found := false
	var err error
	var srcFile, dstFile *os.File

	// Read directories in the debugPath
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	cleanup := func(file *os.File, err error) error {
		_ = file.Close()
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			consulLogPath := filepath.Join(path, entry.Name(), "consul.log")

			// Check if consul.log exists in the directory
			if _, err = os.Stat(consulLogPath); err == nil {

				// Copy consul.log to debugPath
				srcFile, err = os.Open(consulLogPath)
				if err != nil {
					return fmt.Errorf("error opening source consul.log: %w", err)
				}
				defer srcFile.Close()

				dstFile, err = os.Create(filepath.Join(path, "consul.log"))
				if err != nil {
					return fmt.Errorf("error creating destination consul.log: %w", err)
				}
				defer dstFile.Close()

				if _, err = io.Copy(dstFile, srcFile); err != nil {
					return fmt.Errorf("error copying consul.log: %w", err)
				}

				found = true
				break // Stop searching after finding the first instance
			}
		}
	}

	if !found {
		return fmt.Errorf("consul.log not found in any subdirectory")
	}

	if err = srcFile.Close(); err != nil {
		return cleanup(srcFile, err)
	}
	if err = dstFile.Close(); err != nil {
		return cleanup(dstFile, err)
	}

	return nil
}
