package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseRolesBuilder = capellaschema.NewSchemaBuilder("databaseRoles")

// DatabaseRolesSchema returns the schema for the DatabaseRoles data source.
func DatabaseRolesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", databaseRolesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", databaseRolesBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", databaseRolesBuilder, requiredString())

	// Build data item attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", databaseRolesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", databaseRolesBuilder, computedString())

	capellaschema.AddAttr(dataAttrs, "audit", databaseRolesBuilder, computedAudit())

	capellaschema.AddAttr(dataAttrs, "access", databaseRolesBuilder, computedAccessAttribute(databaseRolesBuilder))

	capellaschema.AddAttr(attrs, "data", databaseRolesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to list database user roles for a cluster. Database roles define reusable sets of privileges that can be assigned to database credentials.",
		Attributes:          attrs,
	}
}
