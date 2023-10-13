package types

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Debug struct {
	Agent   Agent
	Members []Member
	Metrics Metrics
	Host    Host
}

func (b *Debug) DecodeJSON(debugPath string) error {
	configs := []string{"agent.json", "members.json", "metrics.json", "host.json"}
	agent, _ := os.Open(fmt.Sprintf("%s/%s", debugPath, configs[0]))
	members, _ := os.Open(fmt.Sprintf("%s/%s", debugPath, configs[1]))
	metrics, _ := os.Open(fmt.Sprintf("%s/%s", debugPath, configs[2]))
	host, _ := os.Open(fmt.Sprintf("%s/%s", debugPath, configs[3]))
	agentDecoder := json.NewDecoder(agent)
	memberDecoder := json.NewDecoder(members)
	metricsDecoder := json.NewDecoder(metrics)
	hostDecoder := json.NewDecoder(host)

	cleanup := func(err error) error {
		_ = agent.Close()
		_ = members.Close()
		_ = metrics.Close()
		_ = host.Close()
		return err
	}

	log.Printf("Parsing %s", agent.Name())
	for {
		var agentConfig Agent
		err := agentDecoder.Decode(&agentConfig)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding %s | file: %v", err, agent.Name())
			return err
		}
		b.Agent = agentConfig
	}

	log.Printf("Parsing %s", members.Name())
	for {
		var membersList []Member
		err := memberDecoder.Decode(&membersList)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding %s | file: %v", err, members.Name())
			return err
		}
		b.Members = membersList
	}

	log.Printf("Parsing %s", metrics.Name())
	for {
		var metric Metric
		err := metricsDecoder.Decode(&metric)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding %s | file: %v", err, metrics.Name())
			return err
		}
		b.Metrics.Metrics = append(b.Metrics.Metrics, metric)
	}

	log.Printf("Parsing %s", host.Name())
	for {
		var hostObject Host
		err := hostDecoder.Decode(&hostObject)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error decoding %s | file: %v", err, metrics.Name())
			return err
		}
		b.Host = hostObject
	}

	if err := agent.Close(); err != nil {
		return cleanup(err)
	}
	if err := members.Close(); err != nil {
		return cleanup(err)
	}
	if err := metrics.Close(); err != nil {
		return cleanup(err)
	}
	if err := host.Close(); err != nil {
		return cleanup(err)
	}

	return nil
}
