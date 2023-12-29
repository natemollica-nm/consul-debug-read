package run

import (
	"consul-debug-read/internal/read"
	"consul-debug-read/internal/read/commands"
	"consul-debug-read/internal/read/commands/flags"
	runtime "consul-debug-read/runtime"
	"flag"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"gopkg.in/yaml.v2"
	"os"
	"os/signal"
	"syscall"
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}

	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string { return commands.Usage(help, c.flags) }

func (c *cmd) Run(args []string) int {
	ui := &cli.PrefixedUi{
		OutputPrefix: "  ==> ",
		InfoPrefix:   "    ", // Note that startupLogger also uses this prefix
		ErrorPrefix:  "  ==> ",
		Ui:           c.ui,
	}
	if err := c.flags.Parse(args); err != nil {
		ui.Error(fmt.Sprintf("Failed to parse flags: %v", err))
		return 1
	}
	if c.verbose && c.silent {
		ui.Error(fmt.Sprintf("Cannot specify both -silent and -verbose"))
		return 1
	}
	c.ui.Output("Starting debug reader as long-running application")

	level := hclog.Info
	if c.verbose {
		level = hclog.Debug
	} else if c.silent {
		level = hclog.Off
	}
	commands.InitLogging(c.ui, level)

	path, err := getPath()
	if err != nil {
		ui.Error(fmt.Sprintf("failed to retrieve path setting error=%v", err))
		return 1
	}
	ui.Output(fmt.Sprintf("current set path: %s", path))

	ui.Output("initializing boltDB")
	ui.Output(fmt.Sprintf("consul-debug-read running | v%s", read.Version))
	var data read.Debug
	data.Backend = read.NewBackend()
	ui.Output("successfully initialized backend boltDB")

	ui.Output("starting bundle serialization to boltDB")
	ui.Output("reading in agent.json")
	if err := data.DecodeJSON(path, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	ui.Output("reading in host.json")
	if err := data.DecodeJSON(path, "host"); err != nil {
		hclog.L().Error("failed to decode host.json", "error", err)
		return 1
	}
	ui.Output("reading in index.json")
	if err := data.DecodeJSON(path, "index"); err != nil {
		hclog.L().Error("failed to decode index.json", "error", err)
		return 1
	}
	ui.Output("reading in metrics.json")
	if err := data.DecodeJSON(path, "metrics"); err != nil {
		hclog.L().Error("failed to decode metrics.json", "error", err)
		return 1
	}
	ui.Output("successfully read in bundle contents to boltDB")

	ui.Output("testing boltDB query")
	values, err := data.GetMetricValuesMemDB("consul.rpc.rate_limit.exceeded", true, false, false)
	if err != nil {
		ui.Error(fmt.Sprintf("failed to retrieve memDB values error=%v", err))
		return 1
	}
	c.ui.Output(values)

	ui.Output("terminate at anytime by sending interrupt sig (ctrl + c)")
	// wait for signal
	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGPIPE)

	for {
		var sig os.Signal
		select {
		case s := <-signalCh:
			sig = s
		case <-runtime.ShutdownChannel():
			sig = os.Interrupt
		}

		switch sig {
		case syscall.SIGPIPE:
			continue
		default:
			ui.Info(fmt.Sprintf("\ncaught terminate signal (signal=%v), exiting...", sig))
			return 0
		}
	}
}

func getPath() (string, error) {
	if path := os.Getenv(read.DebugReadEnvVar); path != "" {
		return path, nil
	}
	cmdYamlCfg, err := os.ReadFile(read.DebugReadConfigFullPath)
	if err != nil {
		return "", fmt.Errorf("error reading consul-debug-read user config file filepath=%s error=%v", read.DebugReadConfigFullPath, err)
	}
	var currentPathSetting read.ReaderConfig
	err = yaml.Unmarshal(cmdYamlCfg, &currentPathSetting)
	if err != nil {
		return "", fmt.Errorf("error deserializing YAML contents filepath=%s error=%v", read.DebugReadConfigFullPath, err)
	}
	if currentPathSetting.DebugDirectoryPath == "" {
		return "<null>", nil
	}
	return currentPathSetting.DebugDirectoryPath, nil
}

const synopsis = "Runs consul-debug-read as long-running application"
const help = `
Usage: consul-debug-read run [options]

  Starts consul-debug-read as a runtime app and runs until an interrupt is received. This
  maintains the in-memory ingestion of debug bundles for quicker querying and response.
`
