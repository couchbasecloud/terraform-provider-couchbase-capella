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
	capellaschema.AddAttr(attrs, "scope_name", scopeBuilder, stringAttribute([]string{required, requiresReplace}))

	collectionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionAttrs, "max_ttl", scopeBuilder, int64Attribute(computed))
	capellaschema.AddAttr(collectionAttrs, "name", scopeBuilder, stringAttribute([]string{computed}))

	capellaschema.AddAttr(attrs, "collections", scopeBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a scope within a bucket.",
		Attributes:          attrs,
	}
}
