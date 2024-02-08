package metrics

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// MetricsCmd represents the metrics command
var (
	keyMetricNames = map[string][]string{
		"Transaction Timing":               transactionTimingMetrics,
		"Leadership Stability":             leaderShipMetrics,
		"Certificate Authority Expiration": certAuthority,
		"Autopilot":                        autoPilotMetrics,
		"Memory Utilization":               memoryMetrics,
		"Network Activity":                 networkMetrics,
		"Raft Thread Saturation":           raftThreadSaturationMetrics,
		"Raft Replication Capacity":        replicationCapacity,
		"BoltDB Performance":               boltDBPerformance,
	}
	transactionTimingMetrics = []string{
		"consul.kvs.apply",
		"consul.txn.apply",
		"consul.raft.apply",
		"consul.raft.commitTime",
	}
	leaderShipMetrics = []string{
		"consul.raft.leader.lastContact",
		"consul.raft.state.candidate",
		"consul.raft.state.leader",
		"consul.server.isLeader",
	}
	certAuthority = []string{
		"consul.mesh.active-root-ca.expiry",
		"consul.mesh.active-signing-ca.expiry",
		"consul.agent.tls.cert.expiry",
	}
	autoPilotMetrics = []string{
		"consul.autopilot.healthy",
		"consul.autopilot.failure_tolerance",
	}
	memoryMetrics = []string{
		"consul.runtime.alloc_bytes",
		"consul.runtime.heap_objects",
		"consul.runtime.sys_bytes",
		"consul.runtime.total_gc_pause_ns",
	}
	networkMetrics = []string{
		"consul.client.rpc",
		"consul.client.rpc.exceeded",
		"consul.client.rpc.failed",
	}
	raftThreadSaturationMetrics = []string{
		"consul.raft.thread.main.saturation",
		"consul.raft.thread.fsm.saturation",
	}
	replicationCapacity = []string{
		"consul.raft.fsm.lastRestoreDuration",
		"consul.raft.leader.oldestLogAge",
		"consul.raft.rpc.installSnapshot",
	}
	boltDBPerformance = []string{
		"consul.raft.boltdb.freelistBytes",
		"consul.raft.boltdb.logsPerBatch",
		"consul.raft.boltdb.storeLogs",
		"consul.raft.boltdb.writeCapacity",
	}
	rateLimitingMetrics = []string{
		"consul.client.rpc",
		"consul.client.rpc.failed",
		"consul.client.rpc.exceeded",
		"consul.rpc.queries",
		"consul.rpc.queries_blocking",
		"consul.rpc.rate_limit.exceeded",
		"consul.rpc.rate_limit.log_dropped",
	}
	dataplaneMetrics = []string{
		"consul.xds.server.streams",
		"consul.xds.server.streamsUnauthenticated",
		"consul.xds.server.idealStreamsMax",
		"consul.xds.server.streamDrained",
		"consul.xds.server.streamStart",
	}
	federationMetrics = []string{
		"consul.leader.replication.acl-policies.status",
		"consul.leader.replication.acl-policies.index",
		"consul.leader.replication.acl-roles.status",
		"consul.leader.replication.acl-roles.index",
		"consul.leader.replication.acl-tokens.status",
		"consul.leader.replication.acl-tokens.index",
		"consul.leader.replication.config-entries.status",
		"consul.leader.replication.config-entries.index",
		"consul.leader.replication.federation-state.status",
		"consul.leader.replication.federation-state.index",
		"consul.leader.replication.namespaces.status",
		"consul.leader.replication.namespaces.index",
	}
	serviceMetrics = []string{
		"consul.state.services",
		"consul.state.service_instances",
		"consul.state.connect_instances",
		"consul.dns.stale_queries",
		"consul.dns.ptr_query",
		"consul.dns.domain_query",
	}
	serfHealthMetrics = []string{
		"consul.serf.member.failed",
		"consul.serf.member.flap",
		"consul.serf.member.join",
		"consul.serf.member.left",
		"consul.memberlist.degraded.probe",
		"consul.memberlist.degraded.timeout",
		"consul.memberlist.msg.dead",
		"consul.memberlist.health.score",
		"consul.memberlist.msg.suspect",
		"consul.memberlist.tcp.accept",
		"consul.memberlist.udp.sent",
		"consul.memberlist.udp.received",
		"consul.memberlist.tcp.connect",
		"consul.memberlist.tcp.sent",
		"consul.memberlist.gossip",
		"consul.memberlist.msg_alive",
		"consul.memberlist.msg_dead",
		"consul.memberlist.msg_suspect",
		"consul.memberlist.node.instances",
		"consul.memberlist.probeNode",
		"consul.memberlist.pushPullNode",
		"consul.memberlist.queue.broadcasts",
		"consul.memberlist.size.local",
		"consul.memberlist.size.remote",
		"consul.catalog.service.query",
		"consul.catalog.service.query-tag",
		"consul.catalog.service.query-tags",
		"consul.catalog.service.not-found",
		"consul.catalog.connect.query",
		"consul.catalog.connect.query-tag",
		"consul.catalog.connect.query-tags",
		"consul.catalog.connect.not-found",
	}
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	name    string
	summary bool

	listAvailableTelemetry bool

	keyMetrics   bool
	host         bool
	memory       bool
	network      bool
	rateLimiting bool

	autopilot         bool
	transactionTiming bool
	leadershipChanges bool

	bolt             bool
	threadSaturation bool

	dataplane        bool
	federationStatus bool

	telegraf bool

	sort    bool
	short   bool
	verify  bool
	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.StringVar(&c.name, "name", "", "Retrieve specific metric timestamped values by name")
	c.flags.BoolVar(&c.summary, "summary", false, "Retrieve metrics summary info from bundle")

	c.flags.BoolVar(&c.listAvailableTelemetry, "list-available-telemetry", false, "List available metric names as retrieved from consul telemetry docs")

	c.flags.BoolVar(&c.keyMetrics, "key-metrics", false, "Retrieve key metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.host, "host", false, "Retrieve Host specific metrics")
	c.flags.BoolVar(&c.memory, "memory", false, "Retrieve key memory metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.network, "network", false, "Retrieve key network metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.rateLimiting, "rate-limiting", false, "Retrieve key rate limit metric values for Consul from debug bundle")

	c.flags.BoolVar(&c.autopilot, "auto-pilot", false, "Retrieve key autopilot related metric values")
	c.flags.BoolVar(&c.transactionTiming, "transaction-timing", false, "Retrieve key transaction timing metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.leadershipChanges, "leadership-health", false, "Retrieve key raft leadership stability metric values for Consul from debug bundle")

	c.flags.BoolVar(&c.bolt, "bolt-db", false, "Retrieve key boltDB related metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.threadSaturation, "raft-thread-health", false, "Retrieve key raft thread saturation metric values for Consul from debug bundle")

	c.flags.BoolVar(&c.dataplane, "dataplane-health", false, "Retrieve key dataplane-related metric values for Consul from debug bundle")
	c.flags.BoolVar(&c.federationStatus, "federation-health", false, "Retrieve key secondary datacenter federation metric values for Consul from debug bundle")

	c.flags.BoolVar(&c.telegraf, "telegraf", false, "Generate telegraf compatible metrics file for ingesting offline metrics")

	c.flags.BoolVar(&c.sort, "sort", false, "Parse metric value by name and sort results by value vice timestamp order")
	c.flags.BoolVar(&c.short, "short", false, "Only print timestamp, value, and labels columns")
	c.flags.BoolVar(&c.verify, "verify", false, "Performs metric name validation with hashicorp docs")

	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(metricsHelp, c.flags) }

func (c *cmd) Synopsis() string { return synopsis }

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		c.ui.Error(fmt.Sprintf("Failed to parse flags: %v", err))
		return 1
	}
	if c.verbose && c.silent {
		c.ui.Error(fmt.Sprintf("Cannot specify both -silent and -verbose"))
		return 1
	}

	level := hclog.Info
	if c.verbose {
		level = hclog.Debug
	} else if c.silent {
		level = hclog.Off
	}

	commands.InitLogging(c.ui, level)
	cmdYamlCfg, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		hclog.L().Error("error reading consul-debug-read user config file", "filepath", read.DebugReadConfigFullPath, "error", err)
		return 1
	}
	var cfg read.ReaderConfig
	err = yaml.Unmarshal(cmdYamlCfg, &cfg)
	if err != nil {
		hclog.L().Error("error deserializing YAML contents", "filepath", read.DebugReadConfigFullPath, "error", err)
		return 1
	}
	if cfg.DebugDirectoryPath == "" {
		hclog.L().Error("empty or null consul-debug-path setting", "error", read.DebugReadConfigFullPath)
		return 1
	}

	var data read.Debug
	var result string

	if c.name != "" {
		hclog.L().Debug("reading in index.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "index"); err != nil {
			hclog.L().Error("failed to decode index.json", "error", err)
			return 1
		}
		hclog.L().Debug("reading in metrics.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "metrics"); err != nil {
			hclog.L().Error("failed to decode metrics.json", "error", err)
			return 1
		}
		hclog.L().Debug("successfully read in bundle contents")
	}
	if c.host {
		hclog.L().Debug("reading in host.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "host"); err != nil {
			hclog.L().Error("failed to decode host.json", "error", err)
			return 1
		}
		hclog.L().Debug("successfully read in host.json bundle contents")
	}
	if c.keyMetrics || c.summary || c.memory || c.network || c.rateLimiting || c.autopilot || c.transactionTiming || c.leadershipChanges || c.bolt || c.dataplane || c.federationStatus || c.telegraf {
		hclog.L().Debug("reading in agent.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "agent"); err != nil {
			hclog.L().Error("failed to decode agent.json", "error", err)
			return 1
		}
		hclog.L().Debug("reading in host.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "host"); err != nil {
			hclog.L().Error("failed to decode host.json", "error", err)
			return 1
		}
		hclog.L().Debug("reading in index.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "index"); err != nil {
			hclog.L().Error("failed to decode index.json", "error", err)
			return 1
		}
		hclog.L().Debug("reading in metrics.json")
		if err := data.DecodeJSON(cfg.DebugDirectoryPath, "metrics"); err != nil {
			hclog.L().Error("failed to decode metrics.json", "error", err)
			return 1
		}
		hclog.L().Debug("successfully read in bundle contents")
	}

	switch {
	case c.summary:
		result = data.Summary()
	case c.listAvailableTelemetry:
		result, err = read.ListMetrics()
		if err != nil {
			hclog.L().Error("failed to retrieve agent telemetry available metrics", "error", err)
			return 1
		}
	case c.host:
		result = data.HostSummary()
	case c.name != "":
		values, err := data.GetMetricValues(c.name, c.verify, c.sort, c.short)
		if err != nil {
			hclog.L().Error("failed to retrieve metric value", "name", c.name, "error", err)
			return 1
		}
		c.ui.Output(values)
	case c.keyMetrics:
		// for repeat ordering and printing
		var keyNames []string
		for k := range keyMetricNames {
			keyNames = append(keyNames, k)
		}
		sort.Strings(keyNames)
		for _, keyMetricTitle := range keyNames {
			fmt.Printf("[Next Key Metric] => %s: press [ENTER] to retrieve values", keyMetricTitle)
			metricNames := keyMetricNames[keyMetricTitle]
			for _, name := range metricNames {
				doneCh := make(chan bool)
				go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
				values, err := data.GetMetricValues(name, false, c.sort, c.short)
				if err != nil {
					hclog.L().Error("failed to retrieve metric value", "name", name, "error", err)
					return 1
				}
				doneCh <- true // Stop the dot printing goroutine
				close(doneCh)
				c.ui.Output(values)
			}
		}
		return 0
	case c.memory:
		ClearScreenPrompt("[Memory Metrics]: press [ENTER] to retrieve values")
		for _, name := range memoryMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.network:
		ClearScreenPrompt("[Network Metrics]: press [ENTER] to retrieve values")
		for _, name := range networkMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.rateLimiting:
		ClearScreenPrompt("[Rate Limiting Metrics]: press [ENTER] to retrieve values")
		for _, name := range rateLimitingMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.autopilot:
		ClearScreenPrompt("[Autopilot Metrics]: press [ENTER] to retrieve values")
		for _, name := range autoPilotMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.transactionTiming:
		ClearScreenPrompt("[Transaction Timing Metrics]: press [ENTER] to retrieve values")
		for _, name := range transactionTimingMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.leadershipChanges:
		ClearScreenPrompt("[Raft Leadership Health Metrics]: press [ENTER] to retrieve values")
		for _, name := range leaderShipMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.bolt:
		ClearScreenPrompt("[BoltDB Metrics]: press [ENTER] to retrieve values")
		for _, name := range boltDBPerformance {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.dataplane:
		ClearScreenPrompt("[Dataplane Metrics]: press [ENTER] to retrieve values")
		for _, name := range dataplaneMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.federationStatus:
		ClearScreenPrompt("[Federation Health Metrics]: press [ENTER] to retrieve values")
		for _, name := range federationMetrics {
			doneCh := make(chan bool)
			go Dots(fmt.Sprintf("==> reading '%s' values", name), doneCh)
			values, err := data.GetMetricValues(name, false, c.sort, c.short)
			if err != nil {
				hclog.L().Error("failed to retrieve metric", "name", name, "error", err)
				return 1
			}
			doneCh <- true // Stop the dot printing goroutine
			close(doneCh)
			c.ui.Output(values)
		}
		return 0
	case c.telegraf:
		if err := data.GenerateTelegrafMetrics(); err != nil {
			hclog.L().Error("failed to generate telegraf metrics files", "error", err)
			return 1
		}
		return 0
	default:
		result = c.Help()
	}
	c.ui.Output(result)
	return 0
}

const synopsis = `Ingest metrics.json from consul debug bundle`
const metricsHelp = `Read metrics information from specified bundle and return timestamped values.
Usage: consul-debug-read metrics [options]
  
  Parses and outputs consul debug bundle metrics.json data in readable format.

	Display summary of bundle capture	
		$ consul-debug-read metrics

	Display full list of queryable metric names
		$ consul-debug-read metrics -list 
	
	Retrieve all timestamped captures of metric
		$ consul-debug-read metrics -name <name_of_metric>
	
	Sort metric capture by value (highest to lowest)
		$ consul-debug-read metrics -name <name_of_metric> -sort
	
	Skip hashidoc metric name validation:
		$ consul-debug-read metrics -name <valid_name_but_not_in_docs> -verify=false`

func ClearScreenPrompt(message string) {
	clearScreen := exec.Command("clear")
	clearScreen.Stdout = os.Stdout
	_ = clearScreen.Run()
	fmt.Printf("\n%s", message)
	_, _ = fmt.Scanln()
}

func Dots(msg string, ch <-chan bool) {
	dots := []string{".", "..", "...", "...."}
	i := 0
	for {
		select {
		case <-ch:
			fmt.Print("\r")                                         // Carriage return to the beginning of the line
			fmt.Print(fmt.Sprintf(strings.Repeat(" ", len(msg)+5))) // Overwrite the line with spaces
			fmt.Print("\r")                                         // Carriage return again to the beginning of the line
			return
		default:
			fmt.Printf("%s", msg)
			fmt.Print(dots[i%len(dots)], "\r")
			i++
			time.Sleep(300 * time.Millisecond)
		}
	}
}
