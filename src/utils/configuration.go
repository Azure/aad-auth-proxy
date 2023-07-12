package utils

import (
	"aad-auth-proxy/constants"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type IConfiguration interface {
	GetAadClientId() string
	GetAadClientCertPath() string
	GetAadTenantId() string
	GetAadTokenRefreshDurationInPercentage() uint8
	GetIdentityType() string
	GetListeningPort() string
	GetAudience() string
	GetTargetHost() string
	GetOtelEndpoint() string
	GetOtelServiceName() string
	GetAdditionalHeaders() map[string]string
}

type configuration struct {
	listeningPort     string
	identityConfig    *identityConfiguration
	hostConfig        *hostConfiguration
	telemetryConfig   *telemetryConfiguration
	additionalHeaders *additionalHeaders
}

type identityConfiguration struct {
	identityType                        string
	aadClientId                         string
	aadClientCertificatePath            string
	aadTenantId                         string
	aadTokenRefreshDurationInPercentage uint8
}

type hostConfiguration struct {
	audience   string
	targetHost string
}

type telemetryConfiguration struct {
	otelEndpoint    string
	otelServiceName string
}

type additionalHeaders struct {
	headers    map[string]string
	headersStr string
}

func NewConfiguration() *configuration {

	config := readConfigurationsFromEnv()

	fields := log.Fields{
		"AAD_CLIENT_ID":                            config.identityConfig.aadClientId,
		"AAD_TENANT_ID":                            config.identityConfig.aadTenantId,
		"AAD_CLIENT_CERTIFICATE_PATH":              config.identityConfig.aadClientCertificatePath,
		"AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE": config.identityConfig.aadTokenRefreshDurationInPercentage,
		"IDENTITY_TYPE":                            config.identityConfig.identityType,
		"LISTENING_PORT":                           config.listeningPort,
		"AUDIENCE":                                 config.hostConfig.audience,
		"TARGET_HOST":                              config.hostConfig.targetHost,
		"OTEL_SERVICE_NAME":                        config.telemetryConfig.otelServiceName,
		"OTEL_GRPC_ENDPOINT":                       config.telemetryConfig.otelEndpoint,
		"OVERRIDE_REQUEST_HEADERS":                 config.additionalHeaders.headersStr,
	}

	if config.listeningPort == "" {
		log.WithFields(fields).Fatalln("Missing required configuration setting: LISTENING_PORT")
	}

	log.WithFields(fields).Infoln("Configuration settings loaded:")

	return config
}

// Reads configurations from environment variables.
func readConfigurationsFromEnv() *configuration {

	aadClientId := os.Getenv("AAD_CLIENT_ID")
	aadClientCertificatePath := os.Getenv("AAD_CLIENT_CERTIFICATE_PATH")
	aadTenantId := os.Getenv("AAD_TENANT_ID")
	aadTokenRefreshDurationInPercentageStr := os.Getenv("AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE")
	identityType := os.Getenv("IDENTITY_TYPE")
	listeningPort := os.Getenv("LISTENING_PORT")
	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")
	otelEndpoint := os.Getenv("OTEL_GRPC_ENDPOINT")
	headersStr := os.Getenv("OVERRIDE_REQUEST_HEADERS")

	// Identity
	identityType = strings.ToLower(identityType)
	var aadTokenRefreshDurationInPercentage uint8 = constants.DEFAULT_TOKEN_REFRESH_PERCENTAGE

	// Parse refresh interval if passed as parameter
	if aadTokenRefreshDurationInPercentageStr != "" {
		value, err := strconv.ParseUint(aadTokenRefreshDurationInPercentageStr, 10, 64)
		if err != nil {
			log.WithField(
				"AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE",
				aadTokenRefreshDurationInPercentageStr).Warningln("failed to parse AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE, using default value of 1% of time before expiry")
		} else if value <= 0 || value >= 100 {
			log.WithField(
				"AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE",
				value).Warningln("AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE should be non-zero number from 1 to 99, using default value of 1% of time before expiry")
		} else {
			aadTokenRefreshDurationInPercentage = uint8(value)
		}
	}

	identityConfig := &identityConfiguration{
		identityType:                        identityType,
		aadClientId:                         aadClientId,
		aadClientCertificatePath:            aadClientCertificatePath,
		aadTenantId:                         aadTenantId,
		aadTokenRefreshDurationInPercentage: aadTokenRefreshDurationInPercentage,
	}

	// Auth
	audience := os.Getenv("AUDIENCE")
	targetHost := os.Getenv("TARGET_HOST")
	// Trim https:// and http:// prefixes, as the scheme will be added to the modified host url
	targetHost = strings.TrimPrefix(strings.TrimPrefix(strings.ToLower(targetHost), constants.HTTPS_PREFIX), constants.HTTP_PREFIX)

	hostConfig := &hostConfiguration{
		audience:   audience,
		targetHost: targetHost,
	}

	// Telemetry
	if otelServiceName == "" {
		otelServiceName = constants.SERVICE_TELEMETRY_KEY
	}

	telemetryConfig := &telemetryConfiguration{
		otelEndpoint:    otelEndpoint,
		otelServiceName: otelServiceName,
	}

	// Headers
	if headersStr == "" {
		headersStr = "{}"
	}

	additionalHeaders := parseHeaders(headersStr)

	for key, value := range additionalHeaders.headers {
		log.WithFields(log.Fields{
			"HEADER_KEY":   key,
			"HEADER_VALUE": value,
		}).Infoln("Additional headers loaded")
	}

	return &configuration{
		listeningPort:     listeningPort,
		identityConfig:    identityConfig,
		hostConfig:        hostConfig,
		telemetryConfig:   telemetryConfig,
		additionalHeaders: additionalHeaders,
	}
}

func parseHeaders(headersStr string) *additionalHeaders {
	headers := make(map[string]string)
	err := json.Unmarshal([]byte(headersStr), &headers)
	if err != nil {
		log.WithField("OVERRIDE_REQUEST_HEADERS", headersStr).Warningln("failed to parse OVERRIDE_REQUEST_HEADERS, using default value of empty map")
		return &additionalHeaders{
			headers:    make(map[string]string),
			headersStr: "{}",
		}
	}

	parsedHeadersStr, err := json.Marshal(headers)
	if err != nil {
		log.WithField("headers", headers).Warningln("failed to marshall")
	}

	return &additionalHeaders{
		headers:    headers,
		headersStr: string(parsedHeadersStr),
	}
}

func (config *configuration) GetAadClientId() string {
	return config.identityConfig.aadClientId
}

func (config *configuration) GetAadClientCertPath() string {
	return config.identityConfig.aadClientCertificatePath
}

func (config *configuration) GetAadTenantId() string {
	return config.identityConfig.aadTenantId
}

func (config *configuration) GetAadTokenRefreshDurationInPercentage() uint8 {
	return config.identityConfig.aadTokenRefreshDurationInPercentage
}

func (config *configuration) GetIdentityType() string {
	return config.identityConfig.identityType
}

func (config *configuration) GetListeningPort() string {
	return config.listeningPort
}

func (config *configuration) GetAudience() string {
	return config.hostConfig.audience
}

func (config *configuration) GetTargetHost() string {
	return config.hostConfig.targetHost
}

func (config *configuration) GetOtelEndpoint() string {
	return config.telemetryConfig.otelEndpoint
}

func (config *configuration) GetOtelServiceName() string {
	return config.telemetryConfig.otelServiceName
}

func (config *configuration) GetAdditionalHeaders() map[string]string {
	return config.additionalHeaders.headers
}
