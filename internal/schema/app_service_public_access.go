package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServicePublicAccess defines the Terraform state for the app service public access resource.
type AppServicePublicAccess struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster on which the app service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the app service whose public access to manage.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// PublicAccess is the requested public access mode: restricted or unrestricted.
	PublicAccess types.String `tfsdk:"public_access"`

	// PublicAccessStatus is the public access mode, including transitional and failed values.
	PublicAccessStatus types.String `tfsdk:"public_access_status"`
}

// NewAppServicePublicAccess creates the Terraform state from the public access that Capella returns. A failed
// change reads as the mode that it requested, so that a failed change does not cause a plan diff.
func NewAppServicePublicAccess(
	organizationId, projectId, clusterId, appServiceId string,
	publicAccess appservice.PublicAccess,
) *AppServicePublicAccess {
	return &AppServicePublicAccess{
		OrganizationId:     types.StringValue(organizationId),
		ProjectId:          types.StringValue(projectId),
		ClusterId:          types.StringValue(clusterId),
		AppServiceId:       types.StringValue(appServiceId),
		PublicAccess:       types.StringValue(string(publicAccess.Mode())),
		PublicAccessStatus: types.StringValue(string(publicAccess)),
	}
}

// Validate validates the resource state and returns parsed IDs.
func (a *AppServicePublicAccess) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
