package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func SnapshotRestoreSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all snapshot restore information for a cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the snapshot restore.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The RFC3339 timestamp representing the time at which snapshot restore was created.",
			},
			"restore_to": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The GUID4 ID of the cluster to restore to.",
			},
			"snapshot": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The RFC3339 timestamp representing the time at which the snapshot was taken.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the snapshot restore.",
			},
		},
	}
}
