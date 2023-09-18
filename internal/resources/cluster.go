package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	"time"

	providerschema "terraform-provider-capella/internal/schema"
	"terraform-provider-capella/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
func (r *Cluster) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema defines the schema for the Cluster resource.
func (r *Cluster) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"cloud_provider": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
					"region": schema.StringAttribute{
						Required: true,
					},
					"cidr": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"couchbase_server": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"service_groups": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.SingleNestedAttribute{
							Required: true,
							Attributes: map[string]schema.Attribute{
								"compute": schema.SingleNestedAttribute{
									Required: true,
									Attributes: map[string]schema.Attribute{
										"cpu": schema.Int64Attribute{
											Required: true,
										},
										"ram": schema.Int64Attribute{
											Required: true,
										},
									},
								},
								"disk": schema.SingleNestedAttribute{
									Required: true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											Required: true,
										},
										"storage": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
										"iops": schema.Int64Attribute{
											Optional: true,
											Computed: true,
										},
									},
								},
							},
						},
						"num_of_nodes": schema.Int64Attribute{
							Required: true,
						},
						"services": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
						},
					},
				},
			},
			"availability": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"support": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"plan": schema.StringAttribute{
						Required: true,
					},
					"timezone": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"current_state": schema.StringAttribute{
				Computed: true,
			},
			"app_service_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"audit": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"created_at": schema.StringAttribute{
						Computed: true,
					},
					"created_by": schema.StringAttribute{
						Computed: true,
					},
					"modified_at": schema.StringAttribute{
						Computed: true,
					},
					"modified_by": schema.StringAttribute{
						Computed: true,
					},
					"version": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
			"if_match": schema.StringAttribute{
				Optional: true,
			},
			"etag": schema.StringAttribute{
				Computed: true,
			},
		},
	}

}

// Create creates a new Cluster.
func (r *Cluster) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	ClusterRequest := api.CreateClusterRequest{
		Name: plan.Name.ValueString(),
		Availability: api.Availability{
			Type: api.AvailabilityType(plan.Availability.Type.ValueString()),
		},
		CloudProvider: api.CloudProvider{
			Cidr:   plan.CloudProvider.Cidr.ValueString(),
			Region: plan.CloudProvider.Region.ValueString(),
			Type:   api.CloudProviderType(plan.CloudProvider.Type.ValueString()),
		},
		Support: api.Support{
			Plan:     api.SupportPlan(plan.Support.Plan.ValueString()),
			Timezone: api.SupportTimezone(plan.Support.Timezone.ValueString()),
		},
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		ClusterRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.CouchbaseServer.Version.IsNull() && !plan.CouchbaseServer.Version.IsUnknown() {
		version := plan.CouchbaseServer.Version.ValueString()
		ClusterRequest.CouchbaseServer = &api.CouchbaseServer{
			Version: &version,
		}
	}

	serviceGroups, err := r.morphToApiServiceGroups(plan)
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

	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", r.HostURL, organizationId, projectId),
		http.MethodPost,
		ClusterRequest,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.Error(),
		)
		return
	}

	ClusterResponse := api.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &ClusterResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, error during unmarshalling:"+err.Error(),
		)
		return
	}

	err = r.checkClusterStatus(ctx, organizationId, projectId, ClusterResponse.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveCluster(ctx, organizationId, projectId, ClusterResponse.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create Cluster, unexpected error: "+err.Error(),
		)
		return
	}

	for i, serviceGroup := range refreshedState.ServiceGroups {
		if utils.AreEqual(plan.ServiceGroups[i].Services, serviceGroup.Services) {
			refreshedState.ServiceGroups[i].Services = plan.ServiceGroups[i].Services
		}
	}

	//refreshedState.ServiceGroups = plan.ServiceGroups

	//need to have proper check since we are passing 7.1 and response is returning 7.1.5
	if !plan.CouchbaseServer.Version.IsNull() && !plan.CouchbaseServer.Version.IsUnknown() {
		refreshedState.CouchbaseServer.Version = plan.CouchbaseServer.Version
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure It adds the provider configured api to the project resource.
func (r *Cluster) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.Data = data
}

