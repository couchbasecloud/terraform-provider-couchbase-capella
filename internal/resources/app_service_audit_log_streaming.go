package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*AppServiceAuditLogStreaming)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceAuditLogStreaming)(nil)
	_ resource.ResourceWithImportState = (*AppServiceAuditLogStreaming)(nil)
)

// AppServiceAuditLogStreaming is the App Service audit log streaming resource implementation.
type AppServiceAuditLogStreaming struct {
	*providerschema.Data
}

// NewAppServiceAuditLogStreaming is a helper function to simplify the provider implementation.
func NewAppServiceAuditLogStreaming() resource.Resource {
	return &AppServiceAuditLogStreaming{}
}

// Metadata returns the resource type name.
func (a *AppServiceAuditLogStreaming) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_log_streaming"
}

// Schema defines the schema for the resource.
func (a *AppServiceAuditLogStreaming) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogStreamingSchema()
}

// Configure adds the provider configured client to the resource.
func (a *AppServiceAuditLogStreaming) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create configures audit log streaming with the full config.
func (a *AppServiceAuditLogStreaming) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceAuditLogStreaming
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := a.put(ctx, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating App Service Audit Log Streaming",
			"Could not configure audit log streaming, unexpected error: "+err.Error(),
		)
		return
	}

	a.waitAndSetState(ctx, plan, &resp.State, &resp.Diagnostics)
}

// Read refreshes the Terraform state with the latest data.
func (a *AppServiceAuditLogStreaming) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceAuditLogStreaming
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Streaming",
			"Could not validate state: "+err.Error(),
		)
		return
	}
	state.OrganizationId = types.StringValue(IDs[providerschema.OrganizationId])
	state.ProjectId = types.StringValue(IDs[providerschema.ProjectId])
	state.ClusterId = types.StringValue(IDs[providerschema.ClusterId])
	state.AppServiceId = types.StringValue(IDs[providerschema.AppServiceId])

	streamingResp, err := a.get(ctx, state)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log Streaming",
			"Could not read audit log streaming, unexpected error: "+api.ParseError(err),
		)
		return
	}

	// A PUT with only streamingEnabled=false removes the whole config, so there is nothing to manage.
	if streamingResp.OutputType == nil {
		tflog.Info(ctx, "audit log streaming is not configured, removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(setAuditLogStreamingState(ctx, &state, streamingResp)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update sends a PATCH when only streaming_enabled changed, else a PUT with the full config.
func (a *AppServiceAuditLogStreaming) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state providerschema.AppServiceAuditLogStreaming
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if plan.OutputType.Equal(state.OutputType) &&
		plan.Credentials.Equal(state.Credentials) &&
		plan.DisabledAppEndpoints.Equal(state.DisabledAppEndpoints) {
		tflog.Debug(ctx, "updating audit log streaming", map[string]interface{}{"method": http.MethodPatch})
		err = a.patch(ctx, plan)
	} else {
		tflog.Debug(ctx, "updating audit log streaming", map[string]interface{}{"method": http.MethodPut})
		err = a.put(ctx, &plan)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating App Service Audit Log Streaming",
			"Could not update audit log streaming, unexpected error: "+err.Error(),
		)
		return
	}

	a.waitAndSetState(ctx, plan, &resp.State, &resp.Diagnostics)
}

// Delete disables audit log streaming and removes its config by sending a PUT with only streamingEnabled=false.
func (a *AppServiceAuditLogStreaming) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceAuditLogStreaming
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		api.EndpointCfg{Url: a.url(state), Method: http.MethodPut, SuccessStatus: http.StatusAccepted},
		apigen.PutAuditLogStreamingRequest{StreamingEnabled: false},
		a.Token,
		nil,
	)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting App Service Audit Log Streaming",
			"Could not disable audit log streaming, unexpected error: "+api.ParseError(err),
		)
		return
	}

	// The config is removed when the GET no longer returns an output type.
	configRemoved := func(r *apigen.GetAuditLogStreamingResponse) bool { return r.OutputType == nil }
	if _, err := a.waitForFinalState(ctx, state, configRemoved); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting App Service Audit Log Streaming",
			"Audit log streaming disable was initiated, but checking its state failed: "+err.Error(),
		)
	}
}

