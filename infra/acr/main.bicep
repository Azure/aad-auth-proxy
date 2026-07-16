targetScope = 'subscription'

@description('Azure region where the resource group and Azure Container Registry will be deployed.')
param location string = 'westus'

@description('Name of the resource group to create for the Azure Container Registry.')
param resourceGroupName string

@description('Name of the Azure Container Registry. Must be globally unique and contain only alphanumeric characters.')
@minLength(5)
@maxLength(50)
param acrName string

@description('IPv4 CIDR ranges allowed through the registry firewall. Supply the trusted build-pool egress IPs at deployment time; do not check them into this public repo.')
param networkRuleSetIpRules string[]

@description('Tags to apply to deployed resources.')
param tags object = {
  workload: 'aad-auth-proxy'
}

@description('Firewall IP rules in AVM/ARM shape.')
var formattedAcrIpRules = [for ip in networkRuleSetIpRules: {
  action: 'Allow'
  value: ip
}]

@description('Firewall IP rules in AVM/ARM shape. Fail closed if no valid IP rules are supplied.')
var acrIpRules = (empty(networkRuleSetIpRules) || contains(networkRuleSetIpRules, ''))
  ? fail('networkRuleSetIpRules must include at least one trusted build-pool egress CIDR.')
  : formattedAcrIpRules

module acrResourceGroup 'br/public:avm/res/resources/resource-group:0.4.3' = {
  name: 'acr-resource-group'
  params: {
    name: resourceGroupName
    location: location
    tags: tags
  }
}

module acr 'br/public:avm/res/container-registry/registry:0.12.1' = {
  name: 'acr'
  scope: resourceGroup(resourceGroupName)
  params: {
    name: acrName
    location: location
    acrSku: 'Premium'
    publicNetworkAccess: 'Enabled'
    networkRuleSetDefaultAction: 'Deny'
    networkRuleSetIpRules: acrIpRules
    networkRuleBypassOptions: 'None'
    tags: tags
  }
  dependsOn: [
    acrResourceGroup
  ]
}

output resourceGroupName string = acrResourceGroup.outputs.name
output resourceGroupResourceId string = acrResourceGroup.outputs.resourceId
output acrName string = acr.outputs.name
output acrResourceId string = acr.outputs.resourceId
output acrLoginServer string = acr.outputs.loginServer