// Read reads project information.
func (r *Cluster) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO
}

// Update updates the Cluster.
func (r *Cluster) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO
}

// Delete deletes the project.
func (r *Cluster) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}

// ImportState imports a remote Cluster that is not created by Terraform.
func (r *Cluster) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *Cluster) getCluster(organizationId, projectId, clusterId string) (*api.GetClusterResponse, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", r.HostURL, organizationId, projectId, clusterId),
		http.MethodGet,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	ClusterResp := api.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &ClusterResp)
	if err != nil {
		return nil, err
	}
	ClusterResp.Etag = response.Response.Header.Get("ETag")
	return &ClusterResp, nil
}

func (r *Cluster) retrieveCluster(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.ClusterResourceModel, error) {
	ClusterResp, err := r.getCluster(organizationId, projectId, clusterId)
	if err != nil {
		return nil, err
	}

	audit := providerschema.CouchbaseAuditData{
		CreatedAt:  types.StringValue(ClusterResp.Audit.CreatedAt.String()),
		CreatedBy:  types.StringValue(ClusterResp.Audit.CreatedBy),
		ModifiedAt: types.StringValue(ClusterResp.Audit.ModifiedAt.String()),
		ModifiedBy: types.StringValue(ClusterResp.Audit.ModifiedBy),
		Version:    types.Int64Value(int64(ClusterResp.Audit.Version)),
	}

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error while audit conversion")
	}

	refreshedState := providerschema.ClusterResourceModel{
		Id:             types.StringValue(ClusterResp.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		Name:           types.StringValue(ClusterResp.Name),
		Description:    types.StringValue(ClusterResp.Description),
		Availability: &providerschema.Availability{
			Type: types.StringValue(string(ClusterResp.Availability.Type)),
		},
		CloudProvider: &providerschema.CloudProvider{
			Cidr:   types.StringValue(ClusterResp.CloudProvider.Cidr),
			Region: types.StringValue(ClusterResp.CloudProvider.Region),
			Type:   types.StringValue(string(ClusterResp.CloudProvider.Type)),
		},
		Support: &providerschema.Support{
			Plan:     types.StringValue(string(ClusterResp.Support.Plan)),
			Timezone: types.StringValue(string(ClusterResp.Support.Timezone)),
		},
		CurrentState: types.StringValue(string(ClusterResp.CurrentState)),
		Audit:        auditObj,
		Etag:         types.StringValue(ClusterResp.Etag),
	}

	if ClusterResp.CouchbaseServer.Version != nil {
		version := *ClusterResp.CouchbaseServer.Version
		refreshedState.CouchbaseServer = &providerschema.CouchbaseServer{
			Version: types.StringValue(version),
		}
	}

	var serviceGroups []providerschema.ServiceGroup
	for _, serviceGroup := range ClusterResp.ServiceGroups {

		serviceGroupData := providerschema.ServiceGroup{
			Node: &providerschema.Node{
				Compute: providerschema.Compute{
					Ram: types.Int64Value(int64(serviceGroup.Node.Compute.Ram)),
					Cpu: types.Int64Value(int64(serviceGroup.Node.Compute.Cpu)),
				},
			},
			NumOfNodes: types.Int64Value(int64(*serviceGroup.NumOfNodes)),
		}

		switch ClusterResp.CloudProvider.Type {
		case api.Aws:
			awsDisk, err := serviceGroup.Node.AsDiskAWS()
			if err != nil {
				return nil, err
			}
			serviceGroupData.Node.Disk = providerschema.Node_Disk{
				Type:    types.StringValue(string(awsDisk.Type)),
				Storage: types.Int64Value(int64(awsDisk.Storage)),
				IOPS:    types.Int64Value(int64(awsDisk.Iops)),
			}
		case api.Azure:
			azureDisk, err := serviceGroup.Node.AsDiskAzure()
			if err != nil {
				return nil, err
			}

			serviceGroupData.Node.Disk = providerschema.Node_Disk{
				Type:    types.StringValue(string(azureDisk.Type)),
				Storage: types.Int64Value(int64(*azureDisk.Storage)),
				IOPS:    types.Int64Value(int64(*azureDisk.Iops)),
			}
		case api.Gcp:
			gcpDisk, err := serviceGroup.Node.AsDiskGCP()
			if err != nil {
				return nil, err
			}
			serviceGroupData.Node.Disk = providerschema.Node_Disk{
				Type:    types.StringValue(string(gcpDisk.Type)),
				Storage: types.Int64Value(int64(gcpDisk.Storage)),
			}
		default:
			return nil, fmt.Errorf("unsupported cloud provider is recieved from server")
		}

		serviceGroupData.NumOfNodes = types.Int64Value(int64(*serviceGroup.NumOfNodes))

		for _, service := range *serviceGroup.Services {
			tfService := types.StringValue(string(service))
			serviceGroupData.Services = append(serviceGroupData.Services, tfService)
		}
		serviceGroups = append(serviceGroups, serviceGroupData)
	}
	refreshedState.ServiceGroups = serviceGroups
	return &refreshedState, nil
}

