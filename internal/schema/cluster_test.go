package schema

import (
	"context"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types"
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

func TestNewClusterZones(t *testing.T) {
	tests := []struct {
		name     string
		zones    []string
		wantNull bool
		want     []string
	}{
		{name: "zones present", zones: []string{"use1-az5", "use1-az1", "use1-az4"}, want: []string{"use1-az1", "use1-az4", "use1-az5"}},
		{name: "zones absent", zones: nil, wantNull: true},
		{name: "zones empty", zones: []string{}, wantNull: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cluster := &clusterapi.GetClusterResponse{
				CloudProvider: clusterapi.CloudProvider{Type: clusterapi.Aws},
				Zones:         test.zones,
			}

			newCluster, err := NewCluster(context.Background(), cluster, "org", "proj", types.ObjectNull(nil))
			assert.NoError(t, err)

			if test.wantNull {
				assert.True(t, newCluster.Zones.IsNull())
			} else {
				var got []string
				assert.False(t, newCluster.Zones.ElementsAs(context.Background(), &got, false).HasError())
				assert.ElementsMatch(t, test.want, got)
			}

			oneCluster, err := NewOneCluster(context.Background(), cluster, "org", "proj", types.ObjectNull(nil))
			assert.NoError(t, err)

			clusterData, err := NewClusterData(cluster, "org", "proj", types.ObjectNull(nil))
			assert.NoError(t, err)

			// The data sources expose zones as a list, sorted so the order is stable between reads.
			for _, zones := range []types.List{oneCluster.Zones, clusterData.Zones} {
				if test.wantNull {
					assert.True(t, zones.IsNull())
					continue
				}
				var got []string
				assert.False(t, zones.ElementsAs(context.Background(), &got, false).HasError())
				assert.Equal(t, test.want, got)
			}
		})
	}
}
