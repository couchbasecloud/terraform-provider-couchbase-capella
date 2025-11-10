package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierBucketBuilder = capellaschema.NewSchemaBuilder("freeTierBucket")

func FreeTierBucketSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", freeTierBucketBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "name", freeTierBucketBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "organization_id", freeTierBucketBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", freeTierBucketBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", freeTierBucketBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "type", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "storage_backend", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "bucket_conflict_resolution", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "durability_level", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "replicas", freeTierBucketBuilder, int64Attribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "flush", freeTierBucketBuilder, boolAttribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "time_to_live_in_seconds", freeTierBucketBuilder, int64Attribute(computed, useStateForUnknown))
	capellaschema.AddAttr(attrs, "eviction_policy", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "memory_allocation_in_mb", freeTierBucketBuilder, int64DefaultAttribute(100, computed, optional))

	statsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statsAttrs, "item_count", freeTierBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "ops_per_second", freeTierBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "disk_used_in_mib", freeTierBucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "memory_used_in_mib", freeTierBucketBuilder, int64Attribute(computed))

	capellaschema.AddAttr(attrs, "stats", freeTierBucketBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statsAttrs,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.UseStateForUnknown(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets of your free tier operational cluster.",
		Attributes:          attrs,
	}
}
