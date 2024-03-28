package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

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

const errorMessageAfterAllowListCreation = "Allow list creation is successful, but encountered an error while checking the current" +
	" state of the allow list. Please run `terraform plan` after 1-2 minutes to know the" +
	" current allow list state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileAllowListCreation = "There is an error during allow list creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

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

// Create creates a new allowlist.
func (r *AllowList) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AllowList
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.validateCreateAllowList(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating allow list",
			"Could not create allow list, unexpected error: "+err.Error(),
		)
		return
	}

	allowListRequest := api.CreateAllowListRequest{
		Cidr:      plan.Cidr.ValueString(),
		Comment:   plan.Comment.ValueString(),
		ExpiresAt: plan.ExpiresAt.ValueString(),
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs",
		r.HostURL,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		allowListRequest,
		r.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error executing request",
			errorMessageWhileAllowListCreation+api.ParseError(err),
		)
		return
	}

	allowListResponse := api.GetAllowListResponse{}
	err = json.Unmarshal(response.Body, &allowListResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating allow list",
			errorMessageWhileAllowListCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeAllowListWithPlanAndId(plan, allowListResponse.Id.String()))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.refreshAllowList(ctx, plan.OrganizationId.ValueString(), plan.ProjectId.ValueString(), plan.ClusterId.ValueString(), allowListResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Capella AllowList",
			errorMessageAfterAllowListCreation+api.ParseError(err),
		)
		return
	}

	// This is added to workaround any timezone conversions that the API does automatically and
	// may cause an issue in the state file.
	if plan.ExpiresAt != refreshedState.ExpiresAt {
		refreshedState.ExpiresAt = plan.ExpiresAt
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
	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allow list: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		allowListId    = IDs[providerschema.Id]
	)

	// refresh the existing allow list
	refreshedState, err := r.refreshAllowList(ctx, organizationId, projectId, clusterId, allowListId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allowListID "+allowListId+": "+errString,
		)
		return
	}

	if state.ExpiresAt != refreshedState.ExpiresAt {
		refreshedState.ExpiresAt = state.ExpiresAt
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the allowlist.
func (r *AllowList) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Couchbase Capella's v4 does not support a PUT endpoint for allowlists.
	// Allowlists can only be created, read and deleted.
	// https://docs.couchbase.com/cloud/management-api-reference/index.html#tag/allowedCIDRs(Cluster)
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

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Capella Allow List",
			"Could not delete Capella allow list: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		allowListId    = IDs[providerschema.Id]
	)
	// Execute request to delete existing allowlist
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		allowListId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella AllowList",
			"Could not read Capella allowListID "+allowListId+": "+errString,
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

// getAllowList is used to retrieve an existing allow list.
func (r *AllowList) getAllowList(ctx context.Context, organizationId, projectId, clusterId, allowListId string) (*api.GetAllowListResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s",
		r.HostURL,
		organizationId,
		projectId,
		clusterId,
		allowListId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		r.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	allowListResp := api.GetAllowListResponse{}
	err = json.Unmarshal(response.Body, &allowListResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &allowListResp, nil
}

// refreshAllowList is used to pass an existing AllowList to the refreshed state.
func (r *AllowList) refreshAllowList(ctx context.Context, organizationId, projectId, clusterId, allowListId string) (*providerschema.OneAllowList, error) {
	allowListResp, err := r.getAllowList(ctx, organizationId, projectId, clusterId, allowListId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrNotFound, err)
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

func (r *AllowList) validateCreateAllowList(plan providerschema.AllowList) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdCannotBeEmpty
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}

	return r.validateAllowListAttributesTrimmed(plan)
}

func (r *AllowList) validateAllowListAttributesTrimmed(plan providerschema.AllowList) error {
	if (!plan.Comment.IsNull() && !plan.Comment.IsUnknown()) && !providerschema.IsTrimmed(plan.Comment.ValueString()) {
		return fmt.Errorf("comment %s", errors.ErrNotTrimmed)
	}
	return nil
}

// initializeAllowListWithPlanAndId initializes an instance of providerschema.AllowList
// with the specified plan and ID. It marks all computed fields as null.
func initializeAllowListWithPlanAndId(plan providerschema.AllowList, id string) providerschema.AllowList {
	plan.Id = types.StringValue(id)
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	if plan.Comment.IsNull() || plan.Comment.IsUnknown() {
		plan.Comment = types.StringNull()
	}
	return plan
}
