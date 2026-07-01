package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*EventingFunction)(nil)
	_ resource.ResourceWithConfigure   = (*EventingFunction)(nil)
	_ resource.ResourceWithImportState = (*EventingFunction)(nil)
)

// EventingFunction is the eventing function resource implementation.
type EventingFunction struct {
	*providerschema.Data
}

// NewEventingFunction is a helper function to simplify the provider implementation.
func NewEventingFunction() resource.Resource {
	return &EventingFunction{}
}

// Metadata returns the eventing function resource type name.
func (e *EventingFunction) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventing_function"
}

// Schema defines the schema for the eventing function resource.
func (e *EventingFunction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = EventingFunctionSchema()
}

// Create creates a new eventing function. The function is created in the undeployed state; if the
// plan requests a deployment state, the activationState endpoint is called to reconcile it.
func (e *EventingFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.EventingFunctionResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.State.ValueString() == eventingStatePaused {
		resp.Diagnostics.AddError(
			"Invalid eventing function state",
			"Cannot create a paused eventing function",
		)
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		name           = plan.Name.ValueString()
	)

	plannedSettings, sdiags := eventingSettingsFromObject(ctx, plan.Settings)
	resp.Diagnostics.Append(sdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createBindings, err := bindingsToAPI(ctx, plan.Bindings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating eventing function",
			"Could not convert eventing function bindings: "+err.Error(),
		)
		return
	}

	createReq := eventingapi.CreateEventingFunctionRequest{
		Name:                 name,
		Code:                 plan.Code.ValueStringPointer(),
		EventSource:          keyspaceToAPI(plan.EventSource),
		EventMetadataStorage: keyspaceToAPI(plan.EventMetadataStorage),
		Settings:             settingsToAPI(plannedSettings),
		Bindings:             createBindings,
	}

	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueStringPointer()
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions", e.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	if _, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, createReq, e.Token, nil); err != nil {
		resp.Diagnostics.AddError(
			"Error creating eventing function",
			"Could not create eventing function, unexpected error: "+api.ParseError(err),
		)
		return
	}

	if !plan.State.IsNull() && plan.State.ValueString() != eventingStateUndeployed {
		if err := e.setActivationState(ctx, organizationId, projectId, clusterId, name, plan.State.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error setting state of eventing function after create",
				"Eventing function was created but its activation state could not be set: "+api.ParseError(err),
			)
			return
		}
	}

	// Read the function back to populate computed attributes with their server-assigned values.
	// The plan is passed forward so URL binding secrets and the State verb are preserved.
	refreshedState, err := e.retrieveEventingFunction(ctx, organizationId, projectId, clusterId, name, &plan)
	if err != nil {
		// The function was created, so do not error out and orphan it; fall back to the plan.
		resp.Diagnostics.AddWarning(
			"Error reading eventing function after create",
			"Eventing function was created but could not be read back: "+api.ParseError(err),
		)

		resp.Diagnostics.Append(setEventingFunctionComputedAttributesToNull(ctx, &plan)...)
		if resp.Diagnostics.HasError() {
			return
		}

		refreshedState = &plan
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// setEventingFunctionComputedAttributesToNull sets the computed attributes on the plan to null. It is
// used when setting state after create if the post-create read fails, so the resulting state holds no
// unknown values.
func setEventingFunctionComputedAttributesToNull(ctx context.Context, plan *providerschema.EventingFunctionResource) diag.Diagnostics {
	var diags diag.Diagnostics

	nullKeyspaceComputedAttributes(plan.EventSource)
	nullKeyspaceComputedAttributes(plan.EventMetadataStorage)

	attrTypes := providerschema.EventingFunctionSettings{}.AttributeTypes()
	if plan.Settings.IsNull() || plan.Settings.IsUnknown() {
		plan.Settings = types.ObjectNull(attrTypes)
	} else {
		var s providerschema.EventingFunctionSettings
		diags.Append(plan.Settings.As(ctx, &s, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)

		settings, d := types.ObjectValue(attrTypes, map[string]attr.Value{
			"worker_count":           s.WorkerCount,
			"script_timeout":         s.ScriptTimeout,
			"sql_consistency":        s.SqlConsistency,
			"language_compatibility": s.LanguageCompatibility,
			"feed_boundary":          s.FeedBoundary,
			"max_timer_context_size": s.MaxTimerContextSize,
			"allow_sync_documents":   s.AllowSyncDocuments,
			"cursor_aware":           s.CursorAware,
		})
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		plan.Settings = settings
	}

	if plan.Bindings != nil {
		for i := range plan.Bindings.Buckets {
			plan.Bindings.Buckets[i].Scope = types.StringNull()
			plan.Bindings.Buckets[i].Collection = types.StringNull()
			plan.Bindings.Buckets[i].Permission = types.StringNull()
		}
		for i := range plan.Bindings.Urls {
			plan.Bindings.Urls[i].AllowCookies = types.BoolNull()
			plan.Bindings.Urls[i].ValidateTLSCertificate = types.BoolNull()

			plan.Bindings.Urls[i].Authentication = types.ObjectNull(
				providerschema.
					EventingFunctionURLBindingAuthentication{}.
					AttributeTypes(),
			)
		}
	}

	return diags
}

// nullKeyspaceComputedAttributes sets the computed scope and collection of a keyspace to null. It is a
// no-op for a nil keyspace.
func nullKeyspaceComputedAttributes(k *providerschema.EventingFunctionKeyspace) {
	if k == nil {
		return
	}
	k.Scope = types.StringNull()
	k.Collection = types.StringNull()
}

// eventingValueChanged determines if scalar values are different.
func eventingValueChanged(plan, state attr.Value) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}
	return !plan.Equal(state)
}