func (r *Cluster) checkClusterStatus(ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		ClusterResp *api.GetClusterResponse
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
			ClusterResp, err = r.getCluster(organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if utils.IsFinalState(ClusterResp.CurrentState) {
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

func (r *Cluster) morphToApiServiceGroups(plan providerschema.ClusterResourceModel) ([]api.ServiceGroup, error) {
	var serviceGroups []api.ServiceGroup
	for _, serviceGroup := range plan.ServiceGroups {
		numOfNodes := int(serviceGroup.NumOfNodes.ValueInt64())
		serviceGroupData := api.ServiceGroup{
			Node: &api.Node{
				Compute: api.Compute{
					Ram: int(serviceGroup.Node.Compute.Ram.ValueInt64()),
					Cpu: int(serviceGroup.Node.Compute.Cpu.ValueInt64()),
				},
			},
			NumOfNodes: &numOfNodes,
		}

		switch plan.CloudProvider.Type.ValueString() {
		case string(api.Aws):
			node := api.Node{}
			diskAws := api.DiskAWS{
				Type: api.DiskAWSType(serviceGroup.Node.Disk.Type.ValueString()),
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
			serviceGroupData.Node.Disk = node.Disk

		case string(api.Azure):
			node := api.Node{}

			diskAzure := api.DiskAzure{
				Type: api.DiskAzureType(serviceGroup.Node.Disk.Type.ValueString()),
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.Storage.IsNull() && !serviceGroup.Node.Disk.Storage.IsUnknown() {
				storage := int(serviceGroup.Node.Disk.Storage.ValueInt64())
				diskAzure.Storage = &storage
			}

			if serviceGroup.Node != nil && !serviceGroup.Node.Disk.IOPS.IsNull() && !serviceGroup.Node.Disk.Storage.IsUnknown() {
				iops := int(serviceGroup.Node.Disk.IOPS.ValueInt64())
				diskAzure.Iops = &iops
			}

			err := node.FromDiskAzure(diskAzure)
			if err != nil {
				return nil, err
			}

			serviceGroupData.Node.Disk = node.Disk

		case string(api.Gcp):
			storage := int(serviceGroup.Node.Disk.Storage.ValueInt64())
			node := api.Node{}
			err := node.FromDiskGCP(api.DiskGCP{
				Type:    api.DiskGCPType(serviceGroup.Node.Disk.Type.ValueString()),
				Storage: storage,
			})
			if err != nil {
				return nil, err
			}
			serviceGroupData.Node.Disk = node.Disk
		}
		for _, service := range serviceGroup.Services {
			serviceapi := service.ValueString()
			if serviceGroupData.Services == nil {
				var emptyList []api.Service
				serviceGroupData.Services = &emptyList
			}
			serviceGroupDataServices := append(*serviceGroupData.Services, api.Service(serviceapi))
			serviceGroupData.Services = &serviceGroupDataServices
		}
		serviceGroups = append(serviceGroups, serviceGroupData)
	}
	return serviceGroups, nil
}
