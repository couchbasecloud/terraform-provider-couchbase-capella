package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func SampleBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":                       stringAttribute([]string{required, requiresReplace}),
			"organization_id":            stringAttribute([]string{required, requiresReplace}),
			"project_id":                 stringAttribute([]string{required, requiresReplace}),
			"cluster_id":                 stringAttribute([]string{required, requiresReplace}),
			"type":                       stringAttribute([]string{computed}),
			"storage_backend":            stringAttribute([]string{computed}),
			"memory_allocation_in_mb":    int64Attribute(computed),
			"bucket_conflict_resolution": stringAttribute([]string{computed}),
			"durability_level":           stringAttribute([]string{computed}),
			"replicas":                   int64Attribute(computed),
			"flush":                      boolAttribute(computed),
			"time_to_live_in_seconds":    int64Attribute(computed),
			"eviction_policy":            stringAttribute([]string{computed}),
			"stats": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"item_count":         int64Attribute(computed),
					"ops_per_second":     int64Attribute(computed),
					"disk_used_in_mib":   int64Attribute(computed),
					"memory_used_in_mib": int64Attribute(computed),
				},
			},
		},
	}
}
