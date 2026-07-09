package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseCredentialsBuilder = capellaschema.NewSchemaBuilder("databaseCredentials")

// DatabaseCredentialsSchema returns the schema for the DatabaseCredentials data source.
func DatabaseCredentialsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", databaseCredentialsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", databaseCredentialsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", databaseCredentialsBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", databaseCredentialsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", databaseCredentialsBuilder, computedString())

	capellaschema.AddAttr(dataAttrs, "audit", databaseCredentialsBuilder, computedAudit())

	capellaschema.AddAttr(dataAttrs, "access", databaseCredentialsBuilder, computedAccessAttribute(databaseCredentialsBuilder))

	capellaschema.AddAttr(dataAttrs, "user_roles", databaseCredentialsBuilder, computedStringSet())

	capellaschema.AddAttr(attrs, "data", databaseCredentialsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve database credentials for a cluster. Database credentials provide programmatic and application-level access to data on a database.",
		Attributes:          attrs,
	}
}