// eventingSettingsChanged determines if any of the eventing function settings have changed.
func eventingSettingsChanged(plan, state *providerschema.EventingFunctionSettings) bool {
	if plan == nil {
		return false
	}

	// this should not happen since settings are computed
	if state == nil {
		return true
	}

	return eventingValueChanged(plan.WorkerCount, state.WorkerCount) ||
		eventingValueChanged(plan.ScriptTimeout, state.ScriptTimeout) ||
		eventingValueChanged(plan.SqlConsistency, state.SqlConsistency) ||
		eventingValueChanged(plan.LanguageCompatibility, state.LanguageCompatibility) ||
		eventingValueChanged(plan.FeedBoundary, state.FeedBoundary) ||
		eventingValueChanged(plan.MaxTimerContextSize, state.MaxTimerContextSize) ||
		eventingValueChanged(plan.AllowSyncDocuments, state.AllowSyncDocuments) ||
		eventingValueChanged(plan.CursorAware, state.CursorAware)
}

// eventingBindingsChanged determines if any of the bindings have changed,
// except for secrets (password or bearer token).
func eventingBindingsChanged(ctx context.Context, plan, state *providerschema.EventingFunctionBindingsResource) (bool, error) {
	if plan == nil {
		return false, nil
	}

	if state == nil {
		return true, nil
	}

	if len(plan.Buckets) != len(state.Buckets) ||
		len(plan.Urls) != len(state.Urls) ||
		len(plan.Constants) != len(state.Constants) {
		return true, nil
	}

	for i := range plan.Buckets {
		p, s := plan.Buckets[i], state.Buckets[i]
		if eventingValueChanged(p.Alias, s.Alias) ||
			eventingValueChanged(p.Bucket, s.Bucket) ||
			eventingValueChanged(p.Scope, s.Scope) ||
			eventingValueChanged(p.Collection, s.Collection) ||
			eventingValueChanged(p.Permission, s.Permission) {
			return true, nil
		}
	}

	for i := range plan.Urls {
		p, s := plan.Urls[i], state.Urls[i]
		if eventingValueChanged(p.Alias, s.Alias) ||
			eventingValueChanged(p.Url, s.Url) ||
			eventingValueChanged(p.AllowCookies, s.AllowCookies) ||
			eventingValueChanged(p.ValidateTLSCertificate, s.ValidateTLSCertificate) {
			return true, nil
		}
		planAuth, err := providerschema.AuthenticationFromObject(ctx, p.Authentication)
		if err != nil {
			return false, err
		}
		stateAuth, err := providerschema.AuthenticationFromObject(ctx, s.Authentication)
		if err != nil {
			return false, err
		}
		if eventingURLAuthChanged(planAuth, stateAuth) {
			return true, nil
		}
	}

	for i := range plan.Constants {
		p, s := plan.Constants[i], state.Constants[i]
		if eventingValueChanged(p.Alias, s.Alias) ||
			eventingValueChanged(p.Value, s.Value) {
			return true, nil
		}
	}

	return false, nil
}

