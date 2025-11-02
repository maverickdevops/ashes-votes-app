package main

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitOTel sets up OpenTelemetry tracing for the backend service
func InitOTel(ctx context.Context) (func(context.Context) error, error) {
	// Get endpoint from env or fallback (Docker Swarm internal service name)
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "otel-collector:4317"
	}

	// Connect to OTLP Collector using insecure gRPC (for local use)
	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	// Create trace exporter
	traceExp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Configure TracerProvider with batching and service name resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ashes-vote-backend"),
		)),
	)

	// Register global tracer
	otel.SetTracerProvider(tp)

	// Cleanup to ensure graceful shutdown
	cleanup := func(ctx context.Context) error {
		err := tp.Shutdown(ctx)
		_ = conn.Close()
		return err
	}

	return cleanup, nil
}
