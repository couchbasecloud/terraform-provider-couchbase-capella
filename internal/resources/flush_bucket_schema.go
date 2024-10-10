package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FlushBucketSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"bucket_id":       stringAttribute([]string{required, requiresReplace}),
		},
	}
}
