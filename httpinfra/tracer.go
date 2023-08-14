package httpinfra

import (
	"context"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/contrib/propagators/ot"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/gophermodz/http/httpinfra"
)

func (s *Server) initTracer(ctx context.Context) {
	if !s.config.TracerEnabled {
		nop := trace.NewNoopTracerProvider()
		s.tracer = nop.Tracer(s.config.TracerServiceName)
		return
	}
	client := otlptracegrpc.NewClient()
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		s.logger.Error("creating OTLP trace exporter", err)
	}

	s.tracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(s.config.TracerServiceName),
			semconv.ServiceVersionKey.String(s.config.Version),
		)),
	)

	otel.SetTracerProvider(s.tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(),
		&jaeger.Jaeger{},
		&ot.OT{},
		&xray.Propagator{},
	))

	s.tracer = s.tracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(s.config.Version),
		trace.WithSchemaURL(semconv.SchemaURL),
	)
}

func NewOpenTelemetryMiddleware(service string) mux.MiddlewareFunc {
	return otelmux.Middleware(service)
}
