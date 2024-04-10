# Getting started

This can be deployed in custom templates using release image or as helm chart. Both contain same paramenters which can be set by users to customize deployment. These parameters are called out [below](#parameters).

## Deployment
### Using release image

[sample-proxy-deployment.yaml](../samples/sample-proxy-deployment.yaml) can be used as a starting file for your proxy. Modify necesary parameters and below command to deploy in "observability" namespace.

`kubectl apply -f sample-proxy-deployment.yaml -n observability`

### Using helm chart

Below sample command can be modified with user specific parameters and deployed as a helm chart.

`helm install aad-auth-proxy oci://mcr.microsoft.com/azuremonitor/auth-proxy/prod/aad-auth-proxy/helmchart/aad-auth-proxy --version 0.1.0-main-08-23-2023-5988d874 -n observability --set targetHost=https://azure-monitor-workspace.eastus.prometheus.monitor.azure.com --set identityType=userAssigned --set aadClientId=a711a6a4-1052-44eb-aec8-182e2b604c7f --set audience=https://monitor.azure.com/.default`


## Parameters

| Image Parameter | Helm chart Parameter name | Description | Supported values | Mandatory |
| --------- | --------- | --------------- | --------- | --------- |
|  TARGET_HOST | targetHost | this is the target host where you want to forward the request to. | | Yes |
|  IDENTITY_TYPE | identityType | this is the identity type which will be used to authenticate requests. This proxy supports 3 types of identities. If this value is not set, it will create [DefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential) | systemassigned, userassigned, aadapplication | No |
| AAD_CLIENT_ID | aadClientId | this is the client_id of the identity used. This is needed for userassigned and aadapplication identity types. Check [Fetch parameters for identities](IDENTITY.md#fetch-parameters-for-identities) on how to fetch client_id | | Yes for userassigned and aadapplication |
| AAD_TENANT_ID | aadTenantId | this is the tenant_id of the identity used. This is needed for aadapplication identity types. Check [Fetch parameters for identities](IDENTITY.md#fetch-parameters-for-identities) on how to fetch tenant_id | | Yes for aadapplication |
| AAD_CLIENT_CERTIFICATE_PATH | aadClientCertificatePath | this is the path where proxy can find the certificate for aadapplication. This path should be accessible by proxy and should be a either a pfx or pem certificate containing private key. Check [CSI driver](IDENTITY.md#set-up-csi-driver-for-certificate-management) for managing certificates. | | Yes for aadapplication |
| AAD_TOKEN_REFRESH_INTERVAL_IN_PERCENTAGE | aadTokenRefreshIntervalInPercentage | token will be refreshed based on the percentage of time till token expiry. Default value is 10% time before expiry. | | No |
| AUDIENCE | audience | this will be the audience for the token | | No |
| LISTENING_PORT | listeningPort | proxy will be listening on this port | | Yes |
| OTEL_SERVICE_NAME | otelServiceName | this will be set as the service name for OTEL traces and metrics. Default value is aad_auth_proxy | | No |
| OTEL_GRPC_ENDPOINT | otelGrpcEndpoint | proxy will push OTEL telemetry to this endpoint. Default values is http://localhost:4317 | | No |
| OVERRIDE_REQUEST_HEADERS | overrideRequestHeaders | (Experimental) proxy will override these headers while forwarding requests. This expects headers to be in JSON, example {"header1": "value1", "header2": "value2" }:  Default values is {} | | No |

## Liveness and readiness probes
Proxy supports readiness and liveness probes. [Sample configuration](../samples/sample-proxy-deployment.yaml) uses these checks to monitor health of the proxy.

## Default Azure credentials
DefaultAzureCredential is intended to simplify getting started by relying on default behaviors of azidentity. Developers who want more control or whose scenario isn't served by the default settings should use other credential types by setting parameter: IDENTITY_TYPE.
Proxy supports workload identity via [DefaultAzureCredential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential). [sample configuration using workload indentity via default Azure credentials](../samples/sample-proxy-using-workload-identity-default.yaml) has example configurations to make workload indetity work.

## Example scenarios
### [Query prometheus metrics for KEDA or Kubecost](EXAMPLE_SCENARIOS.md#query-prometheus-metrics-for-kubecost)
### [Ingest prometheus metrics via prometheus remote write](EXAMPLE_SCENARIOS.md#ingest-prometheus-metrics-via-remote-write)
