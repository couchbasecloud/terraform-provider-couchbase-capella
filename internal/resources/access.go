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

// mapAccessFromAPI converts the API Access slice to the terraform schema Access slice.
func mapAccessFromAPI(apiAccess []api.Access) []providerschema.Access {
	access := make([]providerschema.Access, len(apiAccess))
	for i, acc := range apiAccess {
		access[i] = providerschema.Access{Privileges: make([]types.String, len(acc.Privileges))}
		for j, permission := range acc.Privileges {
			access[i].Privileges[j] = types.StringValue(permission)
		}
		if acc.Resources != nil && acc.Resources.Buckets != nil {
			access[i].Resources = &providerschema.Resources{Buckets: providerschema.MapBucketsFromAPI(acc.Resources.Buckets)}
		}
	}
	return access
}

// reconcileAccess merges the API response access with the prior state to
// prevent perpetual drift. It serves two purposes.
//
// First, ordering: access and its nested resources are List attributes, so
// Terraform treats element order as significant. The API is not guaranteed to
// echo entries back in the order they were sent, so the result is emitted in
// prior state order — entries are matched to the API response by privilege set
// rather than by position. API entries with no state counterpart are appended
// so genuine remote additions still surface as a diff, and state entries with
// no API counterpart are dropped because they no longer exist remotely.
//
// Second, wildcards: for global privileges the V4 API returns a wildcard bucket
// even when the user omitted the resources field. That pattern is detected and
// the prior state's nil resources preserved so Terraform does not report an
// unnecessary diff.
func reconcileAccess(apiAccess, stateAccess []providerschema.Access) []providerschema.Access {
	if stateAccess == nil {
		return apiAccess
	}

	// Index API entries by sorted privileges to allow matching regardless of
	// the order the API returned them in.
	apiByPrivileges := make(map[string][]int, len(apiAccess))
	for i, aa := range apiAccess {
		key := privilegesKey(aa.Privileges)
		apiByPrivileges[key] = append(apiByPrivileges[key], i)
	}

	consumed := make([]bool, len(apiAccess))
	result := make([]providerschema.Access, 0, len(apiAccess))

	for _, stateEntry := range stateAccess {
		key := privilegesKey(stateEntry.Privileges)
		candidates := apiByPrivileges[key]
		if len(candidates) == 0 {
			continue
		}

		// Pop the first matching candidate to handle duplicate privilege sets.
		idx := candidates[0]
		apiByPrivileges[key] = candidates[1:]
		consumed[idx] = true

		apiEntry := apiAccess[idx]
		entry := providerschema.Access{Privileges: apiEntry.Privileges, Resources: apiEntry.Resources}

		// The two privilege slices hold the same values by construction, so
		// adopting the state ordering keeps element order stable across reads.
		if len(stateEntry.Privileges) == len(apiEntry.Privileges) {
			entry.Privileges = stateEntry.Privileges
		}

		switch {
		case stateEntry.Resources == nil && apiEntry.Resources != nil && isWildcardOnlyResourcesSchema(apiEntry.Resources):
			// The user did not specify resources and the API returned only
			// the implicit wildcard — suppress the diff by keeping nil.
			entry.Resources = nil
		case stateEntry.Resources != nil && apiEntry.Resources != nil:
			entry.Resources = &providerschema.Resources{
				Buckets: reconcileBuckets(apiEntry.Resources.Buckets, stateEntry.Resources.Buckets),
			}
		}

		result = append(result, entry)
	}

	for i, aa := range apiAccess {
		if !consumed[i] {
			result = append(result, aa)
		}
	}

	return result
}

// reconcileBuckets reorders the API buckets to follow prior state order,
// matching on the bucket name. Unmatched API buckets are appended.
func reconcileBuckets(apiBuckets, stateBuckets []providerschema.BucketResource) []providerschema.BucketResource {
	if apiBuckets == nil || stateBuckets == nil {
		return apiBuckets
	}

	stateOrder := make([]providerschema.BucketResource, 0, len(stateBuckets))
	consumed := make([]bool, len(apiBuckets))

	for _, stateBucket := range stateBuckets {
		for i, apiBucket := range apiBuckets {
			if consumed[i] || apiBucket.Name.ValueString() != stateBucket.Name.ValueString() {
				continue
			}
			consumed[i] = true
			stateOrder = append(stateOrder, providerschema.BucketResource{
				Name:   apiBucket.Name,
				Scopes: reconcileScopes(apiBucket.Scopes, stateBucket.Scopes),
			})
			break
		}
	}

	for i, apiBucket := range apiBuckets {
		if !consumed[i] {
			stateOrder = append(stateOrder, apiBucket)
		}
	}

	return stateOrder
}

// reconcileScopes reorders the API scopes to follow prior state order, matching
// on the scope name. Unmatched API scopes are appended.
func reconcileScopes(apiScopes, stateScopes []providerschema.ScopeResource) []providerschema.ScopeResource {
	if apiScopes == nil || stateScopes == nil {
		return apiScopes
	}

	stateOrder := make([]providerschema.ScopeResource, 0, len(apiScopes))
	consumed := make([]bool, len(apiScopes))

	for _, stateScope := range stateScopes {
		for i, apiScope := range apiScopes {
			if consumed[i] || apiScope.Name.ValueString() != stateScope.Name.ValueString() {
				continue
			}
			consumed[i] = true
			stateOrder = append(stateOrder, providerschema.ScopeResource{
				Name:        apiScope.Name,
				Collections: reconcileCollections(apiScope.Collections, stateScope.Collections),
			})
			break
		}
	}

	for i, apiScope := range apiScopes {
		if !consumed[i] {
			stateOrder = append(stateOrder, apiScope)
		}
	}

	return stateOrder
}

// reconcileCollections reorders the API collection names to follow prior state
// order. Unmatched API collections are appended.
func reconcileCollections(apiCollections, stateCollections []types.String) []types.String {
	if apiCollections == nil || stateCollections == nil {
		return apiCollections
	}

	stateOrder := make([]types.String, 0, len(apiCollections))
	consumed := make([]bool, len(apiCollections))

	for _, stateColl := range stateCollections {
		for i, apiColl := range apiCollections {
			if consumed[i] || apiColl.ValueString() != stateColl.ValueString() {
				continue
			}
			consumed[i] = true
			stateOrder = append(stateOrder, apiColl)
			break
		}
	}

	for i, apiColl := range apiCollections {
		if !consumed[i] {
			stateOrder = append(stateOrder, apiColl)
		}
	}

	return stateOrder
}

// isWildcardOnlyResourcesSchema returns true when the resources contain only
// a single wildcard ("*") bucket with no scopes.
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
