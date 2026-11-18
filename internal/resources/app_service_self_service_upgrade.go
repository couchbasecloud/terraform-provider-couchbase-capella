package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_self_service_upgrade"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*AppServiceSelfServiceUpgrade)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceSelfServiceUpgrade)(nil)
	_ resource.ResourceWithImportState = (*AppServiceSelfServiceUpgrade)(nil)
)

// The create endpoint answers 202 before the upgrade record exists, and the get endpoint answers 404
// until it does. These are vars so tests can shorten them.
var (
	selfServiceUpgradePollInterval  = 5 * time.Second
	selfServiceUpgradeCreateTimeout = 5 * time.Minute
)

// AppServiceSelfServiceUpgrade is the App Service self-service upgrade resource implementation.
type AppServiceSelfServiceUpgrade struct {
	*providerschema.Data
}

// NewAppServiceSelfServiceUpgrade is a helper function to simplify the provider implementation.
func NewAppServiceSelfServiceUpgrade() resource.Resource {
	return &AppServiceSelfServiceUpgrade{}
}

// Metadata returns the App Service self-service upgrade resource type name.
func (a *AppServiceSelfServiceUpgrade) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_self_service_upgrade"
}

// Schema defines the schema for the App Service self-service upgrade resource.
func (a *AppServiceSelfServiceUpgrade) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceSelfServiceUpgradeSchema()
}

// Create upgrades the App Service now, or schedules the upgrade at queue_at. It waits for the upgrade
// record to exist, not for the upgrade to finish.
func (a *AppServiceSelfServiceUpgrade) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceSelfServiceUpgrade
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := appservice_self_service_upgrade.CreateRequest{VersionId: plan.VersionId.ValueString()}
	if !plan.QueueAt.IsNull() && !plan.QueueAt.IsUnknown() {
		queueAt, diags := plan.QueueAt.ValueRFC3339Time()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createRequest.QueueAt = &queueAt
	}

	IDs := map[providerschema.Attr]string{
		providerschema.OrganizationId: plan.OrganizationId.ValueString(),
		providerschema.ProjectId:      plan.ProjectId.ValueString(),
		providerschema.ClusterId:      plan.ClusterId.ValueString(),
		providerschema.AppServiceId:   plan.AppServiceId.ValueString(),
	}
	appServiceId := IDs[providerschema.AppServiceId]

	cfg := api.EndpointCfg{Url: a.upgradesURL(IDs), Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, createRequest, a.Token, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service self-service upgrade",
			"Could not create self-service upgrade of app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	var createResponse appservice_self_service_upgrade.CreateResponse
	if err = json.Unmarshal(response.Body, &createResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error creating App Service self-service upgrade",
			"Could not read the ID of the self-service upgrade of app service "+appServiceId+": "+errors.ErrUnmarshallingResponse.Error()+": "+err.Error(),
		)
		return
	}
	IDs[providerschema.Id] = createResponse.Id

	// Save the IDs before polling, so that an accepted upgrade stays tracked if the poll fails.
	plan.Id = types.StringValue(createResponse.Id)
	if plan.QueueAt.IsUnknown() {
		plan.QueueAt = timetypes.NewRFC3339Null()
	}
	plan.CurrentState = types.StringNull()
	plan.CompletedAt = timetypes.NewRFC3339Null()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.checkAppServiceSelfServiceUpgradeStatus(ctx, IDs, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for App Service self-service upgrade",
			"Self-service upgrade "+createResponse.Id+" of app service "+appServiceId+" was accepted but could not be read: "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read refreshes the state with the upgrade from Capella. A 404 means the upgrade no longer exists;
// any other error, such as the 403 of an organization without self-service upgrades, keeps the state.
func (a *AppServiceSelfServiceUpgrade) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceSelfServiceUpgrade
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service self-service upgrade",
			"Could not read self-service upgrade "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	upgrade, err := a.getUpgrade(ctx, IDs)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service self-service upgrade",
			"Could not read self-service upgrade "+IDs[providerschema.Id]+": "+errString,
		)
		return
	}

	refreshedState, diags := morphToAppServiceSelfServiceUpgrade(ctx, IDs, upgrade, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update replaces the version and queue time of a pending upgrade. Capella rejects the change once
// the upgrade has started, and that error is returned as is.
func (a *AppServiceSelfServiceUpgrade) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.AppServiceSelfServiceUpgrade
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating App Service self-service upgrade",
			"Could not update self-service upgrade "+state.Id.String()+": "+err.Error(),
		)
		return
	}
	upgradeId := IDs[providerschema.Id]

	// The update is a full replace, so a queue_at that is not planned keeps the stored one.
	queueAtValue := plan.QueueAt
	if queueAtValue.IsNull() || queueAtValue.IsUnknown() {
		queueAtValue = state.QueueAt
	}
	queueAt, diags := queueAtValue.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest := appservice_self_service_upgrade.UpdateRequest{VersionId: plan.VersionId.ValueString(), QueueAt: queueAt}
	cfg := api.EndpointCfg{Url: a.upgradeURL(IDs), Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	if _, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, updateRequest, a.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error updating App Service self-service upgrade",
			"Could not update self-service upgrade "+upgradeId+": "+api.ParseError(err),
		)
		return
	}

	// The update returns no body, so read the upgrade back.
	upgrade, err := a.getUpgrade(ctx, IDs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service self-service upgrade",
			"Self-service upgrade "+upgradeId+" was updated but could not be read: "+api.ParseError(err),
		)
		return
	}

	refreshedState, diags := morphToAppServiceSelfServiceUpgrade(ctx, IDs, upgrade, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Delete cancels the upgrade unless it no longer exists or has finished. On any error the resource
// stays in state, so the next destroy retries; it never leaves state while the upgrade is still pending.
func (a *AppServiceSelfServiceUpgrade) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceSelfServiceUpgrade
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error cancelling App Service self-service upgrade",
			"Could not cancel self-service upgrade "+state.Id.String()+": "+err.Error(),
		)
		return
	}
	upgradeId := IDs[providerschema.Id]

	done, err := a.isUpgradeGoneOrFinished(ctx, IDs)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error cancelling App Service self-service upgrade",
			"Could not read self-service upgrade "+upgradeId+" before cancelling it: "+api.ParseError(err),
		)
		return
	}
	if done {
		tflog.Info(ctx, "self-service upgrade no longer exists or has finished, removing it from state without cancelling")
		return
	}

	cfg := api.EndpointCfg{Url: a.upgradeURL(IDs), Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err == nil {
		return
	}

	resourceNotFound, errString := api.CheckResourceNotFoundError(err)
	if !resourceNotFound {
		resp.Diagnostics.AddError(
			"Error cancelling App Service self-service upgrade",
			"Could not cancel self-service upgrade "+upgradeId+": "+errString,
		)
		return
	}

	// A 404 on cancel only counts once the upgrade is confirmed gone or finished.
	done, err = a.isUpgradeGoneOrFinished(ctx, IDs)
	switch {
	case err != nil:
		resp.Diagnostics.AddError(
			"Error cancelling App Service self-service upgrade",
			"Cancelling self-service upgrade "+upgradeId+" returned not found, and it could not be read back: "+api.ParseError(err),
		)
	case !done:
		resp.Diagnostics.AddError(
			"Error cancelling App Service self-service upgrade",
			"Cancelling self-service upgrade "+upgradeId+" returned not found, but the upgrade is still active. "+
				"It is kept in state so that the next destroy retries the cancel.",
		)
	}
}

