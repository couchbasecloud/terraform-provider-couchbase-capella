package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
)

func TestNewEventingFunctionResourceBindings(t *testing.T) {
	apiConstant := eventingapi.ConstantBinding{Alias: "c1", Value: "1"}
	constant := EventingFunctionConstantBinding{Alias: types.StringValue("c1"), Value: types.StringValue("1")}

	tests := []struct {
		name string
		// apiBindings is the GET response bindings. The API returns {} with no lists when there are no bindings.
		apiBindings *eventingapi.Bindings
		prior       *EventingFunctionResource
		expected    *EventingFunctionBindingsResource
	}{
		{
			name:        "no prior leaves empty bindings null",
			apiBindings: &eventingapi.Bindings{},
			prior:       nil,
			expected:    nil,
		},
		{
			name:        "null prior bindings stay null",
			apiBindings: &eventingapi.Bindings{},
			prior:       &EventingFunctionResource{},
			expected:    nil,
		},
		{
			name:        "empty prior bindings keep null lists",
			apiBindings: &eventingapi.Bindings{},
			prior:       &EventingFunctionResource{Bindings: &EventingFunctionBindingsResource{}},
			expected:    &EventingFunctionBindingsResource{},
		},
		{
			name:        "empty prior lists stay empty",
			apiBindings: &eventingapi.Bindings{},
			prior: &EventingFunctionResource{Bindings: &EventingFunctionBindingsResource{
				Buckets:   []EventingFunctionBucketBinding{},
				Urls:      []EventingFunctionUrlBinding{},
				Constants: []EventingFunctionConstantBinding{},
			}},
			expected: &EventingFunctionBindingsResource{
				Buckets:   []EventingFunctionBucketBinding{},
				Urls:      []EventingFunctionUrlBinding{},
				Constants: []EventingFunctionConstantBinding{},
			},
		},
		{
			name:        "empty prior lists stay empty next to a list with entries",
			apiBindings: &eventingapi.Bindings{Constants: []eventingapi.ConstantBinding{apiConstant}},
			prior: &EventingFunctionResource{Bindings: &EventingFunctionBindingsResource{
				Urls:      []EventingFunctionUrlBinding{},
				Constants: []EventingFunctionConstantBinding{constant},
			}},
			expected: &EventingFunctionBindingsResource{
				Urls:      []EventingFunctionUrlBinding{},
				Constants: []EventingFunctionConstantBinding{constant},
			},
		},
		{
			name:        "prior entries removed on the server become null",
			apiBindings: &eventingapi.Bindings{},
			prior: &EventingFunctionResource{Bindings: &EventingFunctionBindingsResource{
				Constants: []EventingFunctionConstantBinding{constant},
			}},
			expected: &EventingFunctionBindingsResource{},
		},
		{
			name:        "server entries replace an empty prior list",
			apiBindings: &eventingapi.Bindings{Constants: []eventingapi.ConstantBinding{apiConstant}},
			prior: &EventingFunctionResource{Bindings: &EventingFunctionBindingsResource{
				Constants: []EventingFunctionConstantBinding{},
			}},
			expected: &EventingFunctionBindingsResource{
				Constants: []EventingFunctionConstantBinding{constant},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := &eventingapi.GetEventingFunctionResponse{
				Name:                 "fn",
				EventSource:          eventingapi.Keyspace{Bucket: "source"},
				EventMetadataStorage: eventingapi.Keyspace{Bucket: "metadata"},
				Bindings:             test.apiBindings,
			}

			fn, err := NewEventingFunctionResource(context.Background(), resp, "org", "project", "cluster", test.prior)
			require.NoError(t, err)
			assert.Equal(t, test.expected, fn.Bindings)
		})
	}
}
