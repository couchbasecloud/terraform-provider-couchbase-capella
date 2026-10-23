package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// TestClusterExpressScalingFromResponse checks that every state the API reports maps to the
// bool the resource and data sources expose, with the raw state passed through unchanged.
// A transition counts as the state it is heading to, so an apply that enables express
// scaling does not read back as drift while the rebalance runs.
func TestClusterExpressScalingFromResponse(t *testing.T) {
	tests := []struct {
		state       clusterapi.ExpressScalingState
		wantEnabled bool
	}{
		{state: clusterapi.ExpressScalingEnabled, wantEnabled: true},
		{state: clusterapi.ExpressScalingEnabling, wantEnabled: true},
		{state: clusterapi.ExpressScalingDisabled, wantEnabled: false},
		{state: clusterapi.ExpressScalingDisabling, wantEnabled: false},
		{state: "", wantEnabled: false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			resp := &clusterapi.GetClusterResponse{ExpressScaling: tt.state}

			cluster, err := NewCluster(context.Background(), resp, "org", "project", types.ObjectNull(nil))
			require.NoError(t, err)
			assert.Equal(t, types.BoolValue(tt.wantEnabled), cluster.ExpressScaling)
			assert.Equal(t, types.StringValue(string(tt.state)), cluster.ExpressScalingState)

			clusterData, err := NewClusterData(resp, "org", "project", types.ObjectNull(nil))
			require.NoError(t, err)
			assert.Equal(t, types.BoolValue(tt.wantEnabled), clusterData.ExpressScaling)
			assert.Equal(t, types.StringValue(string(tt.state)), clusterData.ExpressScalingState)
		})
	}
}
