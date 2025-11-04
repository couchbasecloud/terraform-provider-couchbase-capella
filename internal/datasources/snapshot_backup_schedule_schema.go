package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SnapshotBackupScheduleSchema() schema.Schema {
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
			"interval": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The interval at which the snapshot backup schedule runs.",
			},
			"retention": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The retention period for the snapshot backup schedule.",
			},
			"start_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The start time for the snapshot backup schedule.",
			},
			"copy_to_regions": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "The region to copy the snapshot backup to.",
				ElementType:         types.StringType,
			},
		},
	}
}
