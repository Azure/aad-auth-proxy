package main

import (
	"aad-auth-proxy/certificate"
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"aad-auth-proxy/handler"
	"aad-auth-proxy/telemetry"
	"aad-auth-proxy/token_provider"
	"aad-auth-proxy/utils"
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Main entry point.
func main() {
	// Handle panics
	defer utils.HandlePanic("main")

	// Telemetry
	logger := telemetry.NewLogger()

	// Read configuration parameters
	configuration := utils.NewConfiguration()
	listeningPort := configuration.GetListeningPort()
	audience := configuration.GetAudience()
	targetHost := configuration.GetTargetHost()

	// Traces
	tracerShutdown, err := telemetry.InitializeTracer(logger, configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tracerShutdown(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	// Metrics
	metricShutdown, err := telemetry.InitializeMetric(logger, configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := metricShutdown(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	// Create handler to fetch tokens
	handler := createHandlerWithTokenProvider(configuration, audience, targetHost, logger)

	// Note: These (Health and Readiness) check whether app can fetch a token with provided identity,
	// it does not evaluate end-to-end authentication and authorization of acquired tokens.
	// Health check handler
	http.HandleFunc("/health", handler.ReadinessCheck)

	// Readiness check handler
	http.HandleFunc("/ready", handler.ReadinessCheck)

	// Reverse proxy handler
	http.HandleFunc("/", handler.ProxyRequest)

	// Listen at specified port
	log.Fatal(http.ListenAndServe(":"+listeningPort, nil))
}

// Creates handler with token provider
func createHandlerWithTokenProvider(configuration utils.IConfiguration, audience string, targetHost string, logger contracts.ILogger) *handler.Handler {
	var certManager *certificate.CertificateManager
	var err error

	// Create cert manager for AAD based authentication
	if configuration.GetIdentityType() == constants.AAD_APPLICATION {
		certManager, err = certificate.NewCerificateManager(configuration.GetAadClientCertPath())
		if err != nil {
			logger.Error("NewCertificateManager failed: ", err)
		}
	}

	// Create TokenProvider
	tokenProvider, err := token_provider.NewTokenProvider(audience, configuration, certManager, logger)
	if err != nil {
		logger.Error("TokenCredential creation failed:", err)
	}

	proxy, err := handler.CreateReverseProxy(targetHost, tokenProvider)
	if err != nil {
		logger.Error("Proxy creation failed:", err)
	}

	// Create handler to return tokens based on audience
	handler, err := handler.NewHandler(proxy, tokenProvider, configuration)
	if err != nil {
		logger.Error("NewHandler failed:", err)
	}
	return handler
}
