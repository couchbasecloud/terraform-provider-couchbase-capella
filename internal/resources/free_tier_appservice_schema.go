package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
				MarkdownDescription: "ID of the app service. This is a unique identifier for the app service.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the capella organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{(stringvalidator.LengthAtLeast(1))},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the capella project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the free-tier capella cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the free-tier app service.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Description of the free-tier app service.",
			},
			"nodes": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Number of nodes in the free-tier app service.",
			},
			"cloud_provider": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cloud provider of the free-tier app service. current supported providers are aws, gcp and azure",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"current_state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current state of the free-tier app service.",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the free-tier app service.",
			},
			"compute": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"cpu": int64Attribute(computed, useStateForUnknown),
					"ram": int64Attribute(computed, useStateForUnknown),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Compute configuration of the free-tier app service.",
			},
			"audit": computedAuditAttribute(),
			"plan": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Plan associated with the free-tier app service.",
			},
			"etag": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ETag of the free-tier app service for concurrency control.",
			},
		},
	}
}
