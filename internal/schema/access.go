package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// MapScopesFromAPI converts an API Scope slice to the provider schema ScopeResource slice.
func MapScopesFromAPI(scopes []api.Scope) []ScopeResource {
	result := make([]ScopeResource, len(scopes))
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

// MapBucketsFromAPI converts an API Bucket slice to the provider schema BucketResource slice.
func MapBucketsFromAPI(buckets []api.Bucket) []BucketResource {
	result := make([]BucketResource, len(buckets))
	for k, bucket := range buckets {
		result[k].Name = types.StringValue(bucket.Name)
		if bucket.Scopes != nil {
			result[k].Scopes = MapScopesFromAPI(bucket.Scopes)
		}
	}
	return result
}
