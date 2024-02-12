package agent

import (
	"consul-debug-read/internal/read/commands"
	"github.com/mitchellh/cli"
)

type Cmd struct{}

func New() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Help() string {
	return commands.Usage(help, nil)
}

func (c *Cmd) Synopsis() string { return synopsis }

func (c *Cmd) Run(args []string) int {
	return cli.RunResultHelp
}

const synopsis = `Returns information related to the debug bundles agent.json file`
const help = `
Usage: 
    consul-debug-read agent <subcommand> [options]

  Run consul-debug-read agent <subcommand> with no arguments for help on that
  subcommand.
`
