package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AppEndpointActivationStatusSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the activation status of an App Endpoint. This resource is used to activate or deactivate an App Endpoint on-demand.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"app_endpoint_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app endpoint."),
			"state":           WithDescription(stringAttribute([]string{required}), "The activation state of the app endpoint. It can be 'active' or 'inactive'."),
		},
	}
}
