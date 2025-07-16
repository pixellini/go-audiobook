package logger

import (
	"fmt"
	"log"
	"os"
)

// StandardLogger implements the Logger interface using standard library logging
type StandardLogger struct {
	logger  *log.Logger
	verbose bool
}

// NewStandardLogger creates a new standard logger
func NewStandardLogger(verbose bool) *StandardLogger {
	return &StandardLogger{
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		verbose: verbose,
	}
}

// Info logs an info message
func (l *StandardLogger) Info(msg string) {
	l.logger.Printf("[INFO] %s", msg)
}

// Infof logs a formatted info message
func (l *StandardLogger) Infof(format string, args ...interface{}) {
	l.logger.Printf("[INFO] "+format, args...)
}

// Warn logs a warning message
func (l *StandardLogger) Warn(msg string) {
	l.logger.Printf("[WARN] %s", msg)
}

// Warnf logs a formatted warning message
func (l *StandardLogger) Warnf(format string, args ...interface{}) {
	l.logger.Printf("[WARN] "+format, args...)
}

// Error logs an error message
func (l *StandardLogger) Error(msg string) {
	l.logger.Printf("[ERROR] %s", msg)
}

// Errorf logs a formatted error message
func (l *StandardLogger) Errorf(format string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+format, args...)
}

// SimpleLogger is a basic logger that prints to stdout (for backward compatibility)
type SimpleLogger struct{}

// NewSimpleLogger creates a simple logger that mimics the current fmt.Println behavior
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{}
}

// Info prints an info message
func (l *SimpleLogger) Info(msg string) {
	fmt.Println(msg)
}

// Infof prints a formatted info message
func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Warn prints a warning message
func (l *SimpleLogger) Warn(msg string) {
	fmt.Printf("WARNING: %s\n", msg)
}

// Warnf prints a formatted warning message
func (l *SimpleLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf("WARNING: "+format+"\n", args...)
}

// Error prints an error message
func (l *SimpleLogger) Error(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
}

// Errorf prints a formatted error message
func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}
