package resources

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func convertScopes(scopes []providerschema.ScopeResource) []api.Scope {
	result := make([]api.Scope, len(scopes))
	for s, scope := range scopes {
		result[s].Name = scope.Name.ValueString()
		if scope.Collections != nil {
			result[s].Collections = make([]string, len(scope.Collections))
			for c, coll := range scope.Collections {
				result[s].Collections[c] = coll.ValueString()
			}
		}
	}
	return result
}

func convertBuckets(buckets []providerschema.BucketResource) []api.Bucket {
	result := make([]api.Bucket, len(buckets))
	for k, bucket := range buckets {
		result[k].Name = bucket.Name.ValueString()
		if bucket.Scopes != nil {
			result[k].Scopes = convertScopes(bucket.Scopes)
		}
	}
	return result
}

// createAccessFromSlice converts the terraform schema Access slice to the API Access slice.
func createAccessFromSlice(access []providerschema.Access) []api.Access {
	result := make([]api.Access, len(access))
	for i, acc := range access {
		result[i] = api.Access{Privileges: make([]string, len(acc.Privileges))}
		for j, priv := range acc.Privileges {
			result[i].Privileges[j] = priv.ValueString()
		}
		if acc.Resources == nil {
			// Workaround: pass an empty list of buckets when no resources are specified
			// to avoid nil pointer issues in the V4 API.
			result[i].Resources = &api.AccessibleResources{Buckets: make([]api.Bucket, 0)}
			continue
		}
		if acc.Resources.Buckets != nil {
			result[i].Resources = &api.AccessibleResources{Buckets: convertBuckets(acc.Resources.Buckets)}
		}
	}
	return result
}

func copyScopeResources(scopes []providerschema.ScopeResource) []providerschema.ScopeResource {
	result := make([]providerschema.ScopeResource, len(scopes))
	for s, scope := range scopes {
		result[s].Name = scope.Name
		if scope.Collections != nil {
			result[s].Collections = make([]types.String, len(scope.Collections))
			copy(result[s].Collections, scope.Collections)
		}
	}
	return result
}

func copyBucketResources(buckets []providerschema.BucketResource) []providerschema.BucketResource {
	result := make([]providerschema.BucketResource, len(buckets))
	for k, bucket := range buckets {
		result[k].Name = bucket.Name
		if bucket.Scopes != nil {
			result[k].Scopes = copyScopeResources(bucket.Scopes)
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
		if acc.Resources == nil {
			continue
		}
		if acc.Resources.Buckets != nil {
			result[i].Resources = &providerschema.Resources{Buckets: copyBucketResources(acc.Resources.Buckets)}
		}
	}
	return result
}
