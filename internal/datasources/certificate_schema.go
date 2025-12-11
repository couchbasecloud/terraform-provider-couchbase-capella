package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var certificateBuilder = capellaschema.NewSchemaBuilder("certificate")

// CertificateSchema returns the schema for the Certificate data source.
func CertificateSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", certificateBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", certificateBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", certificateBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "certificate", certificateBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", certificateBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The Capella certificate data source allows you to retrieve the certificate for a cluster.",
		Attributes:          attrs,
	}
}
