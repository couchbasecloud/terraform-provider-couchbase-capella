package schema

import (
	"context"
	"encoding/json"
	"testing"

	network_peer_api "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/network_peer"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkPeerSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		name                   string
		expectedProjectId      string
		expectedOrganizationId string
		expectedClusterId      string
		expectedPeerId         string
		input                  NetworkPeer
	}

	tests := []test{
		{
			name: "[POSITIVE] project ID, organization ID, and cluster ID are passed via terraform apply",
			input: NetworkPeer{
				Id:             basetypes.NewStringValue("100"),
				ClusterId:      basetypes.NewStringValue("200"),
				ProjectId:      basetypes.NewStringValue("300"),
				OrganizationId: basetypes.NewStringValue("400"),
			},
			expectedPeerId:         "100",
			expectedClusterId:      "200",
			expectedProjectId:      "300",
			expectedOrganizationId: "400",
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: NetworkPeer{
				Id: basetypes.NewStringValue("id=100,cluster_id=200,project_id=300,organization_id=400"),
			},
			expectedPeerId:         "100",
			expectedClusterId:      "200",
			expectedProjectId:      "300",
			expectedOrganizationId: "400",
		},
		{
			name: "[NEGATIVE] only Peer ID is passed via terraform apply",
			input: NetworkPeer{
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

			assert.Equal(t, test.expectedPeerId, IDs[Id])
			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}

func TestMorphToProviderConfig_NilProviderConfig(t *testing.T) {
	resp := &network_peer_api.GetNetworkPeeringRecordResponse{
		ProviderConfig: nil,
	}
	cfg, err := morphToProviderConfig(resp)
	require.NoError(t, err)
	assert.Nil(t, cfg.AWSConfig)
	assert.Nil(t, cfg.GCPConfig)
	assert.Nil(t, cfg.AzureConfig)
}

func TestMorphToProviderConfig_NullProviderConfig(t *testing.T) {
	resp := &network_peer_api.GetNetworkPeeringRecordResponse{
		ProviderConfig: json.RawMessage("null"),
	}
	cfg, err := morphToProviderConfig(resp)
	require.NoError(t, err)
	assert.Nil(t, cfg.AWSConfig)
	assert.Nil(t, cfg.GCPConfig)
	assert.Nil(t, cfg.AzureConfig)
}

func TestNewNetworkPeer_NilReasoning(t *testing.T) {
	state := "failed"
	resp := &network_peer_api.GetNetworkPeeringRecordResponse{
		Id:             uuid.New(),
		Name:           "test-peer",
		ProviderType:   "aws",
		ProviderConfig: nil,
		Status: network_peer_api.PeeringStatus{
			State:     &state,
			Reasoning: nil,
		},
	}

	auditObj, _ := types.ObjectValueFrom(context.Background(), map[string]attr.Type{}, nil)

	peer, err := NewNetworkPeer(context.Background(), resp, "org-1", "proj-1", "clust-1", auditObj)
	require.NoError(t, err)
	assert.Equal(t, "test-peer", peer.Name.ValueString())

	// Status should be set with default failure reasoning
	assert.False(t, peer.Status.IsNull())
	var status PeeringStatus
	diags := peer.Status.As(context.Background(), &status, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError())
	assert.Equal(t, "failed", status.State.ValueString())
	assert.Equal(t, "Network peering failed. Remove this resource with 'terraform destroy' before retrying.", status.Reasoning.ValueString())
}

func TestNewNetworkPeer_NilProviderConfigAndNilReasoning(t *testing.T) {
	state := "failed"
	resp := &network_peer_api.GetNetworkPeeringRecordResponse{
		Id:             uuid.New(),
		Name:           "failed-peer",
		ProviderType:   "aws",
		ProviderConfig: json.RawMessage("null"),
		Status: network_peer_api.PeeringStatus{
			State:     &state,
			Reasoning: nil,
		},
	}

	auditObj, _ := types.ObjectValueFrom(context.Background(), map[string]attr.Type{}, nil)

	peer, err := NewNetworkPeer(context.Background(), resp, "org-1", "proj-1", "clust-1", auditObj)
	require.NoError(t, err)
	assert.NotNil(t, peer.ProviderConfig)
	assert.Nil(t, peer.ProviderConfig.AWSConfig)
	assert.Nil(t, peer.ProviderConfig.GCPConfig)
	assert.Nil(t, peer.ProviderConfig.AzureConfig)
}
