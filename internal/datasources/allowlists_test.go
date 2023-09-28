package datasources

import (
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	providerschema "terraform-provider-capella/internal/schema"
	"testing"

	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func Test_MapResponseBody(t *testing.T) {
	var (
		organizationId = basetypes.NewStringValue(uuid.New().String())
		projectId      = basetypes.NewStringValue(uuid.New().String())
		clusterId      = basetypes.NewStringValue(uuid.New().String())

		cidr      = "0.0.0.0/10"
		comment   = "comment"
		expiresAt = "2023-09-26T19:20:30+01:00"
		id        = uuid.New()

		createdAt  = time.Now()
		createdBy  = "user"
		modifiedAt = time.Now()
		modifiedBy = "user"
		version    = 1

		allowList = api.GetAllowListResponse{
			Audit: api.CouchbaseAuditData{
				CreatedAt:  createdAt,
				CreatedBy:  createdBy,
				ModifiedAt: modifiedAt,
				ModifiedBy: modifiedBy,
				Version:    version,
			},
			Cidr:      cidr,
			Comment:   comment,
			ExpiresAt: expiresAt,
			Id:        id,
		}

		OneAllowList = providerschema.OneAllowList{
			Audit: providerschema.CouchbaseAuditData{
				CreatedAt:  basetypes.NewStringValue(createdAt.String()),
				CreatedBy:  basetypes.NewStringValue(createdBy),
				ModifiedAt: basetypes.NewStringValue(modifiedAt.String()),
				ModifiedBy: basetypes.NewStringValue(modifiedBy),
				Version:    basetypes.NewInt64Value(int64(version)),
			},
			OrganizationId: basetypes.NewStringValue(organizationId.ValueString()),
			ProjectId:      basetypes.NewStringValue(projectId.ValueString()),
			ClusterId:      basetypes.NewStringValue(clusterId.ValueString()),
			Cidr:           basetypes.NewStringValue(cidr),
			Comment:        basetypes.NewStringValue(comment),
			ExpiresAt:      basetypes.NewStringValue(expiresAt),
			Id:             basetypes.NewStringValue(id.String()),
		}
	)

	type test struct {
		desc          string
		response      api.GetAllowListsResponse
		expectedState providerschema.AllowLists
	}

	tests := []test{
		{
			desc: "[POSITIVE] - Fields successfully populated - one allow list in response",
			response: api.GetAllowListsResponse{
				Data: []api.GetAllowListResponse{
					allowList,
				},
			},
			expectedState: providerschema.AllowLists{
				OrganizationId: organizationId,
				ProjectId:      projectId,
				ClusterId:      clusterId,
				Data:           []providerschema.OneAllowList{OneAllowList},
			},
		},
		{
			desc: "[POSITIVE] - Fields successfully populated - multiple allow lists in response",
			response: api.GetAllowListsResponse{
				Data: []api.GetAllowListResponse{
					allowList,
					allowList,
				},
			},
			expectedState: providerschema.AllowLists{
				OrganizationId: organizationId,
				ProjectId:      projectId,
				ClusterId:      clusterId,
				Data: []providerschema.OneAllowList{
					OneAllowList,
					OneAllowList,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			a := &AllowList{}

			state := providerschema.AllowLists{
				OrganizationId: organizationId,
				ProjectId:      projectId,
				ClusterId:      clusterId,
			}
			state = a.mapResponseBody(test.response, &state)

			assert.Equal(t, test.expectedState, state)
		})
	}
}

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
