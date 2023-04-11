package utils

import (
	"aad-auth-proxy/constants"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type IConfiguration interface {
	GetAadClientId() string
	GetAadClientCertPath() string
	GetAadTenantId() string
	GetAadTokenRefreshDurationInMinutes() time.Duration
	GetIdentityType() string
	GetListeningPort() string
	GetAudience() string
	GetTargetHost() string
	GetOtelEndpoint() string
	GetOtelServiceName() string
}

type configuration struct {
	listeningPort   string
	identityConfig  *identityConfiguration
	hostConfig      *hostConfiguration
	telemetryConfig *telemetryConfiguration
}

type identityConfiguration struct {
	identityType                     string
	aadClientId                      string
	aadClientCertificatePath         string
	aadTenantId                      string
	aadTokenRefreshDurationInMinutes time.Duration
}

type hostConfiguration struct {
	audience   string
	targetHost string
}

type telemetryConfiguration struct {
	otelEndpoint    string
	otelServiceName string
}

func NewConfiguration() *configuration {

	config := readConfigurationsFromEnv()

	fields := log.Fields{
		"AAD_CLIENT_ID":                         config.identityConfig.aadClientId,
		"AAD_TENANT_ID":                         config.identityConfig.aadTenantId,
		"AAD_CLIENT_CERTIFICATE_PATH":           config.identityConfig.aadClientCertificatePath,
		"AAD_TOKEN_REFRESH_DURATION_IN_MINUTES": config.identityConfig.aadTokenRefreshDurationInMinutes,
		"IDENTITY_TYPE":                         config.identityConfig.identityType,
		"LISTENING_PORT":                        config.listeningPort,
		"AUDIENCE":                              config.hostConfig.audience,
		"TARGET_HOST":                           config.hostConfig.targetHost,
		"OTEL_SERVICE_NAME":                     config.telemetryConfig.otelServiceName,
		"OTEL_GRPC_ENDPOINT":                    config.telemetryConfig.otelEndpoint,
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
	aadTokenRefreshDurationInMinutesStr := os.Getenv("AAD_TOKEN_REFRESH_INTERVAL_IN_MINUTES")
	identityType := os.Getenv("IDENTITY_TYPE")
	listeningPort := os.Getenv("LISTENING_PORT")
	otelServiceName := os.Getenv("OTEL_SERVICE_NAME")
	otelEndpoint := os.Getenv("OTEL_GRPC_ENDPOINT")

	// Identity
	identityType = strings.ToLower(identityType)
	aadTokenRefreshDurationInMinutes := constants.TIME_60_MINITES

	// Parse refresh interval if passed as parameter
	if aadTokenRefreshDurationInMinutesStr != "" {
		value, err := strconv.ParseInt(aadTokenRefreshDurationInMinutesStr, 10, 32)
		if err != nil {
			log.WithField(
				"AAD_TOKEN_REFRESH_INTERVAL_IN_MINUTES",
				aadTokenRefreshDurationInMinutesStr).Warningln("failed to parse AAD_TOKEN_REFRESH_INTERVAL_IN_MINUTES, using default value of 60 minutes")
		} else {
			aadTokenRefreshDurationInMinutes = time.Minute * time.Duration(value)
		}
	}

	identityConfig := &identityConfiguration{
		identityType:                     identityType,
		aadClientId:                      aadClientId,
		aadClientCertificatePath:         aadClientCertificatePath,
		aadTenantId:                      aadTenantId,
		aadTokenRefreshDurationInMinutes: aadTokenRefreshDurationInMinutes,
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

	return &configuration{
		listeningPort:   listeningPort,
		identityConfig:  identityConfig,
		hostConfig:      hostConfig,
		telemetryConfig: telemetryConfig,
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

func (config *configuration) GetAadTokenRefreshDurationInMinutes() time.Duration {
	return config.identityConfig.aadTokenRefreshDurationInMinutes
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
