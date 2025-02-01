package micromigrations

import (
	"fmt"
	"log"
)

type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type NoopLogger struct{}

func (nl *NoopLogger) Debug(msg string, fields ...any) {}
func (nl *NoopLogger) Info(msg string, fields ...any)  {}
func (nl *NoopLogger) Warn(msg string, fields ...any)  {}
func (nl *NoopLogger) Error(msg string, fields ...any) {}

func NewNoopLogger() Logger {
	return &NoopLogger{}
}

type GenericLogger struct {
	logger *log.Logger
}

func (l *GenericLogger) print(level string, msg string, fields ...any) {
	l.logger.Printf(
		"[%s] %s",
		level,
		fmt.Sprintf(msg, fields...),
	)
}

func (l *GenericLogger) Debug(msg string, fields ...any) {
	l.print("DEBUG", msg, fields...)
}

func (l *GenericLogger) Info(msg string, fields ...any) {
	l.print("INFO", msg, fields...)
}

func (l *GenericLogger) Warn(msg string, fields ...any) {
	l.print("WARN", msg, fields...)
}

func (l *GenericLogger) Error(msg string, fields ...any) {
	l.print("ERROR", msg, fields...)
}

func NewGenericLogger() Logger {
	return &GenericLogger{
		logger: log.Default(),
	}
}

func NewGenericWithLogger(log *log.Logger) Logger {
	return &GenericLogger{
		logger: log,
	}
}
