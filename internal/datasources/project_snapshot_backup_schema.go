package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func ProjectSnapshotBackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a project",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the tenant.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
				MarkdownDescription: "Sets the order of how you would like to sort the results and the key you would like to order by. Valid fields to sort the results are: severity, timestamp.",
			},
			"sort_direction": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The order in which the items will be sorted.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cluster_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the cluster.",
						},
						"cluster_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the cluster at the time of event emission.",
						},
						"creation_date_time": schema.StringAttribute{
							Computed:    true,
							Description: "The RFC3339 timestamp when the resource was created.",
						},
						"created_by": schema.StringAttribute{
							Computed:    true,
							Description: "The user who created the resource.",
						},
						"current_status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The current status of the snapshot backup.",
						},
						"cloud_provider": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The cloud provider of the cluster.",
						},
						"region": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The region of the cluster.",
						},
						"most_recent_snapshot": computedProjectSnapshot,
						"oldest_snapshot":      computedProjectSnapshot,
					},
				},
			},
			"cursor": computedCursorAttribute,
		},
	}
}
