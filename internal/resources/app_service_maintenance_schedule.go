package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = (*AppServiceMaintenanceSchedule)(nil)
	_ resource.ResourceWithConfigure   = (*AppServiceMaintenanceSchedule)(nil)
	_ resource.ResourceWithImportState = (*AppServiceMaintenanceSchedule)(nil)
)

// AppServiceMaintenanceSchedule is the App Service maintenance schedule resource implementation.
type AppServiceMaintenanceSchedule struct {
	*providerschema.Data
}

// NewAppServiceMaintenanceSchedule is a helper function to simplify the provider implementation.
func NewAppServiceMaintenanceSchedule() resource.Resource {
	return &AppServiceMaintenanceSchedule{}
}

// Metadata returns the App Service maintenance schedule resource type name.
func (a *AppServiceMaintenanceSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_service_maintenance_schedule"
}

// Schema defines the schema for the App Service maintenance schedule resource.
func (a *AppServiceMaintenanceSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AppServiceMaintenanceScheduleSchema()
}

// Create adopts an existing maintenance record and reschedules it when its queue time differs
// from the planned one. It never creates a maintenance record.
func (a *AppServiceMaintenanceSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AppServiceMaintenanceSchedule
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		maintenanceId  = plan.MaintenanceId.ValueString()
	)

	queueAt, diags := plan.QueueAt.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	maintenance, err := a.getMaintenance(ctx, organizationId, projectId, clusterId, appServiceId, maintenanceId)
	if err != nil {
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			resp.Diagnostics.AddError(
				"App Service maintenance not found",
				fmt.Sprintf("Maintenance %s does not exist on app service %s. List the maintenance of the app service with the "+
					"couchbase-capella_app_service_maintenances data source: %s", maintenanceId, appServiceId, api.ParseError(err)),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance",
			"Could not read maintenance "+maintenanceId+" of app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	if maintenance.InitiatedBy == appservice_maintenance.InitiatedBySelfService {
		resp.Diagnostics.AddError(
			"App Service maintenance is a self-service upgrade",
			fmt.Sprintf("Maintenance %s of app service %s was initiated by a self-service upgrade. Its queue time is managed "+
				"by the couchbase-capella_app_service_self_service_upgrade resource.", maintenanceId, appServiceId),
		)
		return
	}

	if maintenance.QueueAt == nil || !maintenance.QueueAt.Equal(queueAt) {
		if err = a.rescheduleMaintenance(ctx, organizationId, projectId, clusterId, appServiceId, maintenanceId, queueAt); err != nil {
			resp.Diagnostics.AddError(
				"Error rescheduling App Service maintenance",
				"Could not reschedule maintenance "+maintenanceId+" of app service "+appServiceId+": "+api.ParseError(err),
			)
			return
		}
	}

	a.refreshState(ctx, plan, &resp.State, &resp.Diagnostics)
}

// Read refreshes the state with the maintenance record from Capella. A 404 means the record no
// longer exists; any other error leaves the state untouched.
func (a *AppServiceMaintenanceSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AppServiceMaintenanceSchedule
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	IDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance",
			"Could not read maintenance "+state.MaintenanceId.String()+": "+err.Error(),
		)
		return
	}

	maintenance, err := a.getMaintenance(ctx,
		IDs[providerschema.OrganizationId],
		IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId],
		IDs[providerschema.AppServiceId],
		IDs[providerschema.MaintenanceId],
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading App Service maintenance",
			"Could not read maintenance "+IDs[providerschema.MaintenanceId]+": "+errString,
		)
		return
	}

	refreshedState, diags := mapMaintenanceResponseToState(ctx, maintenance, IDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, refreshedState)...)
}

// Update reschedules the maintenance to the planned queue time. Every other attribute requires replace.
func (a *AppServiceMaintenanceSchedule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AppServiceMaintenanceSchedule
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	queueAt, diags := plan.QueueAt.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		organizationId = plan.OrganizationId.ValueString()
		projectId      = plan.ProjectId.ValueString()
		clusterId      = plan.ClusterId.ValueString()
		appServiceId   = plan.AppServiceId.ValueString()
		maintenanceId  = plan.MaintenanceId.ValueString()
	)

	if err := a.rescheduleMaintenance(ctx, organizationId, projectId, clusterId, appServiceId, maintenanceId, queueAt); err != nil {
		resp.Diagnostics.AddError(
			"Error rescheduling App Service maintenance",
			"Could not reschedule maintenance "+maintenanceId+" of app service "+appServiceId+": "+api.ParseError(err),
		)
		return
	}

	a.refreshState(ctx, plan, &resp.State, &resp.Diagnostics)
}

