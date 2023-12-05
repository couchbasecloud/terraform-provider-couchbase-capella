package schema

import (
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestUserSchemaValidate(t *testing.T) {
	type test struct {
		name                   string
		input                  User
		expectedUserId         string
		expectedOrganizationId string
		expectedErr            error
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID and user ID are passed via terraform apply",
			input: User{
				Id:             basetypes.NewStringValue("100"),
				OrganizationId: basetypes.NewStringValue("200"),
			},
			expectedUserId:         "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: User{
				Id: basetypes.NewStringValue("id=100,organization_id=200"),
			},
			expectedUserId:         "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[NEGATIVE] only allow list ID is passed via terraform apply",
			input: User{
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

			assert.Equal(t, test.expectedUserId, IDs[Id])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
