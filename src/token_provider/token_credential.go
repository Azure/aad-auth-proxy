package token_provider

import (
	"aad-auth-proxy/certificate"
	"aad-auth-proxy/contracts"
	"crypto/x509"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func NewAzureADTokenCredential(tenantId string, clientId string, certManager *certificate.CertificateManager, logger contracts.ILogger) (azcore.TokenCredential, error) {
	cert, _, err := certManager.GetTlsCertificate()
	if err != nil {
		return nil, err
	}
	certOptions := &azidentity.ClientCertificateCredentialOptions{SendCertificateChain: true}
	cred, err := azidentity.NewClientCertificateCredential(tenantId, clientId, []*x509.Certificate{cert.Leaf}, cert.PrivateKey, certOptions)
	if err != nil {
		logger.Error("Client Certificate Credential couldn't be created:", err)
		return nil, err
	}

	return cred, nil
}

func NewManagedIdentityTokenCredential(managedIdentityClientId string, logger contracts.ILogger) (azcore.TokenCredential, error) {
	if logger == nil {
		return nil, errors.New("Required params are missing to create token provider")
	}

	var cred *azidentity.ManagedIdentityCredential
	var err error

	if len(managedIdentityClientId) > 0 {
		clientId := azidentity.ClientID(managedIdentityClientId)
		opts := &azidentity.ManagedIdentityCredentialOptions{ID: clientId}
		cred, err = azidentity.NewManagedIdentityCredential(opts)
	} else {
		cred, err = azidentity.NewManagedIdentityCredential(nil)
	}

	if err != nil {
		logger.Error("ManagedIdentity Credential couldn't be created:", err)
		return nil, err
	}

	return cred, nil
}
