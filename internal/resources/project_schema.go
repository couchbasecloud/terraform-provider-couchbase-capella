package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ProjectSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Resource to create and manage a project in a Capella organization. Projects are used to organize and manage groups of Couchbase databases within organizations.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The ID of the project created.",
			},
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"name":            WithDescription(stringAttribute([]string{required}), "The name of the project (up to 128 characters)."),
			"description":     WithDescription(stringAttribute([]string{optional, computed}), "A short description of the project (up to 256 characters)."),
			"if_match":        WithDescription(stringAttribute([]string{optional}), "A precondition header that specifies the entity tag of a resource."),
			"etag":            stringAttribute([]string{computed}),
			"audit":           computedAuditAttribute(),
		},
	}
}
