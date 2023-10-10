package commands

import (
	"flag"
	"github.com/natemollica-nm/consul-debug-reader/internal"
)

type debugFlags struct {
	debugPath stringValue
	customer  stringValue
	platform  stringValue
	useCases  stringValue
}

func flags() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Var(&)
}

func (f *debugFlags) runConfig() (*Config, error) {
	cfg := debug.DefaultConfig()
	f.mergeOntoConfig(cfg)

	return *cfg
}
func (f *debugFlags) mergeOntoConfig(c *runConfig) {
	f.debugPath.Merge(&c.debugPath)
	f.token.Merge(&c.Token)
	f.tokenFile.Merge(&c.TokenFile)
	f.caFile.Merge(&c.TLSConfig.CAFile)
}
// stringValue provides a flag value that's aware if it has been set.
type stringValue struct {
	v *string
}

// merge will overlay this value if it has been set.
func (s *stringValue) Merge(onto *string) {
	if s.v != nil {
		*onto = *(s.v)
	}
}

// Set implements the flag.Value interface.
func (s *stringValue) Set(v string) error {
	if s.v == nil {
		s.v = new(string)
	}
	*(s.v) = v
	return nil
}

// String implements the flag.Value interface.
func (s *stringValue) String() string {
	var current string
	if s.v != nil {
		current = *(s.v)
	}
	return current
}
