package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

var (
	_ resource.Resource                = &FreeTierCluster{}
	_ resource.ResourceWithConfigure   = &FreeTierCluster{}
	_ resource.ResourceWithImportState = &FreeTierCluster{}
)

type FreeTierCluster struct {
	*providerschema.Data
}

func NewFreeTierCluster() resource.Resource {
	return &FreeTierCluster{}
}

func (f *FreeTierCluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_free_tier_cluster"

}

func (f *FreeTierCluster) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FreeTierClusterSchema()
}

func (f *FreeTierCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	orgUUID, _ := uuid.Parse(plan.OrganizationId.ValueString())
	projUUID, _ := uuid.Parse(plan.ProjectId.ValueString())

	freeTierClusterCreateRequest := apigen.CreateFreeTierClusterRequest{
		Name: plan.Name.ValueString(),
		CloudProvider: apigen.CloudProvider{
			Cidr: func() *string {
				if plan.CloudProvider.Cidr.IsNull() || plan.CloudProvider.Cidr.IsUnknown() {
					return nil
				}
				v := plan.CloudProvider.Cidr.ValueString()
				return &v
			}(),
			Region: plan.CloudProvider.Region.ValueString(),
			Type:   apigen.CloudProviderType(plan.CloudProvider.Type.ValueString()),
		},
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		d := plan.Description.ValueString()
		freeTierClusterCreateRequest.Description = &d
	}

	res, err := f.ClientV2.CreateFreeTierClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), freeTierClusterCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating cluster",
			errors.ErrorMessageWhileFreeTierClusterCreation.Error()+err.Error(),
		)
		return
	}
	if res.JSON202 == nil {
		response.Diagnostics.AddError("Error creating cluster", "unexpected status: "+res.Status())
		return
	}

	diags = response.State.Set(ctx, initializePendingFreeTierClusterWithPlanAndId(plan, res.JSON202.Id.String()))
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	clusterResp, err := f.checkForFreeTierClusterDesiredStatus(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), res.JSON202.Id.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"error getting cluster status",
			"failed to get cluster status, please refresh state later "+err.Error(),
		)
		return
	}

	if clusterResp.CurrentState != apigen.CurrentState("healthy") {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			fmt.Sprintf("Could not create cluster id %s, as current Cluster state: %s", clusterResp.Id, clusterResp.CurrentState),
		)

	}
	morphedState, err := f.morphFreeTierClusterRespToTerraformObj(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), clusterResp)
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error fetching the cluster info",
			errors.ErrorMessageAfterFreeTierClusterCreationInitiation.Error()+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, &morphedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f *FreeTierCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state providerschema.FreeTierCluster
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()

	if err != nil {
		response.Diagnostics.AddError(
			"error reading free tier cluster",
			"could not read free tier cluster "+state.Id.String()+": "+err.Error(),
		)
	}
	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	//get refreshed cluster values from capella.
	refreshedState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error Reading Capella Cluster",
			"Could Not Read Capella Cluster "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	diags = response.State.Set(ctx, &refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f *FreeTierCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)

	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := plan.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)

	updateReq := apigen.UpdateFreeTierClusterRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	cluUUID, _ := uuid.Parse(clusterId)
	_, err = f.ClientV2.UpdateFreeTierClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID), nil, updateReq)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	currentState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating free tier cluster",
			"Could not update cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data.
	diags = response.State.Set(ctx, currentState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (f *FreeTierCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Retrieve values from state.
	var state providerschema.FreeTierCluster
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting free tier cluster",
			"Could not delete cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.Id]
	)

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)

	cluUUID, _ := uuid.Parse(clusterId)
	_, err = f.ClientV2.DeleteFreeTierClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID))
	if err != nil {
		response.Diagnostics.AddError(
			"Error Deleting Free Tier Cluster",
			"Could not delete cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	freeTierClusterResp, err := f.checkForFreeTierClusterDesiredStatus(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.Id.ValueString())

	if err != nil {
		response.Diagnostics.AddError(
			"Error Deleting Capella Cluster",
			"Could not delete cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	if freeTierClusterResp.CurrentState == apigen.CurrentState("destroyFailed") {
		response.Diagnostics.AddError(
			"Error Deleting Free Tier Cluster",
			"Could not delete cluster id "+state.Id.String()+": cluster in destroy failed state",
		)
	}
}

func (f *FreeTierCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	/// Retrieve import ID and save to id attribute.
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

func (f *FreeTierCluster) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}
	data, ok := request.ProviderData.(*providerschema.Data)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected resource configure type",
			fmt.Sprintf("expected *providerschema.FreeTierCluster, got %T", request.ProviderData),
		)
		return
	}
	f.Data = data
}

