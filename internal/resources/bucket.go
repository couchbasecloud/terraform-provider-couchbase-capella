package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"net/http"
	bucketapi "terraform-provider-capella/internal/api/bucket"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Bucket{}
	_ resource.ResourceWithConfigure   = &Bucket{}
	_ resource.ResourceWithImportState = &Bucket{}
)

// Bucket is the bucket resource implementation.
type Bucket struct {
	*providerschema.Data
}

// NewBucket is a helper function to simplify the provider implementation.
func NewBucket() resource.Resource {
	return &Bucket{}
}

// Metadata returns the Cluster resource type name.
func (c *Bucket) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bucket"
}

// Schema defines the schema for the Cluster resource.
func (c *Bucket) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BucketSchema()
}

// Create creates a new Bucket.
func (c *Bucket) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Bucket
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	//BucketRequest := bucketapi.CreateBucketRequest{
	//	Name:                     plan.Name.ValueString(),
	//	Type:                     bucketapi.Type(plan.Type),
	//	StorageBackend:           bucketapi.StorageBackend(plan.StorageBackend),
	//	MemoryAllocationInMb:     plan.MemoryAllocationInMb,
	//	BucketConflictResolution: bucketapi.ConflictResolution(plan.BucketConflictResolution),
	//	DurabilityLevel:          bucketapi.DurabilityLevel(plan.DurabilityLevel),
	//	Replicas:                 plan.Replicas,
	//	Flush:                    plan.Flush,
	//	TimeToLiveInSeconds:      plan.TimeToLiveInSeconds,
	//	EvictionPolicy:           bucketapi.BucketEviction(plan.EvictionPolicy),
	//}
	BucketRequest := bucketapi.CreateBucketRequest{
		Name:                     plan.Name.ValueString(),
		Type:                     plan.Type.ValueString(),
		StorageBackend:           plan.StorageBackend.ValueString(),
		MemoryAllocationInMb:     plan.MemoryAllocationInMb,
		BucketConflictResolution: plan.BucketConflictResolution.ValueString(),
		DurabilityLevel:          plan.DurabilityLevel.ValueString(),
		Replicas:                 plan.Replicas,
		Flush:                    plan.Flush,
		TimeToLiveInSeconds:      plan.TimeToLiveInSeconds,
		//EvictionPolicy:           plan.EvictionPolicy.ValueString(),
	}

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+errors.ErrOrganizationIdCannotBeEmpty.Error(),
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	if plan.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrProjectIdCannotBeEmpty.Error(),
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	if plan.ClusterId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating database credential",
			"Could not create database credential, unexpected error: "+errors.ErrClusterIdCannotBeEmpty.Error(),
		)
		return
	}
	var clusterId = plan.ClusterId.ValueString()

	response, err := c.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", c.HostURL, organizationId, projectId, clusterId),
		http.MethodPost,
		BucketRequest,
		c.Token,
		nil,
	)
	_, err = handleClusterError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating bucket",
			"Could not create bucket, unexpected error: "+string(response.Body),
		)
		return
	}

	ClusterResponse := bucketapi.GetBucketResponse{}
	err = json.Unmarshal(response.Body, &ClusterResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, error during unmarshalling:"+err.Error(),
		)
		return
	}

	/*err = c.checkClusterStatus(ctx, organizationId, projectId, ClusterResponse.Id.String())
	_, err = handleClusterError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := c.retrieveCluster(ctx, organizationId, projectId, ClusterResponse.Id.String())
	_, err = handleClusterError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	for i, serviceGroup := range refreshedState.ServiceGroups {
		if clusterapi.AreEqual(plan.ServiceGroups[i].Services, serviceGroup.Services) {
			refreshedState.ServiceGroups[i].Services = plan.ServiceGroups[i].Services
		}
	}

	//need to have proper check since we are passing 7.1 and response is returning 7.1.5
	c.populateInputServerVersionIfPresent(&plan, refreshedState)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}*/
}

// Configure It adds the provider configured api to the project resource.
func (c *Bucket) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	c.Data = data
}

// Read reads bucket information.
func (c *Bucket) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	/*var state providerschema.Cluster
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	// Get refreshed Cluster value from Capella
	refreshedState, err := c.retrieveCluster(ctx, organizationId, projectId, clusterId)
	resourceNotFound, err := handleClusterError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading cluster",
			"Could not read cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	if len(state.ServiceGroups) == len(refreshedState.ServiceGroups) {
		for i, serviceGroup := range refreshedState.ServiceGroups {
			if clusterapi.AreEqual(state.ServiceGroups[i].Services, serviceGroup.Services) {
				refreshedState.ServiceGroups[i].Services = state.ServiceGroups[i].Services
			}
		}
	}

	//need to have proper check since we are passing 7.1 and response is returning 7.1.5
	c.populateInputServerVersionIfPresent(&state, refreshedState)

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}*/
}

