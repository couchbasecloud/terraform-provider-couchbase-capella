package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var importFilterBuilder = capellaschema.NewSchemaBuilder("importFilter")

// ImportFilterSchema defines the Terraform resource schema for managing App Endpoint Import Filters.
func ImportFilterSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", importFilterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", importFilterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", importFilterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_service_id", importFilterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_endpoint_name", importFilterBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "scope", importFilterBuilder, stringDefaultAttribute("_default", optional, computed, requiresReplace))
	capellaschema.AddAttr(attrs, "collection", importFilterBuilder, stringDefaultAttribute("_default", optional, computed, requiresReplace))
	capellaschema.AddAttr(attrs, "import_filter", importFilterBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "This Import Filter resource allows you to manage the JavaScript import filter for an App Endpoint collection. The import filter specifies which documents in the collection are imported by the App Endpoint.",
		Attributes:          attrs,
	}
}
