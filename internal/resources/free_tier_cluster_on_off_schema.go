package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierClusterOnOffSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("organization id"),
				validator.String(stringvalidator.LengthAtLeast(1))),
			"project_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("project id"),
				validator.String(stringvalidator.LengthAtLeast(1))),
			"cluster_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("cluster id"),
				validator.String(stringvalidator.LengthAtLeast(1))),
			"state": stringAttribute([]string{required}),
		},
	}
}
