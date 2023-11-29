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
			"name":                       stringAttribute(required, requiresReplace),
			"organization_id":            stringAttribute(required),
			"project_id":                 stringAttribute(required),
			"cluster_id":                 stringAttribute(required),
			"type":                       stringDefaultAttribute("couchbase", optional, computed, requiresReplace),
			"storage_backend":            stringAttribute(optional, computed, requiresReplace),
			"memory_allocation_in_mb":    int64DefaultAttribute(100, optional, computed),
			"bucket_conflict_resolution": stringDefaultAttribute("seqno", optional, computed, requiresReplace),
			"durability_level":           stringDefaultAttribute("none", optional, computed),
			"replicas":                   int64DefaultAttribute(1, optional, computed),
			"flush":                      boolDefaultAttribute(false, optional, computed),
			"time_to_live_in_seconds":    int64DefaultAttribute(0, optional, computed),
			"eviction_policy":            stringAttribute(optional, computed, requiresReplace),
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
