package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAlertsQueryParams(t *testing.T) {
	ctx := context.Background()
	stringSet := func(values ...string) types.Set {
		set, diags := types.SetValueFrom(ctx, types.StringType, values)
		require.False(t, diags.HasError())
		return set
	}
	nullSet := types.SetNull(types.StringType)

	t.Run("no filters sends no query parameters", func(t *testing.T) {
		params, err := alertsQueryParams(ctx, providerschema.Alerts{
			Status:         nullSet,
			ClusterIds:     nullSet,
			ProjectIds:     nullSet,
			SeverityLevels: nullSet,
			From:           timetypes.NewRFC3339Null(),
			To:             timetypes.NewRFC3339Null(),
		})
		require.NoError(t, err)
		assert.Empty(t, params)
	})

	t.Run("every filter maps to its query parameter", func(t *testing.T) {
		params, err := alertsQueryParams(ctx, providerschema.Alerts{
			Status:         stringSet("firing"),
			ClusterIds:     stringSet("ffffffff-aaaa-1414-eeee-000000000003"),
			ProjectIds:     stringSet("ffffffff-aaaa-1414-eeee-000000000002"),
			SeverityLevels: stringSet("warning", "critical"),
			From:           timetypes.NewRFC3339ValueMust("2026-10-12T10:00:00Z"),
			To:             timetypes.NewRFC3339ValueMust("2026-10-14T10:00:00Z"),
		})
		require.NoError(t, err)
		assert.Equal(t, []string{"firing"}, params["status"])
		assert.Equal(t, []string{"ffffffff-aaaa-1414-eeee-000000000003"}, params["clusterIds"])
		assert.Equal(t, []string{"ffffffff-aaaa-1414-eeee-000000000002"}, params["projectIds"])
		assert.ElementsMatch(t, []string{"warning", "critical"}, params["severityLevels"])
		assert.Equal(t, []string{"2026-10-12T10:00:00Z"}, params["from"])
		assert.Equal(t, []string{"2026-10-14T10:00:00Z"}, params["to"])
	})
}
