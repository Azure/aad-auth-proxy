package constants

const (
	// Http
	HTTPS_SCHEME = "https"

	// Request Headers
	HEADER_AUTHORIZATION = "Authorization"

	// Identity
	SYSTEM_ASSIGNED = "systemAssigned"
	USER_ASSIGNED   = "userAssigned"
	AAD_APPLICATION = "aadApplication"

	// Host
	AUDIENCE    = "https://monitor.azure.com/"
	TARGET_HOST = "https://monitor.azure.com/"

	// Error Response Headers
	ERROR_PROPERTY_CODE    = "code"
	ERROR_PROPERTY_MESSAGE = "message"

	// Limited reader
	BYTES_2KB = 2048 // 2K
)
