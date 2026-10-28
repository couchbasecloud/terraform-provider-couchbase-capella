package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// AlertConfigurationValues is the alert configuration of one level, shared by the organization,
// project and cluster alert configuration data sources.
type AlertConfigurationValues struct {
	Configs             map[string]AlertConfigurationEntry         `tfsdk:"configs"`
	ResolvedConfigs     map[string]ResolvedAlertConfigurationEntry `tfsdk:"resolved_configs"`
	Recipients          types.Set                                  `tfsdk:"recipients"`
	MutedAlerts         map[string]AlertMuteState                  `tfsdk:"muted_alerts"`
	ResolvedMutedAlerts map[string]AlertMuteState                  `tfsdk:"resolved_muted_alerts"`
}

// ResolvedAlertConfigurationEntry is the configuration in effect for a single alert after
// inheritance. It has no updated_at or updated_by: its fields can come from different levels,
// so the API does not attribute it to a single change.
type ResolvedAlertConfigurationEntry struct {
	FrequencyHours    types.Int64   `tfsdk:"frequency_hours"`
	Threshold         types.Float64 `tfsdk:"threshold"`
	WarningThreshold  types.Float64 `tfsdk:"warning_threshold"`
	CriticalThreshold types.Float64 `tfsdk:"critical_threshold"`
}

// OrganizationAlertConfiguration is the Terraform data source model for an organization's alert configuration.
type OrganizationAlertConfiguration struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	AlertConfigurationValues
}

// ProjectAlertConfiguration is the Terraform data source model for a project's alert configuration.
type ProjectAlertConfiguration struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	AlertConfigurationValues
}

// ClusterAlertConfiguration is the Terraform data source model for a cluster's alert configuration.
type ClusterAlertConfiguration struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	AlertConfigurationValues
}

// NewAlertConfigurationValues converts a get alert configuration response into its Terraform model.
func NewAlertConfigurationValues(ctx context.Context, resp alertsapi.GetAlertConfigurationResponse) (AlertConfigurationValues, diag.Diagnostics) {
	var diags diag.Diagnostics
	var values AlertConfigurationValues
	var d diag.Diagnostics

	values.Configs, d = newAlertConfigurationEntries(resp.Configs)
	diags.Append(d...)
	values.ResolvedConfigs = newResolvedAlertConfigurationEntries(resp.ResolvedConfigs)
	values.Recipients, d = newStringSet(ctx, resp.Recipients)
	diags.Append(d...)
	values.MutedAlerts, d = newAlertMuteStates(resp.MutedAlerts)
	diags.Append(d...)
	values.ResolvedMutedAlerts, d = newAlertMuteStates(resp.ResolvedMutedAlerts)
	diags.Append(d...)

	return values, diags
}

func newAlertConfigurationEntries(entries map[string]alertsapi.AlertConfigurationEntry) (map[string]AlertConfigurationEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := make(map[string]AlertConfigurationEntry, len(entries))
	for kind, entry := range entries {
		model, d := newAlertConfigurationEntry(&entry)
		diags.Append(d...)
		out[kind] = *model
	}
	return out, diags
}

func newResolvedAlertConfigurationEntries(entries map[string]alertsapi.AlertConfigurationEntry) map[string]ResolvedAlertConfigurationEntry {
	out := make(map[string]ResolvedAlertConfigurationEntry, len(entries))
	for kind, entry := range entries {
		out[kind] = ResolvedAlertConfigurationEntry{
			FrequencyHours:    types.Int64PointerValue(entry.FrequencyHours),
			Threshold:         types.Float64PointerValue(entry.Threshold),
			WarningThreshold:  types.Float64PointerValue(entry.WarningThreshold),
			CriticalThreshold: types.Float64PointerValue(entry.CriticalThreshold),
		}
	}
	return out
}

func newAlertMuteStates(mutes map[string]alertsapi.AlertMuteState) (map[string]AlertMuteState, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := make(map[string]AlertMuteState, len(mutes))
	for kind, mute := range mutes {
		model, d := newAlertMuteState(&mute)
		diags.Append(d...)
		out[kind] = *model
	}
	return out, diags
}
