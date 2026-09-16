package schema

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

func testGetReplicationResponse() apigen.GetReplicationResponse {
	priority := apigen.GetReplicationResponsePriorityHigh
	networkUsageLimit := 25
	conflictResolutionType := "seqno"
	projectId := "target-project"
	projectName := "target project"
	regEx := "^order:"

	response := apigen.GetReplicationResponse{
		Id:                "repl-1",
		Status:            apigen.GetReplicationResponseStatus("running"),
		Direction:         apigen.GetReplicationResponseDirectionOneWay,
		Priority:          &priority,
		NetworkUsageLimit: &networkUsageLimit,
		ChangesLeft:       7,
		Audit: apigen.ReplicationAuditData{
			CreatedAt: time.Date(2026, 9, 16, 10, 0, 0, 0, time.UTC),
			CreatedBy: "user-1",
		},
	}

	response.Source.Bucket.Id = "source-bucket-id"
	response.Source.Bucket.Name = "source-bucket"
	response.Source.Bucket.ConflictResolutionType = conflictResolutionType
	response.Source.Cluster.Id = "source-cluster-id"
	response.Source.Cluster.Name = "source-cluster"
	response.Source.Project.Id = "source-project-id"
	response.Source.Project.Name = "source-project"

	response.Target.Bucket.Id = "target-bucket-id"
	response.Target.Bucket.Name = "target-bucket"
	response.Target.Bucket.ConflictResolutionType = &conflictResolutionType
	response.Target.Cluster.Id = "target-cluster-id"
	response.Target.Cluster.Name = "target-cluster"
	response.Target.Type = apigen.ReplicationTargetType("capella")
	response.Target.Project = &struct {
		Id   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}{Id: &projectId, Name: &projectName}

	mappings := apigen.Mappings{{SourceScope: "sales", TargetScope: "sales"}}
	mappings[0].Collections = &[]struct {
		SourceCollection string `json:"sourceCollection"`
		TargetCollection string `json:"targetCollection"`
	}{{SourceCollection: "orders", TargetCollection: "orders"}}
	response.Mappings = &mappings

	response.Filter = &apigen.GetFilter{}
	response.Filter.Expressions = &struct {
		RegEx *string `json:"regEx,omitempty"`
	}{RegEx: &regEx}

	return response
}

func TestNewReplication(t *testing.T) {
	replication := NewReplication(testGetReplicationResponse(), "org-1", "proj-1", "cluster-1")

	assert.Equal(t, "org-1", replication.OrganizationId.ValueString())
	assert.Equal(t, "proj-1", replication.ProjectId.ValueString())
	assert.Equal(t, "cluster-1", replication.ClusterId.ValueString())
	assert.Equal(t, "source-bucket-id", replication.SourceBucket.ValueString())
	assert.Equal(t, "repl-1", replication.Id.ValueString())
	assert.Equal(t, "running", replication.Status.ValueString())
	assert.Equal(t, "oneWay", replication.Direction.ValueString())
	assert.Equal(t, "high", replication.Priority.ValueString())
	assert.Equal(t, int64(25), replication.NetworkUsageLimit.ValueInt64())
	assert.Equal(t, int64(7), replication.ChangesLeft.ValueInt64())
	assert.Equal(t, "user-1", replication.Audit.CreatedBy.ValueString())

	require.NotNil(t, replication.Target)
	assert.Equal(t, "target-cluster-id", replication.Target.Cluster.ValueString())
	assert.Equal(t, "target-bucket-id", replication.Target.Bucket.ValueString())
	assert.Equal(t, "capella", replication.Target.Type.ValueString())

	require.Len(t, replication.Mappings, 1)
	assert.Equal(t, "sales", replication.Mappings[0].SourceScope.ValueString())
	require.Len(t, replication.Mappings[0].Collections, 1)
	assert.Equal(t, "orders", replication.Mappings[0].Collections[0].SourceCollection.ValueString())

	require.NotNil(t, replication.SourceDetails)
	require.NotNil(t, replication.SourceDetails.Project)
	assert.Equal(t, "source-project", replication.SourceDetails.Project.Name.ValueString())
	assert.Equal(t, "seqno", replication.SourceDetails.Bucket.ConflictResolutionType.ValueString())

	require.NotNil(t, replication.TargetDetails)
	require.NotNil(t, replication.TargetDetails.Project)
	assert.Equal(t, "target project", replication.TargetDetails.Project.Name.ValueString())
}

