package commands

import (
	"flag"
	"github.com/mitchellh/cli"
)

type exportCommand struct {
	ui    cli.Ui
	flags *flag.FlagSet

	output  string
	verbose bool
	silent  bool
}

func NewMetrics(ui cli.Ui) (cli.Command, error) {
	c := &exportCommand{
		ui:    ui,
		flags: flag.NewFlagSet("", flag.ContinueOnError),
	}

	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")
	c.flags.StringVar(&c.output, "output", "", "File path to output the data to. Defaults to stdout")

	return c, nil
}
