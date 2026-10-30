package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var (
	_ resource.Resource                = (*AppEndpointAuditLogging)(nil)
	_ resource.ResourceWithConfigure   = (*AppEndpointAuditLogging)(nil)
	_ resource.ResourceWithImportState = (*AppEndpointAuditLogging)(nil)
)

type AppEndpointAuditLogging struct {
	*providerschema.Data
}

// NewAppEndpointAuditLogging is used in (p *capellaProvider) Resources for building the provider.
func NewAppEndpointAuditLogging() resource.Resource {
	return &AppEndpointAuditLogging{}
}

// Metadata returns the App Endpoint audit logging resource type name.
func (a *AppEndpointAuditLogging) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_endpoint_audit_logging"
}

// Schema defines the schema for the App Endpoint audit logging resource.
func (a *AppEndpointAuditLogging) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppEndpointAuditLoggingSchema()
}

// Configure is used to configure the App Endpoint audit logging resource with the provider data.
func (a *AppEndpointAuditLogging) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create sets the audit logging configuration of the App Endpoint.
// There is no create endpoint, so create is an update.
func (a *AppEndpointAuditLogging) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppEndpointAuditLogging
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, diags := a.upsert(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Read reads the audit logging configuration of the App Endpoint.
func (a *AppEndpointAuditLogging) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppEndpointAuditLogging
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Validating App Endpoint Audit Logging State",
			"Could not validate App Endpoint audit logging state: "+err.Error(),
		)
		return
	}

	var (
		organizationId  = IDs[providerschema.OrganizationId]
		projectId       = IDs[providerschema.ProjectId]
		clusterId       = IDs[providerschema.ClusterId]
		appServiceId    = IDs[providerschema.AppServiceId]
		appEndpointName = IDs[providerschema.AppEndpointName]
	)

	refreshedState, diags, err := a.refresh(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		if handled, forbiddenErr := handleAppEndpointForbidden(ctx, err, a.Data, resp, organizationId, projectId, clusterId, appServiceId, appEndpointName); handled {
			return
		} else if forbiddenErr != nil {
			errString = forbiddenErr.Error()
		}
		resp.Diagnostics.AddError(
			"Error Reading App Endpoint Audit Logging",
			"Could not read audit logging configuration of App Endpoint "+appEndpointName+": "+errString,
		)
		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update updates the audit logging configuration of the App Endpoint.
func (a *AppEndpointAuditLogging) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppEndpointAuditLogging
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, diags := a.upsert(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Delete removes the resource from the Terraform state only, because there is no delete endpoint.
// The audit logging configuration in Capella stays in its last state.
func (a *AppEndpointAuditLogging) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// ImportState imports the audit logging configuration of an App Endpoint.
// example: app_endpoint_name=endpoint1,app_service_id=app123,cluster_id=cluster123,project_id=proj123,organization_id=org123
func (a *AppEndpointAuditLogging) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("app_endpoint_name"), req, resp)
}

// upsert sends the known attributes of the plan to the PUT endpoint and returns the refreshed state.
// Unknown attributes are left out of the request, so Capella keeps their current values.
func (a *AppEndpointAuditLogging) upsert(ctx context.Context, plan providerschema.AppEndpointAuditLogging) (*providerschema.AppEndpointAuditLogging, diag.Diagnostics) {
	var diags diag.Diagnostics
	body := apigen.PutAppEndpointAuditLogConfigJSONRequestBody{}

	if !plan.AuditEnabled.IsUnknown() && !plan.AuditEnabled.IsNull() {
		body.AuditEnabled = plan.AuditEnabled.ValueBoolPointer()
	}

	if !plan.EnabledEventIds.IsUnknown() && !plan.EnabledEventIds.IsNull() {
		var ids []int64
		diags.Append(plan.EnabledEventIds.ElementsAs(ctx, &ids, false)...)
		eventIds := make([]struct {
			Id *int `json:"id,omitempty"`
		}, len(ids))
		for i, id := range ids {
			v := int(id)
			eventIds[i].Id = &v
		}
		body.EnabledEventIds = &eventIds
	}

	body.DisabledUsers = toDisabledUserRoles(ctx, plan.DisabledUsers, &diags)
	body.DisabledRoles = toDisabledUserRoles(ctx, plan.DisabledRoles, &diags)
	if diags.HasError() {
		return nil, diags
	}

	var (
		organizationId  = plan.OrganizationId.ValueString()
		projectId       = plan.ProjectId.ValueString()
		clusterId       = plan.ClusterId.ValueString()
		appServiceId    = plan.AppServiceId.ValueString()
		appEndpointName = plan.AppEndpointName.ValueString()
	)

	cfg := api.EndpointCfg{
		Url:           a.auditLogURL(organizationId, projectId, clusterId, appServiceId, appEndpointName),
		Method:        http.MethodPut,
		SuccessStatus: http.StatusNoContent,
	}
	if _, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, body, a.Token, nil); err != nil {
		diags.AddError(
			"Error Setting App Endpoint Audit Logging",
			"Could not set audit logging configuration of App Endpoint "+appEndpointName+": "+api.ParseError(err),
		)
		return nil, diags
	}

	refreshedState, refreshDiags, err := a.refresh(ctx, organizationId, projectId, clusterId, appServiceId, appEndpointName)
	if err != nil {
		diags.AddError(
			"Error Reading App Endpoint Audit Logging",
			"Could not read audit logging configuration of App Endpoint "+appEndpointName+" after it was set: "+api.ParseError(err),
		)
		return nil, diags
	}
	diags.Append(refreshDiags...)

	return refreshedState, diags
}

