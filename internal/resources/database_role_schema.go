package resources

import (
	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var databaseRoleBuilder = capellaschema.NewSchemaBuilder("databaseRole")

// DatabaseRoleSchema defines the Terraform schema for the database_role resource.
func DatabaseRoleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", databaseRoleBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "name", databaseRoleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "description", databaseRoleBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "organization_id", databaseRoleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", databaseRoleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", databaseRoleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "audit", databaseRoleBuilder, computedAuditAttribute())

	scopeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scopeAttrs, "name", databaseRoleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(scopeAttrs, "collections", databaseRoleBuilder, stringSetAttribute(optional))

	bucketAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(bucketAttrs, "name", databaseRoleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(bucketAttrs, "scopes", databaseRoleBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: scopeAttrs,
		},
	})

	resourcesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(resourcesAttrs, "buckets", databaseRoleBuilder, &schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: bucketAttrs,
		},
	})

	accessAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(accessAttrs, "privileges", databaseRoleBuilder, stringSetAttribute(required))
	capellaschema.AddAttr(accessAttrs, "resources", databaseRoleBuilder, &schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: resourcesAttrs,
	})

	capellaschema.AddAttr(attrs, "access", databaseRoleBuilder, &schema.SetNestedAttribute{
		Required: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: accessAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Resource to create and manage a database user role for a cluster. Database roles define reusable sets of privileges that can be assigned to database credentials.",
		Attributes:          attrs,
	}
}
