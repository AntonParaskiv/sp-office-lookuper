package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	DestinationGelf = "gelf"
	RequestIDKey    = "rid"
)

type flusher interface {
	Flush()
}

type Logger struct {
	logger       *logrus.Logger
	serviceName  string
	level        int
	destination  string
	host         string
	port         int
	flusherHooks []flusher
}

func NewLogger(opts ...Option) (*Logger, error) {
	logger := &Logger{
		logger: logrus.New(),
	}

	if len(opts) > 0 {
		for _, o := range opts {
			o(logger)
		}
	}

	switch logger.destination {
	case DestinationGelf:
		return logger.InitGelf()
	default:
		return logger, nil
	}
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Debugf(ctx, format, args...)
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Infof(ctx, format, args...)
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Warnf(ctx, format, args...)
}

func (l *Logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Warningf(ctx, format, args...)
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Errorf(ctx, format, args...)
}

func (l *Logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Fatalf(ctx, format, args...)
}

func (l *Logger) Panicf(ctx context.Context, format string, args ...interface{}) {
	l.CreateEntry().Panicf(ctx, format, args...)
}

func (l *Logger) Trace(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Trace(ctx, args...)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Debug(ctx, args...)
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Info(ctx, args...)
}

func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Warn(ctx, args...)
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Error(ctx, args...)
}

func (l *Logger) Fatal(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Fatal(ctx, args...)
}

func (l *Logger) Panic(ctx context.Context, args ...interface{}) {
	l.CreateEntry().Panic(ctx, args...)
}

func (l *Logger) WithField(name string, value interface{}) *LogEntry {
	return l.CreateEntry().WithField(name, value)
}

func (l *Logger) WithFields(fields map[string]interface{}) *LogEntry {
	return l.CreateEntry().WithFields(fields)
}

func (l *Logger) WithError(err error) *LogEntry {
	return l.CreateEntry().WithError(err)
}

func (l *Logger) Flush() {
	for i := range l.flusherHooks {
		l.flusherHooks[i].Flush()
	}
}
