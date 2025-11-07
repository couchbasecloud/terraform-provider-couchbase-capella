package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var freeTierBucketBuilder = capellaschema.NewSchemaBuilder("freeTierBucket")

func FreeTierBucketSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", freeTierBucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
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

	attrs["memory_allocation_in_mb"] = &schema.Int64Attribute{
		Computed: true,
		Optional: true,
		Default:  int64default.StaticInt64(100),
	}

	attrs["stats"] = schema.SingleNestedAttribute{
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
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets of your free tier operational cluster.",
		Attributes:          attrs,
	}
}