// ImportState imports an existing self-service upgrade. The import ID is a comma separated string
// of the IDs, for example:
// "id=<upgradeId>,app_service_id=<appServiceId>,cluster_id=<clusterId>,project_id=<projId>,organization_id=<orgId>".
func (a *AppServiceSelfServiceUpgrade) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Configure adds the provider configured client to the App Service self-service upgrade resource.
func (a *AppServiceSelfServiceUpgrade) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerschema.Data, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	a.Data = data
}

func (a *AppServiceSelfServiceUpgrade) upgradesURL(IDs map[providerschema.Attr]string) string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/selfServiceUpgrades",
		a.HostURL, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId], IDs[providerschema.ClusterId], IDs[providerschema.AppServiceId])
}

func (a *AppServiceSelfServiceUpgrade) upgradeURL(IDs map[providerschema.Attr]string) string {
	return a.upgradesURL(IDs) + "/" + url.PathEscape(IDs[providerschema.Id])
}

func (a *AppServiceSelfServiceUpgrade) getUpgrade(ctx context.Context, IDs map[providerschema.Attr]string) (*appservice_maintenance.GetMaintenanceResponse, error) {
	cfg := api.EndpointCfg{Url: a.upgradeURL(IDs), Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	var upgrade appservice_maintenance.GetMaintenanceResponse
	if err = json.Unmarshal(response.Body, &upgrade); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	return &upgrade, nil
}

// checkAppServiceSelfServiceUpgradeStatus polls the upgrade until its record exists and returns it
// as state. It does not wait for the upgrade to finish, as a scheduled upgrade can be days away.
// Errors are retried until the timeout, when the last one is returned.
func (a *AppServiceSelfServiceUpgrade) checkAppServiceSelfServiceUpgradeStatus(
	ctx context.Context, IDs map[providerschema.Attr]string, plan providerschema.AppServiceSelfServiceUpgrade,
) (*providerschema.AppServiceSelfServiceUpgrade, error) {
	ctx, cancel := context.WithTimeout(ctx, selfServiceUpgradeCreateTimeout)
	defer cancel()

	ticker := time.NewTicker(selfServiceUpgradePollInterval)
	defer ticker.Stop()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("upgrade %s was not found after %s: %w", IDs[providerschema.Id], selfServiceUpgradeCreateTimeout, lastErr)
		case <-ticker.C:
			upgrade, err := a.getUpgrade(ctx, IDs)
			if err != nil {
				lastErr = err
				tflog.Info(ctx, "retrying after error polling self-service upgrade", map[string]interface{}{"error": err.Error()})
				continue
			}

			refreshedState, diags := morphToAppServiceSelfServiceUpgrade(ctx, IDs, upgrade, plan)
			if diags.HasError() {
				return nil, fmt.Errorf("%w: %v", errors.ErrRefreshingState, diags.Errors())
			}
			return refreshedState, nil
		}
	}
}

