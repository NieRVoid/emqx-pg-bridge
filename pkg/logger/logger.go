package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger provides a simple logging interface
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	level       LogLevel
	format      string
}

// LogLevel represents logging levels
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	ErrorLevel
)

// NewLogger creates a new Logger
func NewLogger(level, format string) *Logger {
	infoLogger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLogger := log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	
	logLevel := InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		logLevel = DebugLevel
	case "info":
		logLevel = InfoLevel
	case "error":
		logLevel = ErrorLevel
	}
	
	return &Logger{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		debugLogger: debugLogger,
		level:       logLevel,
		format:      format,
	}
}

// formatMessage formats a log message with key-value pairs
func formatMessage(msg string, keyValues ...interface{}) string {
	if len(keyValues) == 0 {
		return msg
	}
	
	builder := strings.Builder{}
	builder.WriteString(msg)
	
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			builder.WriteString(", ")
			builder.WriteString(keyValues[i].(string))
			builder.WriteString("=")
			builder.WriteString(fmt.Sprintf("%v", keyValues[i+1]))
		}
	}
	
	return builder.String()
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, keyValues ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLogger.Println(formatMessage(msg, keyValues...))
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, keyValues ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLogger.Println(formatMessage(msg, keyValues...))
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, keyValues ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLogger.Println(formatMessage(msg, keyValues...))
	}
}

// Fatal logs an error message and exits
func (l *Logger) Fatal(msg string, keyValues ...interface{}) {
	l.errorLogger.Fatalln(formatMessage(msg, keyValues...))
}