package log

import (
	"bufio"
	common "consul-debug-read/internal/read"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const RPCMethodRegex = `rpc_server_call: method=([^\s]+)`

var (
	timestampRegex = regexp.MustCompile(common.TimeStampRegex)
	methodRegex    = regexp.MustCompile(RPCMethodRegex)
)

// ParseRPCMethods parses a log file and returns a slice of LogEntry
func ParseRPCMethods(filePath, filterMethod string) ([]Entry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 3 || !timestampRegex.MatchString(parts[0]) {
			continue // skip lines that don't match the expected format
		}

		timestamp, _ := time.Parse(time.RFC3339, parts[0])
		matches := methodRegex.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue // skip lines without a method
		}
		method := matches[1]

		if filterMethod == "" || filterMethod == method {
			entries = append(entries, Entry{
				Timestamp: timestamp,
				Method:    method,
			})
		}
	}

	return entries, scanner.Err()
}

// ParseLog parses a log file for entries of a specified level and source, then returns a slice of LogEntry.
func ParseLog(filePath, levelFilter, sourceFilter string, startTime, endTime time.Time) ([]LogEntry, error) {
	// Default to INFO if no level is specified
	if levelFilter == "" {
		levelFilter = InfoLevel
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logEntries []LogEntry
	scanner := bufio.NewScanner(file)
	// Modify the regex to be prepared for dynamic source filtering, if provided
	logRegexPattern := fmt.Sprintf(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z) \[(%s)\] ([^\:]+): (.+)`, levelFilter)
	logRegex := regexp.MustCompile(logRegexPattern)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := logRegex.FindStringSubmatch(line); matches != nil {
			timestamp, _ := time.Parse(time.RFC3339, matches[1])
			level := matches[2]
			source := matches[3]
			message := matches[4]

			// Time range filtering
			if !startTime.IsZero() && timestamp.Before(startTime) {
				continue // Skip entries before the start time
			}
			if !endTime.IsZero() && timestamp.After(endTime) {
				continue // Skip entries after the end time
			}

			// Check if the current log's source matches the specified source filter
			// If sourceFilter is empty, include all sources. Otherwise, filter by the specified source.
			if sourceFilter == "" || source == sourceFilter {
				logEntries = append(logEntries, LogEntry{
					Timestamp: timestamp,
					Level:     level,
					Source:    source,
					Message:   message,
				})
			}
		}
	}

	return logEntries, scanner.Err()
}
