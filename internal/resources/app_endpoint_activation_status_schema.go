package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// AppEndpointActivationStatusSchema defines the schema for the App Endpoint activation status resource.
func AppEndpointActivationStatusSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the activation status of an App Endpoint. This resource allows you to bring an App Endpoint online (resume) or take it offline (pause).",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"app_endpoint_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app endpoint."),
			"state":           WithDescription(stringAttribute([]string{required}), "The activation state of the app endpoint. It can be 'on' (online/resumed) or 'off' (offline/paused)."),
		},
	}
}
