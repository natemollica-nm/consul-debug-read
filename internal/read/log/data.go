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

type JsonLogEntry struct {
	Timestamp string `json:"@timestamp"`
	Module    string `json:"@module"`
	Level     string `json:"@level"`
	Message   string `json:"@message"`
	Thread    int    `json:"thread"`
}
