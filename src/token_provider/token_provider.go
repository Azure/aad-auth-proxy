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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric/global"
)

type identity struct {
	audience      string
	clientId      string
	tenantId      string
	indentityType string
}

type tokenProvider struct {
	token                            string
	ctx                              context.Context
	lastError                        error
	userConfiguredDurationPercentage uint8
	refreshDuration                  time.Duration
	credentialClient                 azcore.TokenCredential
	options                          *policy.TokenRequestOptions
	identity                         identity
}

func NewTokenProvider(audience string, config utils.IConfiguration, certManager *certificate.CertificateManager, logger contracts.ILogger) (contracts.ITokenProvider, error) {

	if config == nil || logger == nil {
		return nil, errors.New("NewTokenProvider: Required arguments canot be nil")
	}

	identityType := config.GetIdentityType()
	aadClientId := config.GetAadClientId()
	aadTenantId := config.GetAadTenantId()
	userConfiguredDurationPercentage := config.GetAadTokenRefreshDurationInPercentage()

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
		ctx:                              context.Background(),
		token:                            "",
		lastError:                        nil,
		userConfiguredDurationPercentage: userConfiguredDurationPercentage,
		credentialClient:                 cred,
		options:                          &policy.TokenRequestOptions{Scopes: []string{audience}},
		identity: identity{
			audience:      audience,
			clientId:      aadClientId,
			tenantId:      aadTenantId,
			indentityType: identityType,
		},
	}

	err = tokenProvider.refreshAADToken()
	if err != nil {
		return nil, errors.New("Failed to get access token: " + err.Error())
	}

	go tokenProvider.periodicallyRefreshClientToken(logger)
	return tokenProvider, nil
}

func (tokenProvider *tokenProvider) GetAccessToken() (string, error) {
	return tokenProvider.token, tokenProvider.lastError
}

func (tokenProvider *tokenProvider) refreshAADToken() error {
	// Record traces
	ctx, span := otel.Tracer(constants.SERVICE_TELEMETRY_KEY).Start(tokenProvider.ctx, "refreshAADToken")
	defer span.End()

	// Telemetry attributes
	attributes := []attribute.KeyValue{
		attribute.String("audience", tokenProvider.identity.audience),
		attribute.String("client_id", tokenProvider.identity.clientId),
		attribute.String("tenant_id", tokenProvider.identity.tenantId),
		attribute.String("identity_type", tokenProvider.identity.indentityType),
	}

	// Record metrics
	// token_refresh_total{is_success}
	meter := global.Meter(constants.SERVICE_TELEMETRY_KEY)
	intrument, _ := meter.Int64Counter(constants.METRIC_TOKEN_REFRESH_TOTAL)

	accessToken, err := tokenProvider.credentialClient.GetToken(ctx, *tokenProvider.options)
	if err != nil {
		attributes = append(attributes, attribute.Bool("is_success", false))
		span.SetAttributes(attributes...)
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to refresh token")
		intrument.Add(ctx, 1, attributes...)

		// Set last error so that this can be returned back when the token is requested
		tokenProvider.lastError = err

		return err
	}

	// Reset last error
	tokenProvider.lastError = nil

	attributes = append(attributes, attribute.Bool("is_success", true))
	intrument.Add(ctx, 1, attributes...)

	tokenProvider.setToken(ctx, accessToken.Token)
	tokenProvider.updateRefreshDuration(accessToken)

	attributes = append(attributes,
		attribute.String("token.expiry_timestamp", accessToken.ExpiresOn.UTC().String()),
		attribute.String("tokenrefresh.next_refresh_timestamp", time.Now().Add(tokenProvider.refreshDuration).UTC().String()),
		attribute.String("tokenrefresh.refresh_duration", tokenProvider.refreshDuration.String()),
	)
	span.SetAttributes(attributes...)
	return nil
}

func (tokenProvider *tokenProvider) periodicallyRefreshClientToken(logger contracts.ILogger) error {
	defer utils.HandlePanic("periodicallyRefreshClientToken")

	for {
		select {
		case <-tokenProvider.ctx.Done():
			return nil
		case <-time.After(tokenProvider.refreshDuration):
			err := tokenProvider.refreshAADToken()
			if err != nil {
				tokenProvider.refreshDuration = time.Duration(constants.TIME_5_MINUTES)
				logger.Error("Failed to refresh token, retry in 5 minutes", err)
				return errors.New("Failed to refresh token: " + err.Error())
			}
		}
	}
}

func (tokenProvider *tokenProvider) setToken(ctx context.Context, token string) {
	var V atomic.Value
	V.Store(token)
	tokenProvider.token = V.Load().(string)
}

func (tokenProvider *tokenProvider) updateRefreshDuration(accessToken azcore.AccessToken) error {
	earlistTime := tokenProvider.getRefreshDuration(accessToken)
	tokenProvider.refreshDuration = earlistTime.Sub(time.Now().UTC())
	return nil
}

func (tokenProvider *tokenProvider) getRefreshDuration(accessToken azcore.AccessToken) time.Time {
	tokenExpiryTimestamp := accessToken.ExpiresOn.UTC()
	userConfiguredTimeFromNow := time.Now().UTC().Add(time.Duration(100-tokenProvider.userConfiguredDurationPercentage) * accessToken.ExpiresOn.Sub(time.Now()) / 100)
	// 10 seconds before now
	thresholdTimestamp := time.Now().UTC().Add(-10 * time.Second)

	// Some times the token expiry time is less than 10 seconds from now or we received an expired token.
	// In that case, we will refresh the token in 1 minute.
	if userConfiguredTimeFromNow.Before(thresholdTimestamp) {
		return time.Now().UTC().Add(constants.TIME_1_MINUTES)
	} else if userConfiguredTimeFromNow.Before(tokenExpiryTimestamp) {
		// If the user configured time is less than the token expiry time, we will use the user configured time.
		return userConfiguredTimeFromNow
	} else {
		return time.Now().UTC().Add(constants.TIME_1_MINUTES)
	}
}
