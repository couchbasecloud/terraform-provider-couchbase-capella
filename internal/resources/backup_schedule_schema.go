package resources

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var backupScheduleBuilder = capellaschema.NewSchemaBuilder("backupSchedule", "scheduledBackup")

func BackupScheduleSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", backupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", backupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", backupScheduleBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "bucket_id", backupScheduleBuilder, requiredNonEmptyStringAttribute())
	capellaschema.AddAttr(attrs, "type", backupScheduleBuilder, stringAttribute([]string{required, requiresReplace}))

	weeklyScheduleAttrs := make(map[string]schema.Attribute)
	// day_of_week, incremental_every, and retention_time get their OneOf enum
	// validators auto-attached from the OpenAPI spec via AddAttr (see AV-132154).
	// start_at is the exception: the spec models it as an enum of 0..23, but a
	// Between range validator gives a clearer message, so override it here.
	capellaschema.AddAttr(weeklyScheduleAttrs, "day_of_week", backupScheduleBuilder, stringAttribute([]string{required}), "WeeklySchedule")
	capellaschema.AddAttr(weeklyScheduleAttrs, "start_at", backupScheduleBuilder, &schema.Int64Attribute{
		Required: true,
		Validators: []validator.Int64{
			int64validator.Between(0, 23),
		},
	}, "WeeklySchedule")
	capellaschema.AddAttr(weeklyScheduleAttrs, "incremental_every", backupScheduleBuilder, int64Attribute(required), "WeeklySchedule")
	capellaschema.AddAttr(weeklyScheduleAttrs, "retention_time", backupScheduleBuilder, stringAttribute([]string{required}), "WeeklySchedule")
	capellaschema.AddAttr(weeklyScheduleAttrs, "cost_optimized_retention", backupScheduleBuilder, boolAttribute(required), "WeeklySchedule")

	capellaschema.AddAttr(attrs, "weekly_schedule", backupScheduleBuilder, &schema.SingleNestedAttribute{
		Required:   true,
		Attributes: weeklyScheduleAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "Manages the backup schedule resource associated with a bucket for an operational cluster.",
		Attributes:          attrs,
	}
}
