package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// OutboundPrivateEndpointCapacityData is the datasource model for a cluster's outbound private endpoint capacity.
type OutboundPrivateEndpointCapacityData struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`

	// ClusterEndpointLimit is the cluster's endpoints against its own limit.
	ClusterEndpointLimit *OutboundPrivateEndpointCapacityUsage `tfsdk:"cluster_endpoint_limit"`
}

// OutboundPrivateEndpointCapacityUsage is the cluster's endpoints against its own limit.
type OutboundPrivateEndpointCapacityUsage struct {
	Used  types.Int64 `tfsdk:"used"`
	Limit types.Int64 `tfsdk:"limit"`
}
