package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BackupSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": stringAttribute([]string{required}),
			"project_id":      stringAttribute([]string{required}),
			"cluster_id":      stringAttribute([]string{required}),
			"bucket_id":       stringAttribute([]string{required}),
			"cycle_id":        stringAttribute([]string{computed, useStateForUnknown}),
			"date":            stringAttribute([]string{computed, useStateForUnknown}),
			"restore_before":  stringAttribute([]string{optional, computed, useStateForUnknown}),
			"status":          stringAttribute([]string{computed, useStateForUnknown}),
			"method":          stringAttribute([]string{computed, useStateForUnknown}),
			"bucket_name":     stringAttribute([]string{computed, useStateForUnknown}),
			"source":          stringAttribute([]string{computed, useStateForUnknown}),
			"cloud_provider":  stringAttribute([]string{computed, useStateForUnknown}),
			"backup_stats": schema.SingleNestedAttribute{
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
			},
			"elapsed_time_in_seconds": int64Attribute(computed, useStateForUnknown),
			"schedule_info": schema.SingleNestedAttribute{
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
			},
			"restore": schema.SingleNestedAttribute{
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
			},
			"restore_times": numberAttribute(optional),
		},
	}
}
