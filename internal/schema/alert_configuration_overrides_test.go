package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

func TestNewOneAlertConfigurationOverride(t *testing.T) {
	ctx := context.Background()
	str := func(s string) *string { return &s }
	num := func(f float64) *float64 { return &f }
	hours := func(i int64) *int64 { return &i }

	t.Run("cluster override with config, previous config and recipients", func(t *testing.T) {
		got, diags := NewOneAlertConfigurationOverride(ctx, alertsapi.AlertConfigurationOverride{
			Scope:          "cluster",
			OrganizationId: "ffffffff-aaaa-1414-eeee-000000000000",
			ProjectId:      str("ffffffff-aaaa-1414-eeee-000000000001"),
			ResourceId:     str("ffffffff-aaaa-1414-eeee-000000000002"),
			AlertKind:      "cpu_usage_high",
			Config: &alertsapi.AlertConfigurationEntry{
				FrequencyHours:    hours(3),
				WarningThreshold:  num(80),
				CriticalThreshold: num(90),
				UpdatedAt:         str("2026-09-16T10:50:57.793147089Z"),
				UpdatedBy:         str("ffffffff-aaaa-1414-eeee-000000000003"),
			},
			PreviousConfig: &alertsapi.AlertConfigurationEntry{
				WarningThreshold: num(75),
			},
			Recipients:         []string{"alerts@example.com"},
			PreviousRecipients: []string{},
		})
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.Equal(t, "cluster", got.Scope.ValueString())
		assert.Equal(t, "ffffffff-aaaa-1414-eeee-000000000002", got.ResourceId.ValueString())
		assert.Equal(t, "cpu_usage_high", got.AlertKind.ValueString())

		require.NotNil(t, got.Config)
		assert.Equal(t, int64(3), got.Config.FrequencyHours.ValueInt64())
		assert.Equal(t, 80.0, got.Config.WarningThreshold.ValueFloat64())
		assert.Equal(t, 90.0, got.Config.CriticalThreshold.ValueFloat64())
		assert.True(t, got.Config.Threshold.IsNull())
		assert.Equal(t, "2026-09-16T10:50:57.793147089Z", got.Config.UpdatedAt.ValueString())

		require.NotNil(t, got.PreviousConfig)
		assert.Equal(t, 75.0, got.PreviousConfig.WarningThreshold.ValueFloat64())
		assert.True(t, got.PreviousConfig.UpdatedAt.IsNull())

		assert.Nil(t, got.MutedAlert)
		assert.Equal(t, types.SetValueMust(types.StringType, []attr.Value{types.StringValue("alerts@example.com")}), got.Recipients)
		assert.False(t, got.PreviousRecipients.IsNull())
		assert.Empty(t, got.PreviousRecipients.Elements())
	})

	t.Run("organization override that only mutes", func(t *testing.T) {
		got, diags := NewOneAlertConfigurationOverride(ctx, alertsapi.AlertConfigurationOverride{
			Scope:          "organization",
			OrganizationId: "ffffffff-aaaa-1414-eeee-000000000000",
			AlertKind:      "disk_usage_high",
			MutedAlert: &alertsapi.AlertMuteState{
				Action:    "mute",
				UpdatedAt: str("2026-09-20T09:58:12.104Z"),
				UpdatedBy: str("Jane Doe"),
			},
		})
		require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)

		assert.True(t, got.ProjectId.IsNull())
		assert.True(t, got.ResourceId.IsNull())
		assert.Nil(t, got.Config)
		assert.Nil(t, got.PreviousConfig)
		assert.True(t, got.Recipients.IsNull())
		assert.True(t, got.PreviousRecipients.IsNull())

		require.NotNil(t, got.MutedAlert)
		assert.Equal(t, "mute", got.MutedAlert.Action.ValueString())
		assert.True(t, got.MutedAlert.MutedFrom.IsNull())
		assert.True(t, got.MutedAlert.MutedUntil.IsNull())
		assert.Equal(t, "Jane Doe", got.MutedAlert.UpdatedBy.ValueString())
	})

	t.Run("invalid timestamp is reported", func(t *testing.T) {
		_, diags := NewOneAlertConfigurationOverride(ctx, alertsapi.AlertConfigurationOverride{
			Scope:          "organization",
			OrganizationId: "ffffffff-aaaa-1414-eeee-000000000000",
			AlertKind:      "cpu_usage_high",
			Config:         &alertsapi.AlertConfigurationEntry{UpdatedAt: str("not-a-time")},
		})
		assert.True(t, diags.HasError())
	})
}
