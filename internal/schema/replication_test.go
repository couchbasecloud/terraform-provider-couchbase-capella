package schema

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

func testGetReplicationResponse() *replicationapi.GetReplicationResponse {
	return &replicationapi.GetReplicationResponse{
		Id:          "cmVwbC0x",
		Status:      replicationapi.StatusRunning,
		ChangesLeft: 5,
		Source: replicationapi.SourceDetails{
			Project: replicationapi.Reference{Id: "proj-1", Name: "project"},
			Cluster: replicationapi.Reference{Id: "cluster-1", Name: "source"},
			Bucket:  replicationapi.BucketReference{Id: "c291cmNl", Name: "source", ConflictResolutionType: "seqno"},
		},
		Target: replicationapi.TargetDetails{
			Project: replicationapi.Reference{Id: "proj-1", Name: "project"},
			Cluster: replicationapi.Reference{Id: "cluster-2", Name: "target"},
			Bucket:  replicationapi.BucketReference{Id: "dGFyZ2V0", Name: "target", ConflictResolutionType: "seqno"},
			Type:    ReplicationTargetTypeCapella,
		},
		Direction:         "oneWay",
		Priority:          "high",
		NetworkUsageLimit: 0,
		Audit:             replicationapi.ReplicationAuditData{CreatedBy: "user-1", CreatedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)},
	}
}

func allFalseExcludeOptions() *ReplicationDocumentExcludeOptions {
	return &ReplicationDocumentExcludeOptions{
		Deletion:   types.BoolValue(false),
		Expiration: types.BoolValue(false),
		TTL:        types.BoolValue(false),
		Binary:     types.BoolValue(false),
	}
}

func TestNewReplicationFilter(t *testing.T) {
	deletionOnly := &replicationapi.DocumentExcludeOptions{Deletion: true}

	tests := []struct {
		name  string
		api   *replicationapi.Filter
		prior *ReplicationFilter
		want  *ReplicationFilter
	}{
		{
			name: "no filter in the API and no prior filter gives no filter",
		},
		{
			name: "import with a filter maps every part and uses false for skip_restream",
			api: &replicationapi.Filter{
				DocumentExcludeOptions: deletionOnly,
				Expressions:            &replicationapi.FilterExpressions{RegEx: "^a"},
			},
			want: &ReplicationFilter{
				DocumentExcludeOptions: &ReplicationDocumentExcludeOptions{
					Deletion: types.BoolValue(true), Expiration: types.BoolValue(false),
					TTL: types.BoolValue(false), Binary: types.BoolValue(false),
				},
				Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(false)},
			},
		},
		{
			name:  "skip_restream is kept from prior because GET never returns it",
			api:   &replicationapi.Filter{Expressions: &replicationapi.FilterExpressions{RegEx: "^a"}},
			prior: &ReplicationFilter{Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(true)}},
			want:  &ReplicationFilter{Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(true)}},
		},
		{
			name:  "all-false exclude options dropped by the API are kept from prior",
			api:   &replicationapi.Filter{Expressions: &replicationapi.FilterExpressions{RegEx: "^a"}},
			prior: &ReplicationFilter{DocumentExcludeOptions: allFalseExcludeOptions(), Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(false)}},
			want:  &ReplicationFilter{DocumentExcludeOptions: allFalseExcludeOptions(), Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(false)}},
		},
		{
			name:  "a filter with only defaults dropped by the API is kept from prior",
			prior: &ReplicationFilter{DocumentExcludeOptions: allFalseExcludeOptions()},
			want:  &ReplicationFilter{DocumentExcludeOptions: allFalseExcludeOptions()},
		},
		{
			name:  "an empty prior filter is kept as an empty filter",
			prior: &ReplicationFilter{},
			want:  &ReplicationFilter{},
		},
		{
			name:  "exclude options with a true flag that the API dropped show as drift",
			prior: &ReplicationFilter{DocumentExcludeOptions: &ReplicationDocumentExcludeOptions{Deletion: types.BoolValue(true), Expiration: types.BoolValue(false), TTL: types.BoolValue(false), Binary: types.BoolValue(false)}},
			want:  &ReplicationFilter{},
		},
		{
			name:  "expressions that the API dropped show as drift",
			prior: &ReplicationFilter{Expressions: &ReplicationFilterExpressions{RegEx: types.StringValue("^a"), SkipRestream: types.BoolValue(false)}},
			want:  &ReplicationFilter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newReplicationFilter(tt.api, tt.prior))
		})
	}
}

