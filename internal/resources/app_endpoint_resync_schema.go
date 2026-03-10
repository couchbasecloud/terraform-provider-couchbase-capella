package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointResyncBuilder = capellaschema.NewSchemaBuilder("appEndpointResync", "resyncRequest")

func AppEndpointResyncSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointResyncBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", appEndpointResyncBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointResyncBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointResyncBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))

	// Input field from ResyncRequest schema
	capellaschema.AddAttr(attrs, "scopes", appEndpointResyncBuilder, mapAttribute(types.SetType{ElemType: types.StringType}, []string{optional, requiresReplace}...))

	// Output fields from ResyncStatus schema
	capellaschema.AddAttr(attrs, "collections_processing", appEndpointResyncBuilder, mapAttribute(types.SetType{ElemType: types.StringType}, []string{computed}...), "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_changed", appEndpointResyncBuilder, int64Attribute(computed), "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_processed", appEndpointResyncBuilder, int64Attribute(computed), "ResyncStatus")
	capellaschema.AddAttr(attrs, "last_error", appEndpointResyncBuilder, stringAttribute([]string{computed}), "ResyncStatus")
	capellaschema.AddAttr(attrs, "start_time", appEndpointResyncBuilder, stringAttribute([]string{computed}), "ResyncStatus")
	capellaschema.AddAttr(attrs, "state", appEndpointResyncBuilder, stringAttribute([]string{computed}), "ResyncStatus")

	return schema.Schema{
		MarkdownDescription: "Manages an App Endpoints resync operations. This resource allows you to create and manage resync operations for an App Endpoint on an App Service in Couchbase Capella.",
		Attributes:          attrs,
	}
}
