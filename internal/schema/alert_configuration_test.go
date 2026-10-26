package schema

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

func TestNewAlertConfigurationValues(t *testing.T) {
	ctx := context.Background()
	str := func(s string) *string { return &s }
	num := func(f float64) *float64 { return &f }
	hours := func(i int64) *int64 { return &i }

	t.Run("own and resolved values are both converted", func(t *testing.T) {
		got, diags := NewAlertConfigurationValues(ctx, alertsapi.GetAlertConfigurationResponse{
			Configs: map[string]alertsapi.AlertConfigurationEntry{
				"cpu_usage_high": {FrequencyHours: hours(3), UpdatedAt: str("2026-09-16T10:50:57Z"), UpdatedBy: str("ffffffff-aaaa-1414-eeee-000000000003")},
			},
			ResolvedConfigs: map[string]alertsapi.AlertConfigurationEntry{
				"cpu_usage_high": {FrequencyHours: hours(3), WarningThreshold: num(80), CriticalThreshold: num(90)},
			},
			Recipients:         []string{"cluster-oncall@example.com"},
			ResolvedRecipients: []string{"org-oncall@example.com", "cluster-oncall@example.com"},
			MutedAlerts: map[string]alertsapi.AlertMuteState{
				"disk_usage_high": {Action: "mute", MutedUntil: str("2026-11-01T00:00:00Z"), UpdatedBy: str("Jane Doe")},
			},
			ResolvedMutedAlerts: map[string]alertsapi.AlertMuteState{
				"disk_usage_high": {Action: "mute", MutedUntil: str("2026-11-01T00:00:00Z")},
			},
		})
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		require.Contains(t, got.Configs, "cpu_usage_high")
		own := got.Configs["cpu_usage_high"]
		assert.Equal(t, int64(3), own.FrequencyHours.ValueInt64())
		assert.True(t, own.WarningThreshold.IsNull(), "only the frequency is set at this level")
		assert.Equal(t, "2026-09-16T10:50:57Z", own.UpdatedAt.ValueString())
		assert.Equal(t, "ffffffff-aaaa-1414-eeee-000000000003", own.UpdatedBy.ValueString())

		resolved := got.ResolvedConfigs["cpu_usage_high"]
		assert.Equal(t, 80.0, resolved.WarningThreshold.ValueFloat64())
		assert.Equal(t, 90.0, resolved.CriticalThreshold.ValueFloat64())

		assert.Len(t, got.Recipients.Elements(), 1)
		assert.Len(t, got.ResolvedRecipients.Elements(), 2)

		require.Contains(t, got.MutedAlerts, "disk_usage_high")
		mute := got.MutedAlerts["disk_usage_high"]
		assert.Equal(t, "mute", mute.Action.ValueString())
		assert.True(t, mute.MutedFrom.IsNull())
		assert.Equal(t, "2026-11-01T00:00:00Z", mute.MutedUntil.ValueString())
		assert.Equal(t, "Jane Doe", mute.UpdatedBy.ValueString())
		assert.Contains(t, got.ResolvedMutedAlerts, "disk_usage_high")
	})

	t.Run("a level that sets nothing has empty values", func(t *testing.T) {
		got, diags := NewAlertConfigurationValues(ctx, alertsapi.GetAlertConfigurationResponse{
			Configs:             map[string]alertsapi.AlertConfigurationEntry{},
			ResolvedConfigs:     map[string]alertsapi.AlertConfigurationEntry{},
			Recipients:          []string{},
			ResolvedRecipients:  []string{},
			MutedAlerts:         map[string]alertsapi.AlertMuteState{},
			ResolvedMutedAlerts: map[string]alertsapi.AlertMuteState{},
		})
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.Empty(t, got.Configs)
		assert.Empty(t, got.MutedAlerts)
		assert.False(t, got.Recipients.IsNull())
		assert.Empty(t, got.Recipients.Elements())
	})

	t.Run("invalid timestamp is reported", func(t *testing.T) {
		_, diags := NewAlertConfigurationValues(ctx, alertsapi.GetAlertConfigurationResponse{
			MutedAlerts: map[string]alertsapi.AlertMuteState{"disk_usage_high": {Action: "mute", MutedUntil: str("not-a-time")}},
		})
		assert.True(t, diags.HasError())
	})
}
