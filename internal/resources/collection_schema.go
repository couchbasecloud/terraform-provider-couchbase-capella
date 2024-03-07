package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CollectionSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"organization_id": stringAttribute(required, requiresReplace),
			"project_id":      stringAttribute(required, requiresReplace),
			"cluster_id":      stringAttribute(required, requiresReplace),
			"bucket_id":       stringAttribute(required, requiresReplace),
			"scope_name":      stringAttribute(required, requiresReplace),
			"collection_name": stringAttribute(required, requiresReplace),
			"max_ttl":         int64Attribute(required, requiresReplace),
			"uid":             stringAttribute(computed),
		},
	}

}
