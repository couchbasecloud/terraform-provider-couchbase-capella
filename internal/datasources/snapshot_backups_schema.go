package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupsBuilder = capellaschema.NewSchemaBuilder("snapshotBackups")

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
		Attributes: getProgressAttrs(),
	})

	capellaschema.AddAttr(dataAttrs, "cmek", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCmekAttrs(),
		},
	})

	capellaschema.AddAttr(dataAttrs, "server", snapshotBackupsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getServerAttrs(),
	})

	capellaschema.AddAttr(dataAttrs, "cross_region_copies", snapshotBackupsBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCrossRegionCopiesAttrs(),
		},
	})

	return dataAttrs
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
