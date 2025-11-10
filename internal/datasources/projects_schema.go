package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectsBuilder = capellaschema.NewSchemaBuilder("projects")

// ProjectsSchema returns the schema for the Projects data source.
func ProjectsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", projectsBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", projectsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", projectsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", projectsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", projectsBuilder, computedString())

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", projectsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "created_by", projectsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_at", projectsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_by", projectsBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "version", projectsBuilder, computedInt64())

	capellaschema.AddAttr(dataAttrs, "audit", projectsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "if_match", projectsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "etag", projectsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", projectsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Data source to retrieve project details in an organization.",
		Attributes:          attrs,
	}
}
