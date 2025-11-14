package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupBuilder = capellaschema.NewSchemaBuilder("snapshotBackup")

func SnapshotBackupSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "id", snapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "created_at", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "expiration", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "retention", snapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "size", snapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(attrs, "type", snapshotBackupBuilder, computedString())

	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "progress", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: progressAttrs,
	})

	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "cmek", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: cmekAttrs,
		},
	})

	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "server", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: serverAttrs,
	})

	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", snapshotBackupBuilder, computedString())
	capellaschema.AddAttr(attrs, "cross_region_copies", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: crossRegionCopiesAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backup data source retrieves one snapshot backup.",
		Attributes:          attrs,
	}
}
