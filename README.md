# Status
| Step | Build | Release |
| -- | -- | -- |
| Image | [![Build Status](https://github-private.visualstudio.com/azure/_apis/build/status/Azure.aad-auth-proxy?branchName=main&jobName=Build%20image)](https://github-private.visualstudio.com/azure/_build/latest?definitionId=449&branchName=main) | [![Release Status](https://github-private.vsrm.visualstudio.com/_apis/public/Release/badge/2d36c31d-2f89-409f-9a3e-32e4e9699840/80/129)](https://github-private.visualstudio.com/azure/_release?_a=releases&view=mine&definitionId=80)
| Helm chart | [![Build Status](https://github-private.visualstudio.com/azure/_apis/build/status/Azure.aad-auth-proxy?branchName=main&jobName=Package%20helm%20chart)](https://github-private.visualstudio.com/azure/_build/latest?definitionId=440&branchName=main) | [![Release Status](https://github-private.vsrm.visualstudio.com/_apis/public/Release/badge/2d36c31d-2f89-409f-9a3e-32e4e9699840/80/129)](https://github-private.visualstudio.com/azure/_release?_a=releases&view=mine&definitionId=80)

# Project
This project is Azure AD proxy, which is a forward proxy to authenticate requests to backend services (example: Azure Monitor workspaces used to store data for Azure Monitor managed service for Prometheus). Clients can use system identity, or user identity (kubelet identity) or Azure AD app to fetch tokens which will be added as bearer tokens to all forwarded requests.

## Getting Started
This can be deployed in custom templates using release image as a side car or a service. This can be deployed using helm chart as well, which will be deployed as a service. Detailed instructions on how to deploy can be found [here](./docs/getting-started/GETTING_STARTED.md).

## Telemetry
This has been instrumented with [OTEL](https://opentelemetry.io/), it emits traces and metrics, which can be collected using [OTEL Collector](https://github.com/open-telemetry/opentelemetry-collector). A grafana dashboard to visualize metrics is also included.

## Securing traffic
This proxy can be deployed as a side car or as a service. When deployed as a side car, only the containers within the pod can access this proxy, but when deployed as a service without restricting traffic, any container can access this proxy. So there might be a need to secure traffic to proxy pod and can be achieved using [Network policies in Azure Kubernetes Service](https://learn.microsoft.com/azure/aks/use-network-policies).

## Limitations
Only helm v3 is supported.
