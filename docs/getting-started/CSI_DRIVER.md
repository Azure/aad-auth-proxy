# Set up CSI driver for certificate management
*Note: Azure Key Vault CSI driver configuration is just one of the ways to get certificate mounted on the pod. Proxy container only needs a local path to a certificate in the pod for the setting `AAD_CLIENT_CERTIFICATE_PATH` when using AAD application tpe for authenticating requests.*

This step is only required if you didn't enable Azure Key Vault Provider for Secrets Store CSI Driver when you created your cluster.
1. Run the following Azure CLI command to enable Azure Key Vault Provider for Secrets Store CSI Driver for your cluster. `az aks enable-addons --addons azure-keyvault-secrets-provider --name <aks-cluster-name> --resource-group <resource-group-name>`
2. Run the following commands to give the identity access to the key vault.
    ```azurecli
    # show client id of the managed identity of the cluster
    az aks show -g <resource-group> -n <cluster-name> --query addonProfiles.azureKeyvaultSecretsProvider.identity.clientId -o tsv

    # set policy to access keys in your key vault
    az keyvault set-policy -n <keyvault-name> --key-permissions get --spn <identity-client-id>

    # set policy to access secrets in your key vault
    az keyvault set-policy -n <keyvault-name> --secret-permissions get --spn <identity-client-id>
    
    # set policy to access certs in your key vault
    az keyvault set-policy -n <keyvault-name> --certificate-permissions get --spn <identity-client-id>
    ```
3.  Create a *SecretProviderClass* by saving the following YAML to a file named *secretproviderclass.yml*. Replace the values for `userAssignedIdentityID`, `keyvaultName`, `tenantId` and the objects to retrieve from your key vault. See [Provide an identity to access the Azure Key Vault Provider for Secrets Store CSI Driver](https://learn.microsoft.com/azure/aks/csi-secrets-store-identity-access) for details on values to use.

    ```yml
    # This is a SecretProviderClass example using user-assigned identity to access your key vault
    apiVersion: secrets-store.csi.x-k8s.io/v1
    kind: SecretProviderClass
    metadata:
    name: azure-kvname-user-msi
    spec:
    provider: azure
    parameters:
        usePodIdentity: "false"
        useVMManagedIdentity: "true"          # Set to true for using managed identity
        userAssignedIdentityID:  <client-id> # Set the clientID of the user-assigned managed identity to use
        keyvaultName: <key-vault-name> # Set to the name of your key vault
        cloudName: ""                         # [OPTIONAL for Azure] if not provided, the Azure environment defaults to AzurePublicCloud
        objects:  |
        array:
            - |
            objectName: <name-of-cert>
            objectType: secret        # object types: secret, key, or cert
            objectFormat: pfx
            objectEncoding: base64
            objectVersion: ""
        tenantId: <tenant-id> # The tenant ID of the key vault
    ```
4. Apply the *SecretProviderClass* by running the following command on your cluster.

    `
    kubectl apply -f secretproviderclass.yml
    `
