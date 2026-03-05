package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppServiceLogStreamingActivationStatus defines the Terraform state for the
// app service log streaming activation status resource.
type AppServiceLogStreamingActivationStatus struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the app service is deployed on.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the app service for which log streaming activation state is managed for.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// State is the activation state of log streaming: "paused" or "enabled".
	State types.String `tfsdk:"state"`
}

// NewAppServiceLogStreamingActivationStatus creates a new AppServiceLogStreamingActivationStatus
// with properly populated ID fields and the given config state. This ensures that after an import,
// the individual ID fields are correctly set in the Terraform state rather than remaining as the
// raw composite import string.
func NewAppServiceLogStreamingActivationStatus(
	organizationId, projectId, clusterId, appServiceId string,
	state apigen.GetLogStreamingResponseConfigState,
) *AppServiceLogStreamingActivationStatus {
	return &AppServiceLogStreamingActivationStatus{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
		State:          types.StringValue(string(state)),
	}
}

// Validate validates the resource state and returns parsed IDs.
func (a *AppServiceLogStreamingActivationStatus) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}
