package config

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

const synopsis = `Executes consul-debug-read bundle file system operations`
const help = `
Usage: 
    consul-debug-read config <subcommand> [options]

  Run consul-debug-read config <subcommand> with no arguments for help on that
  subcommand.
`
