package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FlushBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)),
			"project_id":      stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)),
			"bucket_id":       stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)),
		},
	}
}
