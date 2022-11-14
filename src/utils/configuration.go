package utils

import (
	"aad-auth-proxy/constants"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type IConfiguration interface {
	GetAadClientId() string
	GetAadClientCertPath() string
	GetAadTenantId() string
	GetIdentityType() string
	GetListeningPort() string
	GetAudience() string
	GetTargetHost() string
}

type configuration struct {
	listeningPort  string
	identityConfig *identityConfiguration
	hostConfig     *hostConfiguration
}

type identityConfiguration struct {
	identityType             string
	aadClientId              string
	aadClientCertificatePath string
	aadTenantId              string
}

type hostConfiguration struct {
	audience   string
	targetHost string
}

func NewConfiguration() *configuration {

	config := readConfigurationsFromEnv()

	fields := log.Fields{
		"AAD_CLIENT_ID":               config.identityConfig.aadClientId,
		"AAD_TENANT_ID":               config.identityConfig.aadTenantId,
		"AAD_CLIENT_CERTIFICATE_PATH": config.identityConfig.aadClientCertificatePath,
		"IDENTITY_TYPE":               config.identityConfig.identityType,
		"LISTENING_PORT":              config.listeningPort,
		"AUDIENCE":                    config.hostConfig.audience,
		"TARGET_HOST":                 config.hostConfig.targetHost,
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
	identityType := os.Getenv("IDENTITY_TYPE")
	listeningPort := os.Getenv("LISTENING_PORT")

	identityType = strings.ToLower(identityType)

	identityConfig := &identityConfiguration{
		identityType:             identityType,
		aadClientId:              aadClientId,
		aadClientCertificatePath: aadClientCertificatePath,
		aadTenantId:              aadTenantId,
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

	return &configuration{
		listeningPort:  listeningPort,
		identityConfig: identityConfig,
		hostConfig:     hostConfig,
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
