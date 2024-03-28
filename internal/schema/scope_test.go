package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

func TestScopeSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		input                  Scope
		name                   string
		expectedOrganizationId string
		expectedProjectId      string
		expectedClusterId      string
		expectedBucketId       string
		expectedScopeName      string
	}

	tests := []test{
		{
			name: "[POSITIVE] organization ID, project ID, cluster ID, bucket ID are passed via terraform apply",
			input: Scope{
				OrganizationId: basetypes.NewStringValue("100"),
				ProjectId:      basetypes.NewStringValue("200"),
				ClusterId:      basetypes.NewStringValue("300"),
				BucketId:       basetypes.NewStringValue("400"),
				Name:           basetypes.NewStringValue("new_terraform_scope"),
			},
			expectedOrganizationId: "100",
			expectedProjectId:      "200",
			expectedClusterId:      "300",
			expectedBucketId:       "400",
			expectedScopeName:      "new_terraform_scope",
		},
		{
			name: "[POSITIVE] Name is passed via terraform import",
			input: Scope{
				Name: basetypes.NewStringValue("scope_name=new_terraform_scope,bucket_id=400,cluster_id=300,project_id=200,organization_id=100"),
			},
			expectedOrganizationId: "100",
			expectedProjectId:      "200",
			expectedClusterId:      "300",
			expectedBucketId:       "400",
			expectedScopeName:      "new_terraform_scope",
		},
		{
			name: "[NEGATIVE] only scope name is passed via terraform apply",
			input: Scope{
				Name: basetypes.NewStringValue("terraform_scope"),
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

			assert.Equal(t, test.expectedScopeName, IDs[ScopeName])
			assert.Equal(t, test.expectedBucketId, IDs[BucketId])
			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
