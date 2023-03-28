# Getting started

This can be deployed in custom templates using release image or as helm chart. Both contain same paramenters which can be set by users to customize deployment. These parameters are called out [below](#parameters).

## Deployment
### Using release image

[sample-proxy-deployment.yaml](../samples/sample-proxy-deployment.yaml) can be used as a starting file for your proxy. Modify necesary parameters and below command to deploy in "observability" namespace.

`kubectl apply -f sample-proxy-deployment.yaml -n observability`

### Using helm chart

Below sample command can be modified with user specific parameters and deployed as a helm chart.

`helm install aad-auth-proxy oci://nagasharedacr.azurecr.io/helm/aad-auth-proxy --version 0.1.0 -n observability --set targetHost=https://azure-monitor-workspace.eastus.prometheus.monitor.azure.com --set identityType=userAssigned --set aadClientId=a711a6a4-1052-44eb-aec8-182e2b604c7f --set audience=https://monitor.azure.com` 


## Parameters

| Image Parameter | Helm chart Parameter name | Description | Supported values | Mandatory |
| --------- | --------- | --------------- | --------- | --------- |
|  TARGET_HOST | targetHost | this is the target host where you want to forward the request to. | | Yes |
|  IDENTITY_TYPE | identityType | this is the identity type which will be used to authenticate requests. This proxy supports 3 types of identities. | systemassigned, userassigned, aadapplication | Yes |
| AAD_CLIENT_ID | aadClientId | this is the client_id of the identity used. This is needed for userassigned and aadapplication identity types. Check [Fetch parameters for identities](IDENTITY.md#fetch-parameters-for-identities) on how to fetch client_id | | Yes for userassigned and aadapplication |
| AAD_TENANT_ID | aadTenantId | this is the tenant_id of the identity used. This is needed for aadapplication identity types. Check [Fetch parameters for identities](IDENTITY.md#fetch-parameters-for-identities) on how to fetch tenant_id | | Yes for aadapplication |
| AAD_CLIENT_CERTIFICATE_PATH | aadClientCertificatePath | this is the path where proxy can find the certificate for aadapplication. This path should be accessible by proxy and should be a pfx certificate. Check [CSI driver](IDENTITY.md#set-up-csi-driver-for-certificate-management) for managing certificates. | | Yes for aadapplication |
| AAD_TOKEN_REFRESH_DURATION_IN_MINUTES | aadTokenRefreshIntervalInMinutes | token will be refreshed based on this duration. Default value is 1 hour. | | No |
| AUDIENCE | audience | this will be the audience for the token | | No |
| LISTENING_PORT | listeningPort | proxy will be listening on this port | | Yes |
| OTEL_SERVICE_NAME | otelServiceName | this will be set as the service name for OTEL traces and metrics. Default value is aad_auth_proxy | | No |
| OTEL_GRPC_ENDPOINT | otelGrpcEndpoint | proxy will push OTEL telemetry to this endpoint. Default values is http://localhost:4317 | | No |

## Liveness and readiness probes
Proxy supports readiness and liveness probes. [Sample configuration](../samples/sample-proxy-deployment.yaml) uses these checks to monitor health of the proxy.

## Example scenarios
### [Query prometheus metrics for KEDA or Kubecost](EXAMPLE_SCENARIOS.md#query-prometheus-metrics-for-kubecost)
### [Ingest prometheus metrics via prometheus remote write](EXAMPLE_SCENARIOS.md#ingest-prometheus-metrics-via-remote-write)
