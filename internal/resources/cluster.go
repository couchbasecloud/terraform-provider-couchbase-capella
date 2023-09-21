package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"terraform-provider-capella/internal/api"
	clusterapi "terraform-provider-capella/internal/api/cluster"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &Cluster{}
	_ resource.ResourceWithConfigure   = &Cluster{}
	_ resource.ResourceWithImportState = &Cluster{}
)

// Cluster is the project resource implementation.
type Cluster struct {
	*providerschema.Data
}

// NewCluster is a helper function to simplify the provider implementation.
func NewCluster() resource.Resource {
	return &Cluster{}
}

// Metadata returns the Cluster resource type name.
func (c *Cluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the Cluster resource.
func (c *Cluster) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterSchema()
}

// Create creates a new Cluster.
func (c *Cluster) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.Cluster
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ClusterRequest := clusterapi.CreateClusterRequest{
		Name: plan.Name.ValueString(),
		Availability: clusterapi.Availability{
			Type: clusterapi.AvailabilityType(plan.Availability.Type.ValueString()),
		},
		CloudProvider: clusterapi.CloudProvider{
			Cidr:   plan.CloudProvider.Cidr.ValueString(),
			Region: plan.CloudProvider.Region.ValueString(),
			Type:   clusterapi.CloudProviderType(plan.CloudProvider.Type.ValueString()),
		},
		Support: clusterapi.Support{
			Plan:     clusterapi.SupportPlan(plan.Support.Plan.ValueString()),
			Timezone: clusterapi.SupportTimezone(plan.Support.Timezone.ValueString()),
		},
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ClusterRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.CouchbaseServer.Version.IsNull() && !plan.CouchbaseServer.Version.IsUnknown() {
		version := plan.CouchbaseServer.Version.ValueString()
		ClusterRequest.CouchbaseServer = &clusterapi.CouchbaseServer{
			Version: &version,
		}
	}

	serviceGroups, err := c.morphToApiServiceGroups(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster : unexpected error "+err.Error(),
		)
		return
	}

	ClusterRequest.ServiceGroups = serviceGroups

	if plan.OrganizationId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var organizationId = plan.OrganizationId.ValueString()

	if plan.ProjectId.IsNull() {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: organization ID cannot be empty.",
		)
		return
	}
	var projectId = plan.ProjectId.ValueString()

	response, err := c.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", c.HostURL, organizationId, projectId),
		http.MethodPost,
		ClusterRequest,
		c.Token,
		nil,
	)
	_, err = handleClusterError(err)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	ClusterResponse := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &ClusterResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, error during unmarshalling:"+err.Error(),
		)
		return
	}

	err = c.checkClusterStatus(ctx, organizationId, projectId, ClusterResponse.Id.String())
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
	}
}

// Configure It adds the provider configured api to the project resource.
func (c *Cluster) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Read reads project information.
func (c *Cluster) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.Cluster
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
	}
}

// Update updates the Cluster.
func (c *Cluster) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO
}

// Delete deletes the project.
func (c *Cluster) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}

// ImportState imports a remote Cluster that is not created by Terraform.
func (c *Cluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO
}

func (c *Cluster) getCluster(organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
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
}

func (c *Cluster) retrieveCluster(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.Cluster, error) {
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
}

func (c *Cluster) checkClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		ClusterResp *clusterapi.GetClusterResponse
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
			ClusterResp, err = c.getCluster(organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(ClusterResp.CurrentState) {
					return nil
				}
				const msg = "waiting for Cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

func (c *Cluster) morphToApiServiceGroups(plan providerschema.Cluster) ([]clusterapi.ServiceGroup, error) {
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
}
