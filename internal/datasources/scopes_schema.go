package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var scopesBuilder = capellaschema.NewSchemaBuilder("scopes")

// ScopesSchema returns the schema for the Scopes data source.
func ScopesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", scopesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", scopesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", scopesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "bucket_id", scopesBuilder, requiredString())

	collectionAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(collectionAttrs, "name", scopesBuilder, computedString())
	capellaschema.AddAttr(collectionAttrs, "max_ttl", scopesBuilder, computedInt64())

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "scope_name", scopesBuilder, computedString())
	capellaschema.AddAttr(scopeAttrs, "collections", scopesBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: collectionAttrs,
		},
	})

	capellaschema.AddAttr(attrs, "scopes", scopesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The scopes data source retrieves the scopes in a bucket.",
		Attributes:          attrs,
	}
}
