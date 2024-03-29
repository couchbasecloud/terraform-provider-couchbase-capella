package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AppServiceOnOffOnDemandSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"state":           stringAttribute(required),
		},
	}
}
