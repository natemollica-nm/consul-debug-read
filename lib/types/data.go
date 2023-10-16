package types

import (
	"encoding/json"
	"fmt"
	"github.com/ryanuber/columnize"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strings"
)

type Debug struct {
	Agent   Agent
	Members []Member
	Metrics Metrics
	Host    Host
}

// ByMemberName sorts members by name with a stable sort.
//
// 1. servers go at the top
type ByMemberName []Member

func (m ByMemberName) Len() int      { return len(m) }
func (m ByMemberName) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m ByMemberName) Less(i, j int) bool {
	tags_i := m[i].Tags
	tags_j := m[j].Tags

	// put role=consul first
	switch {
	case tags_i.Role == "consul" && tags_j.Role != "consul":
		return true
	case tags_i.Role != "consul" && tags_j.Role == "consul":
		return false
	}

	// then by datacenter
	switch {
	case tags_i.Dc < tags_j.Dc:
		return true
	case tags_i.Dc > tags_j.Dc:
		return false
	}

	// finally by name
	return m[i].Name < m[j].Name
}

// MembersStandard is used to dump the most useful information about nodes
// in a more human-friendly format
func (b *Debug) MembersStandard() string {
	result := make([]string, 0, len(b.Members))
	header := "Node\x1fAddress\x1fStatus\x1fType\x1fBuild\x1fProtocol\x1fDC"
	result = append(result, header)
	sort.Sort(ByMemberName(b.Members))
	for _, member := range b.Members {
		tags := member.Tags

		addr := net.TCPAddr{IP: net.ParseIP(member.Addr), Port: int(member.Port)}
		protocol := tags.Vsn
		build := tags.Build
		if build == "" {
			build = "< 0.3"
		} else if idx := strings.Index(build, ":"); idx != -1 {
			build = build[:idx]
		}
		nameIdx := strings.Index(member.Name, ".")
		name := member.Name[:nameIdx]

		var statusString string
		switch {
		case member.Status == 0:
			statusString = "None"
		case member.Status == 1:
			statusString = "Alive"
		case member.Status == 2:
			statusString = "Leaving"
		case member.Status == 3:
			statusString = "Left"
		case member.Status == 4:
			statusString = "Failed"
		}
		switch tags.Role {
		case "node":
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1fclient\x1f%s\x1f%s\x1f%s",
				name, addr.String(), statusString, build, protocol, tags.Dc)
			result = append(result, line)

		case "consul":
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1fserver\x1f%s\x1f%s\x1f%s",
				name, addr.String(), statusString, build, protocol, tags.Dc)
			result = append(result, line)

		default:
			line := fmt.Sprintf("%s\x1f%s\x1f%s\x1funknown\x1f\x1f\x1f",
				name, addr.String(), statusString)
			result = append(result, line)
		}
	}

	output, _ := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}

func (b *Debug) BundleSummary() {
	b.Agent.AgentSummary()
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
