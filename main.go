package main

import (
	"consul-debug-read/cmd"
	"consul-debug-read/cmd/cli"
	"consul-debug-read/internal/read"
	mcli "github.com/mitchellh/cli"
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

	exitStatus, err := appCLI.Run()
	if err != nil {
		ui.Error(err.Error())
	}
	return exitStatus
}
