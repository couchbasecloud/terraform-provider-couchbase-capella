package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			"organization_id": stringAttribute(required),
			"project_id":      stringAttribute(required),
			"cluster_id":      stringAttribute(required),
			"bucket_id":       stringAttribute(required),
			"cycle_id":        stringAttribute(computed),
			"date":            stringAttribute(computed),
			"restore_before":  stringAttribute(optional, computed),
			"status":          stringAttribute(computed),
			"method":          stringAttribute(computed),
			"bucket_name":     stringAttribute(computed),
			"source":          stringAttribute(computed),
			"cloud_provider":  stringAttribute(computed),
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
			},
			"elapsed_time_in_seconds": int64Attribute(computed),
			//"schedule_info": schema.SingleNestedAttribute{
			//	Computed: true,
			//	Attributes: map[string]schema.Attribute{
			//		"backup_type": stringAttribute(computed),
			//		"backup_time": stringAttribute(computed),
			//		"increment":   int64Attribute(computed),
			//		"retention":   stringAttribute(computed),
			//	},
			//},
			//"type": stringAttribute(optional),
			//"weekly_schedule": schema.SingleNestedAttribute{
			//	Optional: true,
			//	Attributes: map[string]schema.Attribute{
			//		"day_of_week":              stringAttribute(optional),
			//		"start_at":                 int64Attribute(optional),
			//		"incremental_every":        int64Attribute(optional),
			//		"retention_time":           stringAttribute(optional),
			//		"cost_optimized_retention": boolAttribute(optional),
			//	},
			//},
		},
	}
}