// Delete deletes the project.
func (r *Bucket) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	/*// Retrieve values from state
	var state providerschema.Cluster
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			"Could not delete cluster id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		projectId      = resourceIDs[providerschema.ProjectId]
		clusterId      = resourceIDs[providerschema.ClusterId]
	)

	// Delete existing Cluster
	_, err = r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", r.HostURL, organizationId, projectId, clusterId),
		http.MethodDelete,
		nil,
		r.Token,
		nil,
	)
	resourceNotFound, err := handleClusterError(err)
	if resourceNotFound {
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			"Could not delete cluster id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	err = r.checkClusterStatus(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.Id.ValueString())
	resourceNotFound, err = handleClusterError(err)
	switch err {
	case nil:
		// This case will only occur when cluster deletion has failed,
		// and the cluster record still exists in the cp metadata. Therefore,
		// no error will be returned when performing a GET call.
		cluster, err := r.retrieveCluster(ctx, state.OrganizationId.ValueString(), state.ProjectId.ValueString(), state.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting cluster",
				fmt.Sprintf("Could not delete cluster id %s: %s", state.Id.String(), err.Error()),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting cluster",
			fmt.Sprintf("Could not delete cluster id %s, as current Cluster state: %s", state.Id.String(), cluster.CurrentState),
		)
		return
	default:
		if !resourceNotFound {
			resp.Diagnostics.AddError(
				"Error deleting cluster",
				"Could not delete cluster id "+state.Id.String()+": "+err.Error(),
			)
			return
		}
	}*/
}

// ImportState imports a remote cluster that is not created by Terraform.
func (c *Bucket) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getCluster retrieves cluster information from the specified organization and project
// using the provided cluster ID by open-api call
/*func (c *Cluster) getCluster(organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	response, err := c.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", c.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
		nil,
		c.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	clusterResp := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}
	clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}*/

// retrieveCluster retrieves cluster information for a specified organization, project, and cluster ID.
/*func (c *Cluster) retrieveCluster(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.Cluster, error) {
	ClusterResp, err := c.getCluster(organizationId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	audit := providerschema.NewCouchbaseAuditData(ClusterResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error while audit conversion")
	}

	refreshedState, err := providerschema.NewCluster(ClusterResp, organizationId, projectId, auditObj)
	if err != nil {
		return nil, err
	}
	return refreshedState, nil
}*/

// checkClusterStatus monitors the status of a cluster creation, update and deletion operation for a specified
// organization, project, and cluster ID. It periodically fetches the cluster status using the `getCluster`
// function and waits until the cluster reaches a final state or until a specified timeout is reached.
// The function returns an error if the operation times out or encounters an error during status retrieval.
/*func (c *Cluster) checkClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *clusterapi.GetClusterResponse
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
			const msg = "cluster creation status transition timed out after initiation"
			return fmt.Errorf(msg)

		case <-timer.C:
			clusterResp, err = c.getCluster(organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterResp.CurrentState) {
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
}*/

