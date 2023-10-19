package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func BucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id":            stringAttribute(required),
			"project_id":                 stringAttribute(required),
			"cluster_id":                 stringAttribute(required),
			"name":                       stringAttribute(required),
			"type":                       stringAttribute(optional),
			"storage_backend":            stringAttribute(optional),
			"memory_allocation_in_mb":    int64Attribute(optional),
			"bucket_conflict_resolution": stringAttribute(optional),
			"durability_level":           stringAttribute(optional),
			"replicas":                   int64Attribute(optional),
			"flush":                      boolAttribute(optional),
			"time_to_live_in_seconds":    int64Attribute(optional),
			"eviction_policy":            stringAttribute(computed, optional),
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
