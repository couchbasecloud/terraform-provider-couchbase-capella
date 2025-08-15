package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func AccessFunctionSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This Access Function resource allows you to manage access control and validation functions for App Endpoints in your Capella organization. Access functions are JavaScript functions that specify access control policies applied to documents in collections.",
		Attributes: map[string]schema.Attribute{
			"organization_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the organization."),
			"project_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the project."),
			"cluster_id":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the cluster."),
			"app_service_id":  WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app service."),
			"app_endpoint_id": WithDescription(stringAttribute([]string{required, requiresReplace}), "The GUID4 ID of the app endpoint."),
			"scope":           WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the scope containing the collection."),
			"collection":      WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the collection for which the access function is defined."),
			"function":        WithDescription(stringAttribute([]string{required}), "The JavaScript function that is used to specify the access control policies to be applied to documents in this collection. Every document update is processed by this function. The default access control function is 'function(doc){channel(doc.channels);}' for the default collection and 'function(doc){channel(collectionName);}' for named collections."),
			"audit":           computedAuditAttribute(),
		},
	}
}
