package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var backupScheduleBuilder = capellaschema.NewSchemaBuilder("backupSchedule")

func BackupScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "project_id", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "cluster_id", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "bucket_id", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))
	capellaschema.AddAttr(attrs, "type", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	attrs["weekly_schedule"] = schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"day_of_week":              stringAttribute([]string{required}),
			"start_at":                 int64Attribute(required),
			"incremental_every":        int64Attribute(required),
			"retention_time":           stringAttribute([]string{required}),
			"cost_optimized_retention": boolAttribute(required),
		},
	}

	return schema.Schema{
		MarkdownDescription: "Manages the backup schedule resource associated with a bucket for an operational cluster.",
		Attributes:          attrs,
	}
}
