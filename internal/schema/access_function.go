package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AccessFunction represents the Terraform schema for an Access Control and Validation function
// associated with an App Endpoint. This schema is used for the getAccessFunction API endpoint
// which retrieves the JavaScript function used to specify access control policies.
type AccessFunction struct {
	// OrganizationId is the ID of the organization to which the App Endpoint belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Endpoint belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Endpoint belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the App Endpoint belongs.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpointId is the ID of the App Endpoint.
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`

	// Scope is the name of the scope containing the collection.
	Scope types.String `tfsdk:"scope"`

	// Collection is the name of the collection for which the access function is defined.
	Collection types.String `tfsdk:"collection"`

	// Function is the JavaScript function that is used to specify the access control policies
	// to be applied to documents in this collection. Every document update is processed by this function.
	// The default access control function is 'function(doc){channel(doc.channels);}'
	// for the default collection and 'function(doc){channel(collectionName);}' for named collections.
	AccessControlFunction types.String `tfsdk:"access_control_function"`
}

// AttributeTypes returns a map of attribute types for the AccessFunction schema.
// This method is required to implement the types.ObjectTypable interface for Terraform.
func (a AccessFunction) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":         types.StringType,
		"project_id":              types.StringType,
		"cluster_id":              types.StringType,
		"app_service_id":          types.StringType,
		"app_endpoint_name":       types.StringType,
		"scope":                   types.StringType,
		"collection":              types.StringType,
		"access_control_function": types.StringType,
	}
}

func (a *AccessFunction) Validate() (map[string]string, error) {
	ids := make(map[string]string)

	if a.OrganizationId.IsNull() || a.OrganizationId.IsUnknown() {
		return nil, fmt.Errorf("organization_id cannot be empty")
	}
	ids["organization_id"] = a.OrganizationId.ValueString()

	if a.ProjectId.IsNull() || a.ProjectId.IsUnknown() {
		return nil, fmt.Errorf("project_id cannot be empty")
	}
	ids["project_id"] = a.ProjectId.ValueString()

	if a.ClusterId.IsNull() || a.ClusterId.IsUnknown() {
		return nil, fmt.Errorf("cluster_id cannot be empty")
	}
	ids["cluster_id"] = a.ClusterId.ValueString()

	if a.AppServiceId.IsNull() || a.AppServiceId.IsUnknown() {
		return nil, fmt.Errorf("app_service_id cannot be empty")
	}
	ids["app_service_id"] = a.AppServiceId.ValueString()

	if a.AppEndpointName.IsNull() || a.AppEndpointName.IsUnknown() {
		return nil, fmt.Errorf("app_endpoint_name cannot be empty")
	}
	ids["app_endpoint_name"] = a.AppEndpointName.ValueString()

	if a.Scope.IsNull() || a.Scope.IsUnknown() {
		return nil, fmt.Errorf("scope cannot be empty")
	}
	ids["scope"] = a.Scope.ValueString()

	if a.Collection.IsNull() || a.Collection.IsUnknown() {
		return nil, fmt.Errorf("collection cannot be empty")
	}
	ids["collection"] = a.Collection.ValueString()

	if a.AccessControlFunction.IsNull() || a.AccessControlFunction.IsUnknown() {
		return nil, fmt.Errorf("access_control_function cannot be empty")
	}
	ids["access_control_function"] = a.AccessControlFunction.ValueString()

	return ids, nil
}
