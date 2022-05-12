package log

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Logger for Cloud Function
// https://cloud.google.com/functions/docs/monitoring/logging?hl=ja

const (
	GCP_TRACE_HEADER = "X-Cloud-Trace-Context"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	SUPRESS
)

var log *Logger
var logLevel int

func init() {
	appPhase := os.Getenv("PHASE")
	if appPhase == "" {
		panic("$PHASE required")
	}

	switch appPhase {
	case "local":
		logLevel = DEBUG
	case "develop":
		logLevel = INFO
	case "production":
		logLevel = WARN
	case "test":
		logLevel = SUPRESS
	default:
		logLevel = DEBUG
	}
}

type Logger struct {
	req   *http.Request
	level int
}

func NewLogger(req *http.Request) *Logger {
	return &Logger{
		req:   req,
		level: logLevel,
	}
}

type Entry struct {
	Message        string `json:"message"`
	Severity       string `json:"severity,omitempty"`
	Trace          string `json:"logging.googleapis.com/trace,omitempty"`
	WebhookEventID string `json:"webhook_event_id,omitempty"`

	// Logs Explorer allows filtering and display of this as `jsonPayload.component`.
	Component string `json:"component,omitempty"`
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if DEBUG >= l.level {
		l.writeToStdout("DEBUG", fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if INFO >= l.level {
		l.writeToStdout("INFO", fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if WARN >= l.level {
		l.writeToStdout("WARN", fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if ERROR >= l.level {
		l.writeToStdout("ERROR", fmt.Sprintf(format, v...))
	}
}

func (l *Logger) writeToStdout(severity, msg string) {
	logEntry := Entry{
		Severity: severity,
		Message:  msg,
		Trace:    l.req.Header.Get(GCP_TRACE_HEADER),
	}

	b, _ := json.Marshal(logEntry)
	fmt.Println(string(b))
}
