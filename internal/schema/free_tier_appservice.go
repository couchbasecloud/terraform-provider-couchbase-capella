package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppService defines the response as received from V4 Capella Public API when asked to create a new app service.
// To learn more about App Services, see https://docs.couchbase.com/cloud/app-services/index.html
type FreeTierAppService struct {
	// Compute is the CPU and RAM configuration of the app service.
	Compute types.Object `tfsdk:"compute"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Description is the description for the app service (up to 256 characters).
	Description types.String `tfsdk:"description"`

	// CloudProvider is the cloud provider where the app service will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider types.String `tfsdk:"cloud_provider"`

	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name types.String `tfsdk:"name"`

	// Id is a UUID of the app service.
	Id types.String `tfsdk:"id"`

	// ProjectId is the projectId of the cluster.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the clusterId of the cluster.
	ClusterId types.String `tfsdk:"cluster_id"`

	// CurrentState defines the current state of app service.
	CurrentState types.String `tfsdk:"current_state"`

	// Version defines the version of the app service server
	Version types.String `tfsdk:"version"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`

	// Etag represents the version of the document.
	Etag types.String `tfsdk:"etag"`

	// Nodes is the number of nodes configured for the app service.
	Nodes types.Int64 `tfsdk:"nodes"`

	// Plan is the plan of the free-tier app service.
	Plan types.String `tfsdk:"plan"`
}

// NewFreeTierAppService creates a new instance of an App Service.
func NewFreeTierAppService(
	appService *appservice.GetAppServiceResponse,
	organizationId, projectId string,
	auditObject basetypes.ObjectValue,
	computeObject basetypes.ObjectValue,
) *FreeTierAppService {
	newFreeTierAppService := FreeTierAppService{
		Id:             types.StringValue(appService.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(appService.Name),
		Description:    types.StringValue(appService.Description),
		CloudProvider:  types.StringValue(appService.CloudProvider),
		Nodes:          types.Int64Value(int64(appService.Nodes)),
		Compute:        computeObject,
		ClusterId:      types.StringValue(appService.ClusterId),
		CurrentState:   types.StringValue(string(appService.CurrentState)),
		Version:        types.StringValue(appService.Version),
		Audit:          auditObject,
		Etag:           types.StringValue(appService.Etag),
		Plan:           types.StringValue(appService.Plan),
	}
	return &newFreeTierAppService
}

// Validate is used to verify that IDs have been properly imported.
func (f FreeTierAppService) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: f.OrganizationId,
		ProjectId:      f.ProjectId,
		ClusterId:      f.ClusterId,
		Id:             f.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}
