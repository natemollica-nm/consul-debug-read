package main

import (
	"consul-debug-read/cmd"
	"consul-debug-read/cmd/cli"
	"consul-debug-read/internal/read"
	"fmt"
	"github.com/hashicorp/go-hclog"
	mcli "github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(io.Discard)

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
	cmds := cmd.RegisteredCommands(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}
	appCLI := mcli.NewCLI("consul-debug-read", read.Version)
	appCLI.Args = os.Args[1:]
	appCLI.Commands = cmds
	appCLI.HelpFunc = mcli.FilteredHelpFunc(names, mcli.BasicHelpFunc("consul-debug-read"))
	appCLI.HelpWriter = os.Stdout
	appCLI.ErrorWriter = os.Stderr
	if ok := baseConfigs(); !ok {
		err := fmt.Errorf("failed to establish consul-debug-read configuration directory")
		ui.Error(err.Error())
		return 1
	}
	exitStatus, err := appCLI.Run()
	if err != nil {
		ui.Error(err.Error())
	}
	return exitStatus
}

func baseConfigs() bool {
	if _, err := os.Stat(read.DebugReadConfigDirPath); os.IsNotExist(err) {
		hclog.L().Info("default configuration filepath not found, attempting to create and populate", "file", read.DebugReadConfigFullPath)
		err = os.MkdirAll(read.DebugReadConfigDirPath, 0755)
		if err != nil {
			hclog.L().Error("failed to create directory", "error", err)
			return false
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
				return false
			}
			err = os.WriteFile(read.DebugReadConfigFullPath, defaultCfgBytes, 0755)
			if err != nil {
				hclog.L().Error("failed to create write to configuration file", "error", err)
				return false
			}
		}
	}
	return true
}
