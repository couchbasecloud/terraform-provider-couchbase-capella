package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var collectionBuilder = capellaschema.NewSchemaBuilder("collection")

func CollectionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", collectionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", collectionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", collectionBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "bucket_id", collectionBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["scope_name"] = WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the scope.")
	attrs["collection_name"] = WithDescription(stringAttribute([]string{required, requiresReplace}), "The name of the collection.")
	attrs["max_ttl"] = WithDescription(int64Attribute(optional, computed), "The maximum Time To Live (TTL) for documents in the collection.")

	return schema.Schema{
		MarkdownDescription: "This resource allows you to manage a collection within a scope in a bucket.",
		Attributes:          attrs,
	}
}
