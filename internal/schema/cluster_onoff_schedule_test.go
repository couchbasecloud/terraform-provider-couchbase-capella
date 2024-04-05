package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestNewClusterOnOffScheduleSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		input                  ClusterOnOffSchedule
		name                   string
		expectedOrganizationId string
		expectedProjectId      string
		expectedClusterId      string
		expectedTimeZone       string
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID, project ID, cluster ID are passed via terraform apply",
			input: ClusterOnOffSchedule{
				OrganizationId: basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				ClusterId:      basetypes.NewStringValue("300"),
			},
			expectedOrganizationId: "100",
			expectedProjectId:      "200",
			expectedClusterId:      "300",
		},
		{
			name: "[POSITIVE] project ID, organization ID, and cluster ID, timezone, are passed via terraform apply",
			input: ClusterOnOffSchedule{
				ClusterId:      basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				OrganizationId: basetypes.NewStringValue("300"),
				Timezone:       basetypes.NewStringValue("US/Pacific"),
			},
			expectedClusterId:      "100",
			expectedProjectId:      "200",
			expectedOrganizationId: "300",
			expectedTimeZone:       "US/Pacific",
		},
		{
			name: "[NEGATIVE] only timezone is passed via terraform apply",
			input: ClusterOnOffSchedule{
				Timezone: basetypes.NewStringValue("US/Pacific"),
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
