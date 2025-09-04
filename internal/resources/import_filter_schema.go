package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ImportFilterSchema defines the Terraform resource schema for managing App Endpoint Import Filters.
func ImportFilterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This Import Filter resource allows you to manage the JavaScript import filter for an App Endpoint collection. The import filter specifies which documents in the collection are imported by the App Endpoint.",
		Attributes: map[string]schema.Attribute{
			"organization_id":   WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":        WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":        WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":    WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"app_endpoint_name": WithDescription(stringAttribute([]string{required, requiresReplace}), "The app endpoint name."),
			"scope":             WithDescription(stringDefaultAttribute("_default", optional, computed, requiresReplace), "Scope is the scope within the keyspace where the collection resides."),
			"collection":        WithDescription(stringDefaultAttribute("_default", optional, computed, requiresReplace), "Collection is the collection within the scope where the documents to be imported reside."),
			"import_filter":     WithDescription(stringAttribute([]string{required}), "The JavaScript function that specifies which documents to import from this collection. By default, all documents are imported."),
		},
	}
}
