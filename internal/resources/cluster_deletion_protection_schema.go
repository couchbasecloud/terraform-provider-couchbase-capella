package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var clusterDeletionProtectionBuilder = capellaschema.NewSchemaBuilder("clusterDeletionProtection")

func ClusterDeletionProtectionSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", clusterDeletionProtectionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", clusterDeletionProtectionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", clusterDeletionProtectionBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "deletion_protection", clusterDeletionProtectionBuilder, boolAttribute(required))

	return schema.Schema{
		MarkdownDescription: "Manages the deletion protection state of a Couchbase Capella cluster.",
		Attributes:          attrs,
	}
}
