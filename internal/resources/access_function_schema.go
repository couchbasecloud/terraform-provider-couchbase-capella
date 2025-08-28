package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AccessFunctionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This Access Function resource allows you to manage access control and validation functions for App Endpoints in your Capella organization. Access functions are JavaScript functions that specify access control policies applied to documents in collections.",
		Attributes: map[string]schema.Attribute{
			"organization_id":         WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":              WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":              WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":          WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"keyspace":                WithDescription(stringAttribute([]string{required, requiresReplace}), "The keyspace of the access function."),
			"access_control_function": WithDescription(stringAttribute([]string{required}), "The JavaScript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function."),
		},
	}
}
