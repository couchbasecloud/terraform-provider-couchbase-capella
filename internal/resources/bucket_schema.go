package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var bucketBuilder = capellaschema.NewSchemaBuilder("bucket")

func BucketSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", bucketBuilder, stringAttribute([]string{computed, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "name", bucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "organization_id", bucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", bucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", bucketBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "type", bucketBuilder, stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "storage_backend", bucketBuilder, stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "memory_allocation_in_mb", bucketBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "vbuckets", bucketBuilder, &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	})
	capellaschema.AddAttr(attrs, "bucket_conflict_resolution", bucketBuilder, stringDefaultAttribute("seqno", computed, optional, requiresReplace, useStateForUnknown))
	capellaschema.AddAttr(attrs, "durability_level", bucketBuilder, stringAttribute([]string{computed, optional}))
	capellaschema.AddAttr(attrs, "replicas", bucketBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "flush", bucketBuilder, boolDefaultAttribute(false, optional, computed))
	capellaschema.AddAttr(attrs, "time_to_live_in_seconds", bucketBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "eviction_policy", bucketBuilder, stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}))

	statsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(statsAttrs, "item_count", bucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "ops_per_second", bucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "disk_used_in_mib", bucketBuilder, int64Attribute(computed))
	capellaschema.AddAttr(statsAttrs, "memory_used_in_mib", bucketBuilder, int64Attribute(computed))

	capellaschema.AddAttr(attrs, "stats", bucketBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: statsAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets for an operational cluster.",
		Attributes:          attrs,
	}
}
