package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databasePrivilegesBuilder = capellaschema.NewSchemaBuilder("databasePrivileges", "GetCapellaPrivilegeResponse")

// DatabasePrivilegesSchema returns the schema for the DatabasePrivileges data source.
func DatabasePrivilegesSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", databasePrivilegesBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", databasePrivilegesBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", databasePrivilegesBuilder, requiredStringWithValidator())

	// Build data item attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "name", databasePrivilegesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "group", databasePrivilegesBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "resources", databasePrivilegesBuilder, computedResourcesAttribute(databasePrivilegesBuilder))

	capellaschema.AddAttr(attrs, "data", databasePrivilegesBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Lists the available Capella privileges that can be assigned to database credentials and database user roles for a given cluster.",
		Attributes:          attrs,
	}
}