// eventingURLAuthChanged determines if authentication type, username or secret changed.
func eventingURLAuthChanged(plan, state *providerschema.EventingFunctionURLBindingAuthentication) bool {
	if plan == nil {
		return false
	}

	// should not happen as auth is computed
	if state == nil {
		return true
	}

	return eventingValueChanged(plan.Type, state.Type) ||
		eventingValueChanged(plan.Username, state.Username) ||
		eventingSecretChanged(plan, state)
}

// eventingSecretChanged determines if the URL binding secret changed. The secret compared
// depends on the auth type: basic and digest use the password, bearer uses the token, and
// none has no secret.
func eventingSecretChanged(plan, state *providerschema.EventingFunctionURLBindingAuthentication) bool {
	switch plan.Type.ValueString() {
	case "basic", "digest":
		return eventingValueChanged(plan.Password, state.Password)
	case "bearer":
		return eventingValueChanged(plan.BearerToken, state.BearerToken)
	default:
		return false
	}
}

// eventingFunctionChanged determines if any of the following has changed for eventing function:
// description, application code, source/metadata keyspace, settings or bindings.
func eventingFunctionChanged(
	ctx context.Context,
	plan, state *providerschema.EventingFunctionResource,
	plannedSettings, stateSettings *providerschema.EventingFunctionSettings,
) (bool, error) {

	bindingsChanged, err := eventingBindingsChanged(ctx, plan.Bindings, state.Bindings)
	if err != nil {
		return false, err
	}

	return eventingValueChanged(plan.Description, state.Description) ||
		eventingValueChanged(plan.Code, state.Code) ||
		eventingKeyspaceChanged(plan.EventSource, state.EventSource) ||
		eventingKeyspaceChanged(plan.EventMetadataStorage, state.EventMetadataStorage) ||
		eventingSettingsChanged(plannedSettings, stateSettings) ||
		bindingsChanged, nil
}

// eventingKeyspaceChanged determines whether the plan changes the keyspace relative to the prior state.
func eventingKeyspaceChanged(plan, state *providerschema.EventingFunctionKeyspace) bool {
	if plan == nil {
		return false
	}

	// this should not happen since keepspace is required
	if state == nil {
		return true
	}

	return eventingValueChanged(plan.Bucket, state.Bucket) ||
		eventingValueChanged(plan.Scope, state.Scope) ||
		eventingValueChanged(plan.Collection, state.Collection)
}

// Read reads the eventing function information.
func (e *EventingFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.EventingFunctionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating eventing function for read",
			"Could not read eventing function "+state.Name.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		name           = IDs[providerschema.FunctionName]
	)

	refreshedState, err := e.retrieveEventingFunction(ctx, organizationId, projectId, clusterId, name, &state)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading eventing function",
			"Could not read eventing function "+name+": "+errString,
		)
		return
	}

	diags = resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the eventing function. If the desired deployment state changed it is applied first