// ImportState imports a remote resource that is not managed by Terraform.
func (a *AppServiceAuditLogStreaming) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

func (a *AppServiceAuditLogStreaming) url(m providerschema.AppServiceAuditLogStreaming) string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLogStreaming",
		a.HostURL,
		m.OrganizationId.ValueString(),
		m.ProjectId.ValueString(),
		m.ClusterId.ValueString(),
		m.AppServiceId.ValueString(),
	)
}

// put sends the full config. disabledAppEndpoints is always sent, so removing it from the config clears it.
func (a *AppServiceAuditLogStreaming) put(ctx context.Context, plan *providerschema.AppServiceAuditLogStreaming) error {
	creds, diags := plan.AsLogStreamingCredentials(ctx)
	if diags.HasError() {
		return fmt.Errorf("failed to extract credentials: %s", diags.Errors())
	}
	if creds == nil {
		return fmt.Errorf("credentials are required")
	}
	v, err := logStreamingCredentials(ctx, plan.OutputType.ValueString(), creds)
	if err != nil {
		return err
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}
	var apiCreds apigen.PutAuditLogStreamingRequest_Credentials
	if err := apiCreds.UnmarshalJSON(b); err != nil {
		return fmt.Errorf("failed to set credentials: %w", err)
	}

	// ElementsAs turns a null set into a nil slice, which marshals to null, so only convert a set that is present.
	disabledAppEndpoints := []string{}
	if !plan.DisabledAppEndpoints.IsNull() {
		if diags := plan.DisabledAppEndpoints.ElementsAs(ctx, &disabledAppEndpoints, false); diags.HasError() {
			return fmt.Errorf("failed to read disabled_app_endpoints: %s", diags.Errors())
		}
	}

	outputType := apigen.PutAuditLogStreamingRequestOutputType(plan.OutputType.ValueString())
	body := apigen.PutAuditLogStreamingRequest{
		StreamingEnabled:     plan.StreamingEnabled.ValueBool(),
		OutputType:           &outputType,
		Credentials:          &apiCreds,
		DisabledAppEndpoints: &disabledAppEndpoints,
	}

	_, err = a.ClientV1.ExecuteWithRetry(
		ctx,
		api.EndpointCfg{Url: a.url(*plan), Method: http.MethodPut, SuccessStatus: http.StatusAccepted},
		body,
		a.Token,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s", api.ParseError(err))
	}
	return nil
}

func (a *AppServiceAuditLogStreaming) patch(ctx context.Context, plan providerschema.AppServiceAuditLogStreaming) error {
	body := apigen.PatchAuditLogStreamingRequest{
		Op:    apigen.Update,
		Path:  "/streamingEnabled",
		Value: plan.StreamingEnabled.ValueBool(),
	}
	_, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		api.EndpointCfg{Url: a.url(plan), Method: http.MethodPatch, SuccessStatus: http.StatusAccepted},
		body,
		a.Token,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s", api.ParseError(err))
	}
	return nil
}