// morphToApiServiceGroups converts a provider cluster serviceGroups to an API-compatible list of service groups.
/*func (c *Cluster) morphToApiServiceGroups(plan providerschema.Cluster) ([]clusterapi.ServiceGroup, error) {
	var newServiceGroups []clusterapi.ServiceGroup
	for _, serviceGroup := range plan.ServiceGroups {
		numOfNodes := int(serviceGroup.NumOfNodes.ValueInt64())
		newServiceGroup := clusterapi.ServiceGroup{
			Node: &clusterapi.Node{
				Compute: clusterapi.Compute{
					Ram: int(serviceGroup.Node.Compute.Ram.ValueInt64()),
					Cpu: int(serviceGroup.Node.Compute.Cpu.ValueInt64()),
				},
			},
			NumOfNodes: &numOfNodes,
		}

		switch plan.CloudProvider.Type.ValueString() {
		case string(clusterapi.Aws):
			node := clusterapi.Node{}
			diskAws := clusterapi.DiskAWS{
				Type: clusterapi.DiskAWSType(serviceGroup.Node.Disk.Type.ValueString()),
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.Storage.IsNull() {
				diskAws.Storage = int(serviceGroup.Node.Disk.Storage.ValueInt64())
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.IOPS.IsNull() {
				diskAws.Iops = int(serviceGroup.Node.Disk.IOPS.ValueInt64())
			}

			err := node.FromDiskAWS(diskAws)
			if err != nil {
				return nil, err
			}
			newServiceGroup.Node.Disk = node.Disk

		case string(clusterapi.Azure):
			node := clusterapi.Node{}
			diskAzure := clusterapi.DiskAzure{
				Type: clusterapi.DiskAzureType(serviceGroup.Node.Disk.Type.ValueString()),
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.Storage.IsNull() && !serviceGroup.Node.Disk.Storage.IsUnknown() {
				storage := int(serviceGroup.Node.Disk.Storage.ValueInt64())
				diskAzure.Storage = &storage
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.IOPS.IsNull() && !serviceGroup.Node.Disk.Storage.IsUnknown() {
				iops := int(serviceGroup.Node.Disk.IOPS.ValueInt64())
				diskAzure.Iops = &iops
			}
			if err := node.FromDiskAzure(diskAzure); err != nil {
				return nil, err
			}
			newServiceGroup.Node.Disk = node.Disk

		case string(clusterapi.Gcp):
			storage := int(serviceGroup.Node.Disk.Storage.ValueInt64())
			node := clusterapi.Node{}
			err := node.FromDiskGCP(clusterapi.DiskGCP{
				Type:    clusterapi.DiskGCPType(serviceGroup.Node.Disk.Type.ValueString()),
				Storage: storage,
			})
			if err != nil {
				return nil, err
			}
			newServiceGroup.Node.Disk = node.Disk
		}
		var newServices []clusterapi.Service
		for _, service := range serviceGroup.Services {
			newService := service.ValueString()
			newServices = append(newServices, clusterapi.Service(newService))
		}
		newServiceGroup.Services = &newServices
		newServiceGroups = append(newServiceGroups, newServiceGroup)
	}
	return newServiceGroups, nil
}

// need to have proper check since we are passing 7.1 and response is returning 7.1.5
func (c *Cluster) populateInputServerVersionIfPresent(stateOrPlanCluster *providerschema.Cluster, refreshStateCluster *providerschema.Cluster) {
	if stateOrPlanCluster.CouchbaseServer != nil &&
		refreshStateCluster.CouchbaseServer != nil &&
		!stateOrPlanCluster.CouchbaseServer.Version.IsNull() &&
		!stateOrPlanCluster.CouchbaseServer.Version.IsUnknown() {
		refreshStateCluster.CouchbaseServer.Version = stateOrPlanCluster.CouchbaseServer.Version
	}
}

// validateClusterUpdate checks if specific fields in a cluster can be updated and returns an error if not.
func (c *Cluster) validateClusterUpdate(plan, state providerschema.Cluster) error {
	var planOrganizationId, stateOrganizationId string
	if !plan.OrganizationId.IsNull() {
		planOrganizationId = plan.OrganizationId.ValueString()
	}

	if !state.OrganizationId.IsNull() {
		stateOrganizationId = state.OrganizationId.ValueString()
	}

	if planOrganizationId != stateOrganizationId {
		return fmt.Errorf("organizationId can't be updated")
	}

	var planProjectId, stateProjectId string
	if !plan.ProjectId.IsNull() {
		planProjectId = plan.ProjectId.ValueString()
	}

	if !state.ProjectId.IsNull() {
		stateProjectId = state.ProjectId.ValueString()
	}

	if planProjectId != stateProjectId {
		return fmt.Errorf("projectId can't be updated")
	}

	var planCouchbaseServerVersion, stateCouchbaseServerVersion string
	if plan.CouchbaseServer != nil && !plan.CouchbaseServer.Version.IsNull() {
		planCouchbaseServerVersion = plan.CouchbaseServer.Version.ValueString()
	}
	if state.CouchbaseServer != nil && !state.CouchbaseServer.Version.IsNull() {
		stateCouchbaseServerVersion = state.CouchbaseServer.Version.ValueString()
	}

	if planCouchbaseServerVersion != stateCouchbaseServerVersion {
		return fmt.Errorf("couchbase server version can't be updated")
	}

	var planAvailabilityType, stateAvailabilityType string
	if plan.Availability != nil && !plan.Availability.Type.IsNull() {
		planAvailabilityType = plan.Availability.Type.ValueString()
	}
	if state.Availability != nil && !state.Availability.Type.IsNull() {
		stateAvailabilityType = state.Availability.Type.ValueString()
	}

	if planAvailabilityType != stateAvailabilityType {
		return fmt.Errorf("availability type can't be updated")
	}

	var planCloudProvider, stateCloudProvider providerschema.CloudProvider
	if plan.CloudProvider != nil {
		planCloudProvider = *plan.CloudProvider
	}
	if state.CloudProvider != nil {
		stateCloudProvider = *state.CloudProvider
	}

	if !reflect.DeepEqual(planCloudProvider, stateCloudProvider) {
		return fmt.Errorf("cloud provider can't be updated")
	}

	return nil
}

// this func extract error message if error is api.Error and also checks whether error is
// resource not found
func handleClusterError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}*/

func (c *Bucket) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
