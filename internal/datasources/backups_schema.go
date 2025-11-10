package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var backupBuilder = capellaschema.NewSchemaBuilder("backup")

func BackupSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", backupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", backupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "cluster_id", backupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "bucket_id", backupBuilder, requiredString())

	backupStatsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(backupStatsAttrs, "size_in_mb", backupBuilder, &schema.Float64Attribute{Computed: true})
	capellaschema.AddAttr(backupStatsAttrs, "items", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "mutations", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "tombstones", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "gsi", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "fts", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "cbas", backupBuilder, computedInt64())
	capellaschema.AddAttr(backupStatsAttrs, "event", backupBuilder, computedInt64())

	scheduleInfoAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(scheduleInfoAttrs, "backup_type", backupBuilder, computedString())
	capellaschema.AddAttr(scheduleInfoAttrs, "backup_time", backupBuilder, computedString())
	capellaschema.AddAttr(scheduleInfoAttrs, "increment", backupBuilder, computedInt64())
	capellaschema.AddAttr(scheduleInfoAttrs, "retention", backupBuilder, computedString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "organization_id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "project_id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "bucket_id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cycle_id", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "date", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "restore_before", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "status", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "method", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "bucket_name", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "source", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", backupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "backup_stats", backupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: backupStatsAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "elapsed_time_in_seconds", backupBuilder, computedInt64())
	capellaschema.AddAttr(dataAttrs, "schedule_info", backupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: scheduleInfoAttrs,
	})

	capellaschema.AddAttr(attrs, "data", backupBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	return schema.Schema{
		MarkdownDescription: "The backups data source retrieves backups associated with a bucket for an operational cluster.",
		Attributes:          attrs,
	}
}
