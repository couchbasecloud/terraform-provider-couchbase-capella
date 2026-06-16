package resources

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// createAccessFromSlice converts the terraform schema Access slice to the API Access slice.
func createAccessFromSlice(access []providerschema.Access) []api.Access {
	result := make([]api.Access, len(access))
	for i, acc := range access {
		result[i] = api.Access{Privileges: make([]string, len(acc.Privileges))}
		for j, priv := range acc.Privileges {
			result[i].Privileges[j] = priv.ValueString()
		}
		if acc.Resources != nil {
			if acc.Resources.Buckets != nil {
				result[i].Resources = &api.AccessibleResources{Buckets: make([]api.Bucket, len(acc.Resources.Buckets))}
				for k, bucket := range acc.Resources.Buckets {
					result[i].Resources.Buckets[k].Name = bucket.Name.ValueString()
					if bucket.Scopes != nil {
						result[i].Resources.Buckets[k].Scopes = make([]api.Scope, len(bucket.Scopes))
						for s, scope := range bucket.Scopes {
							result[i].Resources.Buckets[k].Scopes[s].Name = scope.Name.ValueString()
							if scope.Collections != nil {
								result[i].Resources.Buckets[k].Scopes[s].Collections = make([]string, len(scope.Collections))
								for c, coll := range scope.Collections {
									result[i].Resources.Buckets[k].Scopes[s].Collections[c] = coll.ValueString()
								}
							}
						}
					}
				}
			}
		} else {
			// Workaround: pass an empty list of buckets when no resources are specified
			// to avoid nil pointer issues in the V4 API.
			result[i].Resources = &api.AccessibleResources{Buckets: make([]api.Bucket, 0)}
		}
	}
	return result
}

// mapAccessFromSlice creates a copy of the terraform schema Access slice for state storage.
func mapAccessFromSlice(access []providerschema.Access) []providerschema.Access {
	result := make([]providerschema.Access, len(access))
	for i, acc := range access {
		result[i] = providerschema.Access{Privileges: make([]types.String, len(acc.Privileges))}
		copy(result[i].Privileges, acc.Privileges)
		if acc.Resources != nil {
			if acc.Resources.Buckets != nil {
				result[i].Resources = &providerschema.Resources{Buckets: make([]providerschema.BucketResource, len(acc.Resources.Buckets))}
				for k, bucket := range acc.Resources.Buckets {
					result[i].Resources.Buckets[k].Name = bucket.Name
					if bucket.Scopes != nil {
						result[i].Resources.Buckets[k].Scopes = make([]providerschema.ScopeResource, len(bucket.Scopes))
						for s, scope := range bucket.Scopes {
							result[i].Resources.Buckets[k].Scopes[s].Name = scope.Name
							if scope.Collections != nil {
								result[i].Resources.Buckets[k].Scopes[s].Collections = make([]types.String, len(scope.Collections))
								copy(result[i].Resources.Buckets[k].Scopes[s].Collections, scope.Collections)
							}
						}
					}
				}
			}
		}
	}
	return result
}
