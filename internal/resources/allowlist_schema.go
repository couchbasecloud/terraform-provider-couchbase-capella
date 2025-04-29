package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func AllowlistsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the Allowed IP addresses to connect to databases for Couchbase Capella.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the allowed CIDR.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella project.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The trusted CIDR to allow the database connections from.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A short description of the allowed CIDR.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expires_at": schema.StringAttribute{
				Optional: true,
				MarkdownDescription: "An RFC3339 timestamp determining when the allowed CIDR will expire. " +
					"If this field is omitted then the allowed CIDR is permanent and will never automatically expire.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
