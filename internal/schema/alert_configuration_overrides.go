package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// AlertConfigurationOverrides is the Terraform data source model for listing the alert
// configuration overrides in an organization.
type AlertConfigurationOverrides struct {
	OrganizationId types.String                    `tfsdk:"organization_id"`
	Data           []OneAlertConfigurationOverride `tfsdk:"data"`
}

// OneAlertConfigurationOverride is one alert's override at one level.
type OneAlertConfigurationOverride struct {
	Scope              types.String             `tfsdk:"scope"`
	OrganizationId     types.String             `tfsdk:"organization_id"`
	ProjectId          types.String             `tfsdk:"project_id"`
	ResourceId         types.String             `tfsdk:"resource_id"`
	AlertKind          types.String             `tfsdk:"alert_kind"`
	Config             *AlertConfigurationEntry `tfsdk:"config"`
	PreviousConfig     *AlertConfigurationEntry `tfsdk:"previous_config"`
	MutedAlert         *AlertMuteState          `tfsdk:"muted_alert"`
	Recipients         types.Set                `tfsdk:"recipients"`
	PreviousRecipients types.Set                `tfsdk:"previous_recipients"`
}

// AlertConfigurationEntry is the configuration for a single alert at one level.
type AlertConfigurationEntry struct {
	FrequencyHours    types.Int64       `tfsdk:"frequency_hours"`
	Threshold         types.Float64     `tfsdk:"threshold"`
	WarningThreshold  types.Float64     `tfsdk:"warning_threshold"`
	CriticalThreshold types.Float64     `tfsdk:"critical_threshold"`
	UpdatedAt         timetypes.RFC3339 `tfsdk:"updated_at"`
	UpdatedBy         types.String      `tfsdk:"updated_by"`
}

// AlertMuteState is the mute override set for one alert at one level.
type AlertMuteState struct {
	Action     types.String      `tfsdk:"action"`
	MutedFrom  timetypes.RFC3339 `tfsdk:"muted_from"`
	MutedUntil timetypes.RFC3339 `tfsdk:"muted_until"`
	UpdatedAt  timetypes.RFC3339 `tfsdk:"updated_at"`
	UpdatedBy  types.String      `tfsdk:"updated_by"`
}

// NewOneAlertConfigurationOverride converts an alert configuration override from the API into its
// Terraform model.
func NewOneAlertConfigurationOverride(ctx context.Context, override alertsapi.AlertConfigurationOverride) (OneAlertConfigurationOverride, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := OneAlertConfigurationOverride{
		Scope:          types.StringValue(override.Scope),
		OrganizationId: types.StringValue(override.OrganizationId),
		ProjectId:      types.StringPointerValue(override.ProjectId),
		ResourceId:     types.StringPointerValue(override.ResourceId),
		AlertKind:      types.StringValue(override.AlertKind),
	}

	var d diag.Diagnostics
	model.Config, d = newAlertConfigurationEntry(override.Config)
	diags.Append(d...)
	model.PreviousConfig, d = newAlertConfigurationEntry(override.PreviousConfig)
	diags.Append(d...)
	model.MutedAlert, d = newAlertMuteState(override.MutedAlert)
	diags.Append(d...)
	model.Recipients, d = newStringSet(ctx, override.Recipients)
	diags.Append(d...)
	model.PreviousRecipients, d = newStringSet(ctx, override.PreviousRecipients)
	diags.Append(d...)

	return model, diags
}

func newAlertConfigurationEntry(entry *alertsapi.AlertConfigurationEntry) (*AlertConfigurationEntry, diag.Diagnostics) {
	if entry == nil {
		return nil, nil
	}

	updatedAt, diags := timetypes.NewRFC3339PointerValue(entry.UpdatedAt)

	return &AlertConfigurationEntry{
		FrequencyHours:    types.Int64PointerValue(entry.FrequencyHours),
		Threshold:         types.Float64PointerValue(entry.Threshold),
		WarningThreshold:  types.Float64PointerValue(entry.WarningThreshold),
		CriticalThreshold: types.Float64PointerValue(entry.CriticalThreshold),
		UpdatedAt:         updatedAt,
		UpdatedBy:         types.StringPointerValue(entry.UpdatedBy),
	}, diags
}

func newAlertMuteState(mute *alertsapi.AlertMuteState) (*AlertMuteState, diag.Diagnostics) {
	if mute == nil {
		return nil, nil
	}

	var diags diag.Diagnostics
	mutedFrom, d := timetypes.NewRFC3339PointerValue(mute.MutedFrom)
	diags.Append(d...)
	mutedUntil, d := timetypes.NewRFC3339PointerValue(mute.MutedUntil)
	diags.Append(d...)
	updatedAt, d := timetypes.NewRFC3339PointerValue(mute.UpdatedAt)
	diags.Append(d...)

	return &AlertMuteState{
		Action:     types.StringValue(mute.Action),
		MutedFrom:  mutedFrom,
		MutedUntil: mutedUntil,
		UpdatedAt:  updatedAt,
		UpdatedBy:  types.StringPointerValue(mute.UpdatedBy),
	}, diags
}

func newStringSet(ctx context.Context, values []string) (types.Set, diag.Diagnostics) {
	if values == nil {
		return types.SetNull(types.StringType), nil
	}
	return types.SetValueFrom(ctx, types.StringType, values)
}
