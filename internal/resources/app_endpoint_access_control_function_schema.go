package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AccessControlFunctionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This Access Function resource allows you to manage access control and validation functions for App Endpoints in your Capella organization. Access functions are JavaScript functions that specify access control policies applied to documents in collections.",
		Attributes: map[string]schema.Attribute{
			"organization_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Couchbase Capella organization."),
			"project_id":              WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Couchbase Capella project."),
			"cluster_id":              WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Couchbase Capella Cluster."),
			"app_service_id":          WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the Couchbase Capella App Service."),
			"app_endpoint_name":       WithDescription(stringAttribute([]string{required, requiresReplace}), "The Couchbase Capella App endpoint name."),
			"scope":                   WithDescription(stringDefaultAttribute("_default", optional, computed, requiresReplace), "The scope name of the keyspace (app_endpoint_name.scope.collection) containing documents to be processed with this access function."),
			"collection":              WithDescription(stringDefaultAttribute("_default", optional, computed, requiresReplace), "The collection name of the keyspace (app_endpoint_name.scope.collection) containing documents to be processed with this access function."),
			"access_control_function": WithDescription(stringAttribute([]string{required}), "The JavaScript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function."),
		},
	}
}
