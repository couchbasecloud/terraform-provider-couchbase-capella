package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServiceBuilder = capellaschema.NewSchemaBuilder("appService")

func AppServiceSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "id", appServiceBuilder, &schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	})
	capellaschema.AddAttr(attrs, "organization_id", appServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", appServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", appServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "name", appServiceBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "description", appServiceBuilder, stringDefaultAttribute("", optional, computed, requiresReplace))
	capellaschema.AddAttr(attrs, "nodes", appServiceBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "cloud_provider", appServiceBuilder, stringAttribute([]string{optional, computed}))
	capellaschema.AddAttr(attrs, "current_state", appServiceBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "version", appServiceBuilder, stringAttribute([]string{computed}))
	capellaschema.AddAttr(attrs, "audit", appServiceBuilder, computedAuditAttribute())
	capellaschema.AddAttr(attrs, "if_match", appServiceBuilder, stringAttribute([]string{optional}))
	capellaschema.AddAttr(attrs, "etag", appServiceBuilder, stringAttribute([]string{computed}))

	computeAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computeAttrs, "cpu", appServiceBuilder, int64Attribute(required))
	capellaschema.AddAttr(computeAttrs, "ram", appServiceBuilder, int64Attribute(required))

	capellaschema.AddAttr(attrs, "compute", appServiceBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: computeAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "This resource allows you to create and manage an App Service in Capella. App Service is a fully managed application backend designed to provide data synchronization between mobile or IoT applications running Couchbase Lite and your Couchbase Capella database.",
		Attributes:          attrs,
	}
}
