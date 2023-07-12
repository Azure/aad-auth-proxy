package handler

import (
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/global"
)

// Creates proxy for incoming requests
func CreateReverseProxy(targetHost string, tokenProvider contracts.ITokenProvider) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Director = func(request *http.Request) {
		modifyRequest(request, targetHost, tokenProvider)
	}
	proxy.ErrorHandler = handleError
	proxy.ModifyResponse = modifyResponse

	return proxy, nil
}

// This modifies incoming requests and changes host to targetHost
func modifyRequest(request *http.Request, targetHost string, tokenProvider contracts.ITokenProvider) {
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(request.Context(), "modifyRequest")
	defer span.End()

	span.SetAttributes(
		attribute.String("request.target_scheme", constants.HTTPS_SCHEME),
		attribute.String("request.target_host", targetHost),
	)

	request.URL.Scheme = constants.HTTPS_SCHEME
	request.URL.Host = targetHost
	request.Host = targetHost

	// Record metrics
	// request_bytes_total{target_host, method, path, user_agent}
	metricAttributes := []attribute.KeyValue{
		attribute.String("target_host", request.URL.Host),
		attribute.String("method", request.Method),
		attribute.String("path", request.URL.Path),
		attribute.String("user_agent", request.Header.Get(constants.HEADER_USER_AGENT)),
	}

	meter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
	intrument, err := meter.Int64Counter(constants.METRIC_REQUEST_BYTES_TOTAL)
	if err == nil {
		intrument.Add(ctx, request.ContentLength, metricAttributes...)
	}
}

// This will be called when there is an error in forwarding the request
func handleError(response http.ResponseWriter, request *http.Request, err error) {
	// Record traces
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(request.Context(), "handleError")
	defer span.End()

	attributes := []attribute.KeyValue{
		attribute.String("response.status_code", response.Header().Get(constants.HEADER_STATUS_CODE)),
		attribute.String("response.content_type", response.Header().Get(constants.HEADER_CONTENT_TYPE)),
		attribute.String("response.content_encoding", response.Header().Get(constants.HEADER_CONTENT_ENCODING)),
		attribute.String("response.request_id", response.Header().Get(constants.HEADER_REQUEST_ID)),
		attribute.String("response.error.message", err.Error()),
	}

	span.SetAttributes(attributes...)
	span.RecordError(err)
	span.SetStatus(codes.Error, "failed to forward request")

	// Log error
	log.WithFields(log.Fields{
		"Request": request.URL.String(),
	}).Errorln("Request failed", err)

	// Record metrics
	// requests_total{target_host, method, path, user_agent, status_code}
	status_code, err := strconv.ParseInt(response.Header().Get(constants.HEADER_STATUS_CODE), 10, 32)
	if err != nil {
		log.Errorln("Failed to parse status code", err)
		status_code = 0
	}

	metricAttributes := []attribute.KeyValue{
		attribute.String("target_host", request.URL.Host),
		attribute.String("method", request.Method),
		attribute.String("path", request.URL.Path),
		attribute.String("user_agent", request.Header.Get(constants.HEADER_USER_AGENT)),
		attribute.Int("status_code", int(status_code)),
	}

	requestCountMeter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
	requestCountIntrument, err := requestCountMeter.Int64Counter(constants.METRIC_REQUESTS_TOTAL)
	if err == nil {
		requestCountIntrument.Add(ctx, 1, metricAttributes...)
	}
}

// This will be called once we receive response from targetHost
func modifyResponse(response *http.Response) (err error) {
	// Record traces
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(response.Request.Context(), "modifyResponse")
	defer span.End()

	traceAttributes := []attribute.KeyValue{
		attribute.Int("response.status_code", response.StatusCode),
		attribute.String("response.content_length", response.Header.Get(constants.HEADER_CONTENT_LENGTH)),
		attribute.String("response.content_type", response.Header.Get(constants.HEADER_CONTENT_TYPE)),
		attribute.String("response.content_encoding", response.Header.Get(constants.HEADER_CONTENT_ENCODING)),
		attribute.String("response.request_id", response.Header.Get(constants.HEADER_REQUEST_ID)),
	}

	span.SetAttributes(traceAttributes...)

	// Metric attributes
	metricAttributes := []attribute.KeyValue{
		attribute.String("target_host", response.Request.URL.Host),
		attribute.String("method", response.Request.Method),
		attribute.String("path", response.Request.URL.Path),
		attribute.String("user_agent", response.Request.Header.Get(constants.HEADER_USER_AGENT)),
		attribute.Int("status_code", response.StatusCode),
	}

	// Record metrics
	// requests_total{target_host, method, path, user_agent, status_code}
	requestCountMeter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
	requestCountIntrument, err := requestCountMeter.Int64Counter(constants.METRIC_REQUESTS_TOTAL)
	if err == nil {
		requestCountIntrument.Add(ctx, 1, metricAttributes...)
	}

	// Record metrics
	// response_bytes_total{target_host, method, path, user_agent, status_code}
	responseBytesMeter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
	responseBytesIntrument, err := responseBytesMeter.Int64Counter(constants.METRIC_RESPONSE_BYTES_TOTAL)
	if err == nil {
		responseBytesIntrument.Add(ctx, response.ContentLength, metricAttributes...)
	}

	// Log response
	log.WithFields(log.Fields{
		"Request":       response.Request.URL.String(),
		"StatusCode":    response.StatusCode,
		"ContentLength": response.ContentLength,
	}).Infoln("Successfully sent request, returning response back.")

	response.Header.Set("Status-Code", strconv.Itoa(response.StatusCode))

	// If server returned error, log response as well
	if response.StatusCode >= http.StatusBadRequest {
		err = errors.New("Non 2xx response from target host: " + strconv.Itoa(response.StatusCode))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Non 2xx response from target host")

		// Read 2KB of data
		limitedReader := &io.LimitedReader{R: response.Body, N: constants.BYTES_2KB}
		responseBody, err := ioutil.ReadAll(limitedReader)
		if err != nil {
			return err
		}

		log.Println("Error response body: ", string(responseBody[:]))
	}

	return nil
}
