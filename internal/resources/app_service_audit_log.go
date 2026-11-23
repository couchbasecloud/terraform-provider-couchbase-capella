package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*AppServiceAuditLog)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceAuditLog)(nil)
	_ resource.ResourceWithImportState = (*AppServiceAuditLog)(nil)
)

type AppServiceAuditLog struct {
	*providerschema.Data
}

// NewAppServiceAuditLog is used in (p *capellaProvider) Resources for building the provider.
func NewAppServiceAuditLog() resource.Resource {
	return &AppServiceAuditLog{}
}

// Metadata returns the App Service audit log resource type name.
func (a *AppServiceAuditLog) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_audit_logging"
}

// Schema defines the schema for the App Service audit log resource.
func (a *AppServiceAuditLog) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceAuditLogSchema()
}

// Configure is used to configure the App Service audit log resource with the provider data.
func (a *AppServiceAuditLog) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create enables audit logging on the App Service.
func (a *AppServiceAuditLog) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceAuditLog
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.putAuditLogState(
		ctx,
		plan.OrganizationId.ValueString(),
		plan.ProjectId.ValueString(),
		plan.ClusterId.ValueString(),
		plan.AppServiceId.ValueString(),
		true,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Enabling App Service Audit Logging",
			"Could not enable audit logging on App Service "+plan.AppServiceId.ValueString()+": "+api.ParseError(err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read checks that audit logging is still enabled on the App Service.
// If audit logging is disabled, the resource is removed from state so that the next apply enables it again.
func (a *AppServiceAuditLog) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceAuditLog
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating App Service Audit Log State",
			"Could not validate App Service audit log state: "+err.Error(),
		)
		return
	}

	var (
		organizationId = IDs[providerschema.OrganizationId]
		projectId      = IDs[providerschema.ProjectId]
		clusterId      = IDs[providerschema.ClusterId]
		appServiceId   = IDs[providerschema.AppServiceId]
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLog",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading App Service Audit Log State",
			"Could not read audit log state of App Service "+appServiceId+": "+errString,
		)
		return
	}

	auditLog := apigen.CreateAppServiceAuditLogRequest{}
	if err = json.Unmarshal(response.Body, &auditLog); err != nil {
		resp.Diagnostics.AddError(
			"Error Unmarshalling App Service Audit Log State",
			"Could not unmarshal App Service audit log state: "+err.Error(),
		)
		return
	}

	if !auditLog.AuditEnabled {
		tflog.Info(ctx, "audit logging is disabled on the App Service, removing resource from state file")
		resp.State.RemoveResource(ctx)
		return
	}

	// Set the IDs from Validate, since after an import the state holds only the import string.
	refreshedState := providerschema.AppServiceAuditLog{
		OrganizationId: types.StringValue(organizationId),
		ProjectId:      types.StringValue(projectId),
		ClusterId:      types.StringValue(clusterId),
		AppServiceId:   types.StringValue(appServiceId),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update is not supported. All attributes require replacement,
// so the framework recreates the resource and never calls Update.
func (a *AppServiceAuditLog) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete disables audit logging on the App Service.
func (a *AppServiceAuditLog) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AppServiceAuditLog
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Disabling App Service Audit Logging",
			"Could not validate App Service audit log state: "+err.Error(),
		)
		return
	}

	appServiceId := IDs[providerschema.AppServiceId]
	err = a.putAuditLogState(
		ctx,
		IDs[providerschema.OrganizationId],
		IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId],
		appServiceId,
		false,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			return
		}
		resp.Diagnostics.AddError(
			"Error Disabling App Service Audit Logging",
			"Could not disable audit logging on App Service "+appServiceId+": "+errString,
		)
	}
}

// ImportState imports the audit log state of an App Service.
// example: app_service_id=app123,cluster_id=cluster123,project_id=proj123,organization_id=org123
func (a *AppServiceAuditLog) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_service_id"), req, resp)
}

func (a *AppServiceAuditLog) putAuditLogState(ctx context.Context, organizationId, projectId, clusterId, appServiceId string, enabled bool) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/auditLog",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, apigen.CreateAppServiceAuditLogRequest{AuditEnabled: enabled}, a.Token, nil)
	return err
}
