package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProjectEvents represents a structure for querying project-specific events.
type ProjectEvents struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterIds     types.Set    `tfsdk:"cluster_ids"`
	UserIds        types.Set    `tfsdk:"user_ids"`
	SeverityLevels types.Set    `tfsdk:"severity_levels"`
	Tags           types.Set    `tfsdk:"tags"`
	From           types.String `tfsdk:"from"`
	To             types.String `tfsdk:"to"`
	Page           types.Int64  `tfsdk:"page"`
	PerPage        types.Int64  `tfsdk:"per_page"`
	SortBy         types.String `tfsdk:"sort_by"`
	SortDirection  types.String `tfsdk:"sort_direction"`
	Data           []EventItem  `tfsdk:"data"`
	Cursor         *Cursor      `tfsdk:"cursor"`
}
