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

const (
	RPCMethodRegex = `rpc_server_call: method=([^\s]+)`
)

var (
	timestampLayouts = []string{
		time.RFC3339,
		time.UnixDate,
		"2006-01-02T15:04:05Z0700",            // Without colon in timezone offset
		"2006-01-02T15:04:05",                 // Without timezone info
		"2006-01-02T15:04:05.999Z07:00",       // With milliseconds and timezone
		"2006-01-02T15:04:05.999Z0700",        // With milliseconds and without colon in timezone
		"2006-01-02T15:04:05.999999Z07:00",    // With microseconds and timezone
		"2006-01-02T15:04:05.999999Z0700",     // With microseconds and without colon in timezone
		"2006-01-02T15:04:05.999999999Z07:00", // With nanoseconds and timezone
		"2006-01-02T15:04:05.999999999Z0700",  // With nanoseconds and without colon in timezone
		"2006-01-02T15:04Z07:00",              // Without seconds, with timezone
		"2006-01-02T15:04Z0700",               // Without seconds, without colon in timezone
		"2006-01-02",                          // Date only, no time
	}
	// timestampRegex is used for RPC Counts and `$` is manually added as the
	// common.TimeStampRegex purposefully leaves this appended reg key out to support
	// parsing full log entry timestamps.
	timestampRegex = regexp.MustCompile(common.TimeStampRegex + "$")
	methodRegex    = regexp.MustCompile(RPCMethodRegex)
)

// ParseRPCMethods parses a log file and returns a slice of LogEntry
func ParseRPCMethods(filePath, filterMethod string) ([]Entry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cleanup := func(err error) error {
		_ = file.Close()
		return err
	}
	var entries []Entry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var timestamp time.Time
		line := scanner.Text()
		parts := strings.Fields(line)

		timestamp, err = parseTimestamp(parts[0])
		if err != nil {
			return []Entry{}, err
		}
		if len(parts) < 3 || !timestampRegex.MatchString(parts[0]) {
			continue // skip lines that don't match the expected format
		}
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
	if err = file.Close(); err != nil {
		return []Entry{}, cleanup(err)
	}
	return entries, scanner.Err()
}

func parseTimestamp(timestampStr string) (time.Time, error) {
	var timestamp time.Time
	var err error
	for _, layout := range timestampLayouts {
		timestamp, err = time.Parse(layout, timestampStr)
		if err == nil {
			return timestamp, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestampStr)
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
	cleanup := func(err error) error {
		_ = file.Close()
		return err
	}

	var logEntries []LogEntry
	scanner := bufio.NewScanner(file)
	// Modify the regex to be prepared for dynamic source filtering, if provided
	logRegexPattern := fmt.Sprintf(`%s \[(%s)\] ([^\:]+): (.+)`, common.TimeStampRegex, levelFilter)
	logRegex := regexp.MustCompile(logRegexPattern)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := logRegex.FindStringSubmatch(line); matches != nil {
			var timestamp time.Time
			timestamp, err = parseTimestamp(matches[1])
			if err != nil {
				return []LogEntry{}, err
			}
			level := matches[3]
			source := matches[4]
			message := matches[5]

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
	if err = file.Close(); err != nil {
		return []LogEntry{}, cleanup(err)
	}
	return logEntries, scanner.Err()
}
