package schema

import "github.com/hashicorp/terraform-plugin-framework/types"

// OutboundPrivateEndpoints is the data source model for every outbound private endpoint on a cluster.
// Each item carries the resource's own fields, so the two read the same.
type OutboundPrivateEndpoints struct {
	OrganizationId types.String              `tfsdk:"organization_id"`
	ProjectId      types.String              `tfsdk:"project_id"`
	ClusterId      types.String              `tfsdk:"cluster_id"`
	Data           []OutboundPrivateEndpoint `tfsdk:"data"`
}
