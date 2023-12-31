package agent

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
)

type cmd struct {
	ui        cli.Ui
	flags     *flag.FlagSet
	pathFlags *flags.DebugReadFlags

	summary bool
	config  bool
	verbose bool
	silent  bool
}

func New(ui cli.Ui) (cli.Command, error) {
	c := &cmd{
		ui:        ui,
		pathFlags: &flags.DebugReadFlags{},
		flags:     flag.NewFlagSet("", flag.ContinueOnError),
	}
	c.flags.BoolVar(&c.config, "config", false, "Retrieve agent configuration in HCL format")
	c.flags.BoolVar(&c.summary, "summary", false, "Retrieve agent configuration summary details")
	c.flags.BoolVar(&c.silent, "silent", false, "Disables all normal log output")
	c.flags.BoolVar(&c.verbose, "verbose", false, "Enable verbose debugging output")

	flags.FlagMerge(c.flags, c.pathFlags.Flags())

	return c, nil
}

func (c *cmd) Help() string { return commands.Usage(agentCommandHelp, c.flags) }

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
	if err := data.DecodeJSON(cfg.DebugDirectoryPath, "agent"); err != nil {
		hclog.L().Error("failed to decode agent.json", "error", err)
		return 1
	}
	if err := data.DecodeJSON(cfg.DebugDirectoryPath, "members"); err != nil {
		hclog.L().Error("failed to decode members.json", "error", err)
		return 1
	}
	hclog.L().Debug("successfully read in agent cmd information from bundle")

	var result string

	switch {
	case c.summary:
		result = agentSummary(data)
	case c.config:
		result = agentConfig(data)
	default:
		result = c.Help()
	}
	c.ui.Output(result)
	return 0
}

func agentSummary(bundle read.Debug) string {
	return bundle.Agent.Summary()
}

func agentConfig(bundle read.Debug) string {
	return bundle.Agent.AgentConfigFull()
}

const synopsis = `Debug bundle agent.json information parsing`
const agentCommandHelp = `The agent flag will ingest the agent.json and parse for additional information pertaining to the agent.
This includes:
  - Consul Versioning
  - Server Agent Status
  - Client Agent Status
  - Known Serf Members
  - Current Raft Configuration`

// Cmd is the agent subcommand
//var (
//	agentCmd = &cobra.Command{
//		Use:   "agent",
//		Short: "Debug bundle agent.json information parsing.",
//		Long: `The agent flag will ingest the agent.json and parse for additional information pertaining to the agent.
//This includes:
//  - Consul Versioning
//  - Server Agent Status
//  - Client Agent Status
//  - Known Serf Members
//  - Current Raft Configuration`,
//		Args: func(agentCmd *cobra.Command, args []string) error {
//			if len(args) != 0 {
//				return fmt.Errorf("unknown subcommand %q", args[0])
//			}
//			return nil
//		},
//		PreRunE: func(cmd *cobra.Command, args []string) error {
//			if _, ok := os.LookupEnv(config.EnvDebugPath); ok {
//				envPath := os.Getenv(config.EnvDebugPath)
//				envPath = strings.TrimSuffix(envPath, "/")
//				if _, err := os.Stat(envPath); os.IsNotExist(err) {
//					return fmt.Errorf("directory does not exists: %s - %v\n", envPath, err)
//				} else {
//					config.DebugPath = envPath
//					if config.Verbose {
//						log.Printf("using environment variable CONSUL_DEBUG_PATH - %s\n", config.DebugPath)
//					}
//				}
//			} else {
//				config.DebugPath = viper.GetString("config.DebugPath")
//				if config.Verbose {
//					log.Printf("using config.yaml debug path setting - %s\n", config.DebugPath)
//				}
//			}
//			if config.DebugPath != "" {
//				if config.Verbose {
//					log.Printf("debug-path:  '%s'\n", config.DebugPath)
//				}
//				if err := config.DebugBundle.DecodeJSON(config.DebugPath, "agent"); err != nil {
//					return fmt.Errorf("failed to decode agent.json %v", err)
//				}
//				if err := config.DebugBundle.DecodeJSON(config.DebugPath, "members"); err != nil {
//					return fmt.Errorf("failed to decode members.json %v", err)
//				}
//				if config.Verbose {
//					log.Printf("successfully read-in bundle from:  '%s'\n", config.DebugPath)
//				}
//			} else {
//				return fmt.Errorf("debug-path is null")
//			}
//			return nil
//		},
//		RunE: func(cmd *cobra.Command, args []string) error {
//			summary, _ := cmd.Flags().GetBool("summary")
//			c, _ := cmd.Flags().GetBool("config")
//			// Get Metrics object
//			// var members = config.DebugBundle.Members
//			var agentFile = fmt.Sprintf(config.DebugPath + "/agent.json")
//			switch {
//			case summary:
//				if config.Verbose {
//					log.Printf("agent summary: configuration rendered from: %s\n", agentFile)
//				}
//				fmt.Println(config.DebugBundle.Agent.AgentSummary())
//			case c:
//				if config.Verbose {
//					log.Printf("agent hcl config: configuration rendered from: %s\n", agentFile)
//				}
//				fmt.Println(config.DebugBundle.Agent.AgentConfigFull())
//			default:
//				//fmt.Printf("Agent Configuration Summary:\n")
//				//fmt.Println("----------------------")
//				//fmt.Println("Serf Member Count (wan members):", len(members))
//				//agent.AgentSummary()
//				//fmt.Printf("debug file: %s\n", agentFile)
//				if err := cmd.Usage(); err != nil {
//					return err
//				}
//			}
//			return nil
//		},
//	}
//)
//
//func init() {
//	agentCmd.Flags().BoolP("summary", "s", false, "Retrieve agent configuration summary.")
//	agentCmd.Flags().Bool("config", false, "Retrieve agent configuration in HCL format")
//}
