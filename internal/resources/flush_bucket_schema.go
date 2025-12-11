package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var flushBucketBuilder = capellaschema.NewSchemaBuilder("flushBucket")

func FlushBucketSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", flushBucketBuilder, stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)))
	capellaschema.AddAttr(attrs, "project_id", flushBucketBuilder, stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)))
	capellaschema.AddAttr(attrs, "cluster_id", flushBucketBuilder, stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)))
	capellaschema.AddAttr(attrs, "bucket_id", flushBucketBuilder, stringAttribute([]string{required, requiresReplace}, stringvalidator.LengthAtLeast(1)))

	return schema.Schema{
		MarkdownDescription: "This resource allows you to flush a Capella bucket. Flushing a bucket causes all documents in the bucket to be deleted by the system at the earliest. This operation can only be performed if the bucket has been configured with flushEnabled to 'true'.",
		Attributes:          attrs,
	}
}
