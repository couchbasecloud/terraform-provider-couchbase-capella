package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FlushBucketSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resources allows you to flush a Capella bucket. Flushing a bucket causes all documents in the bucket to be deleted by the system at the earliest. This operation can only be performed if the bucket has been configured with flushEnabled to true.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the Capella organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the cluster."),
			"bucket_id":       WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The ID of the bucket. It is the URL-compatible base64 encoding of the bucket name."),
		},
	}
}
