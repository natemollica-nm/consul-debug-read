package main

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands/agent"
	"consul-debug-read/internal/read/commands/agent/members"
	"consul-debug-read/internal/read/commands/agent/raft"
	"consul-debug-read/internal/read/commands/metrics"
	set "consul-debug-read/internal/read/commands/set-debug-path"
	show "consul-debug-read/internal/read/commands/show-debug-path"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	ui := &cli.ColoredUi{
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		OutputColor: cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		InfoColor:   cli.UiColorBlue,
		WarnColor:   cli.UiColorYellow,
	}

	app := cli.NewCLI("consul-debug-read", read.Version)
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"show-debug-path":          func() (cli.Command, error) { return show.New(ui) },
		"set-debug-path":           func() (cli.Command, error) { return set.New(ui) },
		"agent":                    func() (cli.Command, error) { return agent.New(ui) },
		"agent members":            func() (cli.Command, error) { return members.New(ui) },
		"agent raft-configuration": func() (cli.Command, error) { return raft.New(ui) },
		"metrics":                  func() (cli.Command, error) { return metrics.New(ui) },
	}

	exitStatus, err := app.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
