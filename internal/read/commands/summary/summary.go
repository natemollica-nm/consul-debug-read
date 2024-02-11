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

func (c *cmd) Help() string { return commands.Usage(metricsHelp, c.flags) }

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
	var result string
	hclog.L().Debug("reading in agent.json")
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading members.json")
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "members"); err != nil {
		hclog.L().Error("failed to decode members.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in host.json")
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "host"); err != nil {
		hclog.L().Error("failed to decode host.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in index.json")
	if err = data.DecodeJSON(cfg.DebugDirectoryPath, "index"); err != nil {
		hclog.L().Error("failed to decode index.json", "error", err)
		return 1
	}
	hclog.L().Debug("reading in metrics.json")
	if err := data.DecodeJSON(cfg.DebugDirectoryPath, "metrics"); err != nil {
		hclog.L().Error("failed to decode metrics.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in bundle contents")

	files, err := getLogFiles(cfg.DebugDirectoryPath)
	if err != nil {
		return 0
	}
	captureTime, err := getTimestamp(cfg.DebugDirectoryPath)
	if err != nil {
		return 0
	}

	result = fmt.Sprintf("Consul Debug Bundle (%s): %s\nDebug Command Log Level: %s (Default) %s\n%s\n%s\n%s\n",
		captureTime,
		cfg.DebugDirectoryPath,
		data.Agent.LogLevel(),
		formatIndentedList(files, 1),
		data.Agent.Summary(),
		data.Summary(),
		data.HostSummary(),
	)

	c.ui.Output(result)
	return 0
}

const synopsis = `Parses bundle contents and provides outlined summary of debug capture.`
const metricsHelp = `Parses bundle contents and provides outlined summary of debug capture.

Usage: consul-debug-read summary [options]
  
  Parses and prints results of overall health of Consul from captured 'consul debug' command.

	Display summary of bundle capture	
		$ consul-debug-read summary`

// getLogFiles retrieves all .log files from a directory
func getLogFiles(dir string) ([]string, error) {
	// Retrieve directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	// Slice to store the paths of .log files
	var logFiles []string
	conv := read.ByteConverter{}
	// Iterate through directory contents
	for _, file := range files {
		info, _ := file.Info()
		logFileSize := conv.ConvertToReadableBytes(info.Size())
		// Check if the file is a regular file and has a .log extension
		if !file.IsDir() && filepath.Ext(file.Name()) == ".log" {
			logPath := filepath.Join(dir, file.Name())
			// Construct the full path of the log file and add it to the slice
			logFiles = append(logFiles, fmt.Sprintf("%s (%s)", logPath, logFileSize))
		}
	}

	return logFiles, nil
}

func formatIndentedList(list []string, indentLevel int) string {
	// Calculate the indentation string
	indent := strings.Repeat("  ", indentLevel)

	// Create a string builder to efficiently build the formatted string
	var builder strings.Builder

	// Iterate over the list and append each element with indentation to the builder
	for _, item := range list {
		builder.WriteString(fmt.Sprintf("\n%s* %s\n", indent, item))
	}

	// Return the formatted string
	return builder.String()
}

func getTimestamp(path string) (string, error) {
	// Retrieve file or directory information
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// Extract and format modification time
	modTime := fileInfo.ModTime()
	formattedTime := modTime.Format("2006-01-02 15:04:05")

	return formattedTime, nil
}
