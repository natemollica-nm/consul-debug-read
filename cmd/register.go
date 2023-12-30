package cmd

import (
	"consul-debug-read/cmd/cli"
	"consul-debug-read/internal/read/commands/agent"
	"consul-debug-read/internal/read/commands/agent/members"
	"consul-debug-read/internal/read/commands/agent/raft"
	"consul-debug-read/internal/read/commands/load"
	"consul-debug-read/internal/read/commands/metrics"
	"consul-debug-read/internal/read/commands/run"
	set "consul-debug-read/internal/read/commands/set-debug-path"
	show "consul-debug-read/internal/read/commands/show-debug-path"
	"fmt"
	mcli "github.com/mitchellh/cli"
)

func RegisteredCommands(ui cli.Ui) map[string]mcli.CommandFactory {
	registry := map[string]mcli.CommandFactory{}
	registerCommands(ui, registry,
		entry{"show-debug-path", func(cli.Ui) (mcli.Command, error) { return show.New(ui) }},
		entry{"set-debug-path", func(cli.Ui) (mcli.Command, error) { return set.New(ui) }},
		entry{"agent", func(cli.Ui) (mcli.Command, error) { return agent.New(ui) }},
		entry{"agent members", func(cli.Ui) (mcli.Command, error) { return members.New(ui) }},
		entry{"agent raft-configuration", func(cli.Ui) (mcli.Command, error) { return raft.New(ui) }},
		entry{"metrics", func(cli.Ui) (mcli.Command, error) { return metrics.New(ui) }},
		entry{"run", func(cli.Ui) (mcli.Command, error) { return run.New(ui) }},
		entry{"load", func(cli.Ui) (mcli.Command, error) { return load.New(ui) }},
	)
	return registry
}

// factory is a function that returns a new instance of a CLI-sub command.
type factory func(cli.Ui) (mcli.Command, error)

// entry is a struct that contains a command's name and a factory for that command.
type entry struct {
	name string
	fn   factory
}

func registerCommands(ui cli.Ui, m map[string]mcli.CommandFactory, cmdEntries ...entry) {
	for _, ent := range cmdEntries {
		thisFn := ent.fn
		if _, ok := m[ent.name]; ok {
			panic(fmt.Sprintf("duplicate command: %q", ent.name))
		}
		m[ent.name] = func() (mcli.Command, error) {
			return thisFn(ui)
		}
	}
}
