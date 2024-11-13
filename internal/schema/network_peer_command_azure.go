package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// AzureVNetPeeringCommandRequest to retrieve the role assignment command or script to be executed in the Azure CLI to assign a new network contributor role.
type AzureVNetPeeringCommandRequest struct {

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the network peer needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ResourceGroup The resource group name holding the resource you’re connecting with Capella.
	ResourceGroup types.String `tfsdk:"resource_group"`

	// SubscriptionId Subscription ID is a GUID that uniquely identifies your subscription to use Azure services. To find your subscription ID, see [Find your Azure subscription](https://learn.microsoft.com/en-us/azure/azure-portal/get-subscription-tenant-id#find-your-azure-subscription).
	SubscriptionId types.String `tfsdk:"subscription_id"`

	// TenantId The Azure tenant ID. To find your tenant ID, see [How to find your Azure Active Directory tenant ID](https://learn.microsoft.com/en-us/entra/fundamentals/how-to-find-tenant).
	TenantId types.String `tfsdk:"tenant_id"`

	// VnetId The VNet ID is the name of the virtual network in Azure.
	VnetId types.String `tfsdk:"vnet_id"`

	// VnetPeeringServicePrincipal The enterprise application object ID for the Capella service principal. You can find the enterprise application object ID in Azure by selecting Azure Active Directory -> Enterprise applications. Next, select the application name, the object ID is in the Object ID box.
	VnetPeeringServicePrincipal types.String `tfsdk:"vnet_peering_service_principal"`

	// Command is the Azure cli command.
	Command types.String `tfsdk:"command"`
}