// via the activationState endpoint, then the function definition is updated.
//
// The eventing function must be undeployed/paused before any changes can be made.
// Users can change state to undeployed/paused and make changes at the same time.
func (e *EventingFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.EventingFunctionResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plannedSettings, d := eventingSettingsFromObject(ctx, plan.Settings)
	resp.Diagnostics.Append(d...)
	stateSettings, sd := eventingSettingsFromObject(ctx, state.Settings)
	resp.Diagnostics.Append(sd...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating eventing function for update",
			"Could not update eventing function "+plan.Name.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		name           = IDs[providerschema.FunctionName]
	)

	eventingFunctionHasChanged, err := eventingFunctionChanged(ctx, &plan, &state, plannedSettings, stateSettings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error detecting eventing function change",
			err.Error(),
		)
		return
	}

	// eventing function can only be changed while the function is undeployed or paused.
	// in other states block any changes up front to prevent "inconsistent result after apply" errors.
	if (plan.State.ValueString() != eventingStateUndeployed && plan.State.ValueString() != eventingStatePaused) && eventingFunctionHasChanged {
		resp.Diagnostics.AddError(
			"Cannot change eventing function while deployed",
			"Eventing function "+name+" must be in an undeployed or paused state in order to be changed. "+
				"You can change the state to undeployed/paused and the eventing function at the same time.",
		)
		return
	}

	// Apply a deployment state change first if the desired state was set and differs from the prior state.
	if !plan.State.IsNull() && plan.State.ValueString() != state.State.ValueString() {
		if err := e.setActivationState(ctx, organizationId, projectId, clusterId, name, plan.State.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error setting state of eventing function for update",
				"Could not set activation state for eventing function "+name+": "+api.ParseError(err),
			)
			return
		}

		state.State = plan.State
		diags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// if only eventing function state was changed skip update step.
	if plan.State.ValueString() == eventingStateDeployed ||
		plan.State.ValueString() == eventingStateResumed ||
		!eventingFunctionHasChanged {
		return
	}

	updateBindings, err := bindingsToAPI(ctx, plan.Bindings)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating eventing function",
			"Could not convert eventing function bindings: "+err.Error(),
		)
		return
	}

	updateReq := eventingapi.UpdateEventingFunctionRequest{
		Description:          plan.Description.ValueStringPointer(),
		Code:                 plan.Code.ValueStringPointer(),
		EventSource:          keyspaceToAPIPtr(plan.EventSource),
		EventMetadataStorage: keyspaceToAPIPtr(plan.EventMetadataStorage),
		Settings:             settingsToAPI(plannedSettings),
		Bindings:             updateBindings,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err = e.ClientV1.ExecuteWithRetry(ctx, cfg, updateReq, e.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error updating eventing function",
			"Could not update eventing function "+name+": "+errString,
		)
		return
	}

	refreshedState, err := e.retrieveEventingFunction(ctx, organizationId, projectId, clusterId, name, &plan)
	if err != nil {
		// The function was created, so do not error out and orphan it; fall back to the plan.
		resp.Diagnostics.AddWarning(
			"Error reading eventing function after create",
			"Eventing function was created but could not be read back: "+api.ParseError(err),
		)

		setEventingFunctionComputedAttributesToNull(ctx, &plan)

		refreshedState = &plan
	}

	diags := resp.State.Set(ctx, refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the eventing function. The function must be undeployed prior to deletion; if it is
// not, the API error is surfaced to the user.
func (e *EventingFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.EventingFunctionResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating eventing function for delete",
			"Could not delete eventing function "+state.Name.String()+": "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		name           = IDs[providerschema.FunctionName]
	)

	// A function must be undeployed before it can be deleted. If it is currently in any other state,
	// undeploy it first.
	if state.State.ValueString() != eventingStateUndeployed {
		if err := e.setActivationState(ctx, organizationId, projectId, clusterId, name, eventingStateUndeployed); err != nil {
			resp.Diagnostics.AddError(
				"Error undeploying eventing function before delete",
				"Could not undeploy eventing function "+name+" prior to deletion: "+api.ParseError(err),
			)
			return
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = e.ClientV1.ExecuteWithRetry(ctx, cfg, nil, e.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server")
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting eventing function",
			"Could not delete eventing function "+name+": "+errString,
		)
		return
	}
}

// Configure adds the provider configured api to the eventing function resource.
func (e *EventingFunction) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	e.Data = data
}

