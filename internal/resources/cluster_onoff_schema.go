package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ClusterOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id":            stringAttribute([]string{required, requiresReplace}),
			"project_id":                 stringAttribute([]string{required, requiresReplace}),
			"cluster_id":                 stringAttribute([]string{required, requiresReplace}),
			"state":                      stringAttribute([]string{required}),
			"turn_on_linked_app_service": boolAttribute(optional, computed),
		},
	}
}
