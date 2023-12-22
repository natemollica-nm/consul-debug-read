package read

type ReaderConfig struct {
	DebugDirectoryPath string `yaml:"current-debug-path"`
}

type Debug struct {
	Agent   Agent
	Members []Member
	Metrics Metrics
	Host    Host
	Index   Index
}
