package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func SnapshotBackupScheduleSchema() schema.Schema {
	snapshotBackupScheduleBuilder := capellaschema.NewSchemaBuilder("snapshotBackupSchedule", "GetCloudSnapshotBackupScheduleResponse")

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupScheduleBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupScheduleBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupScheduleBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "interval", snapshotBackupScheduleBuilder, requiredInt64())
	capellaschema.AddAttr(attrs, "retention", snapshotBackupScheduleBuilder, requiredInt64())
	capellaschema.AddAttr(attrs, "start_time", snapshotBackupScheduleBuilder, computedString())
	capellaschema.AddAttr(attrs, "copy_to_regions", snapshotBackupScheduleBuilder, computedStringSet())
	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves the snapshot backup schedule for a cluster.",
		Attributes:          attrs,
	}
}
