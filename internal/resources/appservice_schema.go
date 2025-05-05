package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AppServiceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the App Service created.",
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Capella organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"name":            WithDescription(stringAttribute([]string{required, requiresReplace}), "Name of the cluster (up to 256 characters)."),
			"description":     WithDescription(stringDefaultAttribute("", optional, computed, requiresReplace), "A short description of the App Service."),
			"nodes":           WithDescription(int64Attribute(optional, computed), "Number of nodes configured for the App Service. Number of nodes configured for the App Service. The number of nodes can range from 2 to 12."),
			"cloud_provider":  WithDescription(stringAttribute([]string{optional, computed}), "Provider is the cloud service provider for the App Service."),
			"current_state":   WithDescription(stringAttribute([]string{computed}), "The current state of the App Service."),
			"version":         WithDescription(stringAttribute([]string{computed}), "The version of the App Service server. If left empty, it will be defaulted to the latest available version."),
			"compute": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "The CPU and RAM configuration of the App Service.",
				Attributes: map[string]schema.Attribute{
					"cpu": WithDescription(int64Attribute(required), "CPU units (cores)."),
					"ram": WithDescription(int64Attribute(required), "RAM units (GB)."),
				},
			},
			"audit":    computedAuditAttribute(),
			"if_match": WithDescription(stringAttribute([]string{optional}), "A precondition header that specifies the entity tag of a resource."),
			"etag":     stringAttribute([]string{computed}),
		},
	}
}
