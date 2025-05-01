package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func FreeTierAppServiceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)),
				"Organization ID is the unique identifier for the organization. It is used to group resources and manage access within the organization."),
			"project_id": stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)),
			"cluster_id": stringAttribute([]string{required, requiresReplace}, validator.String(
				stringvalidator.LengthAtLeast(1),
			)),
			"name":           stringAttribute([]string{required}),
			"description":    stringDefaultAttribute("", optional, computed),
			"nodes":          int64Attribute(computed, useStateForUnknown),
			"cloud_provider": stringAttribute([]string{computed, useStateForUnknown}),
			"current_state":  stringAttribute([]string{computed}),
			"version":        stringAttribute([]string{computed}),
			"compute": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"cpu": int64Attribute(computed, useStateForUnknown),
					"ram": int64Attribute(computed, useStateForUnknown),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"audit": computedAuditAttribute(),
			"plan":  stringAttribute([]string{computed}),
			"etag":  stringAttribute([]string{computed}),
		},
	}
}
