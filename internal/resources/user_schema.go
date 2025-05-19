package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The User resource allows you to manage users in your Capella organization. You can create, update, and delete users, as well as manage their roles and permissions within the organization.",
		Attributes: map[string]schema.Attribute{
			"id":                   WithDescription(stringAttribute([]string{computed}), "The UUID of the user created."),
			"name":                 WithDescription(stringAttribute([]string{optional, computed}), "The name of the user."),
			"status":               WithDescription(stringAttribute([]string{computed}), "Status depicts user status whether they are verified or not."+"It can be one of the following values: verified, not-verified, pending-primary."),
			"inactive":             WithDescription(boolAttribute(computed), "Inactive depicts whether the user has accepted the invite for the organization."),
			"email":                WithDescription(stringAttribute([]string{required, requiresReplace}), "Email of the user."),
			"organization_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"organization_roles":   WithDescription(stringListAttribute(required), "The organization roles associated to the user. They determines the privileges user possesses in the organization."),
			"last_login":           WithDescription(stringAttribute([]string{computed}), "The Time(UTC) at which user last logged in."),
			"region":               WithDescription(stringAttribute([]string{computed}), "The region of the user"),
			"time_zone":            WithDescription(stringAttribute([]string{computed}), "The Time zone of the user."),
			"enable_notifications": WithDescription(boolAttribute(computed), "After enabling email notifications for your account, you will start receiving email notification alerts from all databases in projects you are a part of."),
			"expires_at":           WithDescription(stringAttribute([]string{computed}), "Time at which the user expires."),
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type":  WithDescription(stringDefaultAttribute("project", optional, computed), "Type of the resource."),
						"id":    WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the project."),
						"roles": WithDescription(stringSetAttribute(required), "Project Roles associated with the User."),
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
