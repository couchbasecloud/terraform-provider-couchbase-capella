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

	weeklyScheduleAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(weeklyScheduleAttrs, "day_of_week", backupScheduleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(weeklyScheduleAttrs, "start_at", backupScheduleBuilder, int64Attribute(required))
	capellaschema.AddAttr(weeklyScheduleAttrs, "incremental_every", backupScheduleBuilder, int64Attribute(required))
	capellaschema.AddAttr(weeklyScheduleAttrs, "retention_time", backupScheduleBuilder, stringAttribute([]string{required}))
	capellaschema.AddAttr(weeklyScheduleAttrs, "cost_optimized_retention", backupScheduleBuilder, boolAttribute(required))

	attrs["weekly_schedule"] = &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: weeklyScheduleAttrs,
	}

	return schema.Schema{
		MarkdownDescription: "Manages the backup schedule resource associated with a bucket for an operational cluster.",
		Attributes:          attrs,
	}
}
