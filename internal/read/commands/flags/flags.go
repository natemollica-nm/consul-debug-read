package flags

import (
	"flag"
)

type DebugReadFlags struct {
	DebugFilePath stringValue
}

func (f *DebugReadFlags) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	return fs
}

func FlagMerge(dst, src *flag.FlagSet) {
	if dst == nil {
		panic("dst cannot be nil")
	}
	if src == nil {
		return
	}
	src.VisitAll(func(f *flag.Flag) {
		dst.Var(f.Value, f.Name, f.Usage)
	})
}

// stringValue provides a flag value that's aware if it has been set.
type stringValue struct {
	v *string
}

// Merge will overlay this value if it has been set.
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
