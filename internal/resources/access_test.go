package resources

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func privileges(values ...string) []types.String {
	privs := make([]types.String, len(values))
	for i, v := range values {
		privs[i] = types.StringValue(v)
	}
	return privs
}

func collections(values ...string) []types.String {
	return privileges(values...)
}

func bucket(name string, scopes ...providerschema.ScopeResource) providerschema.BucketResource {
	return providerschema.BucketResource{Name: types.StringValue(name), Scopes: scopes}
}

func scope(name string, colls ...string) providerschema.ScopeResource {
	s := providerschema.ScopeResource{Name: types.StringValue(name)}
	if colls != nil {
		s.Collections = collections(colls...)
	}
	return s
}

func resources(buckets ...providerschema.BucketResource) *providerschema.Resources {
	return &providerschema.Resources{Buckets: buckets}
}

func TestReconcileAccess(t *testing.T) {
	tests := []struct {
		name       string
		apiAccess  []providerschema.Access
		stateEntry []providerschema.Access
		want       []providerschema.Access
	}{
		{
			name: "api returns access entries in reversed order",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
			},
		},
		{
			name: "api reorders buckets, scopes and collections within an entry",
			apiAccess: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b2", scope("s1", "c1")),
						bucket("b1", scope("s2", "c2"), scope("s1", "c2", "c1")),
					),
				},
			},
			stateEntry: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b1", scope("s1", "c1", "c2"), scope("s2", "c2")),
						bucket("b2", scope("s1", "c1")),
					),
				},
			},
			want: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b1", scope("s1", "c1", "c2"), scope("s2", "c2")),
						bucket("b2", scope("s1", "c1")),
					),
				},
			},
		},
		{
			name: "privilege ordering follows state regardless of api ordering",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataWrite", "dataRead")},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead", "dataWrite")},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead", "dataWrite")},
			},
		},
		{
			name: "api entry with no state counterpart is appended last",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
			},
		},
		{
			name: "state entry missing from the api is dropped",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
		},
		{
			name: "nil state access returns the api access untouched",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			stateEntry: nil,
			want: []providerschema.Access{
				{Privileges: privileges("dataWrite"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
		},
		{
			name: "wildcard only resources are suppressed when state omitted them",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("analyticsAdmin"), Resources: resources(bucket("*"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("analyticsAdmin")},
			},
			want: []providerschema.Access{
				{Privileges: privileges("analyticsAdmin")},
			},
		},
		{
			name: "wildcard resources are kept when state declared them",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("*"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("*"))},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("*"))},
			},
		},
		{
			name: "duplicate privilege sets are matched one at a time",
			apiAccess: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			stateEntry: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
			want: []providerschema.Access{
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b2"))},
				{Privileges: privileges("dataRead"), Resources: resources(bucket("b1"))},
			},
		},
		{
			name: "nil and empty nested collections round trip unchanged",
			apiAccess: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b1", providerschema.ScopeResource{Name: types.StringValue("s1")}),
						bucket("b2", providerschema.ScopeResource{Name: types.StringValue("s2"), Collections: []types.String{}}),
						providerschema.BucketResource{Name: types.StringValue("b3")},
					),
				},
			},
			stateEntry: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b1", providerschema.ScopeResource{Name: types.StringValue("s1")}),
						bucket("b2", providerschema.ScopeResource{Name: types.StringValue("s2"), Collections: []types.String{}}),
						providerschema.BucketResource{Name: types.StringValue("b3")},
					),
				},
			},
			want: []providerschema.Access{
				{
					Privileges: privileges("dataRead"),
					Resources: resources(
						bucket("b1", providerschema.ScopeResource{Name: types.StringValue("s1")}),
						bucket("b2", providerschema.ScopeResource{Name: types.StringValue("s2"), Collections: []types.String{}}),
						providerschema.BucketResource{Name: types.StringValue("b3")},
					),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiCopy := deepCopyAccess(tt.apiAccess)

			got := reconcileAccess(tt.apiAccess, tt.stateEntry)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconcileAccess() =\n%s\nwant\n%s", formatAccess(got), formatAccess(tt.want))
			}
			if !reflect.DeepEqual(tt.apiAccess, apiCopy) {
				t.Errorf("reconcileAccess() mutated its apiAccess argument:\n%s\nwas\n%s", formatAccess(tt.apiAccess), formatAccess(apiCopy))
			}
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
			name:  "reordered collections follow state order",
			api:   collections("c3", "c1", "c2"),
			state: collections("c1", "c2", "c3"),
			want:  collections("c1", "c2", "c3"),
		},
		{
			name:  "new api collection is appended",
			api:   collections("c2", "c1"),
			state: collections("c1"),
			want:  collections("c1", "c2"),
		},
		{
			name:  "removed collection is dropped",
			api:   collections("c1"),
			state: collections("c1", "c2"),
			want:  collections("c1"),
		},
		{
			name:  "nil api stays nil",
			api:   nil,
			state: collections("c1"),
			want:  nil,
		},
		{
			name:  "nil state returns api unchanged",
			api:   collections("c2", "c1"),
			state: nil,
			want:  collections("c2", "c1"),
		},
		{
			name:  "empty api stays empty",
			api:   []types.String{},
			state: collections("c1"),
			want:  []types.String{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reconcileCollections(tt.api, tt.state)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconcileCollections() = %v, want %v", got, tt.want)
			}
		})
	}
}

// deepCopyAccess clones an access slice so a test can assert the input was not
// mutated; mapAccessFromSlice is not reused here because it drops resources
// with nil buckets.
func deepCopyAccess(access []providerschema.Access) []providerschema.Access {
	if access == nil {
		return nil
	}
	result := make([]providerschema.Access, len(access))
	for i, acc := range access {
		result[i].Privileges = append([]types.String(nil), acc.Privileges...)
		if acc.Privileges == nil {
			result[i].Privileges = nil
		}
		if acc.Resources == nil {
			continue
		}
		result[i].Resources = &providerschema.Resources{Buckets: copyBucketResources(acc.Resources.Buckets)}
		if acc.Resources.Buckets == nil {
			result[i].Resources.Buckets = nil
		}
	}
	return result
}

func formatAccess(access []providerschema.Access) string {
	out := "["
	for _, acc := range access {
		out += "\n  {privileges: ["
		for i, p := range acc.Privileges {
			if i > 0 {
				out += " "
			}
			out += p.ValueString()
		}
		out += "]"
		if acc.Resources == nil {
			out += ", resources: nil}"
			continue
		}
		out += ", buckets: ["
		for _, b := range acc.Resources.Buckets {
			out += b.Name.ValueString() + "("
			for _, s := range b.Scopes {
				out += s.Name.ValueString() + ":"
				for i, c := range s.Collections {
					if i > 0 {
						out += ","
					}
					out += c.ValueString()
				}
				out += " "
			}
			out += ") "
		}
		out += "]}"
	}
	return out + "\n]"
}
