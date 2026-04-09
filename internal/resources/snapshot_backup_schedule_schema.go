package resources

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var snapshotBackupScheduleBuilder = capellaschema.NewSchemaBuilder("snapshotBackupSchedule", "cloudSnapshotBackupSchedule")

func SnapshotBackupScheduleSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", snapshotBackupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", snapshotBackupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", snapshotBackupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "interval", snapshotBackupScheduleBuilder, int64Attribute(required))
	capellaschema.AddAttr(attrs, "retention", snapshotBackupScheduleBuilder, int64Attribute(required))
	capellaschema.AddAttr(attrs, "start_time", snapshotBackupScheduleBuilder, stringDefaultAttribute(time.Now().Truncate(time.Hour).Format(time.RFC3339), optional, computed))
	capellaschema.AddAttr(attrs, "copy_to_regions", snapshotBackupScheduleBuilder, stringSetAttribute(optional, computed))

	return schema.Schema{
		MarkdownDescription: "Manages snapshot backup schedule resource",
		Attributes:          attrs,
	}
}
