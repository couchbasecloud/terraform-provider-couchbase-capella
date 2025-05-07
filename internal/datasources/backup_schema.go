package datasources

import "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

func BackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages backup datasource, retrieves backups associated with a bucket for a Couchbase Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the project.",
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
			},
			"bucket_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name.",
			},
			"data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the backup resource.",
						},
						"organization_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the organization.",
						},
						"project_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the project.",
						},
						"cluster_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the cluster.",
						},
						"bucket_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name.",
						},
						"cycle_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the cycle this backup belongs to.",
						},
						"date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The RFC3339 timestamp representing the time at which backup was created.",
						},
						"restore_before": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The RFC3339 timestamp representing the time at which backup will expire",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status of the backup. values being pending, ready and failed",
						},
						"method": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The mechanism of the backup. It can be either incremental or full.",
						},
						"bucket_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the bucket for which the backup belongs to.",
						},
						"source": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The way a backup job was initiated. It can be either manual or scheduled.",
						},
						"cloud_provider": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The cloud provider where the cluster is hosted.",
						},
						"backup_stats": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"size_in_mb": schema.Float64Attribute{
									Computed:            true,
									MarkdownDescription: "Size in MB is the total size of the backup in megabytes. It represents the amount of data that was backed up during the backup operation.",
								},
								"items": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of items saved during the backup.",
								},
								"mutations": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of mutations saved during the backup.",
								},
								"tombstones": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of tombstones saved during the backup.",
								},
								"gsi": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of global secondary indexes saved during the backup.",
								},
								"fts": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of full text search entities saved during the backup.",
								},
								"cbas": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of analytics entities saved during the backup.",
								},
								"event": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The number of eventing entities saved during the backup.",
								},
							},
						},
						"elapsed_time_in_seconds": computedInt64Attribute,
						"schedule_info": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"backup_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Represents whether the backup is a Weekly or Daily backup.",
								},
								"backup_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The timestamp indicating the backup created time.",
								},
								"increment": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "Represents interval in hours for incremental backup.",
								},
								"retention": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Represents retention time in days.",
								},
							},
						},
					},
				},
			},
		},
	}
}
