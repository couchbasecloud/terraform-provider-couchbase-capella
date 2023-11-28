package schema

import (
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestClusterSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		name                   string
		expectedProjectId      string
		expectedOrganizationId string
		expectedClusterId      string
		input                  Cluster
	}

	tests := []test{
		{
			name: "[POSITIVE] project ID, organization ID, and cluster ID are passed via terraform apply",
			input: Cluster{
				Id:             basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				OrganizationId: basetypes.NewStringValue("300"),
			},
			expectedClusterId:      "100",
			expectedProjectId:      "200",
			expectedOrganizationId: "300",
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: Cluster{
				Id: basetypes.NewStringValue("id=100,project_id=200,organization_id=300"),
			},
			expectedClusterId:      "100",
			expectedProjectId:      "200",
			expectedOrganizationId: "300",
		},
		{
			name: "[NEGATIVE] only allow list ID is passed via terraform apply",
			input: Cluster{
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

			assert.Equal(t, test.expectedClusterId, IDs[Id])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
