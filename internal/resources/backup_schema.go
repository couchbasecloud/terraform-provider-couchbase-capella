package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages backup resource associated with a bucket for a Couchbase Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The GUID4 ID of the backup resource.",
			},
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
				MarkdownDescription: "The GUID4 ID of the bucket.",
			},
			"cycle_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The GUID4 ID of the cycle this backup belongs to.",
			},
			"date": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The RFC3339 timestamp representing the time at which backup was created.",
			},
			"restore_before": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The RFC3339 timestamp representing the time at which backup will expire.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The status of the backup. values being pending, ready and failed",
			},
			"method": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The mechanism of the backup. It can be either incremental or full.",
			},
			"bucket_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The name of the bucket for which the backup belongs to.",
			},
			"source": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The way a backup job was initiated. It can be either manual or scheduled.",
			},
			"cloud_provider": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The cloud provider where the cluster is hosted.",
			},
			"backup_stats": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Represents various backup level data that couchbase provides.",
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
						MarkdownDescription: "The number of event entities saved during the backup.",
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"elapsed_time_in_seconds": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The amount of seconds that have elapsed between the creation and completion of the backup.",
			},
			"schedule_info": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Represents the schedule information of the backup.",
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
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"restore": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Represents restore configuration for the backup.",
				Attributes: map[string]schema.Attribute{
					"target_cluster_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The ID of the target cluster to restore to.",
					},
					"source_cluster_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The ID of the source cluster the restore is based on.",
					},
					"services": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Services is a list of services to be restored. It is used to specify the services that should be included in the restore operation.",
					},
					"force_updates": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Forces data in the Couchbase cluster to be overwritten even if the data in the cluster is newer.",
					},
					"auto_remove_collections": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Automatically delete scopes/collections which are known to be deleted in the backup.",
					},
					"filter_keys": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Only restore data where the key matches a particular regular expression.",
					},
					"filter_values": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Only restore data where the value matches a particular regular expression.",
					},
					"include_data": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Restores only the data specified here.",
					},
					"exclude_data": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Skips restoring the data specified here.",
					},
					"map_data": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Specified when you want to restore source data into a different location.",
					},
					"replace_ttl": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "",
					},
					"replace_ttl_with": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Sets a new expiration (time-to-live) value for the specified keys. Values being none, all or expired",
					},
					"status": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Updates the expiration for the keys.",
					},
				},
			},
			"restore_times": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "Number of times the backup to be restored",
			},
		},
	}
}