func (a *AppServiceAuditLogStreaming) get(ctx context.Context, m providerschema.AppServiceAuditLogStreaming) (*apigen.GetAuditLogStreamingResponse, error) {
	response, err := a.ClientV1.ExecuteWithRetry(
		ctx,
		api.EndpointCfg{Url: a.url(m), Method: http.MethodGet, SuccessStatus: http.StatusOK},
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var streamingResp apigen.GetAuditLogStreamingResponse
	if err := json.Unmarshal(response.Body, &streamingResp); err != nil {
		return nil, fmt.Errorf("could not unmarshal audit log streaming response: %w", err)
	}
	return &streamingResp, nil
}

// waitForFinalState polls the GET endpoint until the streaming state is not transitional and done reports the requested outcome.
// It fails on the errored state.
func (a *AppServiceAuditLogStreaming) waitForFinalState(
	ctx context.Context,
	m providerschema.AppServiceAuditLogStreaming,
	done func(*apigen.GetAuditLogStreamingResponse) bool,
) (*apigen.GetAuditLogStreamingResponse, error) {
	const timeout = 3 * time.Minute
	const interval = 3 * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for audit log streaming to reach a final state: %w", ctx.Err())
		case <-ticker.C:
			streamingResp, err := a.get(ctx, m)
			if err != nil {
				return nil, fmt.Errorf("%s", api.ParseError(err))
			}

			var state string
			if streamingResp.LogStreamingState != nil {
				state = *streamingResp.LogStreamingState
			}
			switch state {
			case "enabling", "disabling", "pausing":
			case "errored":
				return nil, fmt.Errorf("audit log streaming reached the errored state")
			default:
				if done(streamingResp) {
					return streamingResp, nil
				}
			}
			tflog.Info(ctx, "waiting for audit log streaming state transition", map[string]interface{}{"state": state})
		}
	}
}

// waitAndSetState waits for the state transition after a PUT or PATCH and saves the refreshed state.
func (a *AppServiceAuditLogStreaming) waitAndSetState(
	ctx context.Context,
	plan providerschema.AppServiceAuditLogStreaming,
	state *tfsdk.State,
	diags *diag.Diagnostics,
) {
	// Save the applied plan first, so a failed wait does not lose the config or the credentials, which the API does not return.
	plan.LogStreamingState = types.StringNull()
	diags.Append(state.Set(ctx, plan)...)
	if diags.HasError() {
		return
	}

	wantEnabled := plan.StreamingEnabled.ValueBool()
	streamingEnabledMatches := func(r *apigen.GetAuditLogStreamingResponse) bool {
		return r.StreamingEnabled != nil && *r.StreamingEnabled == wantEnabled
	}
	streamingResp, err := a.waitForFinalState(ctx, plan, streamingEnabledMatches)
	if err != nil {
		diags.AddError(
			"Error Refreshing App Service Audit Log Streaming",
			"Audit log streaming change was initiated, but checking its state failed. Run `terraform apply --refresh-only` to update the state: "+err.Error(),
		)
		return
	}

	diags.Append(setAuditLogStreamingState(ctx, &plan, streamingResp)...)
	if diags.HasError() {
		return
	}
	diags.Append(state.Set(ctx, plan)...)
}

// setAuditLogStreamingState copies the API response into m. Credentials are not returned by the API, so m keeps them.
func setAuditLogStreamingState(ctx context.Context, m *providerschema.AppServiceAuditLogStreaming, streamingResp *apigen.GetAuditLogStreamingResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	m.StreamingEnabled = types.BoolPointerValue(streamingResp.StreamingEnabled)
	m.OutputType = types.StringPointerValue(streamingResp.OutputType)
	m.LogStreamingState = types.StringPointerValue(streamingResp.LogStreamingState)

	// The API returns no list when nothing is excluded. Keep null in that case to match an omitted attribute.
	if streamingResp.DisabledAppEndpoints == nil || len(*streamingResp.DisabledAppEndpoints) == 0 {
		if m.DisabledAppEndpoints.IsNull() || m.DisabledAppEndpoints.IsUnknown() {
			m.DisabledAppEndpoints = types.SetNull(types.StringType)
			return diags
		}
	}
	endpoints := []string{}
	if streamingResp.DisabledAppEndpoints != nil {
		endpoints = *streamingResp.DisabledAppEndpoints
	}
	m.DisabledAppEndpoints, diags = types.SetValueFrom(ctx, types.StringType, endpoints)
	return diags
}
