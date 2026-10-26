package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceAdminUserBuilder = capellaschema.NewSchemaBuilder("appServiceAdminUser", "AppServiceAdminUser")

// AppServiceAdminUserSchema returns the schema for the App Service admin user resource.
func AppServiceAdminUserSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", appServiceAdminUserBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", appServiceAdminUserBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appServiceAdminUserBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appServiceAdminUserBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appServiceAdminUserBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "name", appServiceAdminUserBuilder, requiredNonEmptyStringAttribute())
	capellaschema.AddAttr(attrs, "password", appServiceAdminUserBuilder, stringAttribute([]string{required, sensitive}), "CreateAppServiceAdminUserRequest")
	capellaschema.AddAttr(attrs, "access_all_endpoints", appServiceAdminUserBuilder, boolAttribute(optional))
	capellaschema.AddAttr(attrs, "endpoints", appServiceAdminUserBuilder, stringSetAttribute(optional))
	capellaschema.AddAttr(attrs, "audit", appServiceAdminUserBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "Manages an admin user of an App Service in Couchbase Capella. Set either `access_all_endpoints = true` or a non-empty `endpoints` set.",
		Attributes:          attrs,
	}
}
