package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var backupBuilder = capellaschema.NewSchemaBuilder("backup")

func BackupSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "organization_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "project_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "cluster_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "bucket_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "cycle_id", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "date", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "restore_before", backupBuilder, stringAttribute([]string{computed, optional, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "status", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "method", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "bucket_name", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "source", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "cloud_provider", backupBuilder, stringAttribute([]string{computed, useStateForUnknown}))

	backupStatsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(backupStatsAttrs, "size_in_mb", backupBuilder, float64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "items", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "mutations", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "tombstones", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "gsi", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "fts", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "cbas", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(backupStatsAttrs, "event", backupBuilder, int64Attribute(computed))

	capellaschema.AddAttr(attrs, "backup_stats", backupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: backupStatsAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	capellaschema.AddAttr(attrs, "elapsed_time_in_seconds", backupBuilder, int64Attribute(computed, useStateForUnknown))

	scheduleInfoAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scheduleInfoAttrs, "backup_type", backupBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(scheduleInfoAttrs, "backup_time", backupBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(scheduleInfoAttrs, "increment", backupBuilder, int64Attribute(computed))
	capellaschema.AddAttr(scheduleInfoAttrs, "retention", backupBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "schedule_info", backupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: scheduleInfoAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	restoreAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(restoreAttrs, "target_cluster_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(restoreAttrs, "source_cluster_id", backupBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(restoreAttrs, "services", backupBuilder, stringListAttribute(required))
	capellaschema.AddAttr(restoreAttrs, "force_updates", backupBuilder, boolAttribute(optional))
	capellaschema.AddAttr(restoreAttrs, "auto_remove_collections", backupBuilder, boolAttribute(optional))
	capellaschema.AddAttr(restoreAttrs, "filter_keys", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "filter_values", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "include_data", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "exclude_data", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "map_data", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "replace_ttl", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "replace_ttl_with", backupBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(restoreAttrs, "status", backupBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "restore", backupBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: restoreAttrs,
	})

	capellaschema.AddAttr(attrs, "restore_times", backupBuilder, numberAttribute(optional))

	return schema.Schema{
		MarkdownDescription: "Manages backup resource associated with a bucket for an operational Capella cluster.",
		Attributes:          attrs,
	}
}
