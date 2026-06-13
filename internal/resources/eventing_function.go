package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
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
	var plan providerschema.EventingFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		name           = plan.Name.ValueString()
	)

	createReq := eventingapi.CreateEventingFunctionRequest{
		Name:                 name,
		Description:          plan.Description.ValueStringPointer(),
		Code:                 plan.Code.ValueStringPointer(),
		EventSource:          keyspaceToAPI(plan.EventSource),
		EventMetadataStorage: keyspaceToAPI(plan.EventMetadataStorage),
		Settings:             settingsToAPI(plan.Settings),
		Bindings:             bindingsToAPI(plan.Bindings),
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions", e.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	_, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, createReq, e.Token, nil)
	if err != nil {
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

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read reads the eventing function information.
func (e *EventingFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.EventingFunction
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
func (e *EventingFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.EventingFunction
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Apply a deployment state change first if the desired state was set and differs from the prior state.
	if !plan.State.IsNull() && plan.State.ValueString() != state.State.ValueString() {
		if err := e.setActivationState(ctx, organizationId, projectId, clusterId, name, plan.State.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error setting state of eventing function for update",
				"Could not set activation state for eventing function "+name+": "+api.ParseError(err),
			)
			return
		}

		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !plan.State.IsNull() &&
		(plan.State.ValueString() == eventingStateDeployed || plan.State.ValueString() == eventingStateResumed) {
		resp.Diagnostics.AddWarning(
			"eventing function settings not applied",
			"eventing function must be undeployed or paused before changes can be applied")

		return
	}

	updateReq := eventingapi.UpdateEventingFunctionRequest{
		Description:          plan.Description.ValueStringPointer(),
		Code:                 plan.Code.ValueStringPointer(),
		EventSource:          keyspaceToAPIPtr(plan.EventSource),
		EventMetadataStorage: keyspaceToAPIPtr(plan.EventMetadataStorage),
		Settings:             settingsToAPI(plan.Settings),
		Bindings:             bindingsToAPI(plan.Bindings),
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

	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the eventing function. The function must be undeployed prior to deletion; if it is
// not, the API error is surfaced to the user.
func (e *EventingFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.EventingFunction
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
	ctx context.Context, organizationId, projectId, clusterId, name string, prior *providerschema.EventingFunction,
) (*providerschema.EventingFunction, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, nil, e.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}

	eventingResp := eventingapi.GetEventingFunctionResponse{}
	if err := json.Unmarshal(response.Body, &eventingResp); err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
	}

	refreshedState := providerschema.NewEventingFunction(&eventingResp, organizationId, projectId, clusterId, prior)
	// `state` mirrors the terminal runtime status, so set it directly for the terminal statuses.
	// Transient statuses (deploying, pausing, undeploying) are left untouched as they are not valid
	// `state` values.
	switch eventingResp.Status {
	case eventingStateUndeployed, eventingStatePaused:
		refreshedState.State = types.StringValue(eventingResp.Status)
	case eventingStateDeployed:
		// deployed status is reached by both deploy and resume; keep the configured distinction
		// so a function applied as resumed does not show false drift against status deployed.
		if prior != nil && prior.State.ValueString() == eventingStateResumed {
			refreshedState.State = types.StringValue(eventingStateResumed)
		} else {
			refreshedState.State = types.StringValue(eventingStateDeployed)
		}
	}
	return refreshedState, nil
}

// activationVerb returns the activationState API action verb for the desired state. The API only
// accepts verbs, while the resource models state as the terminal status it produces.
func activationVerb(state string) string {
	switch state {
	case eventingStateUndeployed:
		return "undeploy"
	case eventingStatePaused:
		return "pause"
	case eventingStateResumed:
		return "resume"
	default: // eventingStateDeployed
		return "deploy"
	}
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

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s", e.HostURL, organizationId, projectId, clusterId, name)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for eventing function %q to reach status %q: %w", name, target, ctx.Err())
		case <-ticker.C:
			response, err := e.ClientV1.ExecuteWithRetry(ctx, cfg, nil, e.Token, nil)
			if err != nil {
				return fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
			}

			eventingResp := eventingapi.GetEventingFunctionResponse{}
			if err := json.Unmarshal(response.Body, &eventingResp); err != nil {
				return fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
			}

			if eventingResp.Status == target {
				return nil
			}

			tflog.Debug(ctx, fmt.Sprintf("eventing function %q status %q, waiting for %q", name, eventingResp.Status, target))
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
		Scope:      k.Scope.ValueStringPointer(),
		Collection: k.Collection.ValueStringPointer(),
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

func settingsToAPI(s *providerschema.EventingFunctionSettings) *eventingapi.Settings {
	if s == nil {
		return nil
	}
	return &eventingapi.Settings{
		WorkerCount:           s.WorkerCount.ValueInt64Pointer(),
		ScriptTimeout:         s.ScriptTimeout.ValueInt64Pointer(),
		SqlConsistency:        s.SqlConsistency.ValueStringPointer(),
		LanguageCompatibility: s.LanguageCompatibility.ValueStringPointer(),
		FeedBoundary:          s.FeedBoundary.ValueStringPointer(),
		MaxTimerContextSize:   s.MaxTimerContextSize.ValueInt64Pointer(),
		AllowSyncDocuments:    s.AllowSyncDocuments.ValueBoolPointer(),
		CursorAware:           s.CursorAware.ValueBoolPointer(),
	}
}

func bindingsToAPI(b *providerschema.EventingFunctionBindings) *eventingapi.Bindings {
	if b == nil {
		return nil
	}

	bindings := &eventingapi.Bindings{}

	for _, bucket := range b.Buckets {
		bindings.Buckets = append(bindings.Buckets, eventingapi.BucketBinding{
			Alias:      bucket.Alias.ValueString(),
			Bucket:     bucket.Bucket.ValueString(),
			Scope:      bucket.Scope.ValueStringPointer(),
			Collection: bucket.Collection.ValueStringPointer(),
			Permission: bucket.Permission.ValueStringPointer(),
		})
	}

	for _, u := range b.Urls {
		urlBinding := eventingapi.UrlBinding{
			Alias:                  u.Alias.ValueString(),
			Url:                    u.Url.ValueString(),
			AllowCookies:           u.AllowCookies.ValueBoolPointer(),
			ValidateTLSCertificate: u.ValidateTLSCertificate.ValueBoolPointer(),
		}
		if u.Authentication != nil {
			urlBinding.Authentication = &eventingapi.URLBindingAuthentication{
				Type:     u.Authentication.Type.ValueString(),
				Username: u.Authentication.Username.ValueStringPointer(),
			}

			if !u.Authentication.Password.IsNull() {
				urlBinding.Authentication.Password = u.Authentication.Password.ValueStringPointer()
			}
			if !u.Authentication.BearerToken.IsNull() {
				urlBinding.Authentication.BearerToken = u.Authentication.BearerToken.ValueStringPointer()
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

	return bindings
}
