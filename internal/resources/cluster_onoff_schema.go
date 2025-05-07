package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ClusterOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the on/off state of a capella cluster. This resource is used to turn on or off cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id":            WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":                 WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"state":                      WithDescription(stringAttribute([]string{required}), "State of the cluster. Possible values are `on` and `off`."),
			"turn_on_linked_app_service": WithDescription(boolAttribute(optional, computed), "Whether to turn on the linked app service when the cluster is turned on."),
		},
	}
}
