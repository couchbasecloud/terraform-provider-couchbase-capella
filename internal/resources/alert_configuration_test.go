package resources

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAlertConfigurationUpdate(t *testing.T) {
	ctx := context.Background()
	str := func(s string) *string { return &s }
	num := func(f float64) *float64 { return &f }
	hours := func(i int64) *int64 { return &i }
	recipients := func(addresses ...string) types.Set {
		set, _ := types.SetValueFrom(ctx, types.StringType, append([]string{}, addresses...))
		return set
	}
	unmanaged := providerschema.AlertConfigurationSettings{Recipients: types.SetNull(types.StringType)}

	tests := []struct {
		name    string
		plan    providerschema.AlertConfigurationSettings
		current alertsapi.GetAlertConfigurationResponse
		check   func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest)
	}{
		{
			name: "only alerts with a changed field are sent",
			plan: providerschema.AlertConfigurationSettings{
				Configs: map[string]providerschema.AlertConfigurationSetting{
					"cpu_usage_high":          {FrequencyHours: types.Int64Value(3), WarningThreshold: types.Int64Value(85)},
					"CBHardOutOfMemoryErrors": {FrequencyHours: types.Int64Value(1)},
				},
				Recipients: types.SetNull(types.StringType),
			},
			current: alertsapi.GetAlertConfigurationResponse{
				Configs: map[string]alertsapi.AlertConfigurationEntry{
					"cpu_usage_high":          {FrequencyHours: hours(3), WarningThreshold: num(80), CriticalThreshold: num(90)},
					"CBHardOutOfMemoryErrors": {FrequencyHours: hours(1)},
				},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				require.Len(t, got.Configs, 1)
				sent := got.Configs["cpu_usage_high"]
				assert.Equal(t, int64(85), *sent.WarningThreshold)
				assert.Equal(t, int64(3), *sent.FrequencyHours)
				assert.Nil(t, sent.CriticalThreshold, "a field the plan does not set is left out")
				assert.Nil(t, got.Recipients)
				assert.Nil(t, got.Alerts)
			},
		},
		{
			name: "recipients are compared ignoring case",
			plan: providerschema.AlertConfigurationSettings{Recipients: recipients("Oncall@Example.com")},
			current: alertsapi.GetAlertConfigurationResponse{
				Recipients: []string{"oncall@example.com"},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				assert.True(t, isEmptyAlertConfigurationUpdate(got))
			},
		},
		{
			name: "an empty recipients list clears the current one",
			plan: providerschema.AlertConfigurationSettings{Recipients: recipients()},
			current: alertsapi.GetAlertConfigurationResponse{
				Recipients: []string{"oncall@example.com"},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				body, err := json.Marshal(got)
				require.NoError(t, err)
				assert.JSONEq(t, `{"recipients":[]}`, string(body), "null would leave the recipients unchanged")
			},
		},
		{
			name: "without a muted_alerts block mutes are left alone",
			plan: unmanaged,
			current: alertsapi.GetAlertConfigurationResponse{
				MutedAlerts: map[string]alertsapi.AlertMuteState{"disk_usage_high": {Action: "mute"}},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				assert.True(t, isEmptyAlertConfigurationUpdate(got))
			},
		},
		{
			name: "with a muted_alerts block, changed mutes are sent and others are cleared",
			plan: providerschema.AlertConfigurationSettings{
				Recipients: types.SetNull(types.StringType),
				MutedAlerts: map[string]providerschema.AlertMuteSetting{
					"disk_usage_high": {
						Action:     types.StringValue("mute"),
						MutedFrom:  timetypes.NewRFC3339Null(),
						MutedUntil: timetypes.NewRFC3339ValueMust("2026-11-01T05:30:00+05:30"),
					},
					"cpu_usage_high": {
						Action:     types.StringValue("unmute"),
						MutedFrom:  timetypes.NewRFC3339Null(),
						MutedUntil: timetypes.NewRFC3339Null(),
					},
				},
			},
			current: alertsapi.GetAlertConfigurationResponse{
				MutedAlerts: map[string]alertsapi.AlertMuteState{
					"disk_usage_high":    {Action: "mute", MutedUntil: str("2026-11-01T00:00:00Z")},
					"node_memory_usage":  {Action: "mute"},
					"indexer_low_memory": {Action: "unmute"},
					"cpu_usage_high":     {Action: "mute"},
				},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				assert.Equal(t, []alertsapi.AlertMuteAction{
					{AlertName: "cpu_usage_high", Action: "unmute"},
					{AlertName: "indexer_low_memory", Action: "clear"},
					{AlertName: "node_memory_usage", Action: "clear"},
				}, got.Alerts, "disk_usage_high ends at the same instant, so it is not resent")
			},
		},
		{
			name: "an empty muted_alerts block clears every mute",
			plan: providerschema.AlertConfigurationSettings{
				Recipients:  types.SetNull(types.StringType),
				MutedAlerts: map[string]providerschema.AlertMuteSetting{},
			},
			current: alertsapi.GetAlertConfigurationResponse{
				MutedAlerts: map[string]alertsapi.AlertMuteState{"disk_usage_high": {Action: "mute"}},
			},
			check: func(t *testing.T, got alertsapi.UpdateAlertConfigurationRequest) {
				assert.Equal(t, []alertsapi.AlertMuteAction{{AlertName: "disk_usage_high", Action: "clear"}}, got.Alerts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := alertConfigurationUpdate(ctx, tt.plan, tt.current)
			require.NoError(t, err)
			tt.check(t, got)
		})
	}
}