// ImportState imports a remote eventing function that was not created by Terraform.
func (e *EventingFunction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// retrieveEventingFunction fetches an eventing function and morphs it into the Terraform state.
// prior carries forward the State action verb and any URL binding secrets the GET response omits.
func (e *EventingFunction) retrieveEventingFunction(
	ctx context.Context, organizationId, projectId, clusterId, name string, prior *providerschema.EventingFunctionResource,
) (*providerschema.EventingFunctionResource, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, nil, e.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	eventingResp := eventingapi.GetEventingFunctionResponse{}
	if err := json.Unmarshal(response.Body, &eventingResp); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}

	fn, err := providerschema.NewEventingFunctionResource(ctx, &eventingResp, organizationId, projectId, clusterId, prior)
	if err != nil {
		return nil, err
	}

	return fn, nil
}

// activationVerb returns the activationState API action verb for the desired state. The API only
// accepts verbs, while the resource models state as the terminal status it produces.
func activationVerb(state string) string {
	action := ""
	switch state {
	case eventingStateUndeployed:
		action = "undeploy"
	case eventingStatePaused:
		action = "pause"
	case eventingStateResumed:
		action = "resume"
	case eventingStateDeployed:
		action = "deploy"
	}

	return action
}

// targetStatus returns the runtime status the given state drives the function to. resumed has no
// distinct status — a resumed function reports deployed — so it polls for deployed.
func targetStatus(state string) string {
	if state == eventingStateResumed {
		return eventingStateDeployed
	}
	return state
}

// setActivationState calls the activationState endpoint with the verb for the target state, then
// polls until the function reaches the status that state produces. The activation is asynchronous,
// so callers (Create, Update, Delete) rely on this to block until it has taken effect.
func (e *EventingFunction) setActivationState(
	ctx context.Context, organizationId, projectId, clusterId, name, target string,
) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s/activationState", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	if _, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, eventingapi.SetFunctionStateRequest{State: activationVerb(target)}, e.Token, nil); err != nil {
		return err
	}

	return e.waitForStatus(ctx, organizationId, projectId, clusterId, name, targetStatus(target))
}

// waitForStatus polls the eventing function every 5 seconds until its runtime status equals target,
// returning an error if the target is not reached within 5 minutes.
func (e *EventingFunction) waitForStatus(
	ctx context.Context, organizationId, projectId, clusterId, name, target string,
) error {
	const timeout = 5 * time.Minute
	const retryInterval = 5 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for eventing function %q to reach status %q: %w", name, target, ctx.Err())
		case <-ticker.C:
			// ok to pass nil as we just need the status
			f, err := e.retrieveEventingFunction(ctx, organizationId, projectId, clusterId, name, nil)
			if err != nil {
				return fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
			}

			if f.State.ValueString() == target {
				return nil
			}

			tflog.Debug(ctx, fmt.Sprintf("eventing function %q status %q, waiting for %q", name, f.State.ValueString(), target))
		}
	}
}

// keyspaceToAPI converts a schema keyspace into the API keyspace value.
func keyspaceToAPI(k *providerschema.EventingFunctionKeyspace) eventingapi.Keyspace {
	if k == nil {
		return eventingapi.Keyspace{}
	}
	return eventingapi.Keyspace{
		Bucket:     k.Bucket.ValueString(),
		Scope:      utils.StringPointerIfKnown(k.Scope),
		Collection: utils.StringPointerIfKnown(k.Collection),
	}
}

