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
	attrs := make(map[string]schema.Attribute)

	// All fields use AddAttr - automatically finds description from OpenAPI or common registry
	capellaschema.AddAttr(attrs, "id", projectBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", projectBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", projectBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(attrs, "description", projectBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "if_match", projectBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(attrs, "etag", projectBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", projectBuilder, computedAuditAttribute())

	return schema.Schema{
		MarkdownDescription: "This resource allows you to create and manage a project in an organization. Projects are used to organize and manage groups of operational clusters within organizations.",
		Attributes:          attrs,
	}
}
