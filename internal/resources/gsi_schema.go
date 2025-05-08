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
		MarkdownDescription: "Manages Query Indexes in Couchbase Capella",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the Capella organization."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the Capella project."),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the Capella cluster where the index will be created."),
			"bucket_name": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The name of the bucket where the index will be created. Specifies the bucket part of the key space."),
			"scope_name": WithDescription(stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			), "The name of the scope where the index will be created. Specifies the scope part of the key space. If unspecified, this will be the default scope."),
			"collection_name": WithDescription(stringDefaultAttribute(
				"_default", optional, computed, useStateForUnknown, requiresReplace,
			), "Specifies the collection part of the key space. If unspecified, this will be the default collection."),
			"index_name": WithDescription(stringAttribute([]string{optional, requiresReplace}),
				"The name of the index."),
			"is_primary": WithDescription(boolAttribute(optional),
				"Whether this is a primary index."),
			"index_keys": WithDescription(stringListAttribute(optional, requiresReplace),
				"List of document fields to index."),
			"where": WithDescription(stringAttribute([]string{optional, requiresReplace}),
				"WHERE clause for the index."),
			"status": WithDescription(stringAttribute([]string{computed, useStateForUnknown}),
				"The current status of the index. For example 'Created', 'Ready', etc."),
			"partition_by": WithDescription(stringListAttribute(optional, requiresReplace),
				"List of fields to partition the index by."),
			"with": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				Default:             objectdefault.StaticValue(defaultObject),
				MarkdownDescription: "Additional index configuration options.",
				Attributes: map[string]schema.Attribute{
					"defer_build": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "If true, the index will not be built immediately after creation.",
						PlanModifiers:       []planmodifier.Bool{custom_plan_modifiers.ImmutableBoolAttribute()},
					},
					"num_replica": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Number of index replicas to create.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"num_partition": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of partitions for the index.",
						Validators: []validator.Int64{
							int64validator.AtLeast(1),
						},
						PlanModifiers: []planmodifier.Int64{custom_plan_modifiers.ImmutableInt64Attribute()},
					},
				},
			},
			"build_indexes": WithDescription(stringSetAttribute(optional),
				"List of index names to build."),
		},
	}
}