// TestNewReplicationLeavesSkipRestreamUnset pins the documented behaviour that skip_restream is
// never derived from a read, because the API's read model has no such field.
func TestNewReplicationLeavesSkipRestreamUnset(t *testing.T) {
	replication := NewReplication(testGetReplicationResponse(), "org-1", "proj-1", "cluster-1")

	require.NotNil(t, replication.Filter)
	require.NotNil(t, replication.Filter.Expressions)
	assert.Equal(t, "^order:", replication.Filter.Expressions.RegEx.ValueString())
	assert.True(t, replication.Filter.Expressions.SkipRestream.IsNull())
}

// TestNewReplicationOptionalFieldsBecomeNull guards against absent API values surfacing as a
// misleading "" or 0 in state.
func TestNewReplicationOptionalFieldsBecomeNull(t *testing.T) {
	response := testGetReplicationResponse()
	response.Priority = nil
	response.NetworkUsageLimit = nil
	response.Error = nil
	response.Mappings = nil
	response.Filter = nil
	response.Target.Project = nil

	replication := NewReplication(response, "org-1", "proj-1", "cluster-1")

	assert.True(t, replication.Priority.IsNull())
	assert.True(t, replication.NetworkUsageLimit.IsNull())
	assert.True(t, replication.Error.IsNull())
	assert.True(t, replication.ReverseReplicationId.IsNull())
	assert.Nil(t, replication.Mappings)
	assert.Nil(t, replication.Filter)
	require.NotNil(t, replication.TargetDetails)
	assert.Nil(t, replication.TargetDetails.Project)
}

func testConfiguredReplication() Replication {
	return Replication{
		SourceBucket: types.StringValue("source-bucket-id"),
		Target: &ReplicationTarget{
			Cluster: types.StringValue("target-cluster-id"),
			Bucket:  types.StringValue("target-bucket-id"),
			Type:    types.StringValue("external"),
		},
		Direction:         types.StringValue("twoWay"),
		Priority:          types.StringValue("medium"),
		NetworkUsageLimit: types.Int64Value(10),
		Mappings: []ReplicationMapping{{
			SourceScope: types.StringValue("sales"),
			TargetScope: types.StringValue("sales-copy"),
			Collections: []ReplicationCollection{{
				SourceCollection: types.StringValue("orders"),
				TargetCollection: types.StringValue("orders-copy"),
			}},
		}},
		Filter: &ReplicationFilter{
			DocumentExcludeOptions: &ReplicationDocumentExcludeOptions{
				Deletion:   types.BoolValue(true),
				Expiration: types.BoolValue(false),
			},
			Expressions: &ReplicationFilterExpressions{
				RegEx:        types.StringValue("^order:"),
				SkipRestream: types.BoolValue(true),
			},
		},
	}
}

