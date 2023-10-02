package schema

import (
	"testing"

	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseCredentialSchemaValidate(t *testing.T) {
	tests := []struct {
		name                         string
		input                        DatabaseCredential
		expectedProjectId            string
		expectedOrganizationId       string
		expectedClusterId            string
		expectedDatabaseCredentialId string
		expectedErr                  error
	}{
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
			name: "[NEGATIVE] IDs follow the right syntax but order is incorrect in terraform import",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("id=100,organization_id=200,project_id=300,cluster_id=400"),
			},
			expectedErr: errors.ErrClusterIdMissing,
		},
		{
			name: "[NEGATIVE] only database credential ID is passed via terraform apply",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrIdMissing,
		},
		{
			name: "[NEGATIVE] only organization ID is passed via terraform apply",
			input: DatabaseCredential{
				OrganizationId: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrDatabaseCredentialIdCannotBeEmpty,
		},
		{
			name: "[NEGATIVE] IDs are incorrectly passed via terraform import",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("100&organization_id=200,projectId=123&cluster_id=900"),
			},
			expectedErr: errors.ErrIdMissing,
		},
		{
			name: "[NEGATIVE] IDs are incorrectly passed via terraform import",
			input: DatabaseCredential{
				Id: basetypes.NewStringValue("id=100,orgId=200,clusterId=300,project_id=900"),
			},
			expectedErr: errors.ErrClusterIdMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbId, clusterId, projectId, organizationId, err := tt.input.Validate()

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
				return
			}

			assert.Equal(t, tt.expectedDatabaseCredentialId, dbId)
			assert.Equal(t, tt.expectedClusterId, clusterId)
			assert.Equal(t, tt.expectedProjectId, projectId)
			assert.Equal(t, tt.expectedOrganizationId, organizationId)
		})
	}
}
