package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var collectionsBuilder = capellaschema.NewSchemaBuilder("collections")

// CollectionsSchema returns the schema for the Collections data source.
func CollectionsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", collectionsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", collectionsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", collectionsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "bucket_id", collectionsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "scope_name", collectionsBuilder, requiredString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "collection_name", collectionsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "max_ttl", collectionsBuilder, computedInt64())

	capellaschema.AddAttr(attrs, "data", collectionsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source for retrieving collections within a scope in a bucket.",
		Attributes:          attrs,
	}
}
