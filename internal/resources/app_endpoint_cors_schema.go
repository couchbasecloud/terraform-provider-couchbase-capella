package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var corsBuilder = capellaschema.NewSchemaBuilder("cors")

// CorsSchema returns the schema for the CORS resource.
func CorsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", corsBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "project_id", corsBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "cluster_id", corsBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_service_id", corsBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "app_endpoint_name", corsBuilder, stringAttribute([]string{required, requiresReplace}, validator.String(stringvalidator.LengthAtLeast(1))))
	capellaschema.AddAttr(attrs, "login_origin", corsBuilder, stringSetAttribute(optional))
	capellaschema.AddAttr(attrs, "headers", corsBuilder, stringSetAttribute(optional))
	capellaschema.AddAttr(attrs, "max_age", corsBuilder, int64Attribute(optional, computed))
	capellaschema.AddAttr(attrs, "disabled", corsBuilder, boolAttribute(optional, computed))

	attrs["origin"] = &schema.SetAttribute{
		ElementType: types.StringType,
		Required:    true,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
		},
	}

	return schema.Schema{
		MarkdownDescription: "Manages CORS (Cross-Origin Resource Sharing) configuration for App Endpoints in Couchbase Capella.",
		Attributes:          attrs,
	}
}
