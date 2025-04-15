package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierClusterOnOffSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"project_id":      stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))),
			"state":           stringAttribute([]string{required}),
		},
	}
}
