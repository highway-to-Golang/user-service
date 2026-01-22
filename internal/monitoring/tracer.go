package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
	tp     *tracesdk.TracerProvider
)

// InitTracer инициализирует OpenTelemetry tracer
func InitTracer(serviceName, jaegerEndpoint string) error {
	if jaegerEndpoint == "" {
		// Если endpoint не указан, используем noop tracer
		tracer = otel.Tracer(serviceName)
		return nil
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	tp = tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = otel.Tracer(serviceName)
	return nil
}

// GetTracer возвращает tracer
func GetTracer() trace.Tracer {
	if tracer == nil {
		tracer = otel.Tracer("user-service")
	}
	return tracer
}

// Shutdown останавливает tracer provider
func ShutdownTracer(ctx context.Context) error {
	if tp != nil {
		return tp.Shutdown(ctx)
	}
	return nil
}
