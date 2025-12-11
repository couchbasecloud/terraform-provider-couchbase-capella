package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var allowlistsBuilder = capellaschema.NewSchemaBuilder("allowlists")

// AllowListsSchema returns the schema for the AllowLists data source.
func AllowListsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", allowlistsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", allowlistsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", allowlistsBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cidr", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "comment", allowlistsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expires_at", allowlistsBuilder, computedString())

	// Build audit attributes
	auditAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(auditAttrs, "created_at", allowlistsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "created_by", allowlistsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_at", allowlistsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "modified_by", allowlistsBuilder, computedString(), "CouchbaseAuditData")
	capellaschema.AddAttr(auditAttrs, "version", allowlistsBuilder, computedInt64(), "CouchbaseAuditData")

	capellaschema.AddAttr(dataAttrs, "audit", allowlistsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: auditAttrs,
	})

	capellaschema.AddAttr(attrs, "data", allowlistsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "Retrieves the allowlists details for a Capella cluster.",
		Attributes:          attrs,
	}
}
