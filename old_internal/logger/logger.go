package logger

import (
	"fmt"
)

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
}

type Log struct{}

func (l *Log) Info(msg string) {
	fmt.Printf("INFO: %s\n", msg)
}

func (l *Log) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("INFO: %s\n", msg)
}

func (l *Log) Warn(msg string) {
	fmt.Printf("WARN: %s\n", msg)
}

func (l *Log) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("WARN: %s\n", msg)
}

func (l *Log) Error(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
}

func (l *Log) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("ERROR: %s\n", msg)
}

func New() Logger {
	return &Log{}
}
