package telemetry

import (
	"aad-auth-proxy/contracts"
	"aad-auth-proxy/utils"
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitializeTracer(logger contracts.ILogger, configuration utils.IConfiguration) (func(context.Context) error, error) {
	// Create a new otlptrace exporter
	exporter, err := otlptrace.New(context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(configuration.GetOtelEndpoint()),
		),
	)
	if err != nil {
		logger.Fatal(err)
	}

	// Create a trace provider for otel
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(NewResource(configuration.GetOtelServiceName())),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}

func InitializeMetric(logger contracts.ILogger, configuration utils.IConfiguration) (func(context.Context) error, error) {
	// Create a new otlpmetric exporter
	exporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(configuration.GetOtelEndpoint()),
	)

	if err != nil {
		logger.Fatal(err)
	}

	// Create a metric provider for otel
	metricProvider := metric.NewMeterProvider(
		metric.WithResource(NewResource(configuration.GetOtelServiceName())),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)

	global.SetMeterProvider(metricProvider)

	return metricProvider.Shutdown, nil
}
