package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/natemollica-nm/consul-debug-reader/debug-read/internal/commands"
	"os"
)

func mainMetrics() {
	debugPath := flag.String("debug-path", "", "File path to extracted debug bundle")
	flag.Parse()

	if *debugPath == "" {
		fmt.Println("Please provide an input path to extracted debug bundle using the -debug-path flag.")
		return
	}

	var metricsFile = *debugPath + "/metrics.json"
	var agentFile = *debugPath + "/agent.json"
	var hostFile = *debugPath + "/host.json"
}

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
	app := cli.NewCLI("consul-debug-parser", debug.version)
	app.Args = os.Args[1:]
	app.Commands = map[string]cli.CommandFactory{
		"metrics": func() (cli.Command, error) { return commands.NewMetrics },
	}
}
