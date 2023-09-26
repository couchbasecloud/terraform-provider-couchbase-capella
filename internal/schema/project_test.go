package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestProjectSchemaValidate(t *testing.T) {
	tests := []struct {
		name                   string
		input                  Project
		expectedProjectId      string
		expectedOrganizationId string
		expectedErr            error
	}{
		{
			name: "[POSITIVE] project ID and organization ID are passed via terraform apply",
			input: Project{
				Id:             basetypes.NewStringValue("100"),
				OrganizationId: basetypes.NewStringValue("200"),
			},
			expectedProjectId:      "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[POSITIVE] project ID and organization ID are passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("id=100,organization_id=200"),
			},
			expectedProjectId:      "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[NEGATIVE] only project ID is passed via terraform apply",
			input: Project{
				Id: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrIdMissing,
		},
		{
			name: "[NEGATIVE] only organization ID is passed via terraform apply",
			input: Project{
				OrganizationId: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrProjectIdCannotBeEmpty,
		},
		{
			name: "[NEGATIVE] project ID and organization ID are incorrectly passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("100&organization_id=200"),
			},
			expectedErr: errors.ErrIdMissing,
		},
		{
			name: "[NEGATIVE] project ID and organization ID are incorrectly passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("id=100,orgId=200"),
			},
			expectedErr: errors.ErrOrganizationIdMissing,
		},
		{
			name: "[NEGATIVE] project ID and organization ID are incorrectly passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("ProjectID=100,organization_id=200"),
			},
			expectedErr: errors.ErrProjectIdMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectId, organizationId, err := tt.input.Validate()

			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
				return
			}

			assert.Equal(t, tt.expectedProjectId, projectId)
			assert.Equal(t, tt.expectedOrganizationId, organizationId)
		})
	}
}
