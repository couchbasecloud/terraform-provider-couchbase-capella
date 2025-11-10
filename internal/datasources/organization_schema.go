package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var organizationBuilder = capellaschema.NewSchemaBuilder("organization")

// OrganizationSchema returns the schema for the Organization data source.
func OrganizationSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", organizationBuilder, requiredString())
	capellaschema.AddAttr(attrs, "name", organizationBuilder, computedString())
	capellaschema.AddAttr(attrs, "description", organizationBuilder, computedString())

	// Build preferences attributes
	preferencesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(preferencesAttrs, "session_duration", organizationBuilder, computedInt64())

	capellaschema.AddAttr(attrs, "preferences", organizationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: preferencesAttrs,
	})

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", organizationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "created_by", organizationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_at", organizationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "modified_by", organizationBuilder, computedString())
	capellaschema.AddAttr(auditAttrs, "version", organizationBuilder, computedInt64())

	capellaschema.AddAttr(attrs, "audit", organizationBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve information about a Capella organization.",
		Attributes:          attrs,
	}
}
