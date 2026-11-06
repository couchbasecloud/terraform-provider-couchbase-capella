package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// Alerts is the Terraform data source model for listing the alerts in an organization.
// Status, ClusterIds, ProjectIds, SeverityLevels, From and To are optional filters that map to
// query parameters on the list alerts endpoint.
type Alerts struct {
	OrganizationId types.String      `tfsdk:"organization_id"`
	Status         types.Set         `tfsdk:"status"`
	ClusterIds     types.Set         `tfsdk:"cluster_ids"`
	ProjectIds     types.Set         `tfsdk:"project_ids"`
	SeverityLevels types.Set         `tfsdk:"severity_levels"`
	From           timetypes.RFC3339 `tfsdk:"from"`
	To             timetypes.RFC3339 `tfsdk:"to"`
	Data           []OneAlert        `tfsdk:"data"`
}

// OneAlert is a single alert on a single resource.
type OneAlert struct {
	AlertId      types.String      `tfsdk:"alert_id"`
	EventId      types.String      `tfsdk:"event_id"`
	Name         types.String      `tfsdk:"name"`
	Kind         types.String      `tfsdk:"kind"`
	ResourceId   types.String      `tfsdk:"resource_id"`
	ResourceType types.String      `tfsdk:"resource_type"`
	ResourceName types.String      `tfsdk:"resource_name"`
	ProjectId    types.String      `tfsdk:"project_id"`
	ProjectName  types.String      `tfsdk:"project_name"`
	Bucket       types.String      `tfsdk:"bucket"`
	NodeHostname types.String      `tfsdk:"node_hostname"`
	Severity     types.String      `tfsdk:"severity"`
	Status       types.String      `tfsdk:"status"`
	StartsAt     timetypes.RFC3339 `tfsdk:"starts_at"`
	EndsAt       timetypes.RFC3339 `tfsdk:"ends_at"`
	EventTime    timetypes.RFC3339 `tfsdk:"event_time"`
}

// NewOneAlert converts an alert from the API into its Terraform model.
func NewOneAlert(alert alertsapi.Alert) (OneAlert, diag.Diagnostics) {
	var diags diag.Diagnostics

	startsAt, d := timetypes.NewRFC3339Value(alert.StartsAt)
	diags.Append(d...)
	endsAt, d := timetypes.NewRFC3339PointerValue(alert.EndsAt)
	diags.Append(d...)
	eventTime, d := timetypes.NewRFC3339Value(alert.EventTime)
	diags.Append(d...)

	return OneAlert{
		AlertId:      types.StringValue(alert.AlertId),
		EventId:      types.StringValue(alert.EventId),
		Name:         types.StringValue(alert.Name),
		Kind:         types.StringValue(alert.Kind),
		ResourceId:   types.StringValue(alert.ResourceId),
		ResourceType: types.StringValue(alert.ResourceType),
		ResourceName: types.StringPointerValue(alert.ResourceName),
		ProjectId:    types.StringPointerValue(alert.ProjectId),
		ProjectName:  types.StringPointerValue(alert.ProjectName),
		Bucket:       types.StringPointerValue(alert.Bucket),
		NodeHostname: types.StringPointerValue(alert.NodeHostname),
		Severity:     types.StringValue(alert.Severity),
		Status:       types.StringValue(alert.Status),
		StartsAt:     startsAt,
		EndsAt:       endsAt,
		EventTime:    eventTime,
	}, diags
}
