package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice_maintenance"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// Fallback descriptions for the App Service maintenance fields, used only until the published
// OpenAPI spec carries them. The leading "\n - " matches the bullet format buildEnhancedDescription
// emits, so the generated docs do not change once the spec takes over.
const (
	MaintenanceIdDescription           = "\n - The ID of the App Service maintenance record, as listed by the `couchbase-capella_app_service_maintenances` data source."
	MaintenanceKindDescription         = "\n - The kind of maintenance, for example `upgrade`."
	MaintenanceInitiatedByDescription  = "\n - Who initiated the maintenance: `couchbase` or `selfService`."
	MaintenanceVersionIdDescription    = "\n - The ID of the release the maintenance moves the App Service to. Null when no release contains the maintenance image."
	MaintenanceCurrentStateDescription = "\n - The current state of the maintenance, for example `pending`. Only `pending` maintenance can be rescheduled."
	MaintenanceQueueAtDescription      = "\n - The RFC3339 timestamp at which the maintenance is queued to start. It must be in the future and within the maintenance scheduling window."
	MaintenanceCompletedAtDescription  = "\n - The RFC3339 timestamp at which the maintenance completed. Null until it completes."
)

// AppServiceMaintenanceSchedule is the queue time of an existing, Couchbase-initiated App Service
// maintenance record.
type AppServiceMaintenanceSchedule struct {
	OrganizationId types.String      `tfsdk:"organization_id"`
	ProjectId      types.String      `tfsdk:"project_id"`
	ClusterId      types.String      `tfsdk:"cluster_id"`
	AppServiceId   types.String      `tfsdk:"app_service_id"`
	MaintenanceId  types.String      `tfsdk:"maintenance_id"`
	QueueAt        timetypes.RFC3339 `tfsdk:"queue_at"`
	Kind           types.String      `tfsdk:"kind"`
	InitiatedBy    types.String      `tfsdk:"initiated_by"`
	VersionId      types.String      `tfsdk:"version_id"`
	CurrentState   types.String      `tfsdk:"current_state"`
	CompletedAt    timetypes.RFC3339 `tfsdk:"completed_at"`
	Audit          types.Object      `tfsdk:"audit"`
}

// Validate is used to verify that IDs have been properly imported.
func (a AppServiceMaintenanceSchedule) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
		MaintenanceId:  a.MaintenanceId,
	}

	IDs, err := validateSchemaState(state, MaintenanceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// AppServiceMaintenances is the list of maintenance records of an App Service.
type AppServiceMaintenances struct {
	OrganizationId types.String                `tfsdk:"organization_id"`
	ProjectId      types.String                `tfsdk:"project_id"`
	ClusterId      types.String                `tfsdk:"cluster_id"`
	AppServiceId   types.String                `tfsdk:"app_service_id"`
	Data           []AppServiceMaintenanceData `tfsdk:"data"`
}

// AppServiceMaintenanceData is a single App Service maintenance record.
type AppServiceMaintenanceData struct {
	Id           types.String      `tfsdk:"id"`
	Kind         types.String      `tfsdk:"kind"`
	InitiatedBy  types.String      `tfsdk:"initiated_by"`
	VersionId    types.String      `tfsdk:"version_id"`
	CurrentState types.String      `tfsdk:"current_state"`
	QueueAt      timetypes.RFC3339 `tfsdk:"queue_at"`
	CompletedAt  timetypes.RFC3339 `tfsdk:"completed_at"`
	Audit        types.Object      `tfsdk:"audit"`
}

// NewAppServiceMaintenanceData maps a maintenance record from the V4 API to its Terraform model.
// Fields the API omits map to null.
func NewAppServiceMaintenanceData(ctx context.Context, resp *appservice_maintenance.GetMaintenanceResponse) (AppServiceMaintenanceData, diag.Diagnostics) {
	audit := NewCouchbaseAuditData(resp.Audit)
	auditObject, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return AppServiceMaintenanceData{}, diags
	}

	return AppServiceMaintenanceData{
		Id:           types.StringValue(resp.Id),
		Kind:         types.StringValue(resp.Type),
		InitiatedBy:  types.StringValue(resp.InitiatedBy),
		VersionId:    types.StringPointerValue(resp.VersionId),
		CurrentState: types.StringValue(resp.CurrentState),
		QueueAt:      timetypes.NewRFC3339TimePointerValue(resp.QueueAt),
		CompletedAt:  timetypes.NewRFC3339TimePointerValue(resp.CompletedAt),
		Audit:        auditObject,
	}, nil
}
