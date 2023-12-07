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
	// Compute is the CPU and RAM configuration of the app service.
	Compute *AppServiceCompute `tfsdk:"compute"`

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

	// IfMatch is a precondition header that specifies the entity tag of a resource.
	IfMatch types.String `tfsdk:"if_match"`

	// Nodes is the number of nodes configured for the app service.
	Nodes types.Int64 `tfsdk:"nodes"`
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
	// Compute is the CPU and RAM configuration of the app service.
	Compute *AppServiceCompute `tfsdk:"compute"`

	// Id is a UUID of the app service.
	Id types.String `tfsdk:"id"`

	// Name is the name of the app service, the name of the app service should follow this naming criteria:
	// An app service name should have at least 2 characters and up to 256 characters.
	Name types.String `tfsdk:"name"`

	// Description is the description for the app service (up to 256 characters).
	Description types.String `tfsdk:"description"`

	// CloudProvider is the cloud provider where the app service will be hosted.
	// To learn more, see:
	// [AWS] https://docs.couchbase.com/cloud/reference/aws.html
	// [GCP] https://docs.couchbase.com/cloud/reference/gcp.html
	// [Azure] https://docs.couchbase.com/cloud/reference/azure.html
	CloudProvider types.String `tfsdk:"cloud_provider"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ClusterId is the clusterId of the cluster.
	ClusterId types.String `tfsdk:"cluster_id"`

	// CurrentState defines the current state of app service.
	CurrentState types.String `tfsdk:"current_state"`

	// Version defines the version of the app service server
	Version types.String `tfsdk:"version"`

	// Audit represents all audit-related fields. It is of types.Object type to avoid conversion error for a nested field.
	Audit types.Object `tfsdk:"audit"`

	// Nodes is the number of nodes configured for the app service.
	Nodes types.Int64 `tfsdk:"nodes"`
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
