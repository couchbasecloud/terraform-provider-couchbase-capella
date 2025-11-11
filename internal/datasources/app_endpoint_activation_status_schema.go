package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointActivationStatusBuilder = capellaschema.NewSchemaBuilder("appEndpointActivationStatus")

// AppEndpointActivationStatusSchema returns the schema for the AppEndpointActivationStatus data source.
func AppEndpointActivationStatusSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointActivationStatusBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appEndpointActivationStatusBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointActivationStatusBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointActivationStatusBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointActivationStatusBuilder, requiredString())
	capellaschema.AddAttr(attrs, "state", appEndpointActivationStatusBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the activation status of an App Endpoint.",
		Attributes:          attrs,
	}
}
