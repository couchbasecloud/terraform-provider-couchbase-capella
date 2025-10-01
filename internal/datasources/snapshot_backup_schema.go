package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func SnapshotBackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a bucket for an operational cluster.",
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
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Lists the snapshotbackups associated with a cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
						"cmek": schema.SetNestedAttribute{
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
						"cross_region_copies": schema.SetNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The cross region copies of the snapshot backup.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"region_code": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The region the snapshot backup has been copied to.",
									},
									"status": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The status of the cross region copy.",
									},
									"time": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The RFC3339 timestamp representing the time at which the status was last updated.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
