package logging

import (
	"context"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
)

const (
	// context keys
	CtxRespStatusKey = "ResponseStatusCode"

	// greylog values keys
	LogResponseCodeKey = "response_code"
	LogRidKey          = "request_id"
	LogTraceIDKey      = "trace_id"
	LogServiceNameKey  = "service_name"
)

type LogEntry struct {
	logger       *logrus.Logger
	serviceName  string
	fields       map[string]interface{}
	flusherHooks []flusher
}

func newLoggerEntry(logger *logrus.Logger, serviceName string, flushers []flusher) *LogEntry {
	lg := &LogEntry{
		logger:       logger,
		serviceName:  serviceName,
		fields:       make(map[string]interface{}),
		flusherHooks: flushers,
	}
	return lg
}

func (l *Logger) CreateEntry() *LogEntry {
	return newLoggerEntry(l.logger, l.serviceName, l.flusherHooks)
}

// logging methods
func (e *LogEntry) Debugf(ctx context.Context, format string, args ...interface{}) {
	e.logMsgf(ctx, logrus.DebugLevel, format, args...)
}

func (e *LogEntry) Infof(ctx context.Context, format string, args ...interface{}) {
	e.logMsgf(ctx, logrus.InfoLevel, format, args...)
}

func (e *LogEntry) Warnf(ctx context.Context, format string, args ...interface{}) {
	e.logMsgf(ctx, logrus.WarnLevel, format, args...)
}

func (e *LogEntry) Warningf(ctx context.Context, format string, args ...interface{}) {
	e.Warnf(ctx, format, args...)
}

func (e *LogEntry) Errorf(ctx context.Context, format string, args ...interface{}) {
	e.logMsgf(ctx, logrus.ErrorLevel, format, args...)
}

func (e *LogEntry) Fatalf(ctx context.Context, format string, args ...interface{}) {
	e.logMsgf(ctx, logrus.FatalLevel, format, args...)
	e.Flush()
	e.logger.Exit(1)
}

func (e *LogEntry) Panicf(ctx context.Context, format string, args ...interface{}) {
	defer e.Flush()
	e.logMsgf(ctx, logrus.PanicLevel, format, args...)
}

func (e *LogEntry) Trace(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.TraceLevel, args...)
}

func (e *LogEntry) Debug(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.DebugLevel, args...)
}

func (e *LogEntry) Info(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.InfoLevel, args...)
}

func (e *LogEntry) Warn(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.WarnLevel, args...)
}

func (e *LogEntry) Error(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.ErrorLevel, args...)
}

func (e *LogEntry) Fatal(ctx context.Context, args ...interface{}) {
	e.logMsg(ctx, logrus.FatalLevel, args...)
	e.Flush()
	e.logger.Exit(1)
}

func (e *LogEntry) Panic(ctx context.Context, args ...interface{}) {
	defer e.Flush()
	e.logMsg(ctx, logrus.PanicLevel, args...)
}

// WithError add error to current entry
func (e *LogEntry) WithError(err error) *LogEntry {
	return e.WithField(logrus.ErrorKey, err)
}

func (e *LogEntry) WithField(name string, value interface{}) *LogEntry {
	cpLogger := newLoggerEntry(e.logger, e.serviceName, e.flusherHooks)
	for key, val := range e.fields {
		cpLogger.fields[key] = val
	}
	cpLogger.fields[name] = value

	return cpLogger
}

func (e *LogEntry) WithFields(fields map[string]interface{}) *LogEntry {
	if fields == nil {
		return e
	}
	cpLogger := newLoggerEntry(e.logger, e.serviceName, e.flusherHooks)
	for name, val := range e.fields {
		cpLogger.fields[name] = val
	}
	for name, val := range fields {
		cpLogger.fields[name] = val
	}

	return cpLogger
}

func (e *LogEntry) logMsg(ctx context.Context, level logrus.Level, args ...interface{}) {
	if !e.logger.IsLevelEnabled(level) {
		return
	}
	entry := e.logger.WithFields(e.getFields(ctx))
	if e.fields != nil {
		entry = entry.WithFields(e.fields)
		defer func() {
			e.fields = make(map[string]interface{})
		}()
	}
	msg, fields, parsed := e.parseArgs(args)
	if parsed {
		entry.WithFields(fields).Log(level, msg)
		return
	}
	entry.Log(level, args...)
}

func (e *LogEntry) logMsgf(ctx context.Context, level logrus.Level, format string, args ...interface{}) {
	if !e.logger.IsLevelEnabled(level) {
		return
	}
	entry := e.logger.WithFields(e.getFields(ctx))
	if e.fields != nil {
		entry = entry.WithFields(e.fields)
		defer func() {
			e.fields = make(map[string]interface{})
		}()
	}
	entry.WithFields(e.getFields(ctx)).Logf(level, format, args...)
}

func (e *LogEntry) parseArgs(args ...interface{}) (message string, fields logrus.Fields, parsed bool) {
	if len(args) < 2 {
		return "", nil, false
	}
	if msg, ok := args[0].(string); ok {
		message = msg
	} else {
		return "", nil, false
	}
	fields = logrus.Fields{}
	for i, f := range args[1:] {
		id := strconv.Itoa(i)
		fields["field_"+id] = f
	}
	return message, fields, true
}

// parse ctx and get necessary fields from it and add to logging entry message
func (e *LogEntry) getFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{LogServiceNameKey: e.serviceName}

	// parse additional fileds from context
	if ctx != nil {
		// check if got this key in context or not
		rid := ctx.Value(RequestIDKey)
		respStatusCode := ctx.Value(CtxRespStatusKey)

		// check if values is not empty and ann them into fields to logging
		if rid != nil {
			fields[LogRidKey] = rid
		}
		if respStatusCode != nil {
			fields[LogResponseCodeKey] = respStatusCode
		}

		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			spanCtx, ok := span.Context().(jaeger.SpanContext)
			if ok {
				fields[LogTraceIDKey] = spanCtx.TraceID().String()
			}
		}
	}
	return fields
}

func (e *LogEntry) Flush() {
	for i := range e.flusherHooks {
		e.flusherHooks[i].Flush()
	}
}