// Delete only removes the resource from the state. The maintenance keeps its last queue time.
func (a *AppServiceMaintenanceSchedule) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// ImportState imports the schedule of an existing maintenance record. The import ID is a comma
// separated string of the IDs, for example:
// "organization_id=<orgId>,project_id=<projId>,cluster_id=<clusterId>,app_service_id=<appServiceId>,maintenance_id=<maintenanceId>".
func (a *AppServiceMaintenanceSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("maintenance_id"), req, resp)
}

// Configure adds the provider configured client to the App Service maintenance schedule resource.
func (a *AppServiceMaintenanceSchedule) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// refreshState reads the maintenance back and stores it, so the state reflects what Capella accepted.
func (a *AppServiceMaintenanceSchedule) refreshState(ctx context.Context, plan providerschema.AppServiceMaintenanceSchedule, state *tfsdk.State, diagnostics *diag.Diagnostics) {
	IDs := map[providerschema.Attr]string{
		providerschema.OrganizationId: plan.OrganizationId.ValueString(),
		providerschema.ProjectId:      plan.ProjectId.ValueString(),
		providerschema.ClusterId:      plan.ClusterId.ValueString(),
		providerschema.AppServiceId:   plan.AppServiceId.ValueString(),
		providerschema.MaintenanceId:  plan.MaintenanceId.ValueString(),
	}

	maintenance, err := a.getMaintenance(ctx, IDs[providerschema.OrganizationId], IDs[providerschema.ProjectId],
		IDs[providerschema.ClusterId], IDs[providerschema.AppServiceId], IDs[providerschema.MaintenanceId])
	if err != nil {
		diagnostics.AddError(
			"Error reading App Service maintenance",
			"Could not read maintenance "+IDs[providerschema.MaintenanceId]+" of app service "+
				IDs[providerschema.AppServiceId]+": "+api.ParseError(err),
		)
		return
	}

	refreshedState, diags := mapMaintenanceResponseToState(ctx, maintenance, IDs)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return
	}
	diagnostics.Append(state.Set(ctx, refreshedState)...)
}

func (a *AppServiceMaintenanceSchedule) maintenanceURL(organizationId, projectId, clusterId, appServiceId, maintenanceId string) string {
	return fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/maintenances/%s",
		a.HostURL, organizationId, projectId, clusterId, appServiceId, url.PathEscape(maintenanceId))
}

func (a *AppServiceMaintenanceSchedule) getMaintenance(ctx context.Context, organizationId, projectId, clusterId, appServiceId, maintenanceId string) (*appservice_maintenance.GetMaintenanceResponse, error) {
	cfg := api.EndpointCfg{Url: a.maintenanceURL(organizationId, projectId, clusterId, appServiceId, maintenanceId), Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, nil, a.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
	}

	var maintenance appservice_maintenance.GetMaintenanceResponse
	if err = json.Unmarshal(response.Body, &maintenance); err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrUnmarshallingResponse, err)
	}
	return &maintenance, nil
}

func (a *AppServiceMaintenanceSchedule) rescheduleMaintenance(ctx context.Context, organizationId, projectId, clusterId, appServiceId, maintenanceId string, queueAt time.Time) error {
	cfg := api.EndpointCfg{Url: a.maintenanceURL(organizationId, projectId, clusterId, appServiceId, maintenanceId), Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
	_, err := a.ClientV1.ExecuteWithRetry(ctx, cfg, appservice_maintenance.RescheduleMaintenanceRequest{QueueAt: queueAt}, a.Token, nil)
	return err
}

func mapMaintenanceResponseToState(
	ctx context.Context, resp *appservice_maintenance.GetMaintenanceResponse, IDs map[providerschema.Attr]string,
) (*providerschema.AppServiceMaintenanceSchedule, diag.Diagnostics) {
	audit := providerschema.NewCouchbaseAuditData(resp.Audit)
	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, diags
	}

	data := providerschema.NewAppServiceMaintenanceData(resp, auditObj)
	return &providerschema.AppServiceMaintenanceSchedule{
		OrganizationId: types.StringValue(IDs[providerschema.OrganizationId]),
		ProjectId:      types.StringValue(IDs[providerschema.ProjectId]),
		ClusterId:      types.StringValue(IDs[providerschema.ClusterId]),
		AppServiceId:   types.StringValue(IDs[providerschema.AppServiceId]),
		MaintenanceId:  types.StringValue(IDs[providerschema.MaintenanceId]),
		QueueAt:        data.QueueAt,
		Kind:           data.Kind,
		InitiatedBy:    data.InitiatedBy,
		VersionId:      data.VersionId,
		CurrentState:   data.CurrentState,
		CompletedAt:    data.CompletedAt,
		Audit:          data.Audit,
	}, nil
}
