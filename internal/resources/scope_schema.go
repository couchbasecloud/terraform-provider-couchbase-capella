package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var scopeBuilder = capellaschema.NewSchemaBuilder("scope")

func ScopeSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", scopeBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", scopeBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", scopeBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "bucket_id", scopeBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["scope_name"] = WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the scope.")
	attrs["collections"] = schema.SetNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The list of collections within this scope.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"max_ttl": schema.Int64Attribute{
					Computed:            true,
					MarkdownDescription: "The maximum Time To Live (TTL) for documents in the collection.",
				},
				"name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The name of the collection.",
				},
			},
		},
	}

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a scope within a bucket.",
		Attributes:          attrs,
	}
}
