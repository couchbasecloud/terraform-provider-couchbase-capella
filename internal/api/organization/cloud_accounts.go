package organization

// GetCloudAccountsResponse is Capella's cloud accounts for an organization.
//
// The generated apigen.CloudAccounts doesn't have azure-capella-subscription-id until the
// published spec does, so the data source decodes into this type instead.
type GetCloudAccountsResponse struct {
	AwsCapellaAccount        string `json:"aws-capella-account"`
	GcpCapellaProject        string `json:"gcp-capella-project"`
	AzureCapellaSubscription string `json:"azure-capella-subscription"`

	// AzureCapellaSubscriptionId is absent when the organization has no Azure account.
	AzureCapellaSubscriptionId *string `json:"azure-capella-subscription-id,omitempty"`
}
