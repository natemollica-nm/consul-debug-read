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
	"path/filepath"
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
	appCLI.Name = "consul-debug-read"
	appCLI.Args = os.Args[1:]
	appCLI.Commands = cmds
	appCLI.Autocomplete = true
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
		fmt.Printf("default configuration filepath not found, attempting to create and populate file=%s\n", read.DebugReadConfigFullPath)
		err = os.MkdirAll(read.DebugReadConfigDirPath, 0755)
		if err != nil {
			hclog.L().Error("failed to create directory", "error", err)
			return false
		}
	}
	var err error
	var defaultConfig []byte
	if _, err = os.Stat(read.DebugReadConfigFullPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("configuring default debug path to current directory dir=%s\n", read.CurrentDir)
			// Create Default Configuration File
			config := read.DefaultReaderConfig()
			if path := os.Getenv(read.DebugReadEnvVar); path != "" {
				fullPath, _ := filepath.Abs(path)
				config.DebugDirectoryPath = fullPath
				config.PathRenderedFrom = "env:CONSUL_DEBUG_PATH"
			} else {
				currentDir, _ := os.Getwd()
				config.DebugDirectoryPath = currentDir
				config.PathRenderedFrom = "file:config.yaml"
			}
			defaultConfig, err = yaml.Marshal(&config)
			if err != nil {
				fmt.Printf("failed to create configuration file error=%v", err)
				return false
			}
			err = os.WriteFile(read.DebugReadConfigFullPath, defaultConfig, 0755)
			if err != nil {
				fmt.Printf("failed to write to configuration file error=%v", err)
				return false
			}
		} else {
			var raw []byte
			var config read.ReaderConfig
			if raw, err = os.ReadFile(read.DebugReadConfigFullPath); err != nil {
				fmt.Printf("failed to read in config file debugConfigPath=%s, error=%v", read.DebugReadConfigFullPath, err)
				return false
			}
			if err = yaml.Unmarshal(raw, &config); err != nil {
				fmt.Printf("failed to unmarshal raw config file debugConfigPath=%s, error=%v", read.DebugReadConfigFullPath, err)
				return false
			}
			// Update the configuration path setting in the struct
			if config.PathRenderedFrom == "env:CONSUL_DEBUG_PATH" {
				config.DebugDirectoryPath = config.DebugEnvVarSetting
				if err = os.Setenv(read.DebugReadEnvVar, config.DebugEnvVarSetting); err != nil {
					hclog.L().Error("error setting debug-read env var", read.DebugReadEnvVar, config.DebugEnvVarSetting, "error", err)
					return false
				}
			}
		}
	}
	return true
}
