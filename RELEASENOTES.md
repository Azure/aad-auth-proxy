# AAD Auth Proxy

## Release 07-17-2023
* Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-07-17-2023-841abb6f`
* Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-07-17-2023-841abb6f`
* Change log:
    * Feature: Decode gzip and zlib response body before logging to console.
    * Feature (experimental): Override request headers while forwarding to host.
    * Enrich traces and metrics.
    * Bump google.golang.org/grpc version from 1.52.3 to 1.53.0.

## Release 05-24-2023

* Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:0.1.0-main-05-24-2023-b911fe1c`
* Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-05-24-2023-b911fe1c`
* Change log:
    * Fix helmchart.

## Release 04-12-2023

* Image: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/images/aad-auth-proxy:aad-auth-proxy-0.1.0-main-04-11-2023-623473b0`
* Helm chart: `mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy:0.1.0-main-04-11-2023-623473b0`
* Change log:
    * Public preview release image for AAD auth proxy.