func TestNewReplicationTarget(t *testing.T) {
	prior := &ReplicationTarget{Cluster: types.StringValue("prior-cluster"), Bucket: types.StringValue("cHJpb3I="), Type: types.StringValue(ReplicationTargetTypeCapella)}

	tests := []struct {
		name   string
		target replicationapi.TargetDetails
		prior  *ReplicationTarget
		want   *ReplicationTarget
	}{
		{
			name: "capella target maps to the cluster ID and the bucket ID",
			target: replicationapi.TargetDetails{
				Cluster: replicationapi.Reference{Id: "cluster-2", Name: "target"},
				Bucket:  replicationapi.BucketReference{Id: "dGFyZ2V0", Name: "target"},
				Type:    ReplicationTargetTypeCapella,
			},
			want: &ReplicationTarget{Cluster: types.StringValue("cluster-2"), Bucket: types.StringValue("dGFyZ2V0"), Type: types.StringValue(ReplicationTargetTypeCapella)},
		},
		{
			name: "external target maps to the reference name and the bucket name",
			target: replicationapi.TargetDetails{
				Cluster: replicationapi.Reference{Id: "ref-id", Name: "my-reference"},
				Bucket:  replicationapi.BucketReference{Id: "bucket-id", Name: "travel-sample"},
				Type:    ReplicationTargetTypeExternal,
			},
			want: &ReplicationTarget{Cluster: types.StringValue("my-reference"), Bucket: types.StringValue("travel-sample"), Type: types.StringValue(ReplicationTargetTypeExternal)},
		},
		{
			name:  "empty target from the API keeps the prior target",
			prior: prior,
			want:  prior,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newReplicationTarget(tt.target, tt.prior))
		})
	}
}

func TestNewReplication(t *testing.T) {
	ctx := context.Background()

	t.Run("builds the state from the response and keeps the reverse replication ID from prior", func(t *testing.T) {
		prior := &Replication{ReverseReplicationId: types.StringValue("cmV2LTE=")}

		got, err := NewReplication(ctx, testGetReplicationResponse(), "org-1", "proj-1", "cluster-1", prior)
		require.NoError(t, err)

		assert.Equal(t, "cmVwbC0x", got.Id.ValueString())
		assert.Equal(t, "org-1", got.OrganizationId.ValueString())
		assert.Equal(t, "c291cmNl", got.SourceBucket.ValueString())
		assert.Equal(t, "cmV2LTE=", got.ReverseReplicationId.ValueString())
		assert.Equal(t, "running", got.Status.ValueString())
		assert.Equal(t, int64(5), got.ChangesLeft.ValueInt64())
		assert.True(t, got.Mappings.IsNull(), "no mappings means all scopes, which is null")
		assert.Nil(t, got.Filter)
		assert.False(t, got.SourceDetails.IsNull())
		assert.False(t, got.TargetDetails.IsNull())
		assert.Equal(t, `"user-1"`, got.Audit.Attributes()["created_by"].String())
	})

	t.Run("error is null when the API reports none, and set when it reports one", func(t *testing.T) {
		got, err := NewReplication(ctx, testGetReplicationResponse(), "org-1", "proj-1", "cluster-1", nil)
		require.NoError(t, err)
		assert.True(t, got.Error.IsNull())

		resp := testGetReplicationResponse()
		resp.Error = "target bucket not found"
		got, err = NewReplication(ctx, resp, "org-1", "proj-1", "cluster-1", nil)
		require.NoError(t, err)
		assert.Equal(t, "target bucket not found", got.Error.ValueString())
	})

	t.Run("import has no reverse replication ID", func(t *testing.T) {
		got, err := NewReplication(ctx, testGetReplicationResponse(), "org-1", "proj-1", "cluster-1", nil)
		require.NoError(t, err)
		assert.True(t, got.ReverseReplicationId.IsNull())
	})

	t.Run("an unknown reverse replication ID in the plan becomes null", func(t *testing.T) {
		prior := &Replication{ReverseReplicationId: types.StringUnknown()}
		got, err := NewReplication(ctx, testGetReplicationResponse(), "org-1", "proj-1", "cluster-1", prior)
		require.NoError(t, err)
		assert.True(t, got.ReverseReplicationId.IsNull())
	})

	t.Run("an empty source bucket from the API keeps the prior source bucket", func(t *testing.T) {
		resp := testGetReplicationResponse()
		resp.Source = replicationapi.SourceDetails{}
		prior := &Replication{SourceBucket: types.StringValue("cHJpb3I=")}

		got, err := NewReplication(ctx, resp, "org-1", "proj-1", "cluster-1", prior)
		require.NoError(t, err)
		assert.Equal(t, "cHJpb3I=", got.SourceBucket.ValueString())
	})

	t.Run("mappings in any order give the same set", func(t *testing.T) {
		sales := replicationapi.ScopeMapping{
			SourceScope: "sales", TargetScope: "sales",
			Collections: []replicationapi.CollectionMapping{
				{SourceCollection: "orders", TargetCollection: "orders"},
				{SourceCollection: "returns", TargetCollection: "returns"},
			},
		}
		inventory := replicationapi.ScopeMapping{SourceScope: "inventory", TargetScope: "stock"}

		first := testGetReplicationResponse()
		first.Mappings = []replicationapi.ScopeMapping{sales, inventory}

		reversedSales := sales
		reversedSales.Collections = []replicationapi.CollectionMapping{sales.Collections[1], sales.Collections[0]}
		second := testGetReplicationResponse()
		second.Mappings = []replicationapi.ScopeMapping{inventory, reversedSales}

		a, err := NewReplication(ctx, first, "org-1", "proj-1", "cluster-1", nil)
		require.NoError(t, err)
		b, err := NewReplication(ctx, second, "org-1", "proj-1", "cluster-1", nil)
		require.NoError(t, err)

		assert.True(t, a.Mappings.Equal(b.Mappings))
		assert.Len(t, a.Mappings.Elements(), 2)
	})
}
