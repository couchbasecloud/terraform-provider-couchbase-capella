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

// reconcileAccess merges the API response access with the prior state so
// Terraform does not report a diff for a logically unchanged role.
//
// access and its nested resources are List attributes, so element order is
// significant and the API does not guarantee it echoes back the order it was
// sent. Entries are therefore paired with prior state by privilege set and
// emitted in state order. A nil prior state means there is nothing to align to,
// as on import.
func reconcileAccess(apiAccess, stateAccess []providerschema.Access) []providerschema.Access {
	return reorderToState(apiAccess, stateAccess,
		func(entry providerschema.Access) string { return privilegesKey(entry.Privileges) },
		mergeAccessEntry)
}

// mergeAccessEntry combines a matched API and prior state access entry, keeping
// the state's ordering and its choice to omit the resources block.
func mergeAccessEntry(apiEntry, stateEntry providerschema.Access) providerschema.Access {
	merged := apiEntry

	// Both slices hold the same privileges by construction, so adopting the
	// state ordering keeps element order stable across reads.
	if len(stateEntry.Privileges) == len(apiEntry.Privileges) {
		merged.Privileges = stateEntry.Privileges
	}

	switch {
	case stateEntry.Resources == nil && isWildcardOnlyResourcesSchema(apiEntry.Resources):
		// For global privileges the V4 API returns an implicit wildcard bucket
		// even when the user omitted resources. Keep nil to suppress the diff.
		merged.Resources = nil
	case stateEntry.Resources != nil && apiEntry.Resources != nil:
		merged.Resources = &providerschema.Resources{
			Buckets: reconcileBuckets(apiEntry.Resources.Buckets, stateEntry.Resources.Buckets),
		}
	}

	return merged
}

func reconcileBuckets(apiBuckets, stateBuckets []providerschema.BucketResource) []providerschema.BucketResource {
	return reorderToState(apiBuckets, stateBuckets,
		func(bucket providerschema.BucketResource) string { return bucket.Name.ValueString() },
		func(apiBucket, stateBucket providerschema.BucketResource) providerschema.BucketResource {
			apiBucket.Scopes = reconcileScopes(apiBucket.Scopes, stateBucket.Scopes)
			return apiBucket
		})
}

func reconcileScopes(apiScopes, stateScopes []providerschema.ScopeResource) []providerschema.ScopeResource {
	return reorderToState(apiScopes, stateScopes,
		func(scope providerschema.ScopeResource) string { return scope.Name.ValueString() },
		func(apiScope, stateScope providerschema.ScopeResource) providerschema.ScopeResource {
			apiScope.Collections = reconcileCollections(apiScope.Collections, stateScope.Collections)
			return apiScope
		})
}

func reconcileCollections(apiCollections, stateCollections []types.String) []types.String {
	return reorderToState(apiCollections, stateCollections,
		types.String.ValueString,
		func(apiCollection, _ types.String) types.String { return apiCollection })
}

// reorderToState returns the API elements in prior state order, pairing them by
// key and combining each matched pair with merge. Elements the API added are
// appended so the addition still surfaces as a diff; state elements absent from
// the API are dropped because they no longer exist remotely. A nil slice on
// either side is returned untouched, preserving Terraform's null/empty
// distinction. The API elements are never modified in place.
func reorderToState[T any](apiItems, stateItems []T, key func(T) string, merge func(apiItem, stateItem T) T) []T {
	if apiItems == nil || stateItems == nil {
		return apiItems
	}

	apiByKey := make(map[string][]int, len(apiItems))
	for i, item := range apiItems {
		k := key(item)
		apiByKey[k] = append(apiByKey[k], i)
	}

	consumed := make([]bool, len(apiItems))
	result := make([]T, 0, len(apiItems))

	for _, stateItem := range stateItems {
		k := key(stateItem)
		candidates := apiByKey[k]
		if len(candidates) == 0 {
			continue
		}

		// Pop the first unconsumed match so duplicate keys pair up one at a time.
		i := candidates[0]
		apiByKey[k] = candidates[1:]
		consumed[i] = true
		result = append(result, merge(apiItems[i], stateItem))
	}

	for i, item := range apiItems {
		if !consumed[i] {
			result = append(result, item)
		}
	}

	return result
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
