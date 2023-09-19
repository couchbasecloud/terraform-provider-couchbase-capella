package schema

import (
	"testing"

	"terraform-provider-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/require"
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
			name: "[HAPPY PATH] project ID and organization ID are passed via terraform apply",
			input: Project{
				Id:             basetypes.NewStringValue("100"),
				OrganizationId: basetypes.NewStringValue("200"),
			},
			expectedProjectId:      "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[HAPPY PATH] project ID and organization ID are passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("id=100,organization_id=200"),
			},
			expectedProjectId:      "100",
			expectedOrganizationId: "200",
		},
		{
			name: "[SAD PATH] only project ID is passed via terraform apply",
			input: Project{
				Id: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrOrganizationIdMissing,
		},
		{
			name: "[SAD PATH] only organization ID is passed via terraform apply",
			input: Project{
				OrganizationId: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrProjectIdCannotBeEmpty,
		},
		{
			name: "[SAD PATH] project ID and organization ID are incorrectly passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("100&organization_id=200"),
			},
			expectedErr: errors.ErrOrganizationIdMissing,
		},
		{
			name: "[SAD PATH] project ID and organization ID are incorrectly passed via terraform import",
			input: Project{
				Id: basetypes.NewStringValue("id=100,orgId=200"),
			},
			expectedErr: errors.ErrOrganizationIdMissing,
		},
		{
			name: "[SAD PATH] project ID and organization ID are incorrectly passed via terraform import",
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
				require.Equal(t, tt.expectedErr, err)
				return
			}

			require.Equal(t, tt.expectedProjectId, projectId)
			require.Equal(t, tt.expectedOrganizationId, organizationId)
		})
	}
}