// initializePendingClusterWithPlanAndId initializes an instance of providerschema.Cluster.
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingFreeTierClusterWithPlanAndId(plan providerschema.FreeTierCluster, id string) providerschema.FreeTierCluster {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}

	plan.EnablePrivateDNSResolution = types.BoolNull()
	plan.CouchbaseServer = types.ObjectNull(providerschema.CouchbaseServer{}.AttributeTypes())
	plan.AppServiceId = types.StringNull()
	plan.ConnectionString = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	plan.Availability = types.ObjectNull(providerschema.Availability{}.AttributeTypes())
	plan.CmekId = types.StringNull()

	plan.ServiceGroups = types.SetNull(types.ObjectType{}.WithAttributeTypes(providerschema.ServiceGroupAttributeTypes()))
	plan.Support = types.ObjectNull(providerschema.Support{}.AttributeTypes())
	plan.Etag = types.StringNull()
	return plan
}

// checkFreeTierClusterStatus monitors the status of a cluster creation/update/deletion.
func (f *FreeTierCluster) checkForFreeTierClusterDesiredStatus(ctx context.Context, organizationId, projectId, ClusterId string) (*struct {
	AppServiceId               *openapi_types.UUID       `json:"appServiceId,omitempty"`
	Audit                      apigen.CouchbaseAuditData `json:"audit"`
	Availability               apigen.Availability       `json:"availability"`
	CloudProvider              apigen.CloudProvider      `json:"cloudProvider"`
	CmekId                     *string                   `json:"cmekId,omitempty"`
	ConfigurationType          apigen.ConfigurationType  `json:"configurationType"`
	ConnectionString           string                    `json:"connectionString"`
	CouchbaseServer            apigen.CouchbaseServer    `json:"couchbaseServer"`
	CurrentState               apigen.CurrentState       `json:"currentState"`
	Description                string                    `json:"description"`
	EnablePrivateDNSResolution *bool                     `json:"enablePrivateDNSResolution,omitempty"`
	Id                         openapi_types.UUID        `json:"id"`
	Name                       string                    `json:"name"`
	ServiceGroups              []apigen.ServiceGroup     `json:"serviceGroups"`
	Support                    struct {
		Plan     apigen.GetFreeTierCluster200SupportPlan `json:"plan"`
		Timezone apigen.SupportTimezone                  `json:"timezone"`
	} `json:"support"`
}, error) {
	var (
		clusterResp *struct {
			AppServiceId               *openapi_types.UUID       `json:"appServiceId,omitempty"`
			Audit                      apigen.CouchbaseAuditData `json:"audit"`
			Availability               apigen.Availability       `json:"availability"`
			CloudProvider              apigen.CloudProvider      `json:"cloudProvider"`
			CmekId                     *string                   `json:"cmekId,omitempty"`
			ConfigurationType          apigen.ConfigurationType  `json:"configurationType"`
			ConnectionString           string                    `json:"connectionString"`
			CouchbaseServer            apigen.CouchbaseServer    `json:"couchbaseServer"`
			CurrentState               apigen.CurrentState       `json:"currentState"`
			Description                string                    `json:"description"`
			EnablePrivateDNSResolution *bool                     `json:"enablePrivateDNSResolution,omitempty"`
			Id                         openapi_types.UUID        `json:"id"`
			Name                       string                    `json:"name"`
			ServiceGroups              []apigen.ServiceGroup     `json:"serviceGroups"`
			Support                    struct {
				Plan     apigen.GetFreeTierCluster200SupportPlan `json:"plan"`
				Timezone apigen.SupportTimezone                  `json:"timezone"`
			} `json:"support"`
		}
		err error
	)

	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)

	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)

	for {
		select {
		case <-ctx.Done():
			return clusterResp, fmt.Errorf("cluster creation status transition timed out after initiation, unexpected error: %w", err)
		case <-ticker.C:
			cluUUID, _ := uuid.Parse(ClusterId)
			res, err := f.ClientV2.GetFreeTierClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID))
			if err != nil {
				return clusterResp, err
			}
			if res.JSON200 != nil {
				clusterResp = res.JSON200
				if clusterResp.CurrentState == apigen.CurrentState("healthy") || clusterResp.CurrentState == apigen.CurrentState("destroyFailed") {
					tflog.Info(ctx, "cluster status is in final state")
					return clusterResp, nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			}
		}
	}
}

// getFreeTierCluster retrieves cluster information using v2 API.
func (f *FreeTierCluster) getFreeTierCluster(ctx context.Context, organizationId, projectId, clusterId string,
) (*struct {
	AppServiceId               *openapi_types.UUID       `json:"appServiceId,omitempty"`
	Audit                      apigen.CouchbaseAuditData `json:"audit"`
	Availability               apigen.Availability       `json:"availability"`
	CloudProvider              apigen.CloudProvider      `json:"cloudProvider"`
	CmekId                     *string                   `json:"cmekId,omitempty"`
	ConfigurationType          apigen.ConfigurationType  `json:"configurationType"`
	ConnectionString           string                    `json:"connectionString"`
	CouchbaseServer            apigen.CouchbaseServer    `json:"couchbaseServer"`
	CurrentState               apigen.CurrentState       `json:"currentState"`
	Description                string                    `json:"description"`
	EnablePrivateDNSResolution *bool                     `json:"enablePrivateDNSResolution,omitempty"`
	Id                         openapi_types.UUID        `json:"id"`
	Name                       string                    `json:"name"`
	ServiceGroups              []apigen.ServiceGroup     `json:"serviceGroups"`
	Support                    struct {
		Plan     apigen.GetFreeTierCluster200SupportPlan `json:"plan"`
		Timezone apigen.SupportTimezone                  `json:"timezone"`
	} `json:"support"`
}, error) {
	orgUUID, _ := uuid.Parse(organizationId)
	projUUID, _ := uuid.Parse(projectId)
	cluUUID, _ := uuid.Parse(clusterId)
	res, err := f.ClientV2.GetFreeTierClusterWithResponse(ctx, apigen.OrganizationId(orgUUID), apigen.ProjectId(projUUID), apigen.ClusterId(cluUUID))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	if res.JSON200 == nil {
		return nil, fmt.Errorf("%s: unexpected status %s", errors.ErrExecutingRequest, res.Status())
	}
	return res.JSON200, nil
}

