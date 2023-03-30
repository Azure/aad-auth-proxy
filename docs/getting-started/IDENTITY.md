# Identity management
Proxy uses identity to authenticate requests and has 3 options: system assigned, user assigned and AAD application. In this section we show how to configure proxy to run on [Azure Kubernetes service](https://learn.microsoft.com/azure/aks/intro-kubernetes). 
[Managed identity](https://learn.microsoft.com/azure/aks/use-managed-identity) can be enabled on [Azure Kubernetes service](https://learn.microsoft.com/azure/aks/intro-kubernetes), which can be used to authenticate requests.

## System identity
1. The node resource group of the AKS cluster contains resources that you will require for other steps in this process. This resource group has the name `MC_<AKS-RESOURCE-GROUP>_<AKS-CLUSTER-NAME>_<REGION>`. You can locate it from the **Resource groups** menu in the Azure portal. Start by making sure that you can locate this resource group since other steps below will refer to it.
![Resource groups](../images/resource-groups.png)
2. Navigate to VMSS with the name `aks-agentpool-<ID>-vmss`. Select **Identity** TOC and **System assigned** tab, and toggle **Status** to **On**. This will enable system assigned identity on the underlying VMSS of AKS cluster.
![Systme assigned](../images/system-identity.png)
3. **Object (principal) ID** should be used as [AAD_CLIENT_ID](GETTING_STARTED.md#parameters) when using [IDENTITY_TYPE](GETTING_STARTED.md#parameters) as **systemassigned**.

## User identity
1. Managed identity can be enabled on while [creating](https://learn.microsoft.com/azure/aks/use-managed-identity#create-an-aks-cluster-using-a-managed-identity) AKS or can be [updated](https://learn.microsoft.com/azure/aks/use-managed-identity#update-an-aks-cluster-to-use-a-managed-identity) at a later point in time.
2. Run command `az aks show -g <AKS-CLUSTER-RESOURCE-GROUP> -n <AKS-CLUSTER-NAME> --query "identityProfile"` and pick `kubeletidentity.clientId`. This should be used as [AAD_CLIENT_ID](GETTING_STARTED.md#parameters) when using [IDENTITY_TYPE](GETTING_STARTED.md#parameters) as **userassigned**.

## AAD application
1. Follow the procedure at [Register an application with Azure AD and create a service principal](https://learn.microsoft.com/azure/active-directory/develop/howto-create-service-principal-portal#register-an-application-with-azure-ad-and-create-a-service-principal) to register an application for Prometheus remote-write and create a service principal.
2. From the **Azure Active Directory** menu in Azure Portal, select **App registrations**. Locate your application and note the client ID. This should be used as [AAD_CLIENT_ID](GETTING_STARTED.md#parameters) when using [IDENTITY_TYPE](GETTING_STARTED.md#parameters) as **aadapplication**.
![AAD application client ID](../images/application-client-id.png)
3. Have a local copy of ceritifate and share the path to it as [AAD_CLIENT_CERTIFICATE_PATH](GETTING_STARTED.md#parameters). If you are running your process in Azure the best and most secure way to manage you cert is to place it in Azure KeyVault and mount to your pod using [CSI dirver](CSI_DRIVER.md#set-up-csi-driver-for-certificate-management).
