package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func SnapshotBackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a bucket for an operational cluster.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the tenant.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the snapshotbackups associated with a cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_service": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the app service.",
						},
						"cluster_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the cluster.",
						},
						"created_at": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The RFC3339 timestamp representing the time at which backup was created.",
						},
						"expiration": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The RFC3339 timestamp representing the time at which the backup will expire.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the snapshot backup resource.",
						},
						"progress": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The progress of the snapshot backup.",
							Attributes: map[string]schema.Attribute{
								"status": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The status of the snapshot backup.",
								},
								"time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The time of the snapshot backup.",
								},
							},
						},
						"retention": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The retention time in hours.",
						},
						"cmek": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The CMEK configuration for the snapshot backup.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The GUID4 ID of the CMEK configuration.",
									},
									"provider_id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The GUID4 ID of the provider.",
									},
								},
							},
						},
						"server": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"version": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The version of the server.",
								},
							},
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the snapshot backup.",
						},
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of the snapshot backup.",
						},
					},
				},
			},
		},
	}
}
