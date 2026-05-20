package tracing

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// InitializeTracer initializes and returns a TracerProvider with Zipkin exporter
func InitializeTracer(serviceName string) (*tracesdk.TracerProvider, error) {
	// Get Zipkin endpoint from environment or use default
	zipkinEndpoint := os.Getenv("OTEL_EXPORTER_ZIPKIN_ENDPOINT")
	if zipkinEndpoint == "" {
		zipkinEndpoint = "http://zipkin:9411/api/v2/spans"
	}

	// Create Zipkin exporter
	exporter, err := zipkin.New(zipkinEndpoint)
	if err != nil {
		return nil, err
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create TracerProvider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	log.Printf("Tracing initialized for service: %s", serviceName)
	return tp, nil
}

// Shutdown gracefully shuts down the tracer provider
func Shutdown(ctx context.Context, tp *tracesdk.TracerProvider) error {
	return tp.Shutdown(ctx)
}
