package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Builders for nested referenced schemas (single backup)
var snapshotBackupCmekBuilder = capellaschema.NewSchemaBuilder("cmek", "ClusterCMEKConfig")
var snapshotBackupProgressBuilder = capellaschema.NewSchemaBuilder("progress", "CloudSnapshotBackupProgress")
var snapshotBackupCrossRegionCopiesBuilder = capellaschema.NewSchemaBuilder("crossRegionCopies", "CloudSnapshotBackupCrossRegionCopies")
var snapshotBackupServerBuilder = capellaschema.NewSchemaBuilder("server", "CouchbaseServer")

func getSnapshotBackupCmekAttrs() map[string]schema.Attribute {
	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", snapshotBackupCmekBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", snapshotBackupCmekBuilder, computedString())
	return cmekAttrs
}

func getSnapshotBackupProgressAttrs() map[string]schema.Attribute {
	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", snapshotBackupProgressBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", snapshotBackupProgressBuilder, computedString())
	return progressAttrs
}

func getSnapshotBackupCrossRegionCopiesAttrs() map[string]schema.Attribute {
	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", snapshotBackupCrossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", snapshotBackupCrossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", snapshotBackupCrossRegionCopiesBuilder, computedString())
	return crossRegionCopiesAttrs
}

func getSnapshotBackupServerAttrs() map[string]schema.Attribute {
	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", snapshotBackupServerBuilder, computedString())
	return serverAttrs
}

func SnapshotBackupSchema() schema.Schema {

	snapshotBackupBuilder := capellaschema.NewSchemaBuilder("snapshotBackup", "GetCloudSnapshotBackupResponse")

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
		Attributes: getSnapshotBackupProgressAttrs(),
	})

	capellaschema.AddAttr(attrs, "cmek", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getSnapshotBackupCmekAttrs(),
		},
	})

	capellaschema.AddAttr(attrs, "server", snapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getSnapshotBackupServerAttrs(),
	})

	capellaschema.AddAttr(attrs, "cross_region_copies", snapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getSnapshotBackupCrossRegionCopiesAttrs(),
		},
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backup data source retrieves one snapshot backup.",
		Attributes:          attrs,
	}
}
