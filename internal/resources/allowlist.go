package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AllowList{}
	_ resource.ResourceWithConfigure   = &AllowList{}
	_ resource.ResourceWithImportState = &AllowList{}
)

// AllowList is the AllowList resource implementation.
type AllowList struct {
	*providerschema.Data
}

func NewAllowList() resource.Resource {
	return &AllowList{}
}

// Metadata returns the allowlist resource type name.
func (r *AllowList) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_allowlist"
}

// Schema defines the schema for the allowlist resource.
func (r *AllowList) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AllowlistsSchema()
}

// Configure set provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (r *AllowList) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create creates a new allowlist
func (r *AllowList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AllowList
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	allowListRequest := api.CreateAllowListRequest{
		Cidr:      plan.Cidr.ValueString(),
		Comment:   plan.Comment.ValueString(),
		ExpiresAt: plan.ExpiresAt.ValueString(),
	}

	response, err := r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs",
			r.HostURL,
			plan.OrganizationId.ValueString(),
			plan.ProjectId.ValueString(),
			plan.ClusterId.ValueString(),
		),
		http.MethodPost,
		allowListRequest,
		r.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			"Could not execute request, unexpected error: "+err.Error(),
		)
		return
	}

	allowListResponse := api.GetAllowListResponse{}
	err = json.Unmarshal(response.Body, &allowListResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating allow list",
			"Could not create allow list, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.refreshAllowList(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), allowListResponse.Id.String())
	switch err := err.(type) {
	case nil:
	case api.Error:
		resp.Diagnostics.AddError(
			"Error reading Capella AllowList",
			"Could not read Capella AllowList "+allowListResponse.Id.String()+": "+err.CompleteError(),
		)
		return
	default:
		resp.Diagnostics.AddError(
			"Error reading Capella AllowList",
			"Could not read Capella AllowList "+allowListResponse.Id.String()+": "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads project information.
func (r *AllowList) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AllowList
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate parameters were successfully imported
	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allow list: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs["organizationId"]
		projectId      = resourceIDs["projectId"]
		clusterId      = resourceIDs["clusterId"]
		allowListId    = resourceIDs["allowListId"]
	)

	// refresh the existing allow list
	refreshedState, err := r.refreshAllowList(ctx, organizationId, projectId, clusterId, allowListId)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Reading Capella AllowList",
				"Could not read Capella allowListID "+allowListId+": "+err.CompleteError(),
			)
			return
		}
		tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	default:
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allowListID "+allowListId+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the allowlist.
func (r *AllowList) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for allowlists.
	// Allowlists can only be created, read and deleted.
	// http://cbc-cp-api.s3-website-us-east-1.amazonaws.com/#tag/allowedCIDRs(Cluster)
	//
	// Note: In this situation, terraform apply will default to deleting and executing a new create.
	// The update implementation should simply be left empty.
	// https://developer.hashicorp.com/terraform/plugin/framework/resources/update
}

// Delete deletes the allow list.
func (r *AllowList) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve existing state
	var state providerschema.AllowList
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella Allow List",
			"Could not delete Capella allow list: "+err.Error(),
		)
		return
	}
	// Execute request to delete existing allowlist
	_, err = r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s",
			r.HostURL,
			resourceIDs["organizationId"],
			resourceIDs["projectId"],
			resourceIDs["clusterId"],
			resourceIDs["allowListId"],
		),
		http.MethodDelete,
		nil,
		r.Token,
		nil,
	)
	switch err := err.(type) {
	case nil:
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Error Deleting Capella Allow List",
				"Could not delete Capella allowListId "+resourceIDs["allowListId"]+": "+err.CompleteError(),
			)
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
	default:
		resp.Diagnostics.AddError(
			"Error Deleting Capella Allow List",
			"Could not delete Capella allowListId "+resourceIDs["allowListId"]+": "+err.Error(),
		)
		return
	}
}

// ImportState imports a remote allowlist that is not created by Terraform.
// Since Capella APIs may require multiple IDs, such as organizationId, projectId, clusterId,
// this function passes the root attribute which is a comma separated string of multiple IDs.
// example: id=cluster123,project_id=proj123,organization_id=org123
// Unfortunately the terraform import CLI doesn't allow us to pass multiple IDs at this point
// and hence this workaround has been applied.
func (r *AllowList) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getAllowList is used to retrieve an existing allow list
func (r *AllowList) getAllowList(ctx context.Context, organizationId, projectId, clusterId, allowListId string) (*api.GetAllowListResponse, error) {
	response, err := r.Client.Execute(
		fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s",
			r.HostURL,
			organizationId,
			projectId,
			clusterId,
			allowListId,
		),
		http.MethodGet,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)
	}

	allowListResp := api.GetAllowListResponse{}
	err = json.Unmarshal(response.Body, &allowListResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %s", err)
	}
	return &allowListResp, nil
}

// refreshAllowList is used to pass an existing AllowList to the refreshed state
func (r *AllowList) refreshAllowList(ctx context.Context, organizationId, projectId, clusterId, allowListId string) (*providerschema.OneAllowList, error) {
	allowListResp, err := r.getAllowList(ctx, organizationId, projectId, clusterId, allowListId)
	if err != nil {
		return nil, err
	}

	refreshedState := providerschema.OneAllowList{
		Id:             types.StringValue(allowListResp.Id.String()),
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		Cidr:           types.StringValue(allowListResp.Cidr),
		Audit: providerschema.CouchbaseAuditData{
			CreatedAt:  types.StringValue(allowListResp.Audit.CreatedAt.String()),
			CreatedBy:  types.StringValue(allowListResp.Audit.CreatedBy),
			ModifiedAt: types.StringValue(allowListResp.Audit.ModifiedAt.String()),
			ModifiedBy: types.StringValue(allowListResp.Audit.ModifiedBy),
			Version:    types.Int64Value(int64(allowListResp.Audit.Version)),
		},
	}

	// Set optional fields
	if allowListResp.Comment != nil {
		refreshedState.Comment = types.StringValue(*allowListResp.Comment)
	}

	if allowListResp.ExpiresAt != nil {
		refreshedState.ExpiresAt = types.StringValue(*allowListResp.ExpiresAt)
	}

	return &refreshedState, nil
}
