package tracer

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"sp-office-lookuper/internal/app"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

const (
	RequestIDHeader = "X-Request-Id"
	RequestIDKey    = "rid"
)

type Config interface {
	GetServiceName() string
	GetHost() string
	GetPort() int
	GetSamplerType() string
	GetSamplerParam() float64
	GetBufferFlushInterval() time.Duration
}

// InitTracer connects to jaeger tracing
func InitTracer(conf Config) error {
	cfg := jaegercfg.Configuration{
		ServiceName: conf.GetServiceName(),
		Sampler: &jaegercfg.SamplerConfig{
			Type:  conf.GetSamplerType(),
			Param: conf.GetSamplerParam(),
		},
		Reporter: &jaegercfg.ReporterConfig{
			BufferFlushInterval: conf.GetBufferFlushInterval(),
			LocalAgentHostPort:  fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort()),
		},
	}
	tracer, _, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger), jaegercfg.Gen128Bit(true))
	if err != nil {
		return err
	}
	opentracing.SetGlobalTracer(tracer)
	return nil
}

type customSpan struct {
	opentracing.Span
}

func (cs customSpan) SetTag(key string, value interface{}) opentracing.Span {
	if key == app.ErrorTag {
		if _, ok := value.(bool); !ok {
			cs.Span.SetTag(app.ErrorTag, true)
			cs.Span.SetTag(app.ErrorMessageTag, value)
			return &cs
		}
	}

	cs.Span.SetTag(key, value)
	return &cs
}

func getCustomSpanFinisher(span opentracing.Span, finish func()) func(...error) {
	return func(params ...error) {
		var err error
		if len(params) > 0 {
			err = params[0]
		}
		if err != nil {
			span.SetTag(app.ErrorTag, true)
			span.SetTag(app.ErrorMessageTag, err)
		}
		finish()
	}
}

type spanWithContextOptions struct {
	operationName     string
	extractFromHeader *http.Header
}

type spanWithContextOptionModifier func(opts *spanWithContextOptions)

func (s spanWithContextOptionModifier) ApplyTo(opts *spanWithContextOptions) {
	s(opts)
}

// ExtractFrom option extracts trace info from header
func ExtractFrom(header http.Header) func(opts *spanWithContextOptions) {
	return func(opts *spanWithContextOptions) {
		opts.extractFromHeader = &header
	}
}

// Trace creates and return new span from passed context
// Common usage:
// var err error
// ctx, span, finish := tracer.Trace(ctx)
// span.SetTag(tag.TagName, value)
// defer func(){ finish(err) }()
func Trace(
	ctx context.Context,
	opts ...spanWithContextOptionModifier,
) (context.Context, opentracing.Span, func(...error)) {
	callOptions := spanWithContextOptions{}
	for _, opt := range opts {
		opt.ApplyTo(&callOptions)
	}

	operation := callOptions.operationName
	if operation == "" {
		pc, file, _, ok := runtime.Caller(1)
		operation = file
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			operation = details.Name()
		}
	}

	var span opentracing.Span
	if callOptions.extractFromHeader != nil {
		parentSpanCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(*callOptions.extractFromHeader),
		)
		if parentSpanCtx != nil && err == nil {
			span = opentracing.GlobalTracer().StartSpan(
				operation, opentracing.ChildOf(parentSpanCtx),
			)
			xRequestID := callOptions.extractFromHeader.Get(RequestIDHeader)
			if xRequestID != "" {
				span.LogFields(log.String(RequestIDHeader, xRequestID))
			}
		}
	}

	if span == nil {
		parentSpan := opentracing.SpanFromContext(ctx)
		options := make([]opentracing.StartSpanOption, 0)
		if parentSpan != nil {
			options = append(options, opentracing.ChildOf(parentSpan.Context()))
		}
		span = opentracing.StartSpan(operation, options...)
	}

	span.LogFields(log.String(RequestIDHeader, GetOrCreateRequestID(ctx)))
	nextCtx := opentracing.ContextWithSpan(ctx, span)

	return nextCtx, customSpan{span}, getCustomSpanFinisher(span, span.Finish)
}
