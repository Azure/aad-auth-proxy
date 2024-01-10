# AAD Auth Proxy

## Release 01-10-2024

- Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-01-10-2024-08b31473`
- Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-01-10-2024-08b31473`
- Change log:
  - Bump golang.org/x/crypto from 0.14.0 to 0.17.0
  - Bump google.golang.org/grpc from 1.53.0 to 1.56.3
  - Bump golang.org/x/net from 0.7.0 to 0.17.0
  - Fixed trivy vulnerabilities CVE-2023-39325 and GHSA-m425-mq94-257g
  - Fixed helm chart apiVersion typo

## Release 08-23-2023

- Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-08-23-2023-5988d874`
- Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-08-23-2023-5988d874`
- Change log:
  - Retry initialization indefinitely on failure
  - Return 503 on failure to establish connection with remote host

## Release 07-17-2023

- Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-07-17-2023-841abb6f`
- Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-07-17-2023-841abb6f`
- Change log:
  - Feature: Decode gzip and zlib response body before logging to console.
  - Feature (experimental): Override request headers while forwarding to host.
  - Enrich traces and metrics.
  - Bump google.golang.org/grpc version from 1.52.3 to 1.53.0.

## Release 05-24-2023

- Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-05-24-2023-b911fe1c`
- Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-05-24-2023-b911fe1c`
- Change log:
  - Fix helmchart.

## Release 04-12-2023

- Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:aad-auth-proxy-0.1.0-main-04-11-2023-623473b0`
- Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-04-11-2023-623473b0`
- Change log:
  - Public preview release image for AAD auth proxy.
