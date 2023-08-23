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
	HEADER_AUTHORIZATION    = "Authorization"
	HEADER_REQUEST_ID       = "X-Request-ID"
	HEADER_CONTENT_TYPE     = "Content-Type"
	HEADER_CONTENT_LENGTH   = "Content-Length"
	HEADER_CONTENT_ENCODING = "Content-Encoding"
	HEADER_USER_AGENT       = "User-Agent"
	HEADER_STATUS_CODE      = "Status-Code"

	// Encoding
	ENCODING_GZIP         = "gzip"
	ENCODING_DEFLATE_ZLIB = "zlib"

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
	TIME_5_SECONDS = time.Second * 5
	TIME_1_MINUTES = time.Minute * 1
	TIME_5_MINUTES = time.Minute * 5

	// Default token refresh percentage
	DEFAULT_TOKEN_REFRESH_PERCENTAGE = 10 // 10% before expiry
)
