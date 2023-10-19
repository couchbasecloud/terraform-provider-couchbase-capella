package schema

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"strings"
	"terraform-provider-capella/internal/api/appservice"
	"terraform-provider-capella/internal/errors"
)

// AppService defines the response as received from V4 Capella Public API when asked to create a new app service.
type AppService struct {
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

	// Nodes is the number of nodes configured for the app service.
	Nodes types.Int64 `tfsdk:"nodes"`

	// Compute is the CPU and RAM configuration of the app service.
	Compute Compute `tfsdk:"compute"`

	// OrganizationId is the organizationId of the capella tenant.
	OrganizationId types.String `tfsdk:"organization_id"`

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
}

// NewAppService creates a new instance of an App Service
func NewAppService(
	appService *appservice.GetAppServiceResponse,
	auditObject basetypes.ObjectValue,
) *AppService {
	newAppService := AppService{
		Id:            types.StringValue(appService.Id.String()),
		Name:          types.StringValue(appService.Name),
		Description:   types.StringValue(appService.Description),
		CloudProvider: types.StringValue(appService.CloudProvider),
		Nodes:         types.Int64Value(int64(appService.Nodes)),
		Compute: Compute{
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

// Validate will split the IDs by a delimiter i.e. comma , in case a terraform import CLI is invoked.
// The format of the terraform import CLI would include the IDs as follows -
// `terraform import capella_bucket.new_bucket id=<uuid>,cluster_id=<uuid>,project_id=<uuid>,organization_id=<uuid>`
func (a AppService) Validate() (appServiceId, clusterId, projectId, organizationId string, err error) {

	const (
		idDelimiter       = ","
		organizationIdSep = "organization_id="
		projectIdSep      = "project_id="
		clusterIdSep      = "cluster_id="
		appServiceIdSep   = "id="
	)

	organizationId = a.OrganizationId.ValueString()
	projectId = a.ProjectId.ValueString()
	clusterId = a.ClusterId.ValueString()
	appServiceId = a.Id.ValueString()
	var found bool

	// check if the id is a comma separated string of multiple IDs, usually passed during the terraform import CLI
	if a.OrganizationId.IsNull() {
		strs := strings.Split(a.Id.ValueString(), idDelimiter)
		if len(strs) != 4 {
			err = errors.ErrIdMissing
			return
		}
		_, appServiceId, found = strings.Cut(strs[0], appServiceIdSep)
		if !found {
			err = errors.ErrDatabaseCredentialIdMissing
			return
		}

		_, clusterId, found = strings.Cut(strs[1], clusterIdSep)
		if !found {
			err = errors.ErrClusterIdMissing
			return
		}

		_, projectId, found = strings.Cut(strs[2], projectIdSep)
		if !found {
			err = errors.ErrProjectIdMissing
			return
		}

		_, organizationId, found = strings.Cut(strs[3], organizationIdSep)
		if !found {
			err = errors.ErrOrganizationIdMissing
			return
		}
	}

	if appServiceId == "" {
		err = errors.ErrAppServiceIdCannotBeEmpty
		return
	}

	if clusterId == "" {
		err = errors.ErrClusterIdCannotBeEmpty
		return
	}

	if projectId == "" {
		err = errors.ErrProjectIdCannotBeEmpty
		return
	}

	if organizationId == "" {
		err = errors.ErrOrganizationIdCannotBeEmpty
		return
	}

	return appServiceId, clusterId, projectId, organizationId, nil
}

func (a AppService) Validate2() (map[Attr]string, error) {
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
