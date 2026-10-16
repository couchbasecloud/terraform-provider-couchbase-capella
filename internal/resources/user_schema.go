package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	capellaschema.AddAttr(attrs, "organization_id", userBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "organization_roles", userBuilder, &schema.ListAttribute{
		Required:    true,
		ElementType: types.StringType,
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.ValueStringsAre(stringvalidator.OneOf(validOrganizationRoles...)),
		},
	})
	capellaschema.AddAttr(attrs, "last_login", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "region", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "time_zone", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "enable_notifications", userBuilder, boolAttribute(computed))
	capellaschema.AddAttr(attrs, "expires_at", userBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", userBuilder, computedAuditAttribute())

	resourceAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourceAttrs, "type", userBuilder, &schema.StringAttribute{
		Optional:   true,
		Computed:   true,
		Default:    stringdefault.StaticString(resourceTypeProject),
		Validators: []validator.String{stringvalidator.OneOf(resourceTypeProject)},
	}, "Resource")
	capellaschema.AddAttr(resourceAttrs, "id", userBuilder, stringAttribute(
		[]string{required},
		validator.String(stringvalidator.LengthAtLeast(1)),
		validator.String(stringvalidator.RegexMatches(apiKeyResourceIDPattern, "resources.id must be a valid UUID")),
	), "Resource")
	capellaschema.AddAttr(resourceAttrs, "roles", userBuilder, &schema.SetAttribute{
		Required:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			setvalidator.ValueStringsAre(stringvalidator.OneOf(validProjectRoles...)),
		},
	}, "Resource")

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
