package constants

import "time"

const (
	// Version
	VERSION = "v0.1.0"

	// Telemetry
	SERVICE_TELEMETRY_KEY = "aad_auth_proxy"

	// Telemetry metric names
	METRIC_REQUESTS_TOTAL                = "aad_auth_proxy_requests_total"
	METRIC_REQUEST_DURATION_MILLISECONDS = "aad_auth_proxy_request_duration_milliseconds"
	METRIC_REQUEST_BYTES_TOTAL           = "aad_auth_proxy_request_bytes_total"
	METRIC_RESPONSE_BYTES_TOTAL          = "aad_auth_proxy_response_bytes_total"
	METRIC_TOKEN_REFRESH_TOTAL           = "aad_auth_proxy_token_refresh_total"

	// Http
	HTTPS_SCHEME = "https"
	HTTP_PREFIX  = "http://"
	HTTPS_PREFIX = "https://"

	// Request Headers
	HEADER_AUTHORIZATION = "Authorization"

	// Identity
	SYSTEM_ASSIGNED = "systemassigned"
	USER_ASSIGNED   = "userassigned"
	AAD_APPLICATION = "aadapplication"

	// Error Response Headers
	ERROR_PROPERTY_CODE    = "code"
	ERROR_PROPERTY_MESSAGE = "message"

	// Limited reader
	BYTES_2KB = 2048 // 2K

	// Time
	TIME_1_MINUTES  = time.Minute * 1
	TIME_5_MINUTES  = time.Minute * 5
	TIME_60_MINITES = time.Minute * 60
)
