package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	freeTierClusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/freeTierCluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &FreeTierCluster{}
	_ resource.ResourceWithConfigure   = &FreeTierCluster{}
	_ resource.ResourceWithImportState = &FreeTierCluster{}
)

const errorMessageAfterFreeTierClusterCreationInitiation = "Cluster creation is initiated, but encountered an error while checking the current" +
	" state of the cluster. Please run `terraform plan` after 4-5 minutes to know the" +
	" current status of the cluster. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileFreeTierClusterCreation = "There is an error during cluster creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

type FreeTierCluster struct {
	*providerschema.Data
}

func NewFreeTierCluster() resource.Resource {
	return &FreeTierCluster{}
}

func (f *FreeTierCluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_free_tier"

}

func (f *FreeTierCluster) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = FreeTierClusterSchema()
}

func (f FreeTierCluster) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan providerschema.FreeTierCluster
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if err := f.validateFreeTierCreateCluster(plan); err != nil {
		response.Diagnostics.AddError(
			"error while validating create free tier cluster",
			"could not create free tier cluster "+err.Error(),
		)
	}
	freeTierClusterCreateRequest := freeTierClusterapi.CreateFreeTierClusterRequest{
		Name: plan.Name.ValueString(),
		CloudProvider: clusterapi.CloudProvider{
			Cidr:   plan.CloudProvider.Cidr.ValueString(),
			Region: plan.CloudProvider.Region.ValueString(),
			Type:   clusterapi.CloudProviderType(plan.CloudProvider.Type.ValueString()),
		},
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		freeTierClusterCreateRequest.Description = plan.Description.ValueStringPointer()
	}
	if plan.OrganizationId.IsNull() {
		response.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()
	if plan.ProjectId.IsNull() {
		response.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier", f.HostURL, organizationId, projectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	res, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		freeTierClusterCreateRequest,
		f.Token,
		nil,
	)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating cluster",
			errorMessageWhileFreeTierClusterCreation+"error during unmarshalling"+err.Error(),
		)
		return
	}
	freeTierClusterResponse := freeTierClusterapi.GetFreeTierClusterResponse{}
	err = json.Unmarshal(res.Body, &freeTierClusterCreateRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Cluster",
			errorMessageWhileFreeTierClusterCreation+"error during unmarshalling:"+err.Error(),
		)
		return
	}
	diags = response.State.Set(ctx, initializePendingFreeTierClusterWithPlanAndId(plan, freeTierClusterResponse.ID.String()))
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	err = f.checkClusterStatus(ctx, organizationId, projectId, freeTierClusterResponse.ID.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			errorMessageAfterFreeTierClusterCreationInitiation+api.ParseError(err),
		)
		return
	}
	refreshedState, err := f.retrieveFreeTierCluster(ctx, organizationId, projectId, freeTierClusterResponse.ID.String())
	if err != nil {
		response.Diagnostics.AddWarning(
			"Error creating cluster",
			errorMessageAfterFreeTierClusterCreationInitiation+api.ParseError(err),
		)
		return
	}
	//if plan.Zones != nil {
	//	refreshedState.Zones = plan.Zones
	//}

	//for i, serviceGroup := range refreshedState.ServiceGroups {
	//	if freeTierClusterapi.AreEqual(plan.ServiceGroups[i].Services, serviceGroup.Services) {
	//		refreshedState.ServiceGroups[i].Services = plan.ServiceGroups[i].Services
	//	}
	//}

	// Set state to fully populated data
	diags = response.State.Set(ctx, refreshedState)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (f FreeTierCluster) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (f FreeTierCluster) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (f FreeTierCluster) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}

func (f FreeTierCluster) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (f *FreeTierCluster) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	//TODO implement me
	panic("implement me")
}

func (f *FreeTierCluster) validateFreeTierCreateCluster(plan providerschema.FreeTierCluster) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	return nil
}

// initializePendingClusterWithPlanAndId initializes an instance of providerschema.Cluster
// with the specified plan and ID. It marks all computed fields as null and state as pending.
func initializePendingFreeTierClusterWithPlanAndId(plan providerschema.FreeTierCluster, id string) providerschema.FreeTierCluster {
	plan.Id = types.StringValue(id)
	plan.CurrentState = types.StringValue("pending")
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringNull()
	}
	//if plan.ConfigurationType.IsNull() || plan.ConfigurationType.IsUnknown() {
	//	plan.ConfigurationType = types.StringNull()
	//}

	if plan.EnablePrivateDNSResolution.IsNull() || plan.EnablePrivateDNSResolution.IsUnknown() {
		plan.EnablePrivateDNSResolution = types.BoolNull()
	}

	if plan.CouchbaseServer.IsNull() || plan.CouchbaseServer.IsUnknown() {
		plan.CouchbaseServer = types.ObjectNull(providerschema.CouchbaseServer{}.AttributeTypes())
	}
	//plan.AppServiceId = types.StringNull()
	plan.ConnectionString = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	//plan.Etag = types.StringNull()

	for _, serviceGroup := range plan.ServiceGroups {
		if serviceGroup.Node != nil && (serviceGroup.Node.Disk.Storage.IsNull() || serviceGroup.Node.Disk.Storage.IsUnknown()) {
			serviceGroup.Node.Disk.Storage = types.Int64Null()
		}
		if serviceGroup.Node != nil && (serviceGroup.Node.Disk.IOPS.IsNull() || serviceGroup.Node.Disk.IOPS.IsUnknown()) {
			serviceGroup.Node.Disk.IOPS = types.Int64Null()
		}
		if serviceGroup.Node != nil && (serviceGroup.Node.Disk.Autoexpansion.IsNull() || serviceGroup.Node.Disk.Autoexpansion.IsUnknown()) {
			serviceGroup.Node.Disk.Autoexpansion = types.BoolNull()
		}
	}
	return plan
}

// checkClusterStatus monitors the status of a cluster creation, update and deletion operation for a specified
// organization, project, and cluster ID. It periodically fetches the cluster status using the `getCluster`
// function and waits until the cluster reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
func (c *FreeTierCluster) checkClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *freeTierClusterapi.GetFreeTierClusterResponse
		err         error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 3

	timer := time.NewTimer(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("cluster creation status transition timed out after initiation, unexpected error: %w", err)
		case <-timer.C:
			clusterResp, err = c.getCluster(ctx, organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterapi.State(clusterResp.CurrentState)) {
					return nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

// getCluster retrieves cluster information from the specified organization and project
// using the provided cluster ID by open-api call.
func (f *FreeTierCluster) getCluster(ctx context.Context, organizationId, projectId, clusterId string) (*freeTierClusterapi.GetFreeTierClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/freeTier/%s", f.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := f.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		f.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	clusterResp := freeTierClusterapi.GetFreeTierClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	//clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}

// retrieveCluster retrieves cluster information for a specified organization, project, and cluster ID.
func (f *FreeTierCluster) retrieveFreeTierCluster(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.FreeTierCluster, error) {
	clusterResp, err := f.getCluster(ctx, organizationId, projectId, clusterId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
	}

	audit := providerschema.NewCouchbaseAuditData(clusterResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnableToConvertAuditData, err)
	}

	refreshedState, err := providerschema.NewFreeTierCluster(ctx, clusterResp, organizationId, projectId, auditObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}
