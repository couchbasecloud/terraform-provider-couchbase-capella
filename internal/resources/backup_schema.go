package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func BackupSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages backup resource associated with a bucket for a Couchbase Capella cluster.",
		Attributes: map[string]schema.Attribute{
			"id":              WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the backup resource."),
			"organization_id": WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the cluster."),
			"bucket_id":       WithDescription(stringAttribute([]string{required}), "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name."),
			"cycle_id":        WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The GUID4 ID of the cycle this backup belongs to."),
			"date":            WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The RFC3339 timestamp representing the time at which backup was created."),
			"restore_before":  WithDescription(stringAttribute([]string{computed, optional, useStateForUnknown}), "The RFC3339 timestamp representing the time at which backup will expire."),
			"status":          WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The status of the backup. values being pending, ready and failed"),
			"method":          WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The mechanism of the backup. It can be either incremental or full."),
			"bucket_name":     WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The name of the bucket for which the backup belongs to."),
			"source":          WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The source of the backup. It can be either cluster or bucket."),
			"cloud_provider":  WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The cloud provider where the cluster is hosted."),
			"backup_stats": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Represents various backup level data that couchbase provides.",
				Attributes: map[string]schema.Attribute{
					"size_in_mb": WithDescription(float64Attribute(computed), "Size in MB is the total size of the backup in megabytes. It represents the amount of data that was backed up during the backup operation."),
					"items":      WithDescription(int64Attribute(computed), "The number of items saved during the backup."),
					"mutations":  WithDescription(int64Attribute(computed), "The number of mutations saved during the backup."),
					"tombstones": WithDescription(int64Attribute(computed), "The number of tombstones saved during the backup."),
					"gsi":        WithDescription(int64Attribute(computed), "The number of global secondary indexes saved during the backup."),
					"fts":        WithDescription(int64Attribute(computed), "The number of full text search entities saved during the backup."),
					"cbas":       WithDescription(int64Attribute(computed), "The number of analytics entities saved during the backup."),
					"event":      WithDescription(int64Attribute(computed), "The number of event entities saved during the backup."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"elapsed_time_in_seconds": WithDescription(int64Attribute(computed, useStateForUnknown), "The amount of seconds that have elapsed between the creation and completion of the backup."),
			"schedule_info": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Represents the schedule information of the backup.",
				Attributes: map[string]schema.Attribute{
					"backup_type": WithDescription(stringAttribute([]string{computed}), "Represents whether the backup is a Weekly or Daily backup."),
					"backup_time": WithDescription(stringAttribute([]string{computed}), "Represents the time at which the backup is scheduled to be taken."),
					"increment":   WithDescription(int64Attribute(computed), "Represents interval in hours for incremental backup."),
					"retention":   WithDescription(stringAttribute([]string{computed}), "Represents retention time in days."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"restore": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Represents restore configuration for the backup.",
				Attributes: map[string]schema.Attribute{
					"target_cluster_id":       WithDescription(stringAttribute([]string{required}), "The ID of the target cluster to restore to."),
					"source_cluster_id":       WithDescription(stringAttribute([]string{required}), "The ID of the source cluster the restore is based on."),
					"services":                WithDescription(stringListAttribute(required), "Services is a list of services to be restored. It is used to specify the services that should be included in the restore operation."),
					"force_updates":           WithDescription(boolAttribute(optional), "Forces data in the Couchbase cluster to be overwritten even if the data in the cluster is newer."),
					"auto_remove_collections": WithDescription(boolAttribute(optional), "Automatically remove collections that are not present in the backup."),
					"filter_keys":             WithDescription(stringAttribute([]string{optional}), "Only restore data where the key matches a particular regular expression."),
					"filter_values":           WithDescription(stringAttribute([]string{optional}), "Only restore data where the value matches a particular regular expression."),
					"include_data":            WithDescription(stringAttribute([]string{optional}), "Restores only the data specified here."),
					"exclude_data":            WithDescription(stringAttribute([]string{optional}), "Skips restoring the data specified here."),
					"map_data":                WithDescription(stringAttribute([]string{optional}), "Specified when you want to restore source data into a different location"),
					"replace_ttl":             WithDescription(stringAttribute([]string{optional}), "Sets a new expiration (time-to-live) value for the specified keys. Values being 'none' 'all' 'expired'"),
					"replace_ttl_with":        WithDescription(stringAttribute([]string{optional}), "Updates the expiration for the keys."),
					"status":                  WithDescription(stringAttribute([]string{computed}), "The status of the restore."),
				},
			},
			"restore_times": WithDescription(numberAttribute(optional), "Number of times the backup to be restored."),
		},
	}
}
