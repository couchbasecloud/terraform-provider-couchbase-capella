package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appEndpointActivationStatusBuilder = capellaschema.NewSchemaBuilder("appEndpointActivationStatus")

// AppEndpointActivationStatusSchema defines the schema for the App Endpoint activation status resource.
func AppEndpointActivationStatusSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appEndpointActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", appEndpointActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", appEndpointActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_service_id", appEndpointActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_endpoint_name", appEndpointActivationStatusBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "state", appEndpointActivationStatusBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "Manages the activation status of an App Endpoint. This resource is used to activate or deactivate an App Endpoint on-demand.",
		Attributes:          attrs,
	}
}
