package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupsBuilder = capellaschema.NewSchemaBuilder("snapshotBackups", "GetCloudSnapshotBackupResponse")

// Use the same builders as snapshot_backup_schema for consistency
var snapshotBackupsCmekBuilder = capellaschema.NewSchemaBuilder("cmek", "ClusterCMEKConfig")
var snapshotBackupsProgressBuilder = capellaschema.NewSchemaBuilder("progress", "CloudSnapshotBackupProgress")
var snapshotBackupsCrossRegionCopiesBuilder = capellaschema.NewSchemaBuilder("crossRegionCopies", "CloudSnapshotBackupCrossRegionCopies")
var snapshotBackupsServerBuilder = capellaschema.NewSchemaBuilder("server", "CouchbaseServer")

func getSnapshotBackupsDataAttrs() map[string]schema.Attribute {
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "created_at", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "expiration", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "id", snapshotBackupsBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "retention", snapshotBackupsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "size", snapshotBackupsBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "type", snapshotBackupsBuilder, computedString())

	capellaschema.AddAttr(dataAttrs, "progress", snapshotBackupsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getSnapshotBackupsProgressAttrs(),
	})

	capellaschema.AddAttr(dataAttrs, "cmek", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getSnapshotBackupsCmekAttrs(),
		},
	})

	capellaschema.AddAttr(dataAttrs, "server", snapshotBackupsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getSnapshotBackupsServerAttrs(),
	})

	capellaschema.AddAttr(dataAttrs, "cross_region_copies", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getSnapshotBackupsCrossRegionCopiesAttrs(),
		},
	})

	return dataAttrs
}

func getSnapshotBackupsCmekAttrs() map[string]schema.Attribute {
	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", snapshotBackupsCmekBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", snapshotBackupsCmekBuilder, computedString())
	return cmekAttrs
}

func getSnapshotBackupsProgressAttrs() map[string]schema.Attribute {
	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", snapshotBackupsProgressBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", snapshotBackupsProgressBuilder, computedString())
	return progressAttrs
}

func getSnapshotBackupsCrossRegionCopiesAttrs() map[string]schema.Attribute {
	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", snapshotBackupsCrossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", snapshotBackupsCrossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", snapshotBackupsCrossRegionCopiesBuilder, computedString())
	return crossRegionCopiesAttrs
}

func getSnapshotBackupsServerAttrs() map[string]schema.Attribute {
	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", snapshotBackupsServerBuilder, computedString())
	return serverAttrs
}

func SnapshotBackupsSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupsBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupsBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupsBuilder, requiredStringWithValidator())

	capellaschema.AddAttr(attrs, "data", snapshotBackupsBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getSnapshotBackupsDataAttrs(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a bucket for an operational cluster.",
		Attributes:          attrs,
		Blocks: map[string]schema.Block{
			"filters": schema.SingleNestedBlock{
				Attributes: getFilterAttrs(),
			},
		},
	}
}
