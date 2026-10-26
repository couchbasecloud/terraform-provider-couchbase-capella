package datasources

import (
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

	// status filters the list in the provider, as the list endpoint has no status query parameter.
	// When omitted, private endpoints in every state are returned.
	capellaschema.AddAttr(attrs, "status", appServicePrivateEndpointsBuilder, &schema.SetAttribute{
		Optional:    true,
		ElementType: types.StringType,
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			setvalidator.ValueStringsAre(stringvalidator.OneOf(appServicePrivateEndpointStatuses...)),
		},
	})

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", appServicePrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "service_name", appServicePrivateEndpointsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", appServicePrivateEndpointsBuilder, computedString())

	capellaschema.AddAttr(attrs, "data", appServicePrivateEndpointsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve the private endpoints for an App Service, optionally filtered by status.",
		Attributes:          attrs,
	}
}
