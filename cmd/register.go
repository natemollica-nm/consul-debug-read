package cmd

import (
	"consul-debug-read/cmd/cli"
	"consul-debug-read/internal/read/commands/agent"
	"consul-debug-read/internal/read/commands/agent/members"
	"consul-debug-read/internal/read/commands/agent/raft"
	show "consul-debug-read/internal/read/commands/get"
	"consul-debug-read/internal/read/commands/log"
	logdebug "consul-debug-read/internal/read/commands/log/parse/debug"
	logerror "consul-debug-read/internal/read/commands/log/parse/error"
	"consul-debug-read/internal/read/commands/log/parse/rpccounts"
	logtrace "consul-debug-read/internal/read/commands/log/parse/trace"
	"consul-debug-read/internal/read/commands/metrics"
	"consul-debug-read/internal/read/commands/set"
	"consul-debug-read/internal/read/commands/summary"
	"fmt"
	mcli "github.com/mitchellh/cli"
)

func RegisteredCommands(ui cli.Ui) map[string]mcli.CommandFactory {
	registry := map[string]mcli.CommandFactory{}
	registerCommands(ui, registry,
		entry{"current-path", func(cli.Ui) (mcli.Command, error) { return show.New(ui) }},
		entry{"set-path", func(cli.Ui) (mcli.Command, error) { return set.New(ui) }},
		entry{"agent", func(cli.Ui) (mcli.Command, error) { return agent.New(ui) }},
		entry{"agent members", func(ui cli.Ui) (mcli.Command, error) { return members.New(ui) }},
		entry{"agent raft-configuration", func(ui cli.Ui) (mcli.Command, error) { return raft.New(ui) }},
		entry{"metrics", func(cli.Ui) (mcli.Command, error) { return metrics.New(ui) }},
		entry{"summary", func(cli.Ui) (mcli.Command, error) { return summary.New(ui) }},
		entry{"log", func(cli.Ui) (mcli.Command, error) { return log.New(), nil }},
		entry{"log parse-rpc-counts", func(ui cli.Ui) (mcli.Command, error) { return rpccounts.New(ui) }},
		entry{"log parse-errors", func(ui cli.Ui) (mcli.Command, error) { return logerror.New(ui) }},
		entry{"log parse-debug", func(ui cli.Ui) (mcli.Command, error) { return logdebug.New(ui) }},
		entry{"log parse-trace", func(ui cli.Ui) (mcli.Command, error) { return logtrace.New(ui) }},
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
