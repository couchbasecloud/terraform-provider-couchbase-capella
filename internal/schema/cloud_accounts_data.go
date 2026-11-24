package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// CloudAccountsData is the datasource model for Capella's cloud accounts in an organization.
type CloudAccountsData struct {
	// OrganizationId is the ID of the organization.
	OrganizationId types.String `tfsdk:"organization_id"`

	// AwsCapellaAccount is Capella's AWS account ID.
	AwsCapellaAccount types.String `tfsdk:"aws_capella_account"`

	// GcpCapellaProject is Capella's GCP project ID.
	GcpCapellaProject types.String `tfsdk:"gcp_capella_project"`

	// AzureCapellaSubscription is the name of Capella's Azure subscription.
	AzureCapellaSubscription types.String `tfsdk:"azure_capella_subscription"`

	// AzureCapellaSubscriptionId is the ID of Capella's Azure subscription. Null when the
	// organization has no Azure account.
	AzureCapellaSubscriptionId types.String `tfsdk:"azure_capella_subscription_id"`
}
