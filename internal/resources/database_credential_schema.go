package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var databaseCredentialBuilder = capellaschema.NewSchemaBuilder("databaseCredential")

// DatabaseCredentialSchema defines the schema for the terraform provider resource - "DatabaseCredential".
// This terraform resource directly maps to the database credential created for a Capella cluster.
// DatabaseCredential resource supports Create, Destroy, Read, Import and List operations.
func DatabaseCredentialSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", databaseCredentialBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "name", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "password", databaseCredentialBuilder, stringAttribute([]string{optional, computed, sensitive, useStateForUnknown}))

	// The credential type cannot be changed once the database credential is created,
	// so a change to it requires replacing the resource.
	credentialTypeAttr := stringDefaultAttribute(credentialTypeBasic, optional, computed, requiresReplace)
	credentialTypeAttr.Validators = append(credentialTypeAttr.Validators, stringvalidator.OneOf(credentialTypeBasic, credentialTypeAdvanced))
	capellaschema.AddAttr(attrs, "credential_type", databaseCredentialBuilder, credentialTypeAttr)

	capellaschema.AddAttr(attrs, "organization_id", databaseCredentialBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", databaseCredentialBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", databaseCredentialBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "audit", databaseCredentialBuilder, computedAuditAttribute())

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "name", databaseCredentialBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(scopeAttrs, "collections", databaseCredentialBuilder, stringSetAttribute(optional))

	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "name", databaseCredentialBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(bucketAttrs, "scopes", databaseCredentialBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "buckets", databaseCredentialBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	accessAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(accessAttrs, "privileges", databaseCredentialBuilder, stringSetAttribute(required))
	capellaschema.AddAttr(accessAttrs, "resources", databaseCredentialBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: resourcesAttrs,
	})

	// Exactly one of access or user_roles must be configured: access for a basic
	// credential type and user_roles for an advanced credential type.
	capellaschema.AddAttr(attrs, "access", databaseCredentialBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
		Validators: []validator.Set{
			setvalidator.ExactlyOneOf(path.MatchRoot("user_roles")),
		},
	})

	userRolesAttr := stringSetAttribute(optional)
	userRolesAttr.Validators = append(userRolesAttr.Validators, setvalidator.SizeAtLeast(1))
	capellaschema.AddAttr(attrs, "user_roles", databaseCredentialBuilder, userRolesAttr)

	return schema.Schema{
		MarkdownDescription: "Resource to create and manage a database credential for a cluster. Database credentials provide programmatic and application-level access to data on a database.",
		Attributes:          attrs,
	}
}
