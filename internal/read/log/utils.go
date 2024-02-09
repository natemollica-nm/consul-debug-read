package log

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"sort"
)

// MethodCount represents the count of a method at a specific minute
type MethodCount struct {
	Method string
	Minute string
	Count  int
}

// AggregateEntries aggregates log entries by method and minute
func AggregateEntries(entries []Entry) map[string]map[string]int {
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
	result := []string{fmt.Sprintf("Method\x1fMinute-Interval\x1fCounts\x1f")}

	var methodCounts []MethodCount
	// Flatten counts into a slice of MethodCount
	for method, minutes := range counts {
		for minute, count := range minutes {
			methodCounts = append(methodCounts, MethodCount{
				Method: method,
				Minute: minute,
				Count:  count,
			})
		}
	}

	// Sort the slice by count descending
	sort.Slice(methodCounts, func(i, j int) bool {
		return methodCounts[i].Count > methodCounts[j].Count
	})

	// append sorted results
	for _, mc := range methodCounts {
		result = append(result, fmt.Sprintf("%s\x1f%s\x1f%d\x1f", mc.Method, mc.Minute, mc.Count))
	}
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}