// refresh gets the audit logging configuration and converts it to Terraform state.
func (a *AppEndpointAuditLogging) refresh(ctx context.Context, organizationId, projectId, clusterId, appServiceId, appEndpointName string) (*providerschema.AppEndpointAuditLogging, diag.Diagnostics, error) {
	cfg := api.EndpointCfg{
		Url:           a.auditLogURL(organizationId, projectId, clusterId, appServiceId, appEndpointName),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, nil, err
	}

	var auditLog apigen.GetAppEndpointAuditLogResponse
	if err = json.Unmarshal(response.Body, &auditLog); err != nil {
		return nil, nil, err
	}

	var diags diag.Diagnostics
	state := providerschema.AppEndpointAuditLogging{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		AuditEnabled:    types.BoolValue(auditLog.AuditEnabled != nil && *auditLog.AuditEnabled),
	}

	ids := []int64{}
	if auditLog.EnabledEventIds != nil {
		for _, e := range *auditLog.EnabledEventIds {
			if e.Id != nil {
				ids = append(ids, int64(*e.Id))
			}
		}
	}
	var d diag.Diagnostics
	state.EnabledEventIds, d = types.SetValueFrom(ctx, types.Int64Type, ids)
	diags.Append(d...)

	state.DisabledUsers = fromDisabledUserRoles(ctx, auditLog.DisabledUsers, &diags)
	state.DisabledRoles = fromDisabledUserRoles(ctx, auditLog.DisabledRoles, &diags)

	return &state, diags, nil
}

func (a *AppEndpointAuditLogging) auditLogURL(organizationId, projectId, clusterId, appServiceId, appEndpointName string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/auditLog",
		a.HostURL,
		organizationId,
		projectId,
		clusterId,
		appServiceId,
		appEndpointName,
	)
}

// toDisabledUserRoles returns nil for an unknown or null set, so the field is left out of the request.
func toDisabledUserRoles(ctx context.Context, set types.Set, diags *diag.Diagnostics) *apigen.DisabledUserRoles {
	if set.IsUnknown() || set.IsNull() {
		return nil
	}

	var elems []providerschema.AppEndpointAuditLogDisabledUserRole
	diags.Append(set.ElementsAs(ctx, &elems, false)...)
	userRoles := make(apigen.DisabledUserRoles, len(elems))
	for i, e := range elems {
		userRoles[i] = apigen.DisabledUserRole{Domain: e.Domain.ValueString(), Name: e.Name.ValueString()}
	}
	return &userRoles
}

func fromDisabledUserRoles(ctx context.Context, userRoles *apigen.DisabledUserRoles, diags *diag.Diagnostics) types.Set {
	elems := []providerschema.AppEndpointAuditLogDisabledUserRole{}
	if userRoles != nil {
		for _, u := range *userRoles {
			elems = append(elems, providerschema.AppEndpointAuditLogDisabledUserRole{
				Domain: types.StringValue(u.Domain),
				Name:   types.StringValue(u.Name),
			})
		}
	}
	set, d := types.SetValueFrom(ctx, providerschema.AppEndpointAuditLogDisabledUserRoleType, elems)
	diags.Append(d...)
	return set
}
