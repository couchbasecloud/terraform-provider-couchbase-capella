package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func UserSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"inactive": schema.BoolAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"organization_roles": schema.StringAttribute{
				Required: true,
			},
			"lastLogin": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Computed: true,
			},
			"time_zone": schema.StringAttribute{
				Computed: true,
			},
			"enable_notifications": schema.BoolAttribute{
				Computed: true,
			},
			"expires_at": schema.StringAttribute{
				Computed: true,
			},
			"resources": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Optional: true,
					},
					"id": schema.StringAttribute{
						Required: true,
					},
					"roles": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
			"if_match": schema.StringAttribute{
				Optional: true,
			},
			"audit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}
}
