package resources

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// newAccess builds one access entry. Omitting buckets leaves resources nil,
// which is how a config that skipped the resources block is stored.
func newAccess(privileges string, buckets ...providerschema.BucketResource) providerschema.Access {
	entry := providerschema.Access{Privileges: newStrings(strings.Split(privileges, ",")...)}
	if buckets != nil {
		entry.Resources = &providerschema.Resources{Buckets: buckets}
	}
	return entry
}

// newBucket builds a bucket. Omitting scopes leaves scopes nil.
func newBucket(name string, scopes ...providerschema.ScopeResource) providerschema.BucketResource {
	return providerschema.BucketResource{Name: types.StringValue(name), Scopes: scopes}
}

// newScope builds a scope. Omitting collections leaves collections nil.
func newScope(name string, collections ...string) providerschema.ScopeResource {
	return providerschema.ScopeResource{Name: types.StringValue(name), Collections: newStrings(collections...)}
}

func newStrings(values ...string) []types.String {
	if values == nil {
		return nil
	}
	out := make([]types.String, len(values))
	for i, v := range values {
		out[i] = types.StringValue(v)
	}
	return out
}

func TestReconcileAccess(t *testing.T) {
	tests := []struct {
		name  string
		api   []providerschema.Access
		state []providerschema.Access
		want  []providerschema.Access
	}{
		{
			name:  "entries returned out of order follow state order",
			api:   []providerschema.Access{newAccess("dataWrite", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
			state: []providerschema.Access{newAccess("dataRead", newBucket("b1")), newAccess("dataWrite", newBucket("b2"))},
			want:  []providerschema.Access{newAccess("dataRead", newBucket("b1")), newAccess("dataWrite", newBucket("b2"))},
		},
		{
			name:  "privileges follow state order within a matched entry",
			api:   []providerschema.Access{newAccess("dataWrite,dataRead")},
			state: []providerschema.Access{newAccess("dataRead,dataWrite")},
			want:  []providerschema.Access{newAccess("dataRead,dataWrite")},
		},
		{
			name:  "entry with no state counterpart is appended last",
			api:   []providerschema.Access{newAccess("dataWrite", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
			state: []providerschema.Access{newAccess("dataRead", newBucket("b1"))},
			want:  []providerschema.Access{newAccess("dataRead", newBucket("b1")), newAccess("dataWrite", newBucket("b2"))},
		},
		{
			name:  "state entry missing from the api is dropped",
			api:   []providerschema.Access{newAccess("dataRead", newBucket("b1"))},
			state: []providerschema.Access{newAccess("dataRead", newBucket("b1")), newAccess("dataWrite", newBucket("b2"))},
			want:  []providerschema.Access{newAccess("dataRead", newBucket("b1"))},
		},
		{
			name:  "nil state returns the api response untouched, as on import",
			api:   []providerschema.Access{newAccess("dataWrite", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
			state: nil,
			want:  []providerschema.Access{newAccess("dataWrite", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
		},
		{
			name:  "implicit wildcard bucket is suppressed when state omitted resources",
			api:   []providerschema.Access{newAccess("analyticsAdmin", newBucket("*"))},
			state: []providerschema.Access{newAccess("analyticsAdmin")},
			want:  []providerschema.Access{newAccess("analyticsAdmin")},
		},
		{
			name:  "wildcard bucket is kept when state declared it",
			api:   []providerschema.Access{newAccess("dataRead", newBucket("*"))},
			state: []providerschema.Access{newAccess("dataRead", newBucket("*"))},
			want:  []providerschema.Access{newAccess("dataRead", newBucket("*"))},
		},
		{
			name:  "duplicate privilege sets are matched one at a time",
			api:   []providerschema.Access{newAccess("dataRead", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
			state: []providerschema.Access{newAccess("dataRead", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
			want:  []providerschema.Access{newAccess("dataRead", newBucket("b2")), newAccess("dataRead", newBucket("b1"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, reconcileAccess(tt.api, tt.state))
		})
	}
}

func TestReconcileAccessDoesNotMutateAPIResponse(t *testing.T) {
	api := []providerschema.Access{
		newAccess("dataWrite", newBucket("b2", newScope("s2", "c2"))),
		newAccess("dataRead", newBucket("b1", newScope("s1", "c1"))),
	}
	unchanged := []providerschema.Access{
		newAccess("dataWrite", newBucket("b2", newScope("s2", "c2"))),
		newAccess("dataRead", newBucket("b1", newScope("s1", "c1"))),
	}

	reconcileAccess(api, []providerschema.Access{newAccess("dataRead", newBucket("b1", newScope("s1", "c1")))})

	assert.Equal(t, unchanged, api, "reconcileAccess() must not mutate the api response it was given")
}

func TestReconcileBuckets(t *testing.T) {
	tests := []struct {
		name  string
		api   []providerschema.BucketResource
		state []providerschema.BucketResource
		want  []providerschema.BucketResource
	}{
		{
			name:  "buckets returned out of order follow state order",
			api:   []providerschema.BucketResource{newBucket("b2"), newBucket("b1")},
			state: []providerschema.BucketResource{newBucket("b1"), newBucket("b2")},
			want:  []providerschema.BucketResource{newBucket("b1"), newBucket("b2")},
		},
		{
			name:  "scopes and collections follow state order within a bucket",
			api:   []providerschema.BucketResource{newBucket("b1", newScope("s2", "c2"), newScope("s1", "c2", "c1"))},
			state: []providerschema.BucketResource{newBucket("b1", newScope("s1", "c1", "c2"), newScope("s2", "c2"))},
			want:  []providerschema.BucketResource{newBucket("b1", newScope("s1", "c1", "c2"), newScope("s2", "c2"))},
		},
		{
			name:  "bucket with no state counterpart is appended last",
			api:   []providerschema.BucketResource{newBucket("b2"), newBucket("b1")},
			state: []providerschema.BucketResource{newBucket("b1")},
			want:  []providerschema.BucketResource{newBucket("b1"), newBucket("b2")},
		},
		{
			name:  "state bucket missing from the api is dropped",
			api:   []providerschema.BucketResource{newBucket("b1")},
			state: []providerschema.BucketResource{newBucket("b1"), newBucket("b2")},
			want:  []providerschema.BucketResource{newBucket("b1")},
		},
		{
			name:  "nil api buckets stay nil",
			api:   nil,
			state: []providerschema.BucketResource{newBucket("b1")},
			want:  nil,
		},
		{
			name:  "nil state buckets return the api buckets unchanged",
			api:   []providerschema.BucketResource{newBucket("b2"), newBucket("b1")},
			state: nil,
			want:  []providerschema.BucketResource{newBucket("b2"), newBucket("b1")},
		},
		{
			name:  "nil scopes are not turned into an empty list",
			api:   []providerschema.BucketResource{newBucket("b1")},
			state: []providerschema.BucketResource{newBucket("b1")},
			want:  []providerschema.BucketResource{newBucket("b1")},
		},
		{
			name:  "empty scopes are not turned into nil",
			api:   []providerschema.BucketResource{{Name: types.StringValue("b1"), Scopes: []providerschema.ScopeResource{}}},
			state: []providerschema.BucketResource{{Name: types.StringValue("b1"), Scopes: []providerschema.ScopeResource{}}},
			want:  []providerschema.BucketResource{{Name: types.StringValue("b1"), Scopes: []providerschema.ScopeResource{}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, reconcileBuckets(tt.api, tt.state))
		})
	}
}

func TestReconcileCollections(t *testing.T) {
	tests := []struct {
		name  string
		api   []types.String
		state []types.String
		want  []types.String
	}{
		{
			name:  "collections returned out of order follow state order",
			api:   newStrings("c3", "c1", "c2"),
			state: newStrings("c1", "c2", "c3"),
			want:  newStrings("c1", "c2", "c3"),
		},
		{
			name:  "collection with no state counterpart is appended last",
			api:   newStrings("c2", "c1"),
			state: newStrings("c1"),
			want:  newStrings("c1", "c2"),
		},
		{
			name:  "state collection missing from the api is dropped",
			api:   newStrings("c1"),
			state: newStrings("c1", "c2"),
			want:  newStrings("c1"),
		},
		{
			name:  "nil api collections stay nil",
			api:   nil,
			state: newStrings("c1"),
			want:  nil,
		},
		{
			name:  "nil state collections return the api collections unchanged",
			api:   newStrings("c2", "c1"),
			state: nil,
			want:  newStrings("c2", "c1"),
		},
		{
			name:  "empty api collections stay empty",
			api:   []types.String{},
			state: newStrings("c1"),
			want:  []types.String{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, reconcileCollections(tt.api, tt.state))
		})
	}
}
