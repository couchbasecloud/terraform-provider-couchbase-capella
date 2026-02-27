package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ClusterStats defines the cluster stats data source schema.
type ClusterStats struct {
	OrganizationId  types.String `tfsdk:"organization_id"`
	ProjectId       types.String `tfsdk:"project_id"`
	ClusterId       types.String `tfsdk:"cluster_id"`
	FreeMemoryInMb  types.Int64  `tfsdk:"free_memory_in_mb"`
	MaxReplicas     types.Int64  `tfsdk:"max_replicas"`
	TotalMemoryInMb types.Int64  `tfsdk:"total_memory_in_mb"`
}