func TestToCreateRequest(t *testing.T) {
	request := testConfiguredReplication().ToCreateRequest()

	require.NotNil(t, request.Mode)
	assert.Equal(t, apigen.Async, *request.Mode, "creation must always use async mode")

	assert.Equal(t, "source-bucket-id", request.SourceBucket)
	assert.Equal(t, "target-cluster-id", request.Target.Cluster)
	assert.Equal(t, "target-bucket-id", request.Target.Bucket)
	require.NotNil(t, request.Target.Type)
	assert.Equal(t, apigen.CreateReplicationRequestTargetTypeExternal, *request.Target.Type)

	require.NotNil(t, request.Direction)
	assert.Equal(t, apigen.CreateReplicationRequestDirectionTwoWay, *request.Direction)
	require.NotNil(t, request.Priority)
	assert.Equal(t, apigen.CreateReplicationRequestPriorityMedium, *request.Priority)
	require.NotNil(t, request.NetworkUsageLimit)
	assert.Equal(t, 10, *request.NetworkUsageLimit)

	require.NotNil(t, request.Mappings)
	require.Len(t, *request.Mappings, 1)
	mapping := (*request.Mappings)[0]
	assert.Equal(t, "sales", mapping.SourceScope)
	assert.Equal(t, "sales-copy", mapping.TargetScope)
	require.NotNil(t, mapping.Collections)
	require.Len(t, *mapping.Collections, 1)
	assert.Equal(t, "orders", (*mapping.Collections)[0].SourceCollection)
	assert.Equal(t, "orders-copy", (*mapping.Collections)[0].TargetCollection)

	require.NotNil(t, request.Filter)
	require.NotNil(t, request.Filter.DocumentExcludeOptions)
	assert.True(t, *request.Filter.DocumentExcludeOptions.Deletion)
	assert.False(t, *request.Filter.DocumentExcludeOptions.Expiration)
	assert.Nil(t, request.Filter.DocumentExcludeOptions.Ttl, "unset booleans must be omitted, not sent as false")
	require.NotNil(t, request.Filter.Expressions)
	assert.Equal(t, "^order:", *request.Filter.Expressions.RegEx)
	require.NotNil(t, request.Filter.Expressions.SkipRestream)
	assert.True(t, *request.Filter.Expressions.SkipRestream)
}

func TestToCreateRequestOmitsUnsetOptionals(t *testing.T) {
	replication := Replication{
		SourceBucket: types.StringValue("source-bucket-id"),
		Target: &ReplicationTarget{
			Cluster: types.StringValue("target-cluster-id"),
			Bucket:  types.StringValue("target-bucket-id"),
		},
	}

	request := replication.ToCreateRequest()

	assert.Nil(t, request.Target.Type)
	assert.Nil(t, request.Direction)
	assert.Nil(t, request.Priority)
	assert.Nil(t, request.NetworkUsageLimit)
	assert.Nil(t, request.Mappings)
	assert.Nil(t, request.Filter)
}

func TestToUpdateRequest(t *testing.T) {
	request := testConfiguredReplication().ToUpdateRequest()

	require.NotNil(t, request.Priority)
	assert.Equal(t, apigen.UpdateReplicationRequestPriority("medium"), *request.Priority)
	require.NotNil(t, request.Mappings)
	require.Len(t, *request.Mappings, 1)
	assert.Nil(t, request.AllScopes, "allScopes must not be sent while mappings are configured")
}

// TestToUpdateRequestSetsAllScopesWhenMappingsRemoved covers switching a replication back to
// whole-bucket replication, which the API expresses as allScopes rather than as empty mappings.
func TestToUpdateRequestSetsAllScopesWhenMappingsRemoved(t *testing.T) {
	replication := testConfiguredReplication()
	replication.Mappings = nil

	request := replication.ToUpdateRequest()

	require.NotNil(t, request.AllScopes)
	assert.True(t, *request.AllScopes)
	assert.Nil(t, request.Mappings)
}

func TestReplicationValidate(t *testing.T) {
	tests := []struct {
		name        string
		input       Replication
		expectedIds map[Attr]string
		expectError bool
	}{
		{
			name: "[POSITIVE] IDs are passed via terraform apply",
			input: Replication{
				Id:             types.StringValue("repl-1"),
				ClusterId:      types.StringValue("cluster-1"),
				ProjectId:      types.StringValue("proj-1"),
				OrganizationId: types.StringValue("org-1"),
			},
			expectedIds: map[Attr]string{
				Id: "repl-1", ClusterId: "cluster-1", ProjectId: "proj-1", OrganizationId: "org-1",
			},
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: Replication{
				Id: types.StringValue("id=repl-1,cluster_id=cluster-1,project_id=proj-1,organization_id=org-1"),
			},
			expectedIds: map[Attr]string{
				Id: "repl-1", ClusterId: "cluster-1", ProjectId: "proj-1", OrganizationId: "org-1",
			},
		},
		{
			name:        "[NEGATIVE] cluster ID is missing",
			input:       Replication{Id: types.StringValue("repl-1")},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := test.input.Validate()
			if test.expectError {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			for attr, want := range test.expectedIds {
				assert.Equal(t, want, IDs[attr])
			}
		})
	}
}
