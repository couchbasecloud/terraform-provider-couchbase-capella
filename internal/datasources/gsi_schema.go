package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var gsiBuilder = capellaschema.NewSchemaBuilder("gsi")

// GsiSchema returns the schema for the Gsi data source.
func GsiSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", gsiBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", gsiBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", gsiBuilder, requiredString())
	capellaschema.AddAttr(attrs, "bucket_name", gsiBuilder, requiredString())
	capellaschema.AddAttr(attrs, "scope_name", gsiBuilder, optionalString())
	capellaschema.AddAttr(attrs, "collection_name", gsiBuilder, optionalString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "index_name", gsiBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "bucket_name", gsiBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "scope_name", gsiBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "collection_name", gsiBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "index_keys", gsiBuilder, computedStringSet())
	capellaschema.AddAttr(dataAttrs, "status", gsiBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", gsiBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve GSI indexes for a cluster.",
		Attributes:          attrs,
	}
}
