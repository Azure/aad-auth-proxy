package main

import (
	"aad-auth-proxy/certificate"
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"aad-auth-proxy/handler"
	"aad-auth-proxy/telemetry"
	"aad-auth-proxy/token_provider"
	"aad-auth-proxy/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Main entry point.
func main() {

	// Handle panics
	defer utils.HandlePanic("main")

	logger := telemetry.NewLogger()

	// Read configuration parameters
	configuration := utils.NewConfiguration()
	listeningPort := configuration.GetListeningPort()
	audience := configuration.GetAudience()
	targetHost := configuration.GetTargetHost()

	// Create handler to fetch tokens
	handler := createHandlerWithTokenProvider(configuration, audience, targetHost, logger)

	// Note: These (Health and Readiness) check whether app can fetch a token with provided identity,
	// it does not evaluate end-to-end authentication and authorization of acquired tokens.
	// Health check handler
	http.HandleFunc("/health", handler.ReadinessCheckHandler)

	// Readiness check handler
	http.HandleFunc("/ready", handler.ReadinessCheckHandler)

	// Reverse proxy handler
	http.HandleFunc("/", handler.ReverseProxyHandler)

	// Listen at specified port
	log.Fatal(http.ListenAndServe(":"+listeningPort, nil))
}

// Creates handler with token provider
func createHandlerWithTokenProvider(configuration utils.IConfiguration, audience string, targetHost string, logger contracts.ILogger) *handler.Handler {
	var certManager *certificate.CertificateManager
	var err error

	// Create cert manager for AAD based authentication
	if configuration.GetIdentityType() == constants.AAD_APPLICATION {
		certManager, err = certificate.NewCerificateManager(configuration.GetAadClientCertPath())
		if err != nil {
			logger.Error("NewCertificateManager failed: ", err)
		}
	}

	// Create TokenProvider
	tokenProvider, err := token_provider.NewTokenProvider(audience, configuration, certManager, logger)
	if err != nil {
		logger.Error("TokenCredential creation failed:", err)
	}

	proxy, err := createReverseProxy(targetHost, tokenProvider)
	if err != nil {
		logger.Error("Proxy creation failed:", err)
	}

	// Create handler to return tokens based on audience
	handler, err := handler.NewHandler(proxy, tokenProvider)
	if err != nil {
		logger.Error("NewHandler failed:", err)
	}
	return handler
}

func createReverseProxy(targetHost string, tokenProvider contracts.ITokenProvider) (*httputil.ReverseProxy, error) {
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

func modifyRequest(request *http.Request, targetHost string, tokenProvider contracts.ITokenProvider) {
	request.URL.Scheme = constants.HTTPS_SCHEME
	request.URL.Host = targetHost
	request.Host = targetHost
}

func handleError(response http.ResponseWriter, request *http.Request, err error) {
	log.WithFields(log.Fields{
		"Request": request.URL.String(),
	}).Errorln("Request failed", err)
}

func modifyResponse(response *http.Response) (err error) {
	log.WithFields(log.Fields{
		"Request":       response.Request.URL.String(),
		"StatusCode":    response.StatusCode,
		"Status":        response.Status,
		"ContentLength": response.ContentLength,
	}).Infoln("Successfully send request, returning response back.")

	// If server returned error, log response as well
	if response.StatusCode >= http.StatusBadRequest {
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
