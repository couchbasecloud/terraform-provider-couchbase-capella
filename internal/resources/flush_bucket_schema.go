package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FlushBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("organization_id is the ID of the Capella organization"),
				stringvalidator.LengthAtLeast(1)),
			"project_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("project_id is the ID of the Capella project"),
				stringvalidator.LengthAtLeast(1)),
			"cluster_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("cluster_id is the ID of the Capella cluster"),
				stringvalidator.LengthAtLeast(1)),
			"bucket_id": stringAttributeWithValidators([]string{required, requiresReplace},
				withMarkdown[*schema.StringAttribute]("bucket_id is the ID of the Capella bucket"),
				stringvalidator.LengthAtLeast(1)),
		},
	}
}
