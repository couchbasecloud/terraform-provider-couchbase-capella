package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAdminUserBuilder = capellaschema.NewSchemaBuilder("appServiceAdminUser")

// AppServiceAdminUserSchema returns the schema for the App Service admin user data source.
func AppServiceAdminUserSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAdminUserBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", appServiceAdminUserBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAdminUserBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAdminUserBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "id", appServiceAdminUserBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "name", appServiceAdminUserBuilder, computedString())
	capellaschema.AddAttr(attrs, "access_all_endpoints", appServiceAdminUserBuilder, computedBool())
	capellaschema.AddAttr(attrs, "endpoints", appServiceAdminUserBuilder, computedStringSet())
	capellaschema.AddAttr(attrs, "audit", appServiceAdminUserBuilder, computedAudit())

	return schema.Schema{
		MarkdownDescription: "Retrieves an admin user of a Capella App Service.",
		Attributes:          attrs,
	}
}
