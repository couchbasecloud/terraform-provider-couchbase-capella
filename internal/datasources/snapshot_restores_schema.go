package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SnapshotRestoresSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all snapshot restore information for a cluster.",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the cluster.",
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
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The GUID4 ID of the organization.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"data": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of snapshot restores.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The GUID4 ID of the snapshot restore.",
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
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SingleNestedBlock{
				MarkdownDescription: "Filter criteria for the Cloud Snapshot Restores. Only filtering by Snapshot Restore Status is supported.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the attribute to filter.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("status"),
						},
					},
					"values": schema.SetAttribute{
						MarkdownDescription: "List of values to match against.",
						Optional:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
				},
			},
		},
	}
}
