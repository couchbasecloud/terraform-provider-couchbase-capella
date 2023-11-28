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
type AppService struct {
	Compute        *AppServiceCompute `tfsdk:"compute"`
	OrganizationId types.String       `tfsdk:"organization_id"`
	Description    types.String       `tfsdk:"description"`
	CloudProvider  types.String       `tfsdk:"cloud_provider"`
	Name           types.String       `tfsdk:"name"`
	Id             types.String       `tfsdk:"id"`
	ProjectId      types.String       `tfsdk:"project_id"`
	ClusterId      types.String       `tfsdk:"cluster_id"`
	CurrentState   types.String       `tfsdk:"current_state"`
	Version        types.String       `tfsdk:"version"`
	Audit          types.Object       `tfsdk:"audit"`
	Etag           types.String       `tfsdk:"etag"`
	IfMatch        types.String       `tfsdk:"if_match"`
	Nodes          types.Int64        `tfsdk:"nodes"`
}

// AppServiceCompute depicts the couchbase compute, following are the supported compute combinations
// for CPU and RAM for different cloud providers.
// To learn more, see:
// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
type AppServiceCompute struct {
	// Cpu depicts cpu units (cores).
	Cpu types.Int64 `tfsdk:"cpu"`

	// Ram depicts ram units (GB).
	Ram types.Int64 `tfsdk:"ram"`
}

// NewAppService creates a new instance of an App Service.
func NewAppService(
	appService *appservice.GetAppServiceResponse,
	organizationId, projectId string,
	auditObject basetypes.ObjectValue,
) *AppService {
	newAppService := AppService{
		Id:             types.StringValue(appService.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(appService.Name),
		Description:    types.StringValue(appService.Description),
		CloudProvider:  types.StringValue(appService.CloudProvider),
		Nodes:          types.Int64Value(int64(appService.Nodes)),
		Compute: &AppServiceCompute{
			Cpu: types.Int64Value(appService.Compute.Cpu),
			Ram: types.Int64Value(appService.Compute.Ram),
		},
		ClusterId:    types.StringValue(appService.ClusterId),
		CurrentState: types.StringValue(string(appService.CurrentState)),
		Version:      types.StringValue(appService.Version),
		Audit:        auditObject,
		Etag:         types.StringValue(appService.Etag),
	}
	return &newAppService
}

// Validate is used to verify that IDs have been properly imported.
func (a AppService) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		Id:             a.Id,
	}

	IDs, err := validateSchemaState(state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}

// AppServices defines structure based on the response received from V4 Capella Public API when asked to list app services.
type AppServices struct {
	// OrganizationId The organizationId of the capella.
	OrganizationId types.String `tfsdk:"organization_id"`

	// Data contains the list of resources.
	Data []AppServiceData `tfsdk:"data"`
}

// AppServiceData defines attributes for a single cluster when fetched from the V4 Capella Public API.
type AppServiceData struct {
	Compute        *AppServiceCompute `tfsdk:"compute"`
	Id             types.String       `tfsdk:"id"`
	Name           types.String       `tfsdk:"name"`
	Description    types.String       `tfsdk:"description"`
	CloudProvider  types.String       `tfsdk:"cloud_provider"`
	OrganizationId types.String       `tfsdk:"organization_id"`
	ClusterId      types.String       `tfsdk:"cluster_id"`
	CurrentState   types.String       `tfsdk:"current_state"`
	Version        types.String       `tfsdk:"version"`
	Audit          types.Object       `tfsdk:"audit"`
	Nodes          types.Int64        `tfsdk:"nodes"`
}

// NewAppServiceData creates a new cluster data object.
func NewAppServiceData(
	appService *appservice.GetAppServiceResponse,
	organizationId string,
	auditObject basetypes.ObjectValue,
) *AppServiceData {
	newAppService := AppServiceData{
		Id:             types.StringValue(appService.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		Name:           types.StringValue(appService.Name),
		Description:    types.StringValue(appService.Description),
		CloudProvider:  types.StringValue(appService.CloudProvider),
		Nodes:          types.Int64Value(int64(appService.Nodes)),
		Compute: &AppServiceCompute{
			Cpu: types.Int64Value(appService.Compute.Cpu),
			Ram: types.Int64Value(appService.Compute.Ram),
		},
		ClusterId:    types.StringValue(appService.ClusterId),
		CurrentState: types.StringValue(string(appService.CurrentState)),
		Version:      types.StringValue(appService.Version),
		Audit:        auditObject,
	}
	return &newAppService
}
