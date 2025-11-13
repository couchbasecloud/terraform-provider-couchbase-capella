package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectSnapshotBackupBuilder = capellaschema.NewSchemaBuilder("projectSnapshotBackup", "cloudProjectSnapshotBackup")

func ProjectSnapshotBackupSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", projectSnapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "project_id", projectSnapshotBackupBuilder, requiredString())
	capellaschema.AddAttr(attrs, "page", projectSnapshotBackupBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "per_page", projectSnapshotBackupBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "sort_by", projectSnapshotBackupBuilder, optionalString())
	capellaschema.AddAttr(attrs, "sort_direction", projectSnapshotBackupBuilder, optionalString())

	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "cluster_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_name", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "creation_date_time", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "created_by", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "current_status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "region", projectSnapshotBackupBuilder, computedString())

	computedProjectSnapshotAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cluster_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "created_at", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "expiration", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "progress", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "project_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "app_service", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cmek", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cross_region_copies", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "retention", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "server", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "database_size", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "organization_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "type", projectSnapshotBackupBuilder, computedString())

	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cross_region_copies", projectSnapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: crossRegionCopiesAttrs,
		},
	})

	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "progress", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: progressAttrs,
	})

	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cmek", projectSnapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: cmekAttrs,
		},
	})

	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "server", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: serverAttrs,
	})

	capellaschema.AddAttr(dataAttrs, "most_recent_snapshot", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computedProjectSnapshotAttrs,
	})
	capellaschema.AddAttr(dataAttrs, "oldest_snapshot", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: computedProjectSnapshotAttrs,
	})

	capellaschema.AddAttr(attrs, "data", projectSnapshotBackupBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: dataAttrs,
		},
	})

	// Build cursor attributes for pagination
	hrefsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(hrefsAttrs, "first", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "last", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "next", projectEventsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "previous", projectEventsBuilder, computedString())

	pagesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(pagesAttrs, "last", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "next", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "page", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "per_page", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "previous", projectEventsBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "total_items", projectEventsBuilder, computedInt64())

	cursorAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cursorAttrs, "hrefs", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: hrefsAttrs,
	})
	capellaschema.AddAttr(cursorAttrs, "pages", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: pagesAttrs,
	})

	capellaschema.AddAttr(attrs, "cursor", projectEventsBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: cursorAttrs,
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a project",
		Attributes:          attrs,
	}
}
