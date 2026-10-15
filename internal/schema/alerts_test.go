package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	alertsapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/alerts"
)

func TestNewOneAlert(t *testing.T) {
	ptr := func(s string) *string { return &s }

	tests := []struct {
		name  string
		alert alertsapi.Alert
		check func(t *testing.T, got OneAlert)
	}{
		{
			name: "resolved per-node bucket alert sets every field",
			alert: alertsapi.Alert{
				AlertId:      "ffffffff-aaaa-1414-eeee-000000000003:CPUUsageHighWarning",
				EventId:      "ffffffff-aaaa-1414-eeee-000000000009",
				Name:         "High CPU Usage",
				Kind:         "CPUUsageHighWarning",
				ResourceId:   "ffffffff-aaaa-1414-eeee-000000000003",
				ResourceType: "cluster",
				ResourceName: ptr("orders-prod"),
				ProjectId:    ptr("ffffffff-aaaa-1414-eeee-000000000002"),
				ProjectName:  ptr("Payments"),
				Bucket:       ptr("orders"),
				NodeHostname: ptr("svc-d-node-001.example.cloud.couchbase.com"),
				Severity:     "warning",
				Status:       "resolved",
				StartsAt:     "2026-10-12T09:41:00Z",
				EndsAt:       ptr("2026-10-12T11:02:00Z"),
				EventTime:    "2026-10-12T11:02:00Z",
			},
			check: func(t *testing.T, got OneAlert) {
				assert.Equal(t, "ffffffff-aaaa-1414-eeee-000000000003:CPUUsageHighWarning", got.AlertId.ValueString())
				assert.Equal(t, "CPUUsageHighWarning", got.Kind.ValueString())
				assert.Equal(t, "orders-prod", got.ResourceName.ValueString())
				assert.Equal(t, "Payments", got.ProjectName.ValueString())
				assert.Equal(t, "orders", got.Bucket.ValueString())
				assert.Equal(t, "svc-d-node-001.example.cloud.couchbase.com", got.NodeHostname.ValueString())
				assert.Equal(t, "warning", got.Severity.ValueString())
				assert.Equal(t, "resolved", got.Status.ValueString())
				assert.Equal(t, "2026-10-12T09:41:00Z", got.StartsAt.ValueString())
				assert.Equal(t, "2026-10-12T11:02:00Z", got.EndsAt.ValueString())
				assert.Equal(t, "2026-10-12T11:02:00Z", got.EventTime.ValueString())
			},
		},
		{
			name: "firing alert leaves the fields the API omits null",
			alert: alertsapi.Alert{
				AlertId:      "ffffffff-aaaa-1414-eeee-000000000004:SGWCPUUsageHigh",
				EventId:      "ffffffff-aaaa-1414-eeee-000000000010",
				Name:         "App Service High CPU Usage",
				Kind:         "SGWCPUUsageHigh",
				ResourceId:   "ffffffff-aaaa-1414-eeee-000000000004",
				ResourceType: "syncgateway",
				Severity:     "critical",
				Status:       "firing",
				StartsAt:     "2026-10-12T09:41:00Z",
				EventTime:    "2026-10-12T09:41:00Z",
			},
			check: func(t *testing.T, got OneAlert) {
				assert.True(t, got.ResourceName.IsNull())
				assert.True(t, got.ProjectId.IsNull())
				assert.True(t, got.ProjectName.IsNull())
				assert.True(t, got.Bucket.IsNull())
				assert.True(t, got.NodeHostname.IsNull())
				assert.True(t, got.EndsAt.IsNull())
				assert.Equal(t, "firing", got.Status.ValueString())
				assert.Equal(t, "syncgateway", got.ResourceType.ValueString())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, diags := NewOneAlert(tt.alert)
			require.False(t, diags.HasError(), "unexpected diagnostics: %v", diags)
			tt.check(t, got)
		})
	}
}

func TestNewOneAlertInvalidTimestamp(t *testing.T) {
	_, diags := NewOneAlert(alertsapi.Alert{StartsAt: "not-a-time", EventTime: "2026-10-12T09:41:00Z"})
	assert.True(t, diags.HasError())
}
