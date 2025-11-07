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
	capellaschema.AddAttr(attrs, "organization_id", userBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["status"] = WithDescription(stringAttribute([]string{computed}), "Depicts the user's status by determining whether they are verified or not. It can be one of the following values: verified, not-verified, pending-primary.")
	attrs["inactive"] = WithDescription(boolAttribute(computed), "Inactive depicts whether the user has accepted the invite for the organization.")
	attrs["email"] = WithDescription(stringAttribute([]string{required, requiresReplace}), "Email of the user.")
	attrs["organization_roles"] = WithDescription(stringListAttribute(required), "The organization roles associated to the user. They determines the privileges user possesses in the organization.")
	attrs["last_login"] = WithDescription(stringAttribute([]string{computed}), "Time(UTC) at which user last logged in.")
	attrs["region"] = WithDescription(stringAttribute([]string{computed}), "The region of the user.")
	attrs["time_zone"] = WithDescription(stringAttribute([]string{computed}), "The time zone of the user.")
	attrs["enable_notifications"] = WithDescription(boolAttribute(computed), "After enabling email notifications for your account, you will start receiving email notification alerts from all databases in projects you are a part of.")
	attrs["expires_at"] = WithDescription(stringAttribute([]string{computed}), "Time at which the user expires.")
	attrs["resources"] = schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type":  WithDescription(stringDefaultAttribute("project", optional, computed), "Resource type."),
				"id":    WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the project."),
				"roles": WithDescription(stringSetAttribute(required), "The project roles associated with the user."),
			},
		},
	}
	capellaschema.AddAttr(attrs, "audit", userBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "This User resource allows you to manage users in your Capella organization. You can create, update, and delete users, as well as manage their roles and permissions within the organization.",
		Attributes:          attrs,
	}
}
