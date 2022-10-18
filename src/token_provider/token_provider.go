package token_provider

import (
	"aad-auth-proxy/certificate"
	"aad-auth-proxy/constants"
	"aad-auth-proxy/contracts"
	"aad-auth-proxy/utils"
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type tokenProvider struct {
	token            string
	ctx              context.Context
	refreshDuration  time.Duration
	credentialClient azcore.TokenCredential
	options          *policy.TokenRequestOptions
}

func NewTokenProvider(audience string, config utils.IConfiguration, certManager *certificate.CertificateManager, logger contracts.ILogger) (contracts.ITokenProvider, error) {

	if config == nil || logger == nil {
		return nil, errors.New("NewTokenProvider: Required arguments canot be nil")
	}

	identityType := config.GetIdentityType()
	aadClientId := config.GetAadClientId()
	aadTenantId := config.GetAadTenantId()

	var cred azcore.TokenCredential
	var err error

	switch identityType {
	case constants.SYSTEM_ASSIGNED:
		cred, err = NewManagedIdentityTokenCredential("", logger)
	case constants.USER_ASSIGNED:
		if len(aadClientId) > 0 {
			cred, err = NewManagedIdentityTokenCredential(aadClientId, logger)
		} else {
			logger.Error("Client ID not found for UserAssignedIdentity Auth", errors.New("No Client ID"))
			return nil, errors.New("No Client ID")
		}
	case constants.AAD_APPLICATION:
		if len(aadClientId) > 0 && len(aadTenantId) > 0 && certManager != nil {
			cred, err = NewAzureADTokenCredential(aadTenantId, aadClientId, certManager, logger)
		} else {
			logger.Error("Required pararms not found for AAD App Auth", errors.New("AAD params missing"))
			return nil, errors.New("AAD params missing")
		}
	default:
		cred, err = azidentity.NewDefaultAzureCredential(nil)
	}

	if err != nil {
		return nil, err
	}

	tokenProvider := &tokenProvider{
		ctx:              context.Background(),
		credentialClient: cred,
		options:          &policy.TokenRequestOptions{Scopes: []string{audience}},
	}

	err = tokenProvider.setAccessToken()
	if err != nil {
		return nil, errors.New("Failed to get access token: " + err.Error())
	}

	go tokenProvider.periodicallyRefreshClientToken(logger)
	return tokenProvider, nil
}

func (tokenProvider *tokenProvider) GetAccessToken() string {
	return tokenProvider.token
}

func (tokenProvider *tokenProvider) setAccessToken() error {
	accessToken, err := tokenProvider.credentialClient.GetToken(tokenProvider.ctx, *tokenProvider.options)
	if err != nil {
		return err
	}

	tokenProvider.setToken(accessToken.Token)
	tokenProvider.updateRefreshDuration(accessToken)
	return nil
}

func (tokenProvider *tokenProvider) periodicallyRefreshClientToken(logger contracts.ILogger) error {
	defer utils.HandlePanic("periodicallyRefreshClientToken")

	for {
		select {
		case <-tokenProvider.ctx.Done():
			return nil
		case <-time.After(tokenProvider.refreshDuration):
			err := tokenProvider.setAccessToken()
			if err != nil {
				logger.Error("Failed to refresh token", err)
				return errors.New("Failed to refresh token: " + err.Error())
			}
		}
	}
}

func (tokenProvider *tokenProvider) setToken(token string) {
	var V atomic.Value
	V.Store(token)
	tokenProvider.token = V.Load().(string)
}

func (tokenProvider *tokenProvider) updateRefreshDuration(accessToken azcore.AccessToken) error {
	tokenExpiryTimestamp := accessToken.ExpiresOn.UTC()
	deltaExpirytime := tokenExpiryTimestamp.UTC().Add(-time.Minute * 5)
	if deltaExpirytime.After(time.Now().UTC()) {
		tokenProvider.refreshDuration = deltaExpirytime.Sub(time.Now().UTC())
	} else {
		return errors.New("Access Token expiry is less than the current time")
	}

	return nil
}
