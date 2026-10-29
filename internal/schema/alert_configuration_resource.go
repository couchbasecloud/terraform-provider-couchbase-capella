package schema

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

// AlertConfigurationSettings is what an alert configuration resource manages at its level. A nil
// map or a null set means the attribute is not in the user's config, so the resource leaves that
// part of the level alone.
type AlertConfigurationSettings struct {
	Configs     map[string]AlertConfigurationSetting `tfsdk:"configs"`
	Recipients  types.Set                            `tfsdk:"recipients"`
	MutedAlerts map[string]AlertMuteSetting          `tfsdk:"muted_alerts"`
}

// AlertConfigurationSetting is the threshold and frequency override of one alert. A null field
// is not managed by Terraform.
type AlertConfigurationSetting struct {
	FrequencyHours    types.Int64 `tfsdk:"frequency_hours"`
	Threshold         types.Int64 `tfsdk:"threshold"`
	WarningThreshold  types.Int64 `tfsdk:"warning_threshold"`
	CriticalThreshold types.Int64 `tfsdk:"critical_threshold"`
}

// AlertMuteSetting is the mute override of one alert.
type AlertMuteSetting struct {
	Action     types.String      `tfsdk:"action"`
	MutedFrom  timetypes.RFC3339 `tfsdk:"muted_from"`
	MutedUntil timetypes.RFC3339 `tfsdk:"muted_until"`
}

// ClusterAlertConfigurationResource is the Terraform resource model for a cluster's alert configuration.
type ClusterAlertConfigurationResource struct {
	OrganizationId types.String `tfsdk:"organization_id"`
	ProjectId      types.String `tfsdk:"project_id"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	AlertConfigurationSettings
}

// Validate checks the IDs in state, including the comma-separated form that import stores in cluster_id.
func (c *ClusterAlertConfigurationResource) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: c.OrganizationId,
		ProjectId:      c.ProjectId,
		ClusterId:      c.ClusterId,
	}

	IDs, err := validateSchemaState(state, ClusterId)
	if err != nil {
		return nil, fmt.Errorf("failed to validate resource state: %w", err)
	}

	return IDs, nil
}

// Fields returns the setting's fields by their Terraform names, in schema order.
func (s AlertConfigurationSetting) Fields() []AlertConfigurationField {
	return []AlertConfigurationField{
		{Name: "frequency_hours", Value: s.FrequencyHours},
		{Name: "threshold", Value: s.Threshold},
		{Name: "warning_threshold", Value: s.WarningThreshold},
		{Name: "critical_threshold", Value: s.CriticalThreshold},
	}
}

// AlertConfigurationField is one field of an AlertConfigurationSetting.
type AlertConfigurationField struct {
	Name  string
	Value types.Int64
}

// NewAlertConfigurationSettings builds the settings to keep in state from the level's current
// configuration. Only what prior tracks is refreshed:
//   - configs: alerts and fields in prior. Alerts set outside Terraform are ignored, because a
//     threshold or frequency override cannot be removed, so tracking them would make every
//     apply fail. An alert or field that is no longer set is dropped, so the next apply puts it back.
//   - recipients: the full list, when prior has it. An address keeps prior's spelling when it
//     matches ignoring case, since Capella stores addresses in lowercase.
//   - muted_alerts: every mute at the level, when prior has the block, so a mute added outside
//     Terraform shows up as drift and the next apply clears it.
func NewAlertConfigurationSettings(ctx context.Context, prior AlertConfigurationSettings, current alertsapi.GetAlertConfigurationResponse) (AlertConfigurationSettings, diag.Diagnostics) {
	var diags diag.Diagnostics
	var out AlertConfigurationSettings

	if prior.Configs != nil {
		out.Configs = make(map[string]AlertConfigurationSetting, len(prior.Configs))
		for name, priorEntry := range prior.Configs {
			entry, ok := current.Configs[name]
			if !ok {
				continue
			}
			out.Configs[name] = AlertConfigurationSetting{
				FrequencyHours:    refreshInt64(priorEntry.FrequencyHours, entry.FrequencyHours),
				Threshold:         refreshInt64(priorEntry.Threshold, float64ToInt64(entry.Threshold)),
				WarningThreshold:  refreshInt64(priorEntry.WarningThreshold, float64ToInt64(entry.WarningThreshold)),
				CriticalThreshold: refreshInt64(priorEntry.CriticalThreshold, float64ToInt64(entry.CriticalThreshold)),
			}
		}
	}

	out.Recipients = types.SetNull(types.StringType)
	if !prior.Recipients.IsNull() && !prior.Recipients.IsUnknown() {
		var priorRecipients []string
		diags.Append(prior.Recipients.ElementsAs(ctx, &priorRecipients, false)...)
		recipients := make([]string, 0, len(current.Recipients))
		for _, address := range current.Recipients {
			recipients = append(recipients, spellingOf(address, priorRecipients))
		}
		var d diag.Diagnostics
		out.Recipients, d = types.SetValueFrom(ctx, types.StringType, recipients)
		diags.Append(d...)
	}

	if prior.MutedAlerts != nil {
		out.MutedAlerts = make(map[string]AlertMuteSetting, len(current.MutedAlerts))
		for name, mute := range current.MutedAlerts {
			priorMute := prior.MutedAlerts[name]
			out.MutedAlerts[name] = AlertMuteSetting{
				Action:     types.StringValue(mute.Action),
				MutedFrom:  refreshTime(priorMute.MutedFrom, mute.MutedFrom, &diags),
				MutedUntil: refreshTime(priorMute.MutedUntil, mute.MutedUntil, &diags),
			}
		}
	}

	return out, diags
}

func refreshInt64(prior types.Int64, current *int64) types.Int64 {
	if prior.IsNull() {
		return types.Int64Null()
	}
	return types.Int64PointerValue(current)
}

// float64ToInt64 converts a threshold read with the data sources' float64 API struct. The API
// only returns whole numbers.
func float64ToInt64(v *float64) *int64 {
	if v == nil {
		return nil
	}
	i := int64(math.Round(*v))
	return &i
}

func spellingOf(address string, prior []string) string {
	for _, p := range prior {
		if strings.EqualFold(p, address) {
			return p
		}
	}
	return address
}

// refreshTime keeps prior's spelling of a time that is equal to current, so a time written with a
// different offset or precision than Capella returns is not reported as a change.
func refreshTime(prior timetypes.RFC3339, current *string, diags *diag.Diagnostics) timetypes.RFC3339 {
	if current == nil {
		return timetypes.NewRFC3339Null()
	}
	currentTime, err := time.Parse(time.RFC3339Nano, *current)
	if err != nil {
		diags.AddError("Error reading alert mute", fmt.Sprintf("Invalid timestamp %q: %s", *current, err))
		return timetypes.NewRFC3339Null()
	}
	if !prior.IsNull() && !prior.IsUnknown() {
		if priorTime, d := prior.ValueRFC3339Time(); !d.HasError() && priorTime.Equal(currentTime) {
			return prior
		}
	}
	return timetypes.NewRFC3339ValueMust(*current)
}
