package logger

import (
	"io"
	"log"
	"os"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type StandardLogger struct {
	logger *log.Logger
}

func NewLogger() Logger {
	return &StandardLogger{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func NewSilentLogger() Logger {
	return &StandardLogger{
		logger: log.New(io.Discard, "", 0),
	}
}

func (l *StandardLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}
