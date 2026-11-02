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

	// The spec types events as an untyped map, so the nested fields carry call-site descriptions.
	dataAttrs := map[string]schema.Attribute{
		"id":          &schema.Int64Attribute{Computed: true, MarkdownDescription: "The ID of the audit event."},
		"name":        &schema.StringAttribute{Computed: true, MarkdownDescription: "The name of the audit event."},
		"description": &schema.StringAttribute{Computed: true, MarkdownDescription: "The description of the audit event."},
		"enabled":     &schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the audit event is enabled."},
		"filterable":  &schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the audit event can be filtered."},
	}

	capellaschema.AddAttr(attrs, "data", appEndpointAuditLogEventIDsBuilder, &schema.SetNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The audit events of the App Endpoint.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the audit log event IDs of an App Endpoint. Use these event IDs to configure App Endpoint audit logging.",
		Attributes:          attrs,
	}
}
