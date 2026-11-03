package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointAuditLogEventIDsBuilder = capellaschema.NewSchemaBuilder("appEndpointAuditLogEventIDs", "AuditLogEventsResponse")

// AppEndpointAuditLogEventIDsSchema returns the schema for the app_endpoint_audit_log_event_ids data source.
func AppEndpointAuditLogEventIDsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointAuditLogEventIDsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appEndpointAuditLogEventIDsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointAuditLogEventIDsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointAuditLogEventIDsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointAuditLogEventIDsBuilder, requiredStringWithValidator())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appEndpointAuditLogEventIDsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "name", appEndpointAuditLogEventIDsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "description", appEndpointAuditLogEventIDsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "enabled", appEndpointAuditLogEventIDsBuilder, computedBool())
	capellaschema.AddAttr(dataAttrs, "filterable", appEndpointAuditLogEventIDsBuilder, computedBool())

	capellaschema.AddAttr(attrs, "data", appEndpointAuditLogEventIDsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit log event IDs of an App Endpoint. Use these event IDs to configure App Endpoint audit logging.",
		Attributes:          attrs,
	}
}
