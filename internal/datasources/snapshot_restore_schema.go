package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func SnapshotRestoreSchema() schema.Schema {

	snapshotRestoreBuilder := capellaschema.NewSchemaBuilder("snapshotRestore", "GetCloudSnapshotRestoreResponse")

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "id", snapshotRestoreBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotRestoreBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", snapshotRestoreBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "organization_id", snapshotRestoreBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "created_at", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "restore_to", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "snapshot", snapshotRestoreBuilder, computedString())
	capellaschema.AddAttr(attrs, "status", snapshotRestoreBuilder, computedString())

	return schema.Schema{
		MarkdownDescription: "The data source to retrieve all snapshot restore information for a cluster.",
		Attributes:          attrs,
	}
}
