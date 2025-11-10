package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources/custom_plan_modifiers"
)

var gsiBuilder = capellaschema.NewSchemaBuilder("gsi", "indexDDLRequest")

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

	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", gsiBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", gsiBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", gsiBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "bucket_name", gsiBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "scope_name", gsiBuilder, stringDefaultAttribute("_default", optional, computed, useStateForUnknown, requiresReplace))
	capellaschema.AddAttr(attrs, "collection_name", gsiBuilder, stringDefaultAttribute("_default", optional, computed, useStateForUnknown, requiresReplace))
	capellaschema.AddAttr(attrs, "index_name", gsiBuilder, stringAttribute([]string{optional, requiresReplace}))
	capellaschema.AddAttr(attrs, "is_primary", gsiBuilder, boolAttribute(optional))
	capellaschema.AddAttr(attrs, "index_keys", gsiBuilder, stringListAttribute(optional, requiresReplace))
	capellaschema.AddAttr(attrs, "where", gsiBuilder, stringAttribute([]string{optional, requiresReplace}))
	capellaschema.AddAttr(attrs, "status", gsiBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "partition_by", gsiBuilder, stringListAttribute(optional, requiresReplace))
	capellaschema.AddAttr(attrs, "build_indexes", gsiBuilder, stringSetAttribute(optional))

	withAttrs := make(map[string]schema.Attribute)
	withAttrs["defer_build"] = &schema.BoolAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.Bool{custom_plan_modifiers.ImmutableBoolAttribute()},
	}
	capellaschema.AddAttr(withAttrs, "num_replica", gsiBuilder, int64Attribute(optional, computed, useStateForUnknown))
	withAttrs["num_partition"] = &schema.Int64Attribute{
		Optional: true,
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
		PlanModifiers: []planmodifier.Int64{custom_plan_modifiers.ImmutableInt64Attribute()},
	}

	capellaschema.AddAttr(attrs, "with", gsiBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Computed:   true,
		Default:    objectdefault.StaticValue(defaultObject),
		Attributes: withAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage Query Indexes in Couchbase Capella.",
		Attributes:          attrs,
	}
}
