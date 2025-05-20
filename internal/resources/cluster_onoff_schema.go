package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ClusterOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the On/Off state of an operational cluster. This resource allows you to turn your cluster on or off.",
		Attributes: map[string]schema.Attribute{
			"organization_id":            WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"state":                      WithDescription(stringAttribute([]string{required}), "State of the cluster. It can be `on` and `off`."),
			"turn_on_linked_app_service": WithDescription(boolAttribute(optional, computed), " Whether to turn on the linked App Service when the cluster is turned on."),
		},
	}
}
