package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var userBuilder = capellaschema.NewSchemaBuilder("user")

func UserSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "name", userBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "status", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "inactive", userBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "email", userBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "organization_id", userBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "organization_roles", userBuilder, stringListAttribute(required))
	capellaschema.AddAttr(attrs, "last_login", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "region", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "time_zone", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "enable_notifications", userBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "expires_at", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", userBuilder, computedAuditAttribute())

	resourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourceAttrs, "type", userBuilder, stringDefaultAttribute("project", optional, computed), "Resource")
	capellaschema.AddAttr(resourceAttrs, "id", userBuilder, stringAttribute([]string{required}), "Resource")
	capellaschema.AddAttr(resourceAttrs, "roles", userBuilder, stringSetAttribute(required), "Resource")

	capellaschema.AddAttr(attrs, "resources", userBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: resourceAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "This User resource allows you to manage users in your Capella organization. You can create, update, and delete users, as well as manage their roles and permissions within the organization.",
		Attributes:          attrs,
	}
}
