package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestNewClusterOnOffSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr              error
		input                    ClusterOnOffOnDemand
		name                     string
		expectedOrganizationId   string
		expectedProjectId        string
		expectedClusterId        string
		expectedState            string
		expectedLinkedAppService bool
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID, project ID, cluster ID are passed via terraform apply",
			input: ClusterOnOffOnDemand{
				OrganizationId: basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				ClusterId:      basetypes.NewStringValue("300"),
			},
			expectedOrganizationId: "100",
			expectedProjectId:      "200",
			expectedClusterId:      "300",
		},
		{
			name: "[POSITIVE] project ID, organization ID, cluster ID, state and turn on linked app service are passed via terraform apply",
			input: ClusterOnOffOnDemand{
				ClusterId:              basetypes.NewStringValue("100"),
				ProjectId:              basetypes.NewStringValue("200"),
				OrganizationId:         basetypes.NewStringValue("300"),
				State:                  basetypes.NewStringValue("on"),
				TurnOnLinkedAppService: basetypes.NewBoolValue(true),
			},
			expectedClusterId:        "100",
			expectedProjectId:        "200",
			expectedOrganizationId:   "300",
			expectedState:            "on",
			expectedLinkedAppService: true,
		},
		{
			name: "[NEGATIVE] only state is passed via terraform apply",
			input: ClusterOnOffOnDemand{
				State: basetypes.NewStringValue("on"),
			},
			expectedErr: errors.ErrInvalidImport,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := test.input.Validate()

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, test.expectedErr.Error())
				return
			}

			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
