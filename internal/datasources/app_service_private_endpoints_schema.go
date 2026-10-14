package datasources

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var appServicePrivateEndpointsBuilder = capellaschema.NewSchemaBuilder("appServicePrivateEndpoints", "PrivateEndpoint")

// appServicePrivateEndpointStatuses are the accepted values for the status filter.
var appServicePrivateEndpointStatuses = []string{"pending", "pendingAcceptance", "linked", "rejected", "failed", "unrecognized"}

// AppServicePrivateEndpointsSchema returns the schema for the app_service_private_endpoints data source.
func AppServicePrivateEndpointsSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", appServicePrivateEndpointsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", appServicePrivateEndpointsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", appServicePrivateEndpointsBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "app_service_id", appServicePrivateEndpointsBuilder, requiredUUIDString())

	// Set directly, not through AddAttr, so the OpenAPI description of the endpoint status does not
	// replace the filter description.
	attrs["status"] = schema.SetAttribute{
		MarkdownDescription: "Filters the private endpoints by status. If not set, all private endpoints are returned. Valid values: `" +
			strings.Join(appServicePrivateEndpointStatuses, "`, `") + "`.",
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			setvalidator.ValueStringsAre(stringvalidator.OneOf(appServicePrivateEndpointStatuses...)),
		},
	}

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServicePrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "service_name", appServicePrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", appServicePrivateEndpointsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", appServicePrivateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The list of private endpoints for the App Service.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the private endpoints for an App Service, optionally filtered by status.",
		Attributes:          attrs,
	}
}
