package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AccessControlFunction represents the Terraform schema for an Access Control and Validation function
// associated with an App Endpoint. This schema is used for the getAccessFunction API endpoint
// which retrieves the JavaScript function used to specify access control policies.
type AccessControlFunction struct {
	// OrganizationId is the ID of the organization to which the App Endpoint belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the App Endpoint belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster to which the App Endpoint belongs.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the App Service to which the App Endpoint belongs.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpoint is the name of the App Endpoint the access function is associated with.
	AppEndpoint types.String `tfsdk:"app_endpoint"`

	// Scope is scope where the collection resides.
	Scope types.String `tfsdk:"scope"`

	// Collection is the collection the access function will operate on.
	Collection types.String `tfsdk:"collection"`

	// Function is the JavaScript function that is used to specify the access control policies
	// to be applied to documents in this collection. Every document update is processed by this function.
	// The default access control function is 'function(doc){channel(doc.channels);}'
	// for the default collection and 'function(doc){channel(collectionName);}' for named collections.
	AccessControlFunction types.String `tfsdk:"access_control_function"`
}

// ValidateState validates base identifiers using the shared validateSchemaState helper,
// enabling consistent terraform import parsing similar to other resources.
func (a *AccessControlFunction) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  a.OrganizationId,
		ProjectId:       a.ProjectId,
		ClusterId:       a.ClusterId,
		AppServiceId:    a.AppServiceId,
		AppEndpointName: a.AppEndpoint,
		ScopeName:       a.Scope,
		CollectionName:  a.Collection,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs, nil
}
