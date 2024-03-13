package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

var (
	_ resource.Resource                = &AuditLogSettings{}
	_ resource.ResourceWithConfigure   = &AuditLogSettings{}
	_ resource.ResourceWithImportState = &AuditLogSettings{}
)

const errorMessageWhileAuditLogSettingsCreation = "There is an error during audit log settings creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

const errorMessageAfterAuditLogSettingsCreation = "Audit log settings creation is successful, but encountered an error while checking the current" +
	" state of the settings. Please run `terraform plan` after 1-2 minutes to know the" +
	" current state of the setting. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

// AuditLogSettings is the audit log settings resource implementation.
type AuditLogSettings struct {
	*providerschema.Data
}

func NewAuditLogSettings() resource.Resource {
	return &AuditLogSettings{}
}

// Metadata returns the audit log settings resource type name.
func (a *AuditLogSettings) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_log_settings"
}

// Schema defines the schema for the audit log settings resource.
func (a *AuditLogSettings) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AuditLogSettingsSchema()
}

// Configure adds the provider configured client to the audit log settings resource.
func (a *AuditLogSettings) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	a.Data = data
}

// Audit Log API does not have create endpoint
// so create is treated as an update
func (a *AuditLogSettings) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ClusterAuditSettings
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validateAuditSettingPlan(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log settings",
			"Could not create audit log settings, unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
	)

	eventIds := make([]int32, len(plan.EnabledEventIDs))
	for i, event := range plan.EnabledEventIDs {
		eventIds[i] = int32(event.ValueInt64())
	}

	disabledUsers := make([]api.AuditSettingsDisabledUser, len(plan.DisabledUsers))
	for i, user := range plan.DisabledUsers {
		u := api.AuditSettingsDisabledUser{
			Domain: user.Domain.ValueStringPointer(),
			Name:   user.Name.ValueStringPointer(),
		}
		disabledUsers[i] = u
	}

	auditLogUpdateRequest := api.UpdateClusterAuditSettingsRequest{
		AuditEnabled:    plan.AuditEnabled.ValueBool(),
		EnabledEventIDs: eventIds,
		DisabledUsers:   disabledUsers,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/auditLog", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusOK}
	_, err = a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		auditLogUpdateRequest,
		a.Token,
		nil,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating audit log settings",
			errorMessageWhileAuditLogSettingsCreation+api.ParseError(err),
		)
		return
	}

	currentState, err := a.refreshAuditLogSettingsState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating audit log settings",
			errorMessageAfterAuditLogSettingsCreation+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read retrieves audit log settings.
func (a *AuditLogSettings) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ClusterAuditSettings
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Settings",
			"Could not read Capella audit log settings: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.Id]
	)

	refreshedState, err := a.refreshAuditLogSettingsState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Capella Audit Log Settings",
			"Could not read Capella audit log settings "+": "+errString,
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

// Update updates the audit log settings.
func (a *AuditLogSettings) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state providerschema.ClusterAuditSettings
	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = state.OrganizationId.ValueString()
		projectId      = state.ProjectId.ValueString()
		clusterId      = state.ClusterId.ValueString()
	)

	eventIds := make([]int32, len(state.EnabledEventIDs))
	for i, event := range state.EnabledEventIDs {
		eventIds[i] = int32(event.ValueInt64())
	}

	disabledUsers := make([]api.AuditSettingsDisabledUser, len(state.DisabledUsers))
	for i, user := range state.DisabledUsers {
		u := api.AuditSettingsDisabledUser{
			Domain: user.Domain.ValueStringPointer(),
			Name:   user.Name.ValueStringPointer(),
		}
		disabledUsers[i] = u
	}

	auditLogUpdateRequest := api.UpdateClusterAuditSettingsRequest{
		AuditEnabled:    state.AuditEnabled.ValueBool(),
		EnabledEventIDs: eventIds,
		DisabledUsers:   disabledUsers,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/auditLog", a.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusOK}
	_, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		auditLogUpdateRequest,
		a.Token,
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
			"Error updating audit log settings",
			"Could not update audit log settings, unexpected error: "+": "+errString,
		)
		return
	}

	currentState, err := a.refreshAuditLogSettingsState(ctx, organizationId, projectId, clusterId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating audit log settings",
			"Could not update audit log settings "+": "+api.ParseError(err),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// AuditLogSettings does not have delete endpoint
func (a *AuditLogSettings) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"delete is not supported audit log settings",
		"delete is not supported for audit log settings",
	)
	return
}

func (a *AuditLogSettings) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("cluster_id"), req, resp)
}

func (a *AuditLogSettings) refreshAuditLogSettingsState(ctx context.Context, organizationId, projectId, clusterId string) (*providerschema.ClusterAuditSettings, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/auditLog",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
	)

	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	auditSettingsResp := api.GetClusterAuditSettingsResponse{}
	err = json.Unmarshal(response.Body, &auditSettingsResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	state := providerschema.ClusterAuditSettings{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AuditEnabled:   types.BoolValue(auditSettingsResp.AuditEnabled),
	}

	eventIds := make([]types.Int64, len(auditSettingsResp.EnabledEventIDs))
	for i, event := range auditSettingsResp.EnabledEventIDs {
		eventIds[i] = types.Int64Value(int64(event))
	}

	disabledUsers := make([]providerschema.AuditSettingsDisabledUser, len(auditSettingsResp.DisabledUsers))
	for i, user := range disabledUsers {
		disabledUsers[i] = user
	}

	state.EnabledEventIDs = eventIds
	state.DisabledUsers = disabledUsers

	return &state, nil
}

func (a *AuditLogSettings) validateAuditSettingPlan(plan providerschema.ClusterAuditSettings) error {
	if plan.OrganizationId.IsNull() {
		return fmt.Errorf("organization Id cannot be empty")
	}
	if plan.ProjectId.IsNull() {
		return fmt.Errorf("project Id cannot be empty")
	}
	if plan.ClusterId.IsNull() {
		return fmt.Errorf("cluster Id cannot be empty")
	}
	if plan.AuditEnabled.IsUnknown() {
		return fmt.Errorf("please specify value for audit_enabled")
	}
	if len(plan.EnabledEventIDs) == 0 {
		return fmt.Errorf("please provide list of event ids")
	}

	return nil
}
