package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CollectionSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute([]string{required, requiresReplace}),
			"project_id":      stringAttribute([]string{required, requiresReplace}),
			"cluster_id":      stringAttribute([]string{required, requiresReplace}),
			"bucket_id":       stringAttribute([]string{required, requiresReplace}),
			"scope_name":      stringAttribute([]string{required, requiresReplace}),
			"collection_name": stringAttribute([]string{required, requiresReplace}),
			"max_ttl":         int64Attribute(optional, computed),
		},
	}

}
