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

	capellaschema.AddAttr(attrs, "organization_id", appEndpointResyncBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", appEndpointResyncBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointResyncBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointResyncBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointResyncBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "collections_processing", appEndpointResyncBuilder, &schema.MapAttribute{
		ElementType: types.SetType{
			ElemType: types.StringType,
		},
		Computed: true,
	}, "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_changed", appEndpointResyncBuilder, computedInt64(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "docs_processed", appEndpointResyncBuilder, computedInt64(), "ResyncStatus")

	// The published OpenAPI spec does not carry these two fields yet, so they supply a fallback
	// description; AddAttr prefers the spec as soon as it has one.
	docsTargeted := computedInt64()
	docsTargeted.MarkdownDescription = capellaschema.DocsTargetedDescription
	capellaschema.AddAttr(attrs, "docs_targeted", appEndpointResyncBuilder, docsTargeted, "ResyncStatus")

	docsErrored := computedInt64()
	docsErrored.MarkdownDescription = capellaschema.DocsErroredDescription
	capellaschema.AddAttr(attrs, "docs_errored", appEndpointResyncBuilder, docsErrored, "ResyncStatus")

	capellaschema.AddAttr(attrs, "last_error", appEndpointResyncBuilder, computedString(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "start_time", appEndpointResyncBuilder, computedString(), "ResyncStatus")
	capellaschema.AddAttr(attrs, "state", appEndpointResyncBuilder, computedString(), "ResyncStatus")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the resync status of an App Endpoint.",
		Attributes:          attrs,
	}
}
