package telemetry

import (
	"aad-auth-proxy/constants"

	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// Creates a new resource for otel traces
func NewResource(serviceName string) *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(constants.VERSION),
		),
	)
	return r
}
