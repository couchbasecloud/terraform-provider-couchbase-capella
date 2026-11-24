package resources

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*AppServicePublicAccess)(nil)
	_ resource.ResourceWithConfigure   = (*AppServicePublicAccess)(nil)
	_ resource.ResourceWithImportState = (*AppServicePublicAccess)(nil)
)

var (
	// appServicePublicAccessTimeout bounds the 409 retries and the status polling of one public access change.
	appServicePublicAccessTimeout = 30 * time.Minute

	appServicePublicAccessPollInterval = 10 * time.Second
)

// AppServicePublicAccess is the app service public access resource implementation.
type AppServicePublicAccess struct {
	*providerschema.Data
}

// NewAppServicePublicAccess is a helper function to simplify the provider implementation.
func NewAppServicePublicAccess() resource.Resource {
	return &AppServicePublicAccess{}
}

// Metadata returns the app service public access resource type name.
func (a *AppServicePublicAccess) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_public_access"
}

// Schema defines the schema for the app service public access resource.
func (a *AppServicePublicAccess) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServicePublicAccessSchema()
}

// Configure adds the provider configured client to the app service public access resource.
func (a *AppServicePublicAccess) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState imports the public access of an App Service into Terraform.
func (a *AppServicePublicAccess) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

// Create sets the requested public access mode. If the change fails, Create saves the state with a warning, so that
// Terraform does not taint the resource and retry a change that only Support can recover.
func (a *AppServicePublicAccess) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServicePublicAccess
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.changePublicAccess(
		ctx,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
		appservice.PublicAccess(plan.PublicAccess.ValueString()),
	)
	switch {
	case stderrors.Is(err, errors.ErrAppServicePublicAccessChangeFailed):
		resp.Diagnostics.AddWarning(publicAccessChangeFailedSummary, publicAccessChangeFailedDetail(refreshedState))
	case err != nil:
		resp.Diagnostics.AddError(
			"Error changing App Service public access",
			"Could not change public access of App Service "+plan.AppServiceId.ValueString()+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read refreshes the Terraform state with the public access that Capella returns.
func (a *AppServicePublicAccess) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServicePublicAccess
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service public access",
			"Could not validate the public access of App Service "+state.AppServiceId.String()+": "+err.Error(),
		)
		return
	}

	refreshedState, err := a.getPublicAccess(
		ctx,
		IDs[providerschema.OrganizationId],
		IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId],
		IDs[providerschema.AppServiceId],
	)
	if err != nil {
		if resourceNotFound, _ := api.CheckResourceNotFoundError(err); resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service public access",
			"Could not read public access of App Service "+IDs[providerschema.AppServiceId]+": "+api.ParseError(err),
		)
		return
	}

	if appservice.PublicAccess(refreshedState.PublicAccessStatus.ValueString()).IsFailed() {
		resp.Diagnostics.AddWarning(publicAccessChangeFailedSummary, publicAccessChangeFailedDetail(refreshedState))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update sets the requested public access mode. Unlike Create, a failed change is an error.
func (a *AppServicePublicAccess) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppServicePublicAccess
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.changePublicAccess(
		ctx,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
		appservice.PublicAccess(plan.PublicAccess.ValueString()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error changing App Service public access",
			"Could not change public access of App Service "+plan.AppServiceId.ValueString()+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Delete removes the resource from the Terraform state only. Public access stays as it is.
func (a *AppServicePublicAccess) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Info(ctx, "removing app service public access from state, public access is not changed")
}

// changePublicAccess requests the public access mode, then polls the App Service until public access reaches it.
// If the change fails, it returns the refreshed state with errors.ErrAppServicePublicAccessChangeFailed.
func (a *AppServicePublicAccess) changePublicAccess(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	publicAccess appservice.PublicAccess,
) (*providerschema.AppServicePublicAccess, error) {
	ctx, cancel := context.WithTimeout(ctx, appServicePublicAccessTimeout)
	defer cancel()

	if err := a.requestPublicAccess(ctx, organizationId, projectId, clusterId, appServiceId, publicAccess); err != nil {
		return nil, err
	}

	return a.checkAppServicePublicAccessStatus(ctx, organizationId, projectId, clusterId, appServiceId, publicAccess)
}

// requestPublicAccess calls POST or DELETE .../publicAccessRestriction. A 409 means that a public access change is already in progress,
// so the call is retried until that change completes.
func (a *AppServicePublicAccess) requestPublicAccess(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	publicAccess appservice.PublicAccess,
) error {
	method := http.MethodDelete
	if publicAccess == appservice.PublicAccessRestricted {
		method = http.MethodPost
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/publicAccessRestriction",
		a.HostURL, organizationId, projectId, clusterId, appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}

	for {
		_, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)

		var apiErr *api.Error
		if !stderrors.As(err, &apiErr) || apiErr.HttpStatusCode != http.StatusConflict {
			return err
		}

		tflog.Info(ctx, "a public access change is already in progress, retrying", map[string]any{
			"app_service_id": appServiceId,
		})

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for an already ongoing public access change to complete: %w", ctx.Err())
		case <-time.After(appServicePublicAccessPollInterval):
		}
	}
}

// checkAppServicePublicAccessStatus polls the App Service until public access is the requested mode or the change
// fails.
func (a *AppServicePublicAccess) checkAppServicePublicAccessStatus(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
	publicAccess appservice.PublicAccess,
) (*providerschema.AppServicePublicAccess, error) {
	for {
		refreshedState, err := a.getPublicAccess(ctx, organizationId, projectId, clusterId, appServiceId)
		if err != nil {
			return nil, err
		}

		status := appservice.PublicAccess(refreshedState.PublicAccessStatus.ValueString())
		switch {
		case status == publicAccess:
			return refreshedState, nil
		case status.IsFailed():
			return refreshedState, errors.ErrAppServicePublicAccessChangeFailed
		}
		tflog.Info(ctx, "waiting for app service public access to change", map[string]any{"status": status})

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("app service public access change timed out after initiation: %w", ctx.Err())
		case <-time.After(appServicePublicAccessPollInterval):
		}
	}
}

// getPublicAccess reads public access from the App Service.
func (a *AppServicePublicAccess) getPublicAccess(
	ctx context.Context,
	organizationId, projectId, clusterId, appServiceId string,
) (*providerschema.AppServicePublicAccess, error) {
	appServiceResp, err := (&AppService{Data: a.Data}).getAppService(ctx, organizationId, projectId, clusterId, appServiceId)
	if err != nil {
		return nil, err
	}

	if appServiceResp.PublicAccess == nil {
		return nil, errors.ErrAppServicePublicAccessNotAvailable
	}

	return providerschema.NewAppServicePublicAccess(organizationId, projectId, clusterId, appServiceId, *appServiceResp.PublicAccess), nil
}

const publicAccessChangeFailedSummary = "App Service public access change failed"

func publicAccessChangeFailedDetail(state *providerschema.AppServicePublicAccess) string {
	return fmt.Sprintf(
		"The last public access change of App Service %s failed (public_access_status is %s). "+
			"The change will not be automatically retried. Please contact Couchbase Capella Support.",
		state.AppServiceId.ValueString(), state.PublicAccessStatus.ValueString(),
	)
}
