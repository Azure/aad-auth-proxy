using './main.bicep'

param location = readEnvironmentVariable('AAD_AUTH_PROXY_ACR_LOCATION', 'westus')
param resourceGroupName = readEnvironmentVariable('AAD_AUTH_PROXY_ACR_RESOURCE_GROUP_NAME')
param acrName = readEnvironmentVariable('AAD_AUTH_PROXY_ACR_NAME')
param networkRuleSetIpRules = split(readEnvironmentVariable('AAD_AUTH_PROXY_ACR_NETWORK_RULE_IPS'), ',')
param tags = {
  workload: 'aad-auth-proxy'
}
