# Azure Container Registry

This template creates a new resource group and a Premium Azure Container Registry for `aad-auth-proxy` by using Azure Verified Modules.

## Deploy

```bash
export AAD_AUTH_PROXY_ACR_LOCATION="<region>"
export AAD_AUTH_PROXY_ACR_RESOURCE_GROUP_NAME="<resource-group-name>"
export AAD_AUTH_PROXY_ACR_NAME="<acr-name>"
export AAD_AUTH_PROXY_ACR_NETWORK_RULE_IPS="<trusted-build-pool-cidr-1>,<trusted-build-pool-cidr-2>"

az deployment sub what-if \
  --location "$AAD_AUTH_PROXY_ACR_LOCATION" \
  --template-file infra/acr/main.bicep \
  --parameters infra/acr/main.bicepparam

az deployment sub create \
  --location "$AAD_AUTH_PROXY_ACR_LOCATION" \
  --template-file infra/acr/main.bicep \
  --parameters infra/acr/main.bicepparam
```

The requested ACR display name `aad-auth-proxy-acr` should be represented with an alphanumeric registry name because Azure Container Registry names must be globally unique and contain only alphanumeric characters.

`AAD_AUTH_PROXY_ACR_NETWORK_RULE_IPS` should be set to the trusted build-pool egress CIDRs. Keep the concrete values in private deployment configuration rather than checking them into this public repository.
