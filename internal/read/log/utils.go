package log

import (
	"fmt"
	"github.com/ryanuber/columnize"
	"sort"
	"strings"
	"time"
)

// RPCMethodCount represents the count of a method at a specific minute
type RPCMethodCount struct {
	Method string
	Minute string
	Count  int
}

type EntryCount struct {
	Minute string
	Key    string
	Count  int
}

type FormattedEntry struct {
	Minute string
	Key    string // This could represent the method, message, or any field used for aggregation
	Source string // The source of the log entry
	Count  int    // The number of occurrences
}

type AggregateEntry struct {
	Count  int
	Source string
}

type EntrySelector func(entry LogEntry) string

func SourceSelect(entry LogEntry) string {
	return entry.Source
}

func MessageSelect(entry LogEntry) string {
	return entry.Message
}

// AggregateRPCEntries aggregates log entries by method and minute
func AggregateRPCEntries(entries []Entry) map[string]map[string]int {
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
	// Build RPC FormatCounts Title
	result := []string{fmt.Sprintf("Method\x1fMinute-Interval\x1fCounts\x1f")}

	var methodCounts []RPCMethodCount
	// Flatten counts into a slice of RPCMethodCount
	for method, minutes := range counts {
		for minute, count := range minutes {
			methodCounts = append(methodCounts, RPCMethodCount{
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

func AggregateLogEntries(entries []LogEntry, level string, selector EntrySelector) map[string][]AggregateEntry {
	aggregated := make(map[string][]AggregateEntry)

	for _, entry := range entries {
		if entry.Level != level {
			continue
		}
		key := selector(entry) + "|" + entry.Timestamp.Format("2006-01-02 15:04")

		// Find or initialize the aggregated entry
		found := false
		for i, agg := range aggregated[key] {
			if agg.Source == entry.Source {
				aggregated[key][i].Count++
				found = true
				break
			}
		}
		if !found {
			aggregated[key] = append(aggregated[key], AggregateEntry{Count: 1, Source: entry.Source})
		}
	}

	return aggregated
}

func FormatCounts(aggregated map[string][]AggregateEntry, selector string) string {
	var result []string
	// Build result with new struct
	entryType := capitalize(selector)
	if entryType == "Message" {
		result = []string{fmt.Sprintf("Timestamp\x1fCounts\x1fSource\x1f%s\x1f", entryType)}
	} else {
		result = []string{fmt.Sprintf("Minute-Interval\x1fCounts\x1f%s\x1f", entryType)}
	}
	var entries []FormattedEntry

	// Flatten counts into a slice of EntryCount
	for key, aggEntries := range aggregated {
		parts := strings.Split(key, "|") // Split key (composite key) to extract key and minute if using composite key approach
		minute := parts[1]               // Adjust indexing based on actual composite key structure
		for _, aggregate := range aggEntries {
			entries = append(entries, FormattedEntry{
				Minute: minute,
				Key:    parts[0], // Or other logic to extract original key from compositeKey
				Source: aggregate.Source,
				Count:  aggregate.Count,
			})
		}
	}

	// Sort by timestamp for message-count and by count for source-count
	if entryType == "Message" {
		sort.Slice(entries, func(i, j int) bool {
			// Sort by timestamp (Minute) if the selector is for messages
			return entries[i].Minute < entries[j].Minute
		})
	} else {
		// Otherwise, sort by count descending
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Count > entries[j].Count
		})
	}

	// Append sorted results
	for _, mc := range entries {
		// Truncate results for display if necessary
		out := mc.Key
		if len(out) > 150 {
			out = out[:47] + "..."
		}
		if mc.Count > 1 {
			if entryType == "Message" {
				result = append(result, fmt.Sprintf("%s\x1f%d\x1f%s\x1f%s\x1f", mc.Minute, mc.Count, mc.Source, out))
			} else {
				result = append(result, fmt.Sprintf("%s\x1f%d\x1f%s\x1f", mc.Minute, mc.Count, out))
			}
		}

	}

	// Use columnize to format the results into a string with columns
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}

// FormatLog generates a formatted string of log entries, truncating messages if they are too long.
func FormatLog(entries []LogEntry) string {
	// Build Error Logs Title
	result := []string{fmt.Sprintf("Timestamp\x1fSource\x1fMessage\x1f")}

	// Define the maximum message length
	maxMessageLength := 200 // Adjust as needed

	// Sort the slice by timestamp descending (or any other criteria you prefer)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	// Append sorted and possibly truncated results
	for _, entry := range entries {
		formattedTimestamp := entry.Timestamp.Format(time.RFC3339) // Format timestamp as you like
		message := entry.Message
		// Truncate the message if it exceeds the maximum length
		if len(message) > maxMessageLength {
			message = message[:maxMessageLength-3] + "..." // Add ellipsis to indicate truncation
		}
		result = append(result, fmt.Sprintf("%s\x1f%s\x1f%s\x1f", formattedTimestamp, entry.Source, message))
	}

	// Use columnize to format the results into a string with columns
	output := columnize.Format(result, &columnize.Config{Delim: string([]byte{0x1f}), Glue: " "})
	return output
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
