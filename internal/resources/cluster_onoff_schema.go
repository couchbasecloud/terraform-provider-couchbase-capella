package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ClusterOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id":            stringAttribute(required, requiresReplace),
			"project_id":                 stringAttribute(required, requiresReplace),
			"cluster_id":                 stringAttribute(required, requiresReplace),
			"state":                      stringAttribute(required),
			"turn_on_linked_app_service": boolAttribute(optional, computed),
		},
	}
}
