package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name":                       stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"organization_id":            stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"project_id":                 stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"cluster_id":                 stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"type":                       stringAttribute([]string{computed, useStateForUnknown}),
			"storage_backend":            stringAttribute([]string{computed, useStateForUnknown}),
			"memory_allocation_in_mb":    int64DefaultAttribute(100, optional, computed),
			"bucket_conflict_resolution": stringAttribute([]string{computed, useStateForUnknown}),
			"durability_level":           stringAttribute([]string{computed, useStateForUnknown}),
			"replicas":                   int64Attribute(computed, useStateForUnknown),
			"flush":                      boolAttribute(computed, useStateForUnknown),
			"time_to_live_in_seconds":    int64Attribute(computed, useStateForUnknown),
			"eviction_policy":            stringAttribute([]string{computed, useStateForUnknown}),
			"stats": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"item_count":         int64Attribute(computed),
					"ops_per_second":     int64Attribute(computed),
					"disk_used_in_mib":   int64Attribute(computed),
					"memory_used_in_mib": int64Attribute(computed),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
