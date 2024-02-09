package log

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Method    string
}

// ParseLogFile parses a log file and returns a slice of LogEntry
func ParseLogFile(filePath, filterMethod string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	timestampRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)
	methodRegex := regexp.MustCompile(`method=([^\s]+)`)

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
			entries = append(entries, LogEntry{
				Timestamp: timestamp,
				Method:    method,
			})
		}
	}

	return entries, scanner.Err()
}
