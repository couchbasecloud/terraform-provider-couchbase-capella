package resources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// TestValidateClusterUpdateExpressScaling checks that express_scaling is create-only on the
// cluster resource: leaving it out of configuration never blocks an update, and only an
// explicit change is rejected, pointing at the dedicated resource.
func TestValidateClusterUpdateExpressScaling(t *testing.T) {
	tests := []struct {
		name    string
		plan    types.Bool
		state   types.Bool
		wantErr error
	}{
		{name: "omitted from configuration", plan: types.BoolUnknown(), state: types.BoolValue(true)},
		{name: "null in plan", plan: types.BoolNull(), state: types.BoolValue(true)},
		{name: "unchanged true", plan: types.BoolValue(true), state: types.BoolValue(true)},
		{name: "unchanged false", plan: types.BoolValue(false), state: types.BoolValue(false)},
		{name: "enable after create", plan: types.BoolValue(true), state: types.BoolValue(false), wantErr: errors.ErrUnableToUpdateExpressScaling},
		{name: "disable after create", plan: types.BoolValue(false), state: types.BoolValue(true), wantErr: errors.ErrUnableToUpdateExpressScaling},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
			state := awsClusterPlan(types.Int64Value(50), types.Int64Value(3000))
			plan.ExpressScaling = tt.plan
			state.ExpressScaling = tt.state

			err := (&Cluster{}).validateClusterUpdate(plan, state)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
