package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotRestoreBuilder = capellaschema.NewSchemaBuilder("snapshotRestore")

func SnapshotRestoreSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "id", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "project_id", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "organization_id", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "created_at", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "restore_to", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "snapshot", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", snapshotRestoreBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all snapshot restore information for a cluster.",
		Attributes:          attrs,
	}
}
