package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// ImportFilterSchema defines the Terraform resource schema for managing App Endpoint Import Filters.
func ImportFilterSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This Import Filter resource allows you to manage the JavaScript import filter for an App Endpoint collection. The import filter specifies which documents in the collection are imported by the App Endpoint.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"keyspace":        WithDescription(stringAttribute([]string{required, requiresReplace}), "The keyspace of the collection, in the format <app_endpoint_name>.<scope_name>.<collection_name>. If only an App Endpoint name is provided this will be interpreted as \"endpoint1._default._default\". If only an App Endpoint name and collection name are provided these will interpreted as a named collection within the default scope, for example \"endpoint1.collection1\" will be interpreted as \"endpoint1._default.collection1\"."),
			"import_filter":   WithDescription(stringAttribute([]string{required}), "The JavaScript function that specifies which documents to import from this collection. By default, all documents are imported."),
		},
	}
}
