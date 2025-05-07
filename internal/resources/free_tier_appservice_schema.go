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
		MarkdownDescription: "Manages free-tier app services resources assosciated with a free-tier cluster",
		Attributes: map[string]schema.Attribute{
			"id": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "ID of the free-tier app service."),
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)),
				"Organization ID is the unique identifier for the organization. It is used to group resources and manage access within the organization."),
			"project_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)), "The ID of the Capella project"),
			"cluster_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)), "The ID of the Capella cluster"),

			"name": WithDescription(stringAttribute([]string{required}), "Name of the free-tier app service."),
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Description of the free-tier app service.",
			},
			"nodes": WithDescription(int64Attribute(computed, useStateForUnknown), "Number of nodes in the free-tier app service."),

			"cloud_provider": WithDescription(stringAttribute([]string{computed, useStateForUnknown}), "Cloud provider of the free-tier app service. current supported providers are aws, gcp and azure"),

			"current_state": WithDescription(stringAttribute([]string{computed}), "Current state of the free-tier app service."),

			"version": WithDescription(stringAttribute([]string{computed}), "Version of the free-tier app service."),

			"compute": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Compute configuration of the free-tier app service.",
				Attributes: map[string]schema.Attribute{
					"cpu": WithDescription(int64Attribute(computed, useStateForUnknown), "Number of CPUs of the free-tier app service node."),

					"ram": WithDescription(int64Attribute(computed, useStateForUnknown), "Amount of RAM of the free-tier app service node."),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"audit": computedAuditAttribute(),
			"plan":  WithDescription(stringAttribute([]string{computed}), "Plan associated with the free-tier app service."),
			"etag":  WithDescription(stringAttribute([]string{computed}), "ETag of the free-tier app service."),
		},
	}
}
