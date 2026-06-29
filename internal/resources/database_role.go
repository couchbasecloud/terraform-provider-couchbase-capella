package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = (*DatabaseRole)(nil)
	_ resource.ResourceWithConfigure   = (*DatabaseRole)(nil)
	_ resource.ResourceWithImportState = (*DatabaseRole)(nil)
)

// DatabaseRole is the database role resource implementation.
type DatabaseRole struct {
	*providerschema.Data
}

func NewDatabaseRole() resource.Resource {
	return &DatabaseRole{}
}

func (r *DatabaseRole) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_role"
}

func (r *DatabaseRole) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DatabaseRoleSchema()
}

func (r *DatabaseRole) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DatabaseRole) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.DatabaseRole
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.validateCreateDatabaseRole(plan); err != nil {
		resp.Diagnostics.AddError(
			"Error creating database role",
			"Could not create database role, unexpected error: "+err.Error(),
		)
		return
	}

	organizationId := plan.OrganizationId.ValueString()
	projectId := plan.ProjectId.ValueString()
	clusterId := plan.ClusterId.ValueString()

	roleRequest := api.CreateDatabaseRoleRequest{
		Name:   plan.Name.ValueString(),
		Access: createAccessFromSlice(plan.Access),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		roleRequest.Description = plan.Description.ValueString()
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/roles", r.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := r.ClientV1.ExecuteWithRetry(ctx, cfg, roleRequest, r.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database role",
			"Could not create database role, unexpected error: "+api.ParseError(err),
		)
		return
	}

	roleResponse := api.CreateDatabaseRoleResponse{}
	err = json.Unmarshal(response.Body, &roleResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating database role",
			"Could not create database role, error during unmarshalling: "+err.Error(),
		)
		return
	}

	refreshedState, err := r.retrieveDatabaseRole(ctx, organizationId, projectId, clusterId, roleResponse.Id.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading database role after creation",
			"Database role was created but could not be read: "+api.ParseError(err),
		)
		return
	}

	refreshedState.Access = mapAccessFromSlice(plan.Access)

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

func (r *DatabaseRole) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.DatabaseRole
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading database role",
			"Could not read database role with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	roleId := IDs[providerschema.Id]

	refreshedState, err := r.retrieveDatabaseRole(ctx, organizationId, projectId, clusterId, roleId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading database role",
			"Could not read database role with ID "+state.Id.String()+": "+errString,
		)
		return
	}

	refreshedState.Access = reconcileAccess(refreshedState.Access, state.Access)

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

func (r *DatabaseRole) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.DatabaseRole
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := plan.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating database role",
			"Could not update database role with ID "+plan.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	roleId := IDs[providerschema.Id]

	updateRequest := api.UpdateDatabaseRoleRequest{
		Access: createAccessFromSlice(plan.Access),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateRequest.Description = plan.Description.ValueString()
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/roles/%s", r.HostURL, organizationId, projectId, clusterId, roleId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = r.ClientV1.ExecuteWithRetry(ctx, cfg, updateRequest, r.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating database role",
			"Could not update database role, unexpected error: "+errString,
		)
		return
	}

	refreshedState, err := r.retrieveDatabaseRole(ctx, organizationId, projectId, clusterId, roleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading database role after update",
			"Could not read database role after update: "+api.ParseError(err),
		)
		return
	}

	refreshedState.Access = mapAccessFromSlice(plan.Access)

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

func (r *DatabaseRole) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.DatabaseRole
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting database role",
			"Could not delete database role with ID "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	organizationId := IDs[providerschema.OrganizationId]
	projectId := IDs[providerschema.ProjectId]
	clusterId := IDs[providerschema.ClusterId]
	roleId := IDs[providerschema.Id]

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/roles/%s", r.HostURL, organizationId, projectId, clusterId, roleId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = r.ClientV1.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting database role",
			"Could not delete database role associated with cluster "+clusterId+": "+errString,
		)
		return
	}
}

// ImportState imports a remote database role that is not created by Terraform.
// example: id=<uuid>,organization_id=<uuid>,project_id=<uuid>,cluster_id=<uuid>
func (r *DatabaseRole) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DatabaseRole) retrieveDatabaseRole(ctx context.Context, organizationId, projectId, clusterId, roleId string) (*providerschema.DatabaseRole, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/roles/%s", r.HostURL, organizationId, projectId, clusterId, roleId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := r.ClientV1.ExecuteWithRetry(ctx, cfg, nil, r.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	roleResp := api.GetDatabaseRoleResponse{}
	err = json.Unmarshal(response.Body, &roleResp)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	audit := providerschema.NewCouchbaseAuditData(roleResp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("%w: %s", errors.ErrUnableToConvertAuditData, diags.Errors())
	}

	refreshedState := providerschema.NewDatabaseRole(
		types.StringValue(roleResp.Id.String()),
		types.StringValue(roleResp.Name),
		types.StringValue(roleResp.Description),
		types.StringValue(organizationId),
		types.StringValue(projectId),
		types.StringValue(clusterId),
		auditObj,
	)

	refreshedState.Access = mapAccessFromAPI(roleResp.Access)

	return refreshedState, nil
}

func (r *DatabaseRole) validateCreateDatabaseRole(plan providerschema.DatabaseRole) error {
	if plan.OrganizationId.IsNull() {
		return errors.ErrOrganizationIdMissing
	}
	if plan.ProjectId.IsNull() {
		return errors.ErrProjectIdMissing
	}
	if plan.ClusterId.IsNull() {
		return errors.ErrClusterIdMissing
	}
	_, err := utils.ParseUUIDs(
		utils.IDField{Name: "organization_id", Value: plan.OrganizationId.ValueString()},
		utils.IDField{Name: "project_id", Value: plan.ProjectId.ValueString()},
		utils.IDField{Name: "cluster_id", Value: plan.ClusterId.ValueString()},
	)
	if err != nil {
		return err
	}
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %w", errors.ErrNotTrimmed)
	}
	return nil
}
