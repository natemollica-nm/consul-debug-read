package read

type ReaderConfig struct {
	ConfigFile         string `yaml:"configFile"`
	DebugDirectoryPath string `yaml:"debugDirectoryPath"`
	PathRenderedFrom   string `yaml:"pathRenderedFrom"`
	DebugEnvVarSetting string `yaml:"CONSUL_DEBUG_PATH"`
}

func DefaultReaderConfig() *ReaderConfig {
	var renderedFrom string
	if EnvVarPathSetting != "" {
		renderedFrom = "env:CONSUL_DEBUG_PATH"
	} else {
		renderedFrom = "default: <current-directory>"
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
	Metrics Metrics
	Host    Host
	Index   Index
}
