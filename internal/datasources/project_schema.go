package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectBuilder = capellaschema.NewSchemaBuilder("project")

// ProjectSchema returns the schema for the Project data source.
func ProjectSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", projectBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "id", projectBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "name", projectBuilder, computedString())
	capellaschema.AddAttr(attrs, "description", projectBuilder, computedString())
	capellaschema.AddAttr(attrs, "etag", projectBuilder, computedString())
	capellaschema.AddAttr(attrs, "if_match", projectBuilder, optionalString())
	capellaschema.AddAttr(attrs, "audit", projectBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve information about a Capella project.",
		Attributes:          attrs,
	}
}
