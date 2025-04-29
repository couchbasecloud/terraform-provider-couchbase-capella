package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the user created.",
			},
			"name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the user.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Status depicts user status whether they are verified or not. " +
					"It can be one of the following values: verified, not-verified, pending-primary.",
			},
			"inactive": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Inactive depicts whether the user has accepted the invite for the organization.",
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email of the user.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Capella organization.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"organization_roles": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The organization roles associated to the user. They determines the privileges user possesses in the organization.",
			},
			"last_login": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Time(UTC) at which user last logged in.",
			},
			"region": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The region of the user",
			},
			"time_zone": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The Time zone of the user.",
			},
			"enable_notifications": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "After enabling email notifications for your account, you will start receiving email notification alerts from all databases in projects you are a part of.",
			},
			"expires_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time at which the user expires.",
			},
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("project"),
							MarkdownDescription: "Type of the resource.",
						},
						"id": schema.StringAttribute{Required: true, MarkdownDescription: "The ID of the project."},
						"roles": schema.SetAttribute{
							ElementType:         types.StringType,
							Required:            true,
							MarkdownDescription: "Project Roles associated with the User.",
						},
					},
				},
			},
			"audit": computedAuditAttribute(),
		},
	}
}
