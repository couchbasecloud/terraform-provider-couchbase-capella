package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ApiKeySchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expiry": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.RequiresReplace(),
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")})),
			},
			"organization_roles": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"resources": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
						"roles": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"rotate": schema.NumberAttribute{
				Optional: true,
				Computed: true,
			},
			"secret": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
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
