package api

// CreateVPCEndpointCommandRequest is a request to get the CLI to create an AWS VPC endpoint.
type CreateVPCEndpointCommandRequest struct {
	SubnetIDs *[]string `json:"subnetIDs,omitempty"`

	// VpcID The ID of your virtual network
	VpcID string `json:"vpcID"`
}

// CreateAzurePrivateEndpointCommandRequest is a request to get the CLI to create an Azure private endpoint.
type CreateAzurePrivateEndpointCommandRequest struct {
	// ResourceGroupName The name of your resource group
	ResourceGroupName string `json:"resourceGroupName"`

	// VirtualNetwork The virtual network and subnet name
	VirtualNetwork string `json:"virtualNetwork"`
}

type CreatePrivateEndpointCommandResponse struct {
	// Command The CLI command or script used to create private endpoint within your CSP.
	Command string `json:"command"`
}
