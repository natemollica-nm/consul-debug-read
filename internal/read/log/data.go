package log

import "time"

const (
	InfoLevel  = "INFO"
	ErrorLevel = "ERROR"
	DebugLevel = "DEBUG"
	TraceLevel = "TRACE"
	WarnLevel  = "WARN"
)

// LogEntry represents a single log entry
type Entry struct {
	Timestamp time.Time
	Method    string
}
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Source    string
	Message   string
}
