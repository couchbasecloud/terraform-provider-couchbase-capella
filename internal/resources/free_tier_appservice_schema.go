package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierAppServiceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage free tier App Services associated with a free tier operational cluster.",
		Attributes: map[string]schema.Attribute{
			"id": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "ID of the free tier App Service."),
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)),
				"The GUID4 ID of the organization."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)), "The GUID4 ID of the project."),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)), " The GUID4 ID of the cluster."),

			"name": WithDescription(stringAttribute([]string{required}), "Name of the free tier App Service."),
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Description of the free tier App Service.",
			},
			"nodes": WithDescription(int64Attribute(computed, useStateForUnknown), "Number of nodes in the free tier App Service."),

			"cloud_provider": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "The Cloud Service Provider of the free tier App Service. The supported Cloud Service Providers are AWS, GCP, and Azure."),

			"current_state": WithDescription(stringAttribute([]string{computed}), "Current state of the free tier App Service."),

			"version": WithDescription(stringAttribute([]string{computed}), "The Server version of the free tier App Service."),

			"compute": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Compute configuration of the free tier App Service.",
				Attributes: map[string]schema.Attribute{
					"cpu": WithDescription(int64Attribute(computed, useStateForUnknown), "The number of CPUs of the free tier App Service node."),

					"ram": WithDescription(int64Attribute(computed, useStateForUnknown), "The amount of RAM of the free tier App Service node."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"audit": computedAuditAttribute(),
			"plan":  WithDescription(stringAttribute([]string{computed}), "The Support plan associated with the free tier App Service. The Support plan options are 'Basic', 'Developer Pro', or 'Enterprise'."),
			"etag":  WithDescription(stringAttribute([]string{computed}), "ETag of the free tier App Service."),
		},
	}
}