func TestAlertConfigurationDeleteUpdate(t *testing.T) {
	current := alertsapi.GetAlertConfigurationResponse{
		Recipients:  []string{"ui-added@example.com"},
		MutedAlerts: map[string]alertsapi.AlertMuteState{"disk_usage_high": {Action: "mute"}},
	}

	t.Run("clears the mutes and recipients state manages", func(t *testing.T) {
		got := alertConfigurationDeleteUpdate(providerschema.AlertConfigurationSettings{
			Recipients:  types.SetValueMust(types.StringType, nil),
			MutedAlerts: map[string]providerschema.AlertMuteSetting{},
		}, current)

		require.NotNil(t, got.Recipients)
		assert.Empty(t, *got.Recipients)
		assert.Equal(t, []alertsapi.AlertMuteAction{{AlertName: "disk_usage_high", Action: "clear"}}, got.Alerts)
		assert.Empty(t, got.Configs, "thresholds and frequencies cannot be removed")
	})

	t.Run("leaves alone what state does not manage, such as after an import", func(t *testing.T) {
		got := alertConfigurationDeleteUpdate(providerschema.AlertConfigurationSettings{
			Recipients: types.SetNull(types.StringType),
		}, current)

		assert.True(t, isEmptyAlertConfigurationUpdate(got))
	})
}

func TestAlertConfigurationWarnings(t *testing.T) {
	state := providerschema.AlertConfigurationSettings{
		Configs: map[string]providerschema.AlertConfigurationSetting{
			"cpu_usage_high":          {FrequencyHours: types.Int64Value(3), WarningThreshold: types.Int64Value(80)},
			"CBHardOutOfMemoryErrors": {FrequencyHours: types.Int64Value(1)},
		},
	}

	t.Run("removed alerts and fields are listed", func(t *testing.T) {
		plan := providerschema.AlertConfigurationSettings{
			Configs: map[string]providerschema.AlertConfigurationSetting{
				"cpu_usage_high": {FrequencyHours: types.Int64Value(3)},
			},
		}

		got := alertConfigurationRemovalWarning("cluster", state, plan)
		assert.Contains(t, got, "  - CBHardOutOfMemoryErrors: frequency_hours = 1\n  - cpu_usage_high: warning_threshold = 80")
		assert.NotContains(t, got, "frequency_hours = 3", "a field still in the plan is not listed")
	})

	t.Run("nothing removed gives no warning", func(t *testing.T) {
		assert.Empty(t, alertConfigurationRemovalWarning("cluster", state, state))
	})

	t.Run("destroy lists everything left behind", func(t *testing.T) {
		got := alertConfigurationLeftBehindWarning("cluster", state)
		assert.Contains(t, got, "  - CBHardOutOfMemoryErrors: frequency_hours = 1\n  - cpu_usage_high: frequency_hours = 3, warning_threshold = 80")
		assert.Empty(t, alertConfigurationLeftBehindWarning("cluster", providerschema.AlertConfigurationSettings{}))
	})
}
