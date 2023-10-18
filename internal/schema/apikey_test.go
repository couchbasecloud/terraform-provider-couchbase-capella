package schema

import (
	"terraform-provider-capella/internal/errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestApiKeySchemaValidate(t *testing.T) {
	type test struct {
		name                   string
		input                  ApiKey
		expectedOrganizationId string
		expectedApiKeyId       string
		expectedErr            error
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID and api key ID are passed via terraform apply",
			input: ApiKey{
				Id:             basetypes.NewStringValue("100"),
				OrganizationId: basetypes.NewStringValue("200"),
			},
			expectedApiKeyId:       "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: ApiKey{
				Id: basetypes.NewStringValue("id=100,organization_id=200"),
			},
			expectedApiKeyId:       "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[NEGATIVE] only allow list ID is passed via terraform apply",
			input: ApiKey{
				Id: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrIdMissing,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := test.input.Validate()

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, test.expectedErr.Error())
				return
			}

			assert.Equal(t, test.expectedApiKeyId, IDs[Id])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])

		})
	}
}
