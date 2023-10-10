package debug_read

type Config struct {
	Customer  string
	DebugPath string
	Platform  string
	UseCases  []string
}

func DefaultConfig() *Config {
	conf := &Config{
		Customer:  "test",
		DebugPath: "customers/test",
		Platform:  "VM",
		UseCases:  []string{"service-discovery", "service-mesh", "key-value"},
	}
	return conf
}
