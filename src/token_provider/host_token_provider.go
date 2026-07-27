// filepath: host_token_provider.go
package token_provider

import (
    "aad-auth-proxy/certificate"
    "aad-auth-proxy/contracts"
    "aad-auth-proxy/utils"
    "fmt"
    "net"
    "strings"
)

type HostTokenProvider struct {
    defaultProvider contracts.ITokenProvider
    perHostProvider map[string]contracts.ITokenProvider
}

func NewHostTokenProvider(
    defaultAudience string,
    hostAudience map[string]string,
    configuration utils.IConfiguration,
    certManager *certificate.CertificateManager,
    logger contracts.ILogger,
) (*HostTokenProvider, error) {
    defaultProvider, err := NewTokenProvider(defaultAudience, configuration, certManager, logger)
    if err != nil {
        return nil, err
    }

    perHost := make(map[string]contracts.ITokenProvider, len(hostAudience))
    for h, aud := range hostAudience {
        p, e := NewTokenProvider(aud, configuration, certManager, logger)
        if e != nil {
            return nil, e
        }
        perHost[normalizeHost(h)] = p
    }

    return &HostTokenProvider{
        defaultProvider: defaultProvider,
        perHostProvider: perHost,
    }, nil
}

// Compat: espone vari nomi possibili usati dal progetto.
func (h *HostTokenProvider) GetToken() (string, error) {
    return getTokenFromProvider(h.defaultProvider)
}

func (h *HostTokenProvider) GetClientToken() (string, error) {
    return getTokenFromProvider(h.defaultProvider)
}

func (h *HostTokenProvider) GetAccessToken() (string, error) {
    return getTokenFromProvider(h.defaultProvider)
}

func (h *HostTokenProvider) GetTokenForHost(host string) (string, error) {
    n := normalizeHost(host)
    if p, ok := h.perHostProvider[n]; ok {
        return getTokenFromProvider(p)
    }
    return getTokenFromProvider(h.defaultProvider)
}

func normalizeHost(host string) string {
    host = strings.TrimSpace(strings.ToLower(host))
    onlyHost, _, err := net.SplitHostPort(host)
    if err == nil {
        return onlyHost
    }
    return host
}

func getTokenFromProvider(p contracts.ITokenProvider) (string, error) {
    // prova le firme più comuni
    if tp, ok := any(p).(interface{ GetAccessToken() (string, error) }); ok {
        return tp.GetAccessToken()
    }
    if tp, ok := any(p).(interface{ GetClientToken() (string, error) }); ok {
        return tp.GetClientToken()
    }
    if tp, ok := any(p).(interface{ GetToken() (string, error) }); ok {
        return tp.GetToken()
    }
    return "", fmt.Errorf("unsupported token getter on provider type %T", p)
}