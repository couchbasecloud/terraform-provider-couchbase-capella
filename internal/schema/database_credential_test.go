package schema

import (
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseCredentialSchemaValidate(t *testing.T) {
	type test struct {
		name                         string
		input                        DatabaseCredential
		expectedProjectId            string
		expectedOrganizationId       string
		expectedClusterId            string
		expectedDatabaseCredentialId string
		expectedErr                  error
	}

	tests := []test{
		{
			name: "[POSITIVE] project ID, organization ID, cluster ID, database credential ID are passed via terraform apply",
			input: DatabaseCredential{
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
			name: "[POSITIVE] IDs are passed via terraform import",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("id=100,cluster_id=200,project_id=300,organization_id=400"),
			},
			expectedDatabaseCredentialId: "100",
			expectedClusterId:            "200",
			expectedProjectId:            "300",
			expectedOrganizationId:       "400",
		},
		{
			name: "[NEGATIVE] only database credential ID is passed via terraform apply",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("100"),
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

			assert.Equal(t, test.expectedDatabaseCredentialId, IDs[Id])
			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])

		})
	}
}
