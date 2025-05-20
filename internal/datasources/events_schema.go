package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// EventsSchema returns the schema for the Events data source.
func EventsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: " Data source to retrieve all events in an organization. Events represent a trail of actions that users performs within Capella at an organization level.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The GUID4 ID of projects to filter on. By default, events corresponding to all projects are returned",
			},
			"cluster_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of clusterIds to filter on. By default, events corresponding to all clusters are returned.",
			},
			"user_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by user UUID. Default is to return events corresponding to all users.",
			},
			"severity_levels": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by severity levels. Default is to return events corresponding to all supported severity levels.",
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Filter by tags. Default is to return events corresponding to all supported tag. The tags are: availability, billing, maintenance, performance, security, alert.",
			},
			"from": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Start date in RFC3339 format. If not provided, events starting from last 24 hours are returned.",
			},
			"to": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "End datetime in the last 24 hours, RFC3339 format. Defaults to Now.",
			},
			"page": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Sets the page you would like to view.",
			},
			"per_page": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Sets the number of results you would like to have on each page.",
			},
			"sort_by": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Sets the order of how you would like to sort the results and the key you would like to order by. The valid fields to sort the results are: severity, timestamp.",
			},
			"sort_direction": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The order in which the items will be sorted. The valid values are: asc, desc.",
			},
			"data":   computedEventAttributes,
			"cursor": computedCursorAttribute,
		},
	}
}
