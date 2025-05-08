package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AppServiceOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the on-demand state of an app service. This resource is used to turn on/off the app service on-demand.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"state":           WithDescription(stringAttribute([]string{required}), "The state of the app service on-demand. It can be on/off."),
		},
	}
}