// isUpgradeGoneOrFinished reports whether the upgrade no longer exists or is in a final state, in
// which case there is nothing left to cancel.
func (a *AppServiceSelfServiceUpgrade) isUpgradeGoneOrFinished(ctx context.Context, IDs map[providerschema.Attr]string) (bool, error) {
	upgrade, err := a.getUpgrade(ctx, IDs)
	if err != nil {
		if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
			return true, nil
		}
		return false, err
	}

	return isTerminalSelfServiceUpgradeState(upgrade.CurrentState), nil
}

// morphToAppServiceSelfServiceUpgrade maps the upgrade to state. prior is the plan or the stored
// state: its queue_at string is kept when it is the same instant, and its version_id is kept when
// the API omits the version, which it does for an image no release contains.
func morphToAppServiceSelfServiceUpgrade(
	ctx context.Context, IDs map[providerschema.Attr]string, upgrade *appservice_maintenance.GetMaintenanceResponse, prior providerschema.AppServiceSelfServiceUpgrade,
) (*providerschema.AppServiceSelfServiceUpgrade, diag.Diagnostics) {
	data, diags := providerschema.NewAppServiceMaintenanceData(ctx, upgrade)
	if diags.HasError() {
		return nil, diags
	}

	versionId := data.VersionId
	if versionId.IsNull() {
		versionId = prior.VersionId
	}

	refreshedState := providerschema.AppServiceSelfServiceUpgrade{
		OrganizationId: types.StringValue(IDs[providerschema.OrganizationId]),
		ProjectId:      types.StringValue(IDs[providerschema.ProjectId]),
		ClusterId:      types.StringValue(IDs[providerschema.ClusterId]),
		AppServiceId:   types.StringValue(IDs[providerschema.AppServiceId]),
		Id:             types.StringValue(IDs[providerschema.Id]),
		VersionId:      versionId,
		QueueAt:        normalizeQueueAt(prior.QueueAt, upgrade.QueueAt),
		CurrentState:   data.CurrentState,
		CompletedAt:    data.CompletedAt,
		Audit:          data.Audit,
	}
	return &refreshedState, nil
}

// normalizeQueueAt keeps the prior queue_at string when it is the same instant as the API's, so a
// different format of the same time (for example "+00:00" for "Z") does not show as a change.
func normalizeQueueAt(prior timetypes.RFC3339, actual *time.Time) timetypes.RFC3339 {
	if actual == nil {
		return timetypes.NewRFC3339Null()
	}

	if !prior.IsNull() && !prior.IsUnknown() {
		priorTime, diags := prior.ValueRFC3339Time()
		if !diags.HasError() && priorTime.Equal(*actual) {
			return prior
		}
	}

	return timetypes.NewRFC3339TimePointerValue(actual)
}

// isTerminalSelfServiceUpgradeState reports whether an upgrade in this state can no longer change.
func isTerminalSelfServiceUpgradeState(state string) bool {
	switch state {
	case "completed", "failed", "cancelled":
		return true
	default:
		return false
	}
}
