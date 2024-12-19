package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
)

func GsiSchema() schema.Schema {
	defaultObject, _ := types.ObjectValue(
		map[string]attr.Type{
			"defer_build":   types.BoolType,
			"num_replica":   types.Int64Type,
			"num_partition": types.Int64Type,
		},
		map[string]attr.Value{
			"defer_build":   types.BoolNull(),
			"num_replica":   types.Int64Null(),
			"num_partition": types.Int64Null(),
		})

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
			"index_name":   stringAttribute([]string{optional, requiresReplace}),
			"is_primary":   boolAttribute(optional),
			"index_keys":   stringListAttribute(optional, requiresReplace),
			"where":        stringAttribute([]string{optional, requiresReplace}),
			"status":       stringAttribute([]string{computed, useStateForUnknown}),
			"partition_by": stringListAttribute(optional, requiresReplace),
			"with": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Default:  objectdefault.StaticValue(defaultObject),
				Attributes: map[string]schema.Attribute{
					"defer_build": schema.BoolAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{custom_plan_modifiers.ImmutableBoolAttribute()},
					},
					"num_replica": schema.Int64Attribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"num_partition": schema.Int64Attribute{
						Optional: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
						PlanModifiers: []planmodifier.Int64{custom_plan_modifiers.ImmutableInt64Attribute()},
					},
				},
			},
			"build_indexes": stringSetAttribute(optional),
		},
	}
}
