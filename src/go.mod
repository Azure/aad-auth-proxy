module aad-auth-proxy

go 1.19

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.1.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.1.0
	github.com/sirupsen/logrus v1.8.1
	go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.36.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.13.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.13.0
	go.opentelemetry.io/otel/metric v0.36.0
	go.opentelemetry.io/otel/sdk v1.13.0
	go.opentelemetry.io/otel/sdk/metric v0.36.0
	software.sslmate.com/src/go-pkcs12 v0.0.0-20210415151418-c5206de65a78
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.11.3 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.13.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.56.3 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.0.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.5.1 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	github.com/t-tomalak/logrus-easy-formatter v0.0.0-20190827215021-c074f06c5816
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace (
	golang.org/x/net => golang.org/x/net v0.17.0
	golang.org/x/text => golang.org/x/text v0.3.8
	google.golang.org/grpc => google.golang.org/grpc v1.56.3
)
