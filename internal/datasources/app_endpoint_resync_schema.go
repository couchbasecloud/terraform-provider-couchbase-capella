package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointResyncBuilder = capellaschema.NewSchemaBuilder("appEndpointResync")

// AppEndpointResyncSchema returns the schema for the AppEndpointResync data source.
func AppEndpointResyncSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointResyncBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appEndpointResyncBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointResyncBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointResyncBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointResyncBuilder, requiredString())
	capellaschema.AddAttr(attrs, "collections_processing", appEndpointResyncBuilder, &schema.MapAttribute{
		ElementType: types.SetType{
			ElemType: types.StringType,
		},
		Computed: true,
	}, "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_changed", appEndpointResyncBuilder, computedInt64(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_processed", appEndpointResyncBuilder, computedInt64(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "last_error", appEndpointResyncBuilder, computedString(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "start_time", appEndpointResyncBuilder, computedString(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "state", appEndpointResyncBuilder, computedString(), "ResyncStatus")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the resync status of an App Endpoint.",
		Attributes:          attrs,
	}
}
