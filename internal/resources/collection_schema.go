package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func CollectionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to manage a collection within a scope in a bucket.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the organization."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the project."),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The GUID4 ID of the cluster."),
			"bucket_id": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name."),
			"scope_name": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The name of the scope."),
			"collection_name": WithDescription(stringAttribute([]string{required, requiresReplace}),
				"The name of the collection."),
			"max_ttl": WithDescription(int64Attribute(optional, computed),
				"The maximum Time To Live (TTL) for documents in the collection."),
		},
	}
}
