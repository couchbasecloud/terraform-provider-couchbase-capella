package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ImportFilter represents the Terraform schema for an Import Filter function
// associated with an App Endpoint. It mirrors the AccessFunction schema shape
// and is used for the getImportFilter/putImportFilter API endpoints.
// See: https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/App-Endpoints/operation/putImportFilter
type ImportFilter struct {
	// OrganizationId is the ID of the organization to which the App Endpoint belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Endpoint belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Endpoint belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the App Endpoint belongs.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpointName is the name of the App Endpoint.
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`

	// Scope is the scope within the keyspace where the collection resides.
	Scope types.String `tfsdk:"scope"`

	// Collection is the collection within the scope where the documents to be
	// imported reside.
	Collection types.String `tfsdk:"collection"`

	// ImportFilter is the JavaScript function used to specify the documents in this
	// collection that are to be imported by the App Endpoint. By default, all
	// documents in the corresponding collection are imported.
	ImportFilter types.String `tfsdk:"import_filter"`
}

// AttributeTypes returns a map of attribute types for the ImportFilter schema.
func (i ImportFilter) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":   types.StringType,
		"project_id":        types.StringType,
		"cluster_id":        types.StringType,
		"app_service_id":    types.StringType,
		"app_endpoint_name": types.StringType,
		"scope":             types.StringType,
		"collection":        types.StringType,
		"import_filter":     types.StringType,
	}
}
