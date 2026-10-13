package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppServiceAdminUser defines the attributes for an App Service admin user.
type AppServiceAdminUser struct {
	// Id is the unique identifier of the admin user.
	Id types.String `tfsdk:"id"`

	// OrganizationId is the Capella tenant id associated with the App Service.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the id of the Capella project associated with the App Service.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the id of the Capella cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the id of the Capella App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// Name is the name of the admin user.
	Name types.String `tfsdk:"name"`

	// Password is the password of the admin user. It is never returned by the API.
	Password types.String `tfsdk:"password"`

	// AccessAllEndpoints gives the user access to all App Endpoints.
	AccessAllEndpoints types.Bool `tfsdk:"access_all_endpoints"`

	// Endpoints is the set of App Endpoints the user has access to.
	Endpoints types.Set `tfsdk:"endpoints"`

	// Audit contains the audit information for the admin user.
	Audit types.Object `tfsdk:"audit"`
}

// Validate validates the state and returns organizationId, projectId, clusterId, appServiceId and userId.
// It also parses the import string when the resource is imported.
func (a *AppServiceAdminUser) Validate() (organizationId, projectId, clusterId, appServiceId, userId string, err error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to validate resource state: %s", err)
	}

	return IDs[OrganizationId], IDs[ProjectId], IDs[ClusterId], IDs[AppServiceId], IDs[Id], nil
}
