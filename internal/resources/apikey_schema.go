package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ApiKeySchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to create and manage API keys in Capella. API keys are used to authenticate and authorize access to Capella resources and services.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID (Access key) of the API key.",
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"name":            WithDescription(stringAttribute([]string{required, requiresReplace}), "Name of the API key."),
			"description":     WithDescription(stringDefaultAttribute("", optional, computed, requiresReplace, useStateForUnknown), "Description for the API key."),
			"expiry":          WithDescription(float64DefaultAttribute(180, optional, computed, requiresReplace, useStateForUnknown), "Expiry of the API key in number of days. If set to -1, the token will not expire."),
			"allowed_cidrs": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0.0.0.0/0")})),
				MarkdownDescription: "List of inbound CIDRs for the API key. The system making a request must come from one of the allowed CIDRs.",
			},
			"organization_roles": stringSetAttribute(required, requiresReplace),
			"resources": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    WithDescription(stringAttribute([]string{required}), "The GUID4 ID of the project."),
						"roles": WithDescription(stringSetAttribute(required), "Project Roles associated with the API key."),
						"type":  WithDescription(stringDefaultAttribute("project", optional, computed), "Resource type."),
					},
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Resources are the resource level permissions associated with the API key.",
			},
			"rotate": schema.NumberAttribute{
				Optional: true,
				Computed: true,
			},
			"secret": WithDescription(stringAttribute([]string{optional, computed, sensitive}), "A secret associated with API key. One has to follow the secret key policy, such as allowed characters and a length of 64 characters. If this field is left empty, a secret will be auto-generated."),
			"token":  WithDescription(stringAttribute([]string{computed, sensitive}), "The Token is a confidential piece of information that is used to authorize requests made to v4 endpoints."),
			"audit":  computedAuditAttribute(),
		},
	}
}
