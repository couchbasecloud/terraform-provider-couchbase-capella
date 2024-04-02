package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AppServiceOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"app_service_id":  stringAttribute([]string{required, requiresReplace}),
			"state":           stringAttribute([]string{required}),
		},
	}
}
