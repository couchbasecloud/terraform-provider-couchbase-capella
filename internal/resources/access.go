package resources

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
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
			// Global privileges (e.g. analyticsAdmin) can be created without a
			// resources field. Pass an empty buckets slice so the API receives a
			// well-formed request body rather than a nil pointer.
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

func mapScopesFromAPI(scopes []api.Scope) []providerschema.ScopeResource {
	result := make([]providerschema.ScopeResource, len(scopes))
	for s, scope := range scopes {
		result[s].Name = types.StringValue(scope.Name)
		if scope.Collections != nil {
			result[s].Collections = make([]types.String, len(scope.Collections))
			for c, coll := range scope.Collections {
				result[s].Collections[c] = types.StringValue(coll)
			}
		}
	}
	return result
}

func mapBucketsFromAPI(buckets []api.Bucket) []providerschema.BucketResource {
	result := make([]providerschema.BucketResource, len(buckets))
	for k, bucket := range buckets {
		result[k].Name = types.StringValue(bucket.Name)
		if bucket.Scopes != nil {
			result[k].Scopes = mapScopesFromAPI(bucket.Scopes)
		}
	}
	return result
}

// mapAccessFromAPI converts the API Access slice to the terraform schema Access slice.
func mapAccessFromAPI(apiAccess []api.Access) []providerschema.Access {
	access := make([]providerschema.Access, len(apiAccess))
	for i, acc := range apiAccess {
		access[i] = providerschema.Access{Privileges: make([]types.String, len(acc.Privileges))}
		for j, permission := range acc.Privileges {
			access[i].Privileges[j] = types.StringValue(permission)
		}
		if acc.Resources != nil && acc.Resources.Buckets != nil {
			access[i].Resources = &providerschema.Resources{Buckets: mapBucketsFromAPI(acc.Resources.Buckets)}
		}
	}
	return access
}

// reconcileAccess merges the API response access with the prior state to
// prevent perpetual drift. For global privileges the V4 API returns a wildcard
// bucket even when the user omitted the resources field. This function detects
// that pattern and preserves the prior state's nil resources so Terraform does
// not report an unnecessary diff.
func reconcileAccess(apiAccess, stateAccess []providerschema.Access) []providerschema.Access {
	if stateAccess == nil {
		return apiAccess
	}

	// Build a lookup of state access entries keyed by sorted privileges to
	// allow matching entries regardless of ordering.
	type accessKey struct {
		index int
		entry providerschema.Access
	}
	stateByPrivileges := make(map[string][]accessKey, len(stateAccess))
	for i, sa := range stateAccess {
		key := privilegesKey(sa.Privileges)
		stateByPrivileges[key] = append(stateByPrivileges[key], accessKey{index: i, entry: sa})
	}

	result := make([]providerschema.Access, len(apiAccess))
	for i, apiEntry := range apiAccess {
		result[i] = apiAccess[i]
		key := privilegesKey(apiEntry.Privileges)
		candidates, ok := stateByPrivileges[key]
		if !ok || len(candidates) == 0 {
			continue
		}

		// Pop the first matching candidate to handle duplicate privilege sets.
		stateEntry := candidates[0].entry
		stateByPrivileges[key] = candidates[1:]

		if stateEntry.Resources == nil && apiEntry.Resources != nil && isWildcardOnlyResourcesSchema(apiEntry.Resources) {
			// The user did not specify resources and the API returned only
			// the implicit wildcard — suppress the diff by keeping nil.
			result[i].Resources = nil
		}
	}
	return result
}

// isWildcardOnlyResourcesSchema mirrors isWildcardOnlyResources for the
// schema representation.
func isWildcardOnlyResourcesSchema(res *providerschema.Resources) bool {
	if res == nil || len(res.Buckets) != 1 {
		return false
	}
	b := res.Buckets[0]
	return b.Name.ValueString() == "*" && len(b.Scopes) == 0
}

// privilegesKey builds a deterministic string from a set of privilege values
// for use as a map key.
func privilegesKey(privs []types.String) string {
	sorted := make([]string, len(privs))
	for i, p := range privs {
		sorted[i] = p.ValueString()
	}
	// Sort for deterministic matching regardless of API ordering.
	sort.Strings(sorted)
	return strings.Join(sorted, "\x00")
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
