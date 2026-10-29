package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

func TestNewAlertConfigurationSettings(t *testing.T) {
	ctx := context.Background()
	str := func(s string) *string { return &s }
	num := func(f float64) *float64 { return &f }
	hours := func(i int64) *int64 { return &i }

	current := alertsapi.GetAlertConfigurationResponse{
		Configs: map[string]alertsapi.AlertConfigurationEntry{
			"cpu_usage_high":  {FrequencyHours: hours(3), WarningThreshold: num(85), CriticalThreshold: num(90)},
			"disk_usage_high": {WarningThreshold: num(70)},
		},
		Recipients: []string{"oncall@example.com", "ui-added@example.com"},
		MutedAlerts: map[string]alertsapi.AlertMuteState{
			"disk_usage_high":   {Action: "mute", MutedUntil: str("2026-11-01T00:00:00Z")},
			"node_memory_usage": {Action: "mute"},
		},
	}

	t.Run("only what prior tracks is refreshed", func(t *testing.T) {
		recipients, _ := types.SetValueFrom(ctx, types.StringType, []string{"OnCall@example.com"})
		prior := AlertConfigurationSettings{
			Configs: map[string]AlertConfigurationSetting{
				"cpu_usage_high":          {WarningThreshold: types.Int64Value(80)},
				"CBHardOutOfMemoryErrors": {FrequencyHours: types.Int64Value(1)},
			},
			Recipients: recipients,
			MutedAlerts: map[string]AlertMuteSetting{
				"disk_usage_high": {
					Action:     types.StringValue("mute"),
					MutedUntil: timetypes.NewRFC3339ValueMust("2026-11-01T05:30:00+05:30"),
				},
			},
		}

		got, diags := NewAlertConfigurationSettings(ctx, prior, current)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		require.Len(t, got.Configs, 1, "disk_usage_high is ignored and CBHardOutOfMemoryErrors is dropped")
		cpu := got.Configs["cpu_usage_high"]
		assert.Equal(t, int64(85), cpu.WarningThreshold.ValueInt64(), "drift is picked up")
		assert.True(t, cpu.FrequencyHours.IsNull(), "a field prior does not track stays unmanaged")
		assert.True(t, cpu.CriticalThreshold.IsNull())

		var addresses []string
		require.False(t, got.Recipients.ElementsAs(ctx, &addresses, false).HasError())
		assert.ElementsMatch(t, []string{"OnCall@example.com", "ui-added@example.com"}, addresses)

		require.Len(t, got.MutedAlerts, 2, "a mute added outside Terraform is tracked so apply can clear it")
		assert.Equal(t, "2026-11-01T05:30:00+05:30", got.MutedAlerts["disk_usage_high"].MutedUntil.ValueString(), "an equal time keeps prior's spelling")
		assert.True(t, got.MutedAlerts["node_memory_usage"].MutedUntil.IsNull())
	})

	t.Run("attributes not in the config stay unmanaged", func(t *testing.T) {
		got, diags := NewAlertConfigurationSettings(ctx, AlertConfigurationSettings{Recipients: types.SetNull(types.StringType)}, current)
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.Nil(t, got.Configs)
		assert.True(t, got.Recipients.IsNull())
		assert.Nil(t, got.MutedAlerts)
	})

	t.Run("invalid timestamp is reported", func(t *testing.T) {
		_, diags := NewAlertConfigurationSettings(ctx, AlertConfigurationSettings{
			Recipients:  types.SetNull(types.StringType),
			MutedAlerts: map[string]AlertMuteSetting{},
		}, alertsapi.GetAlertConfigurationResponse{
			MutedAlerts: map[string]alertsapi.AlertMuteState{"disk_usage_high": {Action: "mute", MutedUntil: str("not-a-time")}},
		})
		assert.True(t, diags.HasError())
	})
}
