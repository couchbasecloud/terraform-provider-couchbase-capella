package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ScopeSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a scope within a bucket.",
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
			"collections": schema.SetNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of collections within this scope.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"max_ttl": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The maximum Time To Live (TTL) for documents in the collection.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the collection.",
						},
					},
				},
			},
		},
	}
}
