package cmd

import (
	"consul-debug-read/cmd/cli"
	"consul-debug-read/internal/read/commands/agent"
	agentconfig "consul-debug-read/internal/read/commands/agent/config"
	"consul-debug-read/internal/read/commands/agent/members"
	"consul-debug-read/internal/read/commands/agent/raft"
	agentsummary "consul-debug-read/internal/read/commands/agent/summary"
	"consul-debug-read/internal/read/commands/config"
	"consul-debug-read/internal/read/commands/config/get"
	"consul-debug-read/internal/read/commands/config/set"
	"consul-debug-read/internal/read/commands/config/show"
	"consul-debug-read/internal/read/commands/log"
	logdebug "consul-debug-read/internal/read/commands/log/parse/debug"
	logerror "consul-debug-read/internal/read/commands/log/parse/error"
	loginfo "consul-debug-read/internal/read/commands/log/parse/info"
	"consul-debug-read/internal/read/commands/log/parse/rpccounts"
	logtrace "consul-debug-read/internal/read/commands/log/parse/trace"
	logwarn "consul-debug-read/internal/read/commands/log/parse/warn"
	"consul-debug-read/internal/read/commands/metrics"
	"consul-debug-read/internal/read/commands/summary"
	"fmt"
	mcli "github.com/mitchellh/cli"
)

func RegisteredCommands(ui cli.Ui) map[string]mcli.CommandFactory {
	registry := map[string]mcli.CommandFactory{}
	registerCommands(ui, registry,
		entry{"config", func(cli.Ui) (mcli.Command, error) { return config.New(), nil }},
		entry{"config current-path", func(ui cli.Ui) (mcli.Command, error) { return get.New(ui) }},
		entry{"config set-path", func(ui cli.Ui) (mcli.Command, error) { return set.New(ui) }},
		entry{"config show", func(ui cli.Ui) (mcli.Command, error) { return show.New(ui) }},
		entry{"agent", func(cli.Ui) (mcli.Command, error) { return agent.New(), nil }},
		entry{"agent summary", func(ui cli.Ui) (mcli.Command, error) { return agentsummary.New(ui) }},
		entry{"agent config", func(ui cli.Ui) (mcli.Command, error) { return agentconfig.New(ui) }},
		entry{"agent members", func(ui cli.Ui) (mcli.Command, error) { return members.New(ui) }},
		entry{"agent raft-configuration", func(ui cli.Ui) (mcli.Command, error) { return raft.New(ui) }},
		entry{"metrics", func(cli.Ui) (mcli.Command, error) { return metrics.New(ui) }},
		entry{"summary", func(cli.Ui) (mcli.Command, error) { return summary.New(ui) }},
		entry{"log", func(cli.Ui) (mcli.Command, error) { return log.New(), nil }},
		entry{"log parse-rpc-counts", func(ui cli.Ui) (mcli.Command, error) { return rpccounts.New(ui) }},
		entry{"log parse-error", func(ui cli.Ui) (mcli.Command, error) { return logerror.New(ui) }},
		entry{"log parse-debug", func(ui cli.Ui) (mcli.Command, error) { return logdebug.New(ui) }},
		entry{"log parse-trace", func(ui cli.Ui) (mcli.Command, error) { return logtrace.New(ui) }},
		entry{"log parse-warn", func(ui cli.Ui) (mcli.Command, error) { return logwarn.New(ui) }},
		entry{"log parse-info", func(ui cli.Ui) (mcli.Command, error) { return loginfo.New(ui) }},
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
