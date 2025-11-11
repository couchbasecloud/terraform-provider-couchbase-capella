package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var auditLogEventIDsBuilder = capellaschema.NewSchemaBuilder("auditLogEventIDs")

// AuditLogEventIDsSchema returns the schema for the AuditLogEventIDs data source.
func AuditLogEventIDsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", auditLogEventIDsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", auditLogEventIDsBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", auditLogEventIDsBuilder, requiredString())

	// Build data attributes
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "description", auditLogEventIDsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "id", auditLogEventIDsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "module", auditLogEventIDsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "name", auditLogEventIDsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", auditLogEventIDsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve audit log event IDs for an operational cluster. These event IDs can be used to filter audit logs and configure audit logging.",
		Attributes:          attrs,
	}
}
