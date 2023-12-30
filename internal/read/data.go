package read

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
)

var (
	DefaultDBPath = fmt.Sprintf("%s/%s", DebugReadConfigDirPath, ConsulDebugDb)
)

type Backend struct {
	id string
	DB *bolt.DB
}

type ReaderConfig struct {
	DebugDirectoryPath string `yaml:"current-debug-path"`
}

type Debug struct {
	Agent   Agent
	Members []Member
	Metrics Metrics
	Host    Host
	Index   Index
	Backend *Backend
}

func NewBackend() *Backend {
	_, uuid := generateUUID()
	store, _ := initMemDB()
	return &Backend{
		id: uuid,
		DB: store,
	}
}

func initMemDB() (*bolt.DB, error) {
	if _, err := os.Stat(DefaultDBPath); err == nil {
		if err = os.Remove(DefaultDBPath); err != nil {
			return nil, err
		}
	}
	// Open the BoltDB file
	// It will be created if it doesn't exist
	db, err := bolt.Open(DefaultDBPath, 0666, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}

	return db, nil
}
