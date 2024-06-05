package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// PrivateEndpoint represents a private endpoint resource.
type PrivateEndpoint struct {
	// EndpointId is the id of the bucket for which the collection needs to be created.
	EndpointId types.String `tfsdk:"endpoint_id"`

	// ClusterId is the ID of the cluster for which the collection needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// PrivateEndpoints defines a structure used by the LIST endpoint for private endpoints.
type PrivateEndpoints struct {
	// Data is an array of private endpoints.
	Data []PrivateEndpointData `tfsdk:"data"`

	// ClusterId is the ID of the cluster for which the collection needs to be created.
	ClusterId types.String `tfsdk:"cluster_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`
}

// PrivateEndpointData defines a single private endpoint.
type PrivateEndpointData struct {
	// Id is the endpoint id.
	Id types.String `tfsdk:"id"`
	// Status is the endpoint status.  Possible values are failed, linked, pending, pendingAcceptance, rejected and unrecognized.
	Status types.String `tfsdk:"status"`
}
