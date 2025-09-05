package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// CorsSchema returns the schema for the CORS resource.
func CorsSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages CORS (Cross-Origin Resource Sharing) configuration for App Endpoints in Couchbase Capella.",
		Attributes: map[string]schema.Attribute{
			"organization_id":   WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The ID of the Capella organization."),
			"project_id":        WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The ID of the Capella project."),
			"cluster_id":        WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The ID of the Capella cluster."),
			"app_service_id":    WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The ID of the Capella App Service."),
			"app_endpoint_name": WithDescription(stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))), "The name of the App Endpoint."),
			"origin":            WithDescription(stringSetAttribute(required), "Set of allowed origins for CORS. Use ['*'] to allow access from everywhere."),
			"login_origin":      WithDescription(stringSetAttribute(optional), "Set of allowed login origins for CORS."),
			"headers":           WithDescription(stringSetAttribute(optional), "Set of allowed headers for CORS."),
			"max_age":           WithDescription(int64Attribute(optional, computed), "Specifies the duration (in seconds) for which the results of a preflight request can be cached."),
			"disabled":          WithDescription(boolAttribute(optional, computed), "Indicates whether CORS is disabled."),
		},
	}
}
