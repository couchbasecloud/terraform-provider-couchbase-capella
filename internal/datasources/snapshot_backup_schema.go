package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func SnapshotBackupSchema() schema.Schema {

	snapshotBackupBuilder := capellaschema.NewSchemaBuilder("snapshotBackup")

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "id", snapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "created_at", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "expiration", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "retention", snapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "size", snapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "type", snapshotBackupBuilder, computedString())

	capellaschema.AddAttr(attrs, "progress", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getProgressAttrs(),
	})

	capellaschema.AddAttr(attrs, "cmek", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCmekAttrs(),
		},
	})

	capellaschema.AddAttr(attrs, "server", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getServerAttrs(),
	})

	capellaschema.AddAttr(attrs, "cross_region_copies", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCrossRegionCopiesAttrs(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backup data source retrieves one snapshot backup.",
		Attributes:          attrs,
	}
}