// keyspaceToAPIPtr converts a schema keyspace into a pointer to the API keyspace value, or nil.
func keyspaceToAPIPtr(k *providerschema.EventingFunctionKeyspace) *eventingapi.Keyspace {
	if k == nil {
		return nil
	}
	keyspace := keyspaceToAPI(k)
	return &keyspace
}

// eventingSettingsFromObject converts the settings object value into the concrete struct, returning
// nil when the block was omitted (null/unknown). Unknown computed fields are nulled so they are not
// sent to the API and do not register as user changes.
func eventingSettingsFromObject(ctx context.Context, obj types.Object) (*providerschema.EventingFunctionSettings, diag.Diagnostics) {
	if obj.IsNull() || obj.IsUnknown() {
		return nil, nil
	}
	var s providerschema.EventingFunctionSettings
	diags := obj.As(ctx, &s, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})
	if diags.HasError() {
		return nil, diags
	}
	return &s, diags
}

func settingsToAPI(s *providerschema.EventingFunctionSettings) *eventingapi.Settings {
	if s == nil {
		return nil
	}
	return &eventingapi.Settings{
		WorkerCount:           utils.Int64PointerIfKnown(s.WorkerCount),
		ScriptTimeout:         utils.Int64PointerIfKnown(s.ScriptTimeout),
		SqlConsistency:        utils.StringPointerIfKnown(s.SqlConsistency),
		LanguageCompatibility: utils.StringPointerIfKnown(s.LanguageCompatibility),
		FeedBoundary:          utils.StringPointerIfKnown(s.FeedBoundary),
		MaxTimerContextSize:   utils.Int64PointerIfKnown(s.MaxTimerContextSize),
		AllowSyncDocuments:    utils.BoolPointerIfKnown(s.AllowSyncDocuments),
		CursorAware:           utils.BoolPointerIfKnown(s.CursorAware),
	}
}

func bindingsToAPI(ctx context.Context, b *providerschema.EventingFunctionBindingsResource) (*eventingapi.Bindings, error) {
	if b == nil {
		return nil, nil
	}

	bindings := &eventingapi.Bindings{}

	for _, bucket := range b.Buckets {
		bindings.Buckets = append(bindings.Buckets, eventingapi.BucketBinding{
			Alias:      bucket.Alias.ValueString(),
			Bucket:     bucket.Bucket.ValueString(),
			Scope:      utils.StringPointerIfKnown(bucket.Scope),
			Collection: utils.StringPointerIfKnown(bucket.Collection),
			Permission: utils.StringPointerIfKnown(bucket.Permission),
		})
	}

	for _, u := range b.Urls {
		urlBinding := eventingapi.UrlBinding{
			Alias:                  u.Alias.ValueString(),
			Url:                    u.Url.ValueString(),
			AllowCookies:           utils.BoolPointerIfKnown(u.AllowCookies),
			ValidateTLSCertificate: utils.BoolPointerIfKnown(u.ValidateTLSCertificate),
		}
		auth, err := providerschema.AuthenticationFromObject(ctx, u.Authentication)
		if err != nil {
			return nil, err
		}
		if auth != nil {
			urlBinding.Authentication = &eventingapi.URLBindingAuthentication{
				Type:        auth.Type.ValueString(),
				Username:    auth.Username.ValueStringPointer(),
				Password:    utils.StringPointerIfKnown(auth.Password),
				BearerToken: utils.StringPointerIfKnown(auth.BearerToken),
			}
		}

		bindings.Urls = append(bindings.Urls, urlBinding)
	}

	for _, c := range b.Constants {
		bindings.Constants = append(bindings.Constants, eventingapi.ConstantBinding{
			Alias: c.Alias.ValueString(),
			Value: c.Value.ValueString(),
		})
	}

	return bindings, nil
}
