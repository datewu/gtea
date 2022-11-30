package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Logger holds the output destination that the log entries
// will be written to, the minimum severity level that log
// entries will be written for, and a mutex for concurrent writes.
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

var defaultLogger = Default()

// New create a new Logger
func New(out io.Writer, minLevel Level) *Logger {
	l := &Logger{
		out:      out,
		minLevel: minLevel,
	}
	defaultLogger = l
	return l
}

// Default write to stdout with the minimum level set to LevelDebug.
func Default() *Logger {
	return &Logger{
		out:      os.Stdout,
		minLevel: LevelDebug,
	}
}

// Debug writes a log entry at LevelDebug to the output destination.
func (l *Logger) Debug(msg string, props ...map[string]any) {
	l.print(LevelDebug, msg, props...)
}

// Info writes a log entry at LevelInfo to the output destination.
func (l *Logger) Info(msg string, props ...map[string]any) {
	l.print(LevelInfo, msg, props...)
}

// Err writes a log entry at LevelError to the output destination.
func (l *Logger) Err(err error, props ...map[string]any) {
	l.print(LevelError, err.Error(), props...)
}

// Trace writes a log entry at TraceLevel to the output destination.
func (l *Logger) Trace(err error, props ...map[string]any) {
	l.print(LevelTrace, err.Error(), props...)
}

// Fatal writes a log entry at LevelFatal to the output destination
// and exit 1.
func (l *Logger) Fatal(err error, props ...map[string]any) {
	l.print(LevelFatal, err.Error(), props...)
	os.Exit(1)
}

func (l *Logger) print(level Level, msg string, props ...map[string]any) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}
	aux := struct {
		Level      string           `json:"level"`
		Time       string           `json:"time"`
		Message    string           `json:"message"`
		Properties []map[string]any `json:"properties,omitempty"`
		Trace      string           `json:"trace,omitempty"`
	}{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    msg,
		Properties: props,
	}

	if level >= LevelTrace {
		aux.Trace = string(debug.Stack())
		l.mu.Lock()
		fmt.Fprintln(l.out, "Trace:")
		fmt.Fprintln(l.out, aux.Trace)
		fmt.Fprintln(l.out)
		l.mu.Unlock()
	}
	var line []byte
	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}

// Write satisfies the io.Writer interface by call writeRaw.
func (l *Logger) Write(msg []byte) (int, error) {
	return l.writeRaw(msg)
}

// writeRaw writes a raw log entry at LevelError to the output destination.
func (l *Logger) writeRaw(msg []byte) (int, error) {
	return l.print(LevelError, string(msg), nil)
}
