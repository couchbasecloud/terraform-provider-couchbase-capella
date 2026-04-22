package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
)

func TestAppEndpointInList(t *testing.T) {
	const (
		e1 = "e1"
		e2 = "e2"
		e3 = "e3"
	)

	tests := []struct {
		name     string
		list     []app_endpoints.GetAppEndpointResponse
		target   string
		expected bool
	}{
		{
			name:     "empty list",
			list:     nil,
			target:   e1,
			expected: false,
		},
		{
			name: "name present",
			list: []app_endpoints.GetAppEndpointResponse{
				{Name: e3},
				{Name: e2},
			},
			target:   e2,
			expected: true,
		},
		{
			name: "name absent",
			list: []app_endpoints.GetAppEndpointResponse{
				{Name: e1},
				{Name: e2},
				{Name: e3},
			},
			target:   "other",
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, appEndpointInList(tc.list, tc.target))
		})
	}
}
