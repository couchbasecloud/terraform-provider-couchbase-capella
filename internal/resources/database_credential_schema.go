package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

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
	capellaschema.AddAttr(attrs, "organization_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", databaseCredentialBuilder, stringAttribute([]string{required, requiresReplace}))
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

	capellaschema.AddAttr(attrs, "access", databaseCredentialBuilder, &schema.SetNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Resource to create and manage a database credential for a cluster. Database credentials provide programmatic and application-level access to data on a database.",
		Attributes:          attrs,
	}
}
