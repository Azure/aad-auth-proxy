# Example scenarios

## Query Prometheus metrics for KEDA or Kubecost
In this example we will deploy proxy to query Prometheus metrics from Azure Monitor Workspace in order to scale out deployments using [KEDA](https://keda.sh/). 
1. [sample-proxy-query.yaml](../samples/sample-proxy-query.yaml) can be used as a template to modify [parameters](GETTING_STARTED.md#parameters).
2. [Create Azure monitor workspace](https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/azure-monitor-workspace-manage?tabs=azure-portal#create-an-azure-monitor-workspace).
3. Use "Query endpoint" from Azure monitor workspace oiverview page as TARGET_HOST.
4. Modify identity parameters based on the identity type you choose. Assign [read](IDENTITY.md#add-read-permissions) permissions to identity.
    - [System identity](IDENTITY.md#system-identity)
    - [User identity](IDENTITY.md#user-identity)
    - [AAD application](IDENTITY.md#aad-application)
5. AUDIENCE has to be "https://prometheus.monitor.azure.com" for querying metrics from Azure Monitor Workspace.
6. Change OTEL_GRPC_ENDPOINT to receive endpoint of OTEL collector if deployed, else remove it.
7. Deploy proxy using command: `kubectl apply -f sample-proxy-query.yaml -n observability`
8. Deploy KEDA and in prometheus scaled object point to this proxy endpoint.

## Ingest Prometheus metrics via remote write
In this example we will deploy proxy to ingest Prometheus metrics to Azure Monitor Workspace via [prometheus remote write](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#remote_write). 
1. [sample-proxy-ingestion.yaml](../samples/sample-proxy-ingestion.yaml) can be used as a template to modify [parameters](GETTING_STARTED.md#parameters).
2. [Create Azure monitor workspace](https://learn.microsoft.com/azure/azure-monitor/essentials/azure-monitor-workspace-manage?tabs=azure-portal#create-an-azure-monitor-workspace).
3. Use "Metrics ingestion endpoint" from Azure monitor workspace oiverview page as TARGET_HOST.
*Note: metrics ingestion endpoint will look like "https://naga-aad-auth-proxy-amw-cuqx.eastus-1.metrics.ingest.monitor.azure.com/dataCollectionRules/dcr-65cb9d21936f43e3b2035d2/streams/Microsoft-PrometheusMetrics/api/v1/write?api-version=2021-11-01-preview", here only pick host part for TARGET_HOST, which is "https://naga-aad-auth-proxy-amw-cuqx.eastus-1.metrics.ingest.monitor.azure.com". The remaining path will be used in step 8 in remote write endpoint.*
4. Modify identity parameters based on the identity type you choose. Assign [write](IDENTITY.md#add-write-permissions) permissions to identity.
    - [System identity](IDENTITY.md#system-identity)
    - [User identity](IDENTITY.md#user-identity)
    - [AAD application](IDENTITY.md#aad-application)
5. AUDIENCE has to be "https://monitor.azure.com" for ingesting metrics to Azure Monitor Workspace.
6. Change OTEL_GRPC_ENDPOINT to receive endpoint of OTEL collector if deployed, else remove it.
7. Deploy proxy using command: `kubectl apply -f sample-proxy-ingestion.yaml -n observability`
8. Update Prometheus' remote write configuration to point to proxy endpoint as host and the path from "Metrics ingestion endpoint" in step 3.

<pre>
remoteWrite:
    - url: "http://azuremonitor-ingestion.observability.svc.cluster.local/dataCollectionRules/dcr-65cb9d21936f43e3b2035d2/streams/Microsoft-PrometheusMetrics/api/v1/write?api-version=2021-11-01-preview"
</pre>