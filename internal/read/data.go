package read

type ReaderConfig struct {
	ConfigFile         string `yaml:"config"`
	DebugDirectoryPath string `yaml:"current-debug-path"`
	PathRenderedFrom   string `yaml:"Using"`
	DebugEnvVarSetting string `yaml:"CONSUL_DEBUG_PATH"`
}

func DefaultReaderConfig() *ReaderConfig {
	var renderedFrom string
	if EnvVarPathSetting != "" {
		renderedFrom = "env:CONSUL_DEBUG_PATH"
	} else {
		renderedFrom = "file:config.yaml"
	}
	return &ReaderConfig{
		// Create Default Configuration File
		ConfigFile:         DebugReadConfigFullPath,
		DebugDirectoryPath: CurrentDir,
		DebugEnvVarSetting: EnvVarPathSetting,
		PathRenderedFrom:   renderedFrom,
	}
}

type Debug struct {
	Agent   Agent
	Members []Member
	Metrics Metrics
	Host    Host
	Index   Index
}
