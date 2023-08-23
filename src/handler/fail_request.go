package handler

import (
	"aad-auth-proxy/constants"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Fail request
func FailRequest(w http.ResponseWriter, r *http.Request, statusCode int, errorCode string, ctx context.Context, err error) {
	_, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(ctx, "failRequest")
	defer span.End()

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
	w.Header().Add(constants.HEADER_CONTENT_TYPE, "application/json")
	w.Header().Add(constants.HEADER_STATUS_CODE, strconv.Itoa(statusCode))
	resp := make(map[string]string)
	resp[constants.ERROR_PROPERTY_CODE] = errorCode
	resp[constants.ERROR_PROPERTY_MESSAGE] = errorMessage
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return
	}
	w.Write(jsonResp)
}
