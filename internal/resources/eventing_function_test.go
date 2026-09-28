package resources

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func bucketBinding(alias string) providerschema.EventingFunctionBucketBinding {
	return providerschema.EventingFunctionBucketBinding{
		Alias:      types.StringValue(alias),
		Bucket:     types.StringValue("travel-sample"),
		Scope:      types.StringValue("_default"),
		Collection: types.StringValue("_default"),
		Permission: types.StringValue("read"),
	}
}

func constantBinding(alias, value string) providerschema.EventingFunctionConstantBinding {
	return providerschema.EventingFunctionConstantBinding{
		Alias: types.StringValue(alias),
		Value: types.StringValue(value),
	}
}

func Test_eventingBindingsChanged(t *testing.T) {
	withEntries := &providerschema.EventingFunctionBindingsResource{
		Buckets:   []providerschema.EventingFunctionBucketBinding{bucketBinding("b1")},
		Constants: []providerschema.EventingFunctionConstantBinding{constantBinding("c1", "1")},
	}

	tests := []struct {
		name     string
		plan     *providerschema.EventingFunctionBindingsResource
		state    *providerschema.EventingFunctionBindingsResource
		expected bool
	}{
		{
			name:     "null plan and null state",
			plan:     nil,
			state:    nil,
			expected: false,
		},
		{
			name:     "null plan and empty state",
			plan:     nil,
			state:    &providerschema.EventingFunctionBindingsResource{},
			expected: false,
		},
		{
			name: "empty lists plan and null state",
			plan: &providerschema.EventingFunctionBindingsResource{
				Buckets:   []providerschema.EventingFunctionBucketBinding{},
				Urls:      []providerschema.EventingFunctionUrlBinding{},
				Constants: []providerschema.EventingFunctionConstantBinding{},
			},
			state:    nil,
			expected: false,
		},
		{
			name:     "null plan removes the state bindings",
			plan:     nil,
			state:    withEntries,
			expected: true,
		},
		{
			name: "empty list removes a state list",
			plan: &providerschema.EventingFunctionBindingsResource{
				Buckets:   []providerschema.EventingFunctionBucketBinding{},
				Constants: []providerschema.EventingFunctionConstantBinding{constantBinding("c1", "1")},
			},
			state:    withEntries,
			expected: true,
		},
		{
			name: "null list removes a state list",
			plan: &providerschema.EventingFunctionBindingsResource{
				Buckets: []providerschema.EventingFunctionBucketBinding{bucketBinding("b1")},
			},
			state:    withEntries,
			expected: true,
		},
		{
			name:     "bindings added to null state",
			plan:     withEntries,
			state:    nil,
			expected: true,
		},
		{
			name:     "unchanged bindings",
			plan:     withEntries,
			state:    withEntries,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			changed, err := eventingBindingsChanged(context.Background(), test.plan, test.state)
			require.NoError(t, err)
			assert.Equal(t, test.expected, changed)
		})
	}
}

func Test_bindingsToAPI(t *testing.T) {
	tests := []struct {
		name     string
		bindings *providerschema.EventingFunctionBindingsResource
		expected string
	}{
		{
			name:     "null bindings send every list empty",
			bindings: nil,
			expected: `{"buckets":[],"urls":[],"constants":[]}`,
		},
		{
			name:     "null lists are sent empty",
			bindings: &providerschema.EventingFunctionBindingsResource{},
			expected: `{"buckets":[],"urls":[],"constants":[]}`,
		},
		{
			name: "empty lists are sent empty",
			bindings: &providerschema.EventingFunctionBindingsResource{
				Buckets:   []providerschema.EventingFunctionBucketBinding{},
				Urls:      []providerschema.EventingFunctionUrlBinding{},
				Constants: []providerschema.EventingFunctionConstantBinding{},
			},
			expected: `{"buckets":[],"urls":[],"constants":[]}`,
		},
		{
			name: "lists with entries are sent with the other lists empty",
			bindings: &providerschema.EventingFunctionBindingsResource{
				Constants: []providerschema.EventingFunctionConstantBinding{constantBinding("c1", "1")},
			},
			expected: `{"buckets":[],"urls":[],"constants":[{"alias":"c1","value":"1"}]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bindings, err := bindingsToAPI(context.Background(), test.bindings)
			require.NoError(t, err)

			body, err := json.Marshal(bindings)
			require.NoError(t, err)
			assert.JSONEq(t, test.expected, string(body))
		})
	}
}

func Test_setEventingFunctionComputedAttributesToNull(t *testing.T) {
	tests := []struct {
		name     string
		code     types.String
		expected types.String
	}{
		{
			name:     "unknown code is set to null",
			code:     types.StringUnknown(),
			expected: types.StringNull(),
		},
		{
			name:     "known code is kept",
			code:     types.StringValue("function OnUpdate(doc, meta) {}"),
			expected: types.StringValue("function OnUpdate(doc, meta) {}"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			plan := &providerschema.EventingFunctionResource{Code: test.code}

			diags := setEventingFunctionComputedAttributesToNull(context.Background(), plan)
			require.False(t, diags.HasError())
			assert.Equal(t, test.expected, plan.Code)
		})
	}
}
