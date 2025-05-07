package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func FreeTierClusterOnOffSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the on/off state resource of a Free Tier Cluster.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)), "The GUID4 ID of the cluster."),
			"state":           WithDescription(stringAttribute([]string{required}), "The on/off state of the Free Tier Cluster."),
		},
	}
}
