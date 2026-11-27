package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAuditLogBuilder = capellaschema.NewSchemaBuilder("appServiceAuditLog", "CreateAppServiceAuditLogRequest")

// AppServiceAuditLogSchema returns the schema for the App Service audit log resource.
func AppServiceAuditLogSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServiceAuditLogBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceAuditLogBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAuditLogBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAuditLogBuilder, requiredUUIDStringAttribute())

	return schema.Schema{
		MarkdownDescription: "Enables audit logging on an App Service in Couchbase Capella. Destroying this resource disables audit logging on the App Service.",
		Attributes:          attrs,
	}
}
