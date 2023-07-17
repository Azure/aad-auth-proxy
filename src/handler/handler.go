package handler

import (
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"aad-auth-proxy/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/global"
)

// This manages token provider handler
type Handler struct {
	targetHost      string
	proxy           *httputil.ReverseProxy
	tokenProvider   contracts.ITokenProvider
	configuration   utils.IConfiguration
	overrideHeaders map[string]string
}

// Creates a new handler
func NewHandler(proxy *httputil.ReverseProxy, tokenProvider contracts.ITokenProvider, configuration utils.IConfiguration) (handler *Handler, err error) {
	if proxy == nil {
		return nil, errors.New("proxy cannot be nil")
	}

	if tokenProvider == nil {
		return nil, errors.New("tokenProvider cannot be nil")
	}

	if configuration == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	var overrideHeaders map[string]string = nil
	additionalheaders := configuration.GetAdditionalHeaders()
	if additionalheaders != nil && len(additionalheaders) > 0 {
		overrideHeaders = additionalheaders
	}

	return &Handler{
		targetHost:      configuration.GetTargetHost(),
		proxy:           proxy,
		tokenProvider:   tokenProvider,
		configuration:   configuration,
		overrideHeaders: overrideHeaders,
	}, nil
}

// Reverse proxy handler
func (handler *Handler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Start tracing
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(r.Context(), "ProxyRequest")
	defer span.End()

	attributes := []attribute.KeyValue{
		attribute.String("request.query_string", r.URL.RawQuery),
		attribute.String("request.path", r.URL.Path),
		attribute.String("request.method", r.Method),
		attribute.Int64("request.content_length", r.ContentLength),
		attribute.String("request.content_type", r.Header.Get("Content-Type")),
		attribute.String("request.user_agent", r.Header.Get("user-Agent")),
		attribute.String("request.content_encoding", r.Header.Get("Content-Encoding")),
	}

	span.SetAttributes(attributes...)

	err := handler.checkTokenProvider(ctx)
	if err != nil {
		// Metric attributes
		metricAttributes := []attribute.KeyValue{
			attribute.String("target_host", r.URL.Host),
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("user_agent", r.Header.Get("User-Agent")),
			attribute.Int("status_code", http.StatusServiceUnavailable),
		}

		// Record metrics
		// requests_total{target_host, method, path, user_agent, status_code}
		requestCountMeter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
		requestCountIntrument, err := requestCountMeter.Int64Counter(constants.METRIC_REQUESTS_TOTAL)
		if err == nil {
			requestCountIntrument.Add(ctx, 1, metricAttributes...)
		}

		handler.failRequest(w, r, ctx, err)
		return
	}

	// Add authorization header
	token, _ := handler.tokenProvider.GetAccessToken()
	r.Header.Set(constants.HEADER_AUTHORIZATION, "Bearer "+token)

	// Add additional headers
	if handler.overrideHeaders != nil {
		for key, value := range handler.overrideHeaders {
			r.Header.Set(key, value)
		}
	}

	// Start timer for calculating request duration
	startTime := time.Now()
	defer func() {
		// Extract duration and status_code
		duration := time.Since(startTime).Milliseconds()
		status_code, err := strconv.ParseInt(w.Header().Get("Status-Code"), 10, 32)
		if err != nil {
			log.Errorln("Failed to parse status code", err)
			status_code = 0
		}
		// Record metrics
		// request_duration_milliseconds{target_host, method, path, user_agent, status_code}
		requestDurationMeter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
		requestDurationIntrument, err := requestDurationMeter.Int64Histogram(constants.METRIC_REQUEST_DURATION_MILLISECONDS)
		if err == nil {
			metricAttributes := []attribute.KeyValue{
				attribute.String("target_host", handler.targetHost),
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.String("user_agent", r.Header.Get("User-Agent")),
				attribute.Int("status_code", int(status_code)),
			}
			requestDurationIntrument.Record(ctx, duration, metricAttributes...)
		}
	}()

	// Handle request
	handler.proxy.ServeHTTP(w, r.WithContext(ctx))
}

// Readiness check handler
func (handler *Handler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(r.Context(), "ReadinessCheck")
	defer span.End()

	// Check token provider
	err := handler.checkTokenProvider(ctx)
	if err != nil {
		handler.failRequest(w, r, ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("response.status_code", http.StatusOK))
	w.WriteHeader(http.StatusOK)
}

// TokenProvider check
// If token provider is not instantiated, cannot fetch token, so fail request
func (handler *Handler) checkTokenProvider(ctx context.Context) error {
	if handler.tokenProvider == nil {
		token, err := handler.tokenProvider.GetAccessToken()
		if len(token) == 0 || err != nil {
			// Start tracing
			_, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(ctx, "checkTokenProvider")
			defer span.End()

			// If we run into a case where we received empty token without any errors
			if err == nil {
				err = errors.New("handler, tokenProvider is not instantiated, cannot forward request")
			}

			span.SetAttributes(attribute.Int("proxy.status_code", http.StatusServiceUnavailable))
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to forward request")

			log.Errorln("failed to forward request", err)

			return err
		}
	}

	return nil
}

// Fail request
// This will fail the request with 503
func (handler *Handler) failRequest(w http.ResponseWriter, r *http.Request, ctx context.Context, err error) {
	_, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(ctx, "failRequest")
	defer span.End()

	statusCode := http.StatusServiceUnavailable
	errorCode := "AuthenticationTokenNotFound"
	errorMessage := err.Error()
	requestId := r.Header.Get(constants.HEADER_REQUEST_ID)

	span.SetAttributes(
		attribute.Int("response.status_code", statusCode),
		attribute.String("response.request_id", requestId),
		attribute.String("response.error.code", errorCode),
		attribute.String("response.error.message", errorMessage),
	)
	span.SetStatus(codes.Error, "failed to forward request")

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	resp := make(map[string]string)
	resp[constants.ERROR_PROPERTY_CODE] = errorCode
	resp[constants.ERROR_PROPERTY_MESSAGE] = errorMessage
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(jsonResp)
}
