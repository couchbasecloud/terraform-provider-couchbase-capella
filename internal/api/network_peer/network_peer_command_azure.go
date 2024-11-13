package network_peer

// GetAzureVNetPeeringCommandRequest request to retrieve the role assignment command or script to be executed in the Azure CLI to assign a new network contributor role.
type GetAzureVNetPeeringCommandRequest struct {
	// ResourceGroup The resource group name holding the resource you’re connecting with Capella.
	ResourceGroup string `json:"resourceGroup"`

	// SubscriptionId Subscription ID is a GUID that uniquely identifies your subscription to use Azure services. To find your subscription ID, see [Find your Azure subscription](https://learn.microsoft.com/en-us/azure/azure-portal/get-subscription-tenant-id#find-your-azure-subscription).
	SubscriptionId string `json:"subscriptionId"`

	// TenantId The Azure tenant ID. To find your tenant ID, see [How to find your Azure Active Directory tenant ID](https://learn.microsoft.com/en-us/entra/fundamentals/how-to-find-tenant).
	TenantId string `json:"tenantId"`

	// VnetId The VNet ID is the name of the virtual network in Azure.
	VnetId string `json:"vnetId"`

	// VnetPeeringServicePrincipal The enterprise application object ID for the Capella service principal. You can find the enterprise application object ID in Azure by selecting Azure Active Directory -> Enterprise applications. Next, select the application name, the object ID is in the Object ID box.
	VnetPeeringServicePrincipal string `json:"vnetPeeringServicePrincipal"`
}

// GetAzureVNetPeeringCommandResponse retrieves the role assignment command or script to be executed in the Azure CLI to assign a new network contributor role.
type GetAzureVNetPeeringCommandResponse struct {
	// Command The command to be run by the customer in is their external azure account in order to grant the service principal a network contributor role that is required for VNET peering.
	Command string `json:"command"`
}
