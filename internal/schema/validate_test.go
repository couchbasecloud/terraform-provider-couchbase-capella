package schema

import (
	"terraform-provider-capella/internal/errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"gotest.tools/assert"
)

func Test_ValidateSchemaState(t *testing.T) {
	type test struct {
		name                         string
		state                        map[Attr]basetypes.StringValue
		expectedProjectId            string
		expectedOrganizationId       string
		expectedClusterId            string
		expectedDatabaseCredentialId string
		expectedErr                  error
	}

	tests := []test{
		{
			name: "[POSITIVE] - Multiple IDs successfully validated via terraform apply",
			state: map[Attr]basetypes.StringValue{
				Id:             basetypes.NewStringValue("100"),
				ClusterId:      basetypes.NewStringValue("200"),
				ProjectId:      basetypes.NewStringValue("300"),
				OrganizationId: basetypes.NewStringValue("400"),
			},
			expectedDatabaseCredentialId: "100",
			expectedClusterId:            "200",
			expectedProjectId:            "300",
			expectedOrganizationId:       "400",
		},
		{
			name: "[POSITIVE] - Multiple IDs successfully validated via terraform import",
			state: map[Attr]basetypes.StringValue{
				Id: basetypes.NewStringValue("id=100,cluster_id=200,project_id=300,organization_id=400"),
			},
			expectedDatabaseCredentialId: "100",
			expectedClusterId:            "200",
			expectedProjectId:            "300",
			expectedOrganizationId:       "400",
		},
		{
			name: "[POSITIVE] - IDs are passed in a different order via terraform import",
			state: map[Attr]basetypes.StringValue{
				Id: basetypes.NewStringValue("cluster_id=200,id=100,organization_id=400,project_id=300"),
			},
			expectedDatabaseCredentialId: "100",
			expectedClusterId:            "200",
			expectedProjectId:            "300",
			expectedOrganizationId:       "400",
		},
		{
			name: "[NEGATIVE] - IDs are passed incorrectly via terraform import",
			state: map[Attr]basetypes.StringValue{
				Id: basetypes.NewStringValue("100&organization_id=200,projectId=123&cluster_id=900"),
			},
			expectedErr: errors.ErrIdMissing,
		},
		{
			name: "[NEGATIVE] - IDs are passed with incorrect names via terraform import",
			state: map[Attr]basetypes.StringValue{
				Id: basetypes.NewStringValue("id=100,orgId=200,clusterId=300,project_id=900"),
			},
			expectedErr: errors.ErrIdMissing,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := validateSchemaState(test.state)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, test.expectedErr.Error())
				return
			}

			assert.Equal(t, test.expectedDatabaseCredentialId, IDs[Id])
			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
