package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level represents log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelStrings = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// Logger is a simple structured logger
type Logger struct {
	level  Level
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New(INFO)
}

// New creates a new logger
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// SetLevel sets the logging level
func SetLevel(level Level) {
	defaultLogger.level = level
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := levelStrings[level]
	message := fmt.Sprintf(format, args...)

	l.logger.Printf("[%s] %s: %s", timestamp, levelStr, message)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	defaultLogger.log(DEBUG, format, args...)
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	defaultLogger.log(INFO, format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	defaultLogger.log(WARN, format, args...)
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	defaultLogger.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func Fatal(format string, args ...interface{}) {
	defaultLogger.log(FATAL, format, args...)
}
