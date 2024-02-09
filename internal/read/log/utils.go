package log

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"sort"
)

// AggregateEntries aggregates log entries by method and minute
func AggregateEntries(entries []LogEntry) map[string]map[string]int {
	counts := make(map[string]map[string]int) // method -> minute -> count

	for _, entry := range entries {
		method := entry.Method
		minute := entry.Timestamp.Format("2006-01-02 15:04")
		if _, ok := counts[method]; !ok {
			counts[method] = make(map[string]int)
		}
		counts[method][minute]++
	}

	return counts
}

// RPCCounts generate the aggregated counts
func RPCCounts(counts map[string]map[string]int) string {
	// Build RPC Counts Title
	result := []string{fmt.Sprintf("Timestamp\x1fMethod\x1fCounts\x1f")}
	for method, minutes := range counts {
		var keys []string
		for k := range minutes {
			keys = append(keys, k)
		}
		sort.Strings(keys) // Sort minutes for consistent output
		for _, k := range keys {
			result = append(result, fmt.Sprintf("%s\x1f%s\x1f%d\x1f", method, k, minutes[k]))
		}
	}
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}
