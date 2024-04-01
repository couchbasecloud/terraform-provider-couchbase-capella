package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestNewAppServiceOnOffSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		input                  AppServiceOnOffOnDemand
		name                   string
		expectedOrganizationId string
		expectedProjectId      string
		expectedClusterId      string
		expectedState          string
		expectedAppServiceId   string
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID, project ID, cluster ID, app service ID are passed via terraform apply",
			input: AppServiceOnOffOnDemand{
				OrganizationId: basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				ClusterId:      basetypes.NewStringValue("300"),
				AppServiceId:   basetypes.NewStringValue("400"),
			},
			expectedOrganizationId: "100",
			expectedProjectId:      "200",
			expectedClusterId:      "300",
			expectedAppServiceId:   "400",
		},
		{
			name: "[POSITIVE] project ID, organization ID, cluster ID, app service ID and state are passed via terraform apply",
			input: AppServiceOnOffOnDemand{
				ClusterId:      basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				OrganizationId: basetypes.NewStringValue("300"),
				AppServiceId:   basetypes.NewStringValue("400"),
				State:          basetypes.NewStringValue("off"),
			},
			expectedClusterId:      "100",
			expectedProjectId:      "200",
			expectedOrganizationId: "300",
			expectedAppServiceId:   "400",
			expectedState:          "off",
		},
		{
			name: "[NEGATIVE] only state is passed via terraform apply",
			input: AppServiceOnOffOnDemand{
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
