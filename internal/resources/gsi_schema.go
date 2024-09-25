package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
)

func GsiSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"bucket_name":     stringAttribute([]string{required, requiresReplace}),
			"scope_name": stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			),
			"collection_name": stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			),
			"index_name":   stringAttribute([]string{optional, useStateForUnknown, requiresReplace}),
			"is_primary":   boolAttribute(optional),
			"index_keys":   stringListAttribute(optional, requiresReplace),
			"where":        stringAttribute([]string{optional, useStateForUnknown, requiresReplace}),
			"partition_by": stringListAttribute(optional, requiresReplace),
			"with": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"defer_build": schema.BoolAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{custom_plan_modifiers.ImmutableBoolAttribute()},
					},
					"num_replica": int64Attribute(optional),
					"num_partition": schema.Int64Attribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Int64{custom_plan_modifiers.ImmutableInt64Attribute()},
					},
				},
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"defer_build":   types.BoolType,
							"num_replica":   types.Int64Type,
							"num_partition": types.Int64Type,
						},
						map[string]attr.Value{
							"defer_build":   types.BoolValue(false),
							"num_replica":   types.Int64Value(0),
							"num_partition": types.Int64Value(8),
						},
					),
				),
			},
			"build_indexes": stringListAttribute(optional),
		},
	}
}
