package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// projectBuilder is the SchemaBuilder instance for the project resource.
// It encapsulates the resource name and provides OpenAPI-aware description methods.
var projectBuilder = capellaschema.NewSchemaBuilder("project")

func ProjectSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This resource allows you to create and manage a project in an organization. Projects are used to organize and manage groups of operational clusters within organizations.",
		Attributes: map[string]schema.Attribute{
			"id": projectBuilder.WithOpenAPIDescription(
				&schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"id",
			).(*schema.StringAttribute),
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"name": projectBuilder.WithOpenAPIDescription(
				stringAttribute([]string{required}),
				"name",
			).(*schema.StringAttribute),
			"description": projectBuilder.WithOpenAPIDescription(
				stringAttribute([]string{optional, computed}),
				"description",
			).(*schema.StringAttribute),
			"if_match": WithDescription(stringAttribute([]string{optional}), "A precondition header that specifies the entity tag of a resource."),
			"etag":     WithDescription(stringAttribute([]string{computed}), "The ETag header value returned by the server, used for optimistic concurrency control."),
			"audit":    computedAuditAttribute(),
		},
	}
}
