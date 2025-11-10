package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicesCidrBuilder = capellaschema.NewSchemaBuilder("appServicesCidr")

// AppServicesCidrSchema returns the schema for the App Service CIDRs data source.
func AppServicesCidrSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "app_service_id", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cidr", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "comment", appServicesCidrBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expires_at", appServicesCidrBuilder, computedString())

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", appServicesCidrBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", appServicesCidrBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", appServicesCidrBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", appServicesCidrBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", appServicesCidrBuilder, computedInt64(), "CouchbaseAuditData")

	capellaschema.AddAttr(dataAttrs, "audit", appServicesCidrBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	capellaschema.AddAttr(attrs, "data", appServicesCidrBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Retrieves the allowed CIDR blocks for a Capella App Service.",
		Attributes:          attrs,
	}
}
