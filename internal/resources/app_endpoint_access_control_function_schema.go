package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var accessControlFunctionBuilder = capellaschema.NewSchemaBuilder("accessControlFunction")

func AccessControlFunctionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", accessControlFunctionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", accessControlFunctionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", accessControlFunctionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_service_id", accessControlFunctionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "app_endpoint_name", accessControlFunctionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "scope", accessControlFunctionBuilder, stringDefaultAttribute("_default", optional, computed, requiresReplace))
	capellaschema.AddAttr(attrs, "collection", accessControlFunctionBuilder, stringDefaultAttribute("_default", optional, computed, requiresReplace))
	capellaschema.AddAttr(attrs, "access_control_function", accessControlFunctionBuilder, stringAttribute([]string{required}))

	return schema.Schema{
		MarkdownDescription: "This Access Function resource allows you to manage access control and validation functions for App Endpoints in your Capella organization. Access functions are JavaScript functions that specify access control policies applied to documents in collections.",
		Attributes:          attrs,
	}
}
