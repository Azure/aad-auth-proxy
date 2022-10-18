package handler

import (
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httputil"
)

// This manages token provider handler
type Handler struct {
	proxy         *httputil.ReverseProxy
	tokenProvider contracts.ITokenProvider
}

// Creates a new handler
func NewHandler(proxy *httputil.ReverseProxy, tokenProvider contracts.ITokenProvider) (handler *Handler, err error) {
	if proxy == nil {
		return nil, errors.New("proxy cannot be nil")
	}

	if tokenProvider == nil {
		return nil, errors.New("tokenProvider cannot be nil")
	}

	return &Handler{
		proxy:         proxy,
		tokenProvider: tokenProvider,
	}, nil
}

// Reverse proxy handler
func (handler *Handler) ReverseProxyHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Set(constants.HEADER_AUTHORIZATION, "Bearer "+handler.tokenProvider.GetAccessToken())
	handler.proxy.ServeHTTP(w, r)
}

// Readiness check handler
func (handler *Handler) ReadinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	if len(handler.tokenProvider.GetAccessToken()) != 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		resp := make(map[string]string)
		resp[constants.ERROR_PROPERTY_CODE] = "AuthenticationTokenNotFound"
		resp[constants.ERROR_PROPERTY_MESSAGE] = "cannot get authentication token"
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Write(jsonResp)
	}
}
