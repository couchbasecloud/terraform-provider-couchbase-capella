package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointResyncBuilder = capellaschema.NewSchemaBuilder("appEndpointResync")

func AppEndpointResyncSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointResyncBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "scopes", appEndpointResyncBuilder, mapAttribute(types.SetType{ElemType: types.StringType}, []string{optional}...))
	capellaschema.AddAttr(attrs, "collections_processing", appEndpointResyncBuilder, mapAttribute(types.SetType{ElemType: types.StringType}, []string{computed}...))
	capellaschema.AddAttr(attrs, "docs_changed", appEndpointResyncBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "docs_processed", appEndpointResyncBuilder, int64Attribute(computed))
	capellaschema.AddAttr(attrs, "last_error", appEndpointResyncBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "start_time", appEndpointResyncBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "state", appEndpointResyncBuilder, stringAttribute([]string{computed}))

	return schema.Schema{
		MarkdownDescription: "Manages App Endpoint Resync operations. This resource allows you to create and manage resync operations for App Endpoints in Couchbase Capella.",
		Attributes:          attrs,
	}
}
