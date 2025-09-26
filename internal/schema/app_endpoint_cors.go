package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Cors defines the attributes for CORS configuration in Terraform.
type Cors struct {
	// OrganizationId is the Capella tenant id associated with the App Endpoint.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the id of the Capella project associated with the App Endpoint.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the id of the Capella cluster associated with the App Endpoint.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the id of the Capella App Service associated with the App Endpoint.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpointName is the name of the App Endpoint associated with the CORS configuration.
	AppEndpointName types.String `tfsdk:"app_endpoint_name"`

	// Origin is a set of allowed origins for CORS.
	Origin types.Set `tfsdk:"origin"`

	// LoginOrigin is a set of allowed login origins for CORS.
	LoginOrigin types.Set `tfsdk:"login_origin"`

	// Headers is a set of allowed headers for CORS.
	Headers types.Set `tfsdk:"headers"`

	// MaxAge specifies the duration (in seconds) for which the results of a preflight request can be cached.
	MaxAge types.Int64 `tfsdk:"max_age"`

	// Disabled indicates whether CORS is disabled.
	Disabled types.Bool `tfsdk:"disabled"`
}

// Validate validates the CORS state and returns the IDs as a map.
// It returns organizationId, projectId, clusterId, appServiceId, and appEndpointName.
func (c *Cors) Validate() (organizationId, projectId, clusterId, appServiceId, appEndpointName string, err error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  c.OrganizationId,
		ProjectId:       c.ProjectId,
		ClusterId:       c.ClusterId,
		AppServiceId:    c.AppServiceId,
		AppEndpointName: c.AppEndpointName,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs[OrganizationId], IDs[ProjectId], IDs[ClusterId], IDs[AppServiceId], IDs[AppEndpointName], nil
}