// retrieveFreeTierCluster retrieves cluster information.
func (f *FreeTierCluster) retrieveFreeTierCluster(
	ctx context.Context,
	organizationId,
	projectId,
	clusterId string,
) (*providerschema.FreeTierCluster, error) {
	freeTierClusterResp, err := f.getFreeTierCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.CouchbaseAuditData{
		CreatedAt:  types.StringValue(freeTierClusterResp.Audit.CreatedAt.String()),
		CreatedBy:  types.StringValue(freeTierClusterResp.Audit.CreatedBy),
		ModifiedAt: types.StringValue(freeTierClusterResp.Audit.ModifiedAt.String()),
		ModifiedBy: types.StringValue(freeTierClusterResp.Audit.ModifiedBy),
		Version:    types.Int64Value(int64(freeTierClusterResp.Audit.Version)),
	}
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
	}

	availability := providerschema.NewAvailability(apigen.Availability{Type: apigen.AvailabilityType(freeTierClusterResp.Availability.Type)})
	availabilityObj, diags := types.ObjectValueFrom(ctx, availability.AttributeTypes(), availability)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert availability data %w", err)
	}

	support := providerschema.NewSupport(struct {
		Plan     apigen.GetFreeTierCluster200SupportPlan
		Timezone apigen.SupportTimezone
	}{Plan: freeTierClusterResp.Support.Plan, Timezone: freeTierClusterResp.Support.Timezone})
	supportObj, diags := types.ObjectValueFrom(ctx, support.AttributeTypes(), support)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert support data %w", err)
	}

	serviceGroups, err := providerschema.NewTerraformServiceGroups(&apigen.GetClusterResponse{
		ServiceGroups: freeTierClusterResp.ServiceGroups,
		CloudProvider: freeTierClusterResp.CloudProvider,
	})
	if diags.HasError() {
		return nil, fmt.Errorf("unable to convert service groups data %w", err)
	}
	serviceGroupObjList, err, diag := providerschema.NewServiceGroups(ctx, serviceGroups)
	if err != nil {
		if diag.HasError() {
			return nil, err
		}
	}

	serviceGroupsObj, diags := types.SetValueFrom(ctx, types.ObjectType{}.WithAttributeTypes(providerschema.ServiceGroupAttributeTypes()), serviceGroupObjList)
	if diags.HasError() {
		return nil, fmt.Errorf("error while converting servicegroups to service group object ")
	}

	refreshedState, err := providerschema.NewFreeTierCluster(ctx, freeTierClusterResp, organizationId, projectId, auditObj, availabilityObj, supportObj, serviceGroupsObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}

func (f *FreeTierCluster) morphFreeTierClusterRespToTerraformObj(
	ctx context.Context,
	organizationId string,
	projectId string,
	freeTierClusterResp *struct {
		AppServiceId               *openapi_types.UUID       `json:"appServiceId,omitempty"`
		Audit                      apigen.CouchbaseAuditData `json:"audit"`
		Availability               apigen.Availability       `json:"availability"`
		CloudProvider              apigen.CloudProvider      `json:"cloudProvider"`
		CmekId                     *string                   `json:"cmekId,omitempty"`
		ConfigurationType          apigen.ConfigurationType  `json:"configurationType"`
		ConnectionString           string                    `json:"connectionString"`
		CouchbaseServer            apigen.CouchbaseServer    `json:"couchbaseServer"`
		CurrentState               apigen.CurrentState       `json:"currentState"`
		Description                string                    `json:"description"`
		EnablePrivateDNSResolution *bool                     `json:"enablePrivateDNSResolution,omitempty"`
		Id                         openapi_types.UUID        `json:"id"`
		Name                       string                    `json:"name"`
		ServiceGroups              []apigen.ServiceGroup     `json:"serviceGroups"`
		Support                    struct {
			Plan     apigen.GetFreeTierCluster200SupportPlan `json:"plan"`
			Timezone apigen.SupportTimezone                  `json:"timezone"`
		} `json:"support"`
	},
) (*providerschema.FreeTierCluster, error) {
	return f.retrieveFreeTierCluster(ctx, organizationId, projectId, freeTierClusterResp.Id.String())
}
