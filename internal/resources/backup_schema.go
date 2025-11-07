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
	capellaschema.AddAttr(attrs, "elapsed_time_in_seconds", backupBuilder, int64Attribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "restore_times", backupBuilder, numberAttribute(optional))

	attrs["backup_stats"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"size_in_mb": float64Attribute(computed),
			"items":      int64Attribute(computed),
			"mutations":  int64Attribute(computed),
			"tombstones": int64Attribute(computed),
			"gsi":        int64Attribute(computed),
			"fts":        int64Attribute(computed),
			"cbas":       int64Attribute(computed),
			"event":      int64Attribute(computed),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["schedule_info"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"backup_type": stringAttribute([]string{computed}),
			"backup_time": stringAttribute([]string{computed}),
			"increment":   int64Attribute(computed),
			"retention":   stringAttribute([]string{computed}),
		},
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["restore"] = schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"target_cluster_id":       stringAttribute([]string{required}),
			"source_cluster_id":       stringAttribute([]string{required}),
			"services":                stringListAttribute(required),
			"force_updates":           boolAttribute(optional),
			"auto_remove_collections": boolAttribute(optional),
			"filter_keys":             stringAttribute([]string{optional}),
			"filter_values":           stringAttribute([]string{optional}),
			"include_data":            stringAttribute([]string{optional}),
			"exclude_data":            stringAttribute([]string{optional}),
			"map_data":                stringAttribute([]string{optional}),
			"replace_ttl":             stringAttribute([]string{optional}),
			"replace_ttl_with":        stringAttribute([]string{optional}),
			"status":                  stringAttribute([]string{computed}),
		},
	}

	return schema.Schema{
		MarkdownDescription: "Manages backup resource associated with a bucket for an operational Capella cluster.",
		Attributes:          attrs,
	}
}
