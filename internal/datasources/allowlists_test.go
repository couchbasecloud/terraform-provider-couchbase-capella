package datasources

import (
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	var (
		organizationId = basetypes.NewStringValue("organizationId")
		projectId      = basetypes.NewStringValue("projectId")
		clusterId      = basetypes.NewStringValue("clusterId")
	)

	type test struct {
		desc        string
		state       providerschema.AllowLists
		expectedErr error
	}

	tests := []test{
		{
			desc: "[POSITIVE] - All fields populated",
			state: providerschema.AllowLists{
				OrganizationId: organizationId,
				ProjectId:      projectId,
				ClusterId:      clusterId,
			},
		},
		{
			desc: "[NEGATIVE] - OrganizationId is missing",
			state: providerschema.AllowLists{
				ProjectId: projectId,
				ClusterId: clusterId,
			},
			expectedErr: errors.ErrOrganizationIdMissing,
		},
		{
			desc: "[NEGATIVE] - ProjectId is missing",
			state: providerschema.AllowLists{
				OrganizationId: organizationId,
				ClusterId:      clusterId,
			},
			expectedErr: errors.ErrProjectIdMissing,
		},
		{
			desc: "[NEGATIVE] - ClusterId is missing",
			state: providerschema.AllowLists{
				OrganizationId: organizationId,
				ProjectId:      projectId,
			},
			expectedErr: errors.ErrClusterIdMissing,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			a := &AllowList{}

			err := a.validate(test.state)

			if test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
