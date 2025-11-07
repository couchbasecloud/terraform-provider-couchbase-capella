package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	capellaschema.AddAttr(attrs, "storage_backend", bucketBuilder, stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}))
	capellaschema.AddAttr(attrs, "eviction_policy", bucketBuilder, stringAttribute([]string{computed, optional, requiresReplace, useStateForUnknown}))

	attrs["type"] = &schema.StringAttribute{
		Computed: true,
		Optional: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["memory_allocation_in_mb"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
	}
	attrs["vbuckets"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	}
	attrs["bucket_conflict_resolution"] = &schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString("seqno"),
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	attrs["durability_level"] = &schema.StringAttribute{
		Computed: true,
		Optional: true,
	}
	attrs["replicas"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
	}
	attrs["flush"] = &schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default:  booldefault.StaticBool(false),
	}
	attrs["time_to_live_in_seconds"] = &schema.Int64Attribute{
		Optional: true,
		Computed: true,
	}
	attrs["stats"] = schema.SingleNestedAttribute{
		Computed: true,
		Attributes: map[string]schema.Attribute{
			"item_count":         int64Attribute(computed),
			"ops_per_second":     int64Attribute(computed),
			"disk_used_in_mib":   int64Attribute(computed),
			"memory_used_in_mib": int64Attribute(computed),
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage the buckets for an operational cluster.",
		Attributes:          attrs,
	}
}
