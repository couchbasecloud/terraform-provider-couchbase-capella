package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointServiceBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpointService")

// AppServicePrivateEndpointServiceSchema returns the schema for the
// app_service_private_endpoint_service data source.
func AppServicePrivateEndpointServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointServiceBuilder, requiredString())
	capellaschema.AddAttr(attrs, "enabled", appServicePrivateEndpointServiceBuilder, computedBool())
	capellaschema.AddAttr(attrs, "state", appServicePrivateEndpointServiceBuilder, computedString(), "GetAppServicePrivateEndpointStateResponse")
	capellaschema.AddAttr(attrs, "target_state", appServicePrivateEndpointServiceBuilder, computedString(), "GetAppServicePrivateEndpointStateResponse")

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve private endpoint service information for an App Service (Sync Gateway).",
		Attributes:          attrs,
	}
}
