package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectSnapshotBackupBuilder = capellaschema.NewSchemaBuilder("projectSnapshotBackup")

func gethrefsAttrs() map[string]schema.Attribute {
	hrefsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(hrefsAttrs, "first", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "last", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "next", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "previous", projectSnapshotBackupBuilder, computedString())

	return hrefsAttrs
}

func getPagesAttrs() map[string]schema.Attribute {
	pagesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(pagesAttrs, "last", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "next", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "page", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "per_page", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "previous", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "total_items", projectSnapshotBackupBuilder, computedInt64())

	return pagesAttrs
}

func getCursorAttrs() map[string]schema.Attribute {
	cursorAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cursorAttrs, "hrefs", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: gethrefsAttrs(),
	})
	capellaschema.AddAttr(cursorAttrs, "pages", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getPagesAttrs(),
	})

	return cursorAttrs
}

func getServerAttrs() map[string]schema.Attribute {
	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", projectSnapshotBackupBuilder, computedString())

	return serverAttrs
}

func getProgressAttrs() map[string]schema.Attribute {
	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", projectSnapshotBackupBuilder, computedString())

	return progressAttrs
}

func getCrossRegionCopiesAttrs() map[string]schema.Attribute {
	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", projectSnapshotBackupBuilder, computedString())

	return crossRegionCopiesAttrs
}

func getCmekAttrs() map[string]schema.Attribute {
	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", projectSnapshotBackupBuilder, computedString())

	return cmekAttrs
}

func getComputedProjectSnapshotAttrs() map[string]schema.Attribute {
	computedProjectSnapshotAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cluster_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "created_at", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "expiration", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "project_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "app_service", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "retention", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "database_size", projectSnapshotBackupBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "organization_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "type", projectSnapshotBackupBuilder, computedString())

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cross_region_copies", projectSnapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCrossRegionCopiesAttrs(),
		},
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "progress", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getProgressAttrs(),
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cmek", projectSnapshotBackupBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCmekAttrs(),
		},
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "server", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getServerAttrs(),
	})

	return computedProjectSnapshotAttrs
}

func getProjectSnapshotBackupsDataAttrs() map[string]schema.Attribute {
	dataAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(dataAttrs, "cluster_id", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cluster_name", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "creation_date_time", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "created_by", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "current_status", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "cloud_provider", projectSnapshotBackupBuilder, computedString())
	capellaschema.AddAttr(dataAttrs, "region", projectSnapshotBackupBuilder, computedString())

	capellaschema.AddAttr(dataAttrs, "most_recent_snapshot", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getComputedProjectSnapshotAttrs(),
	})
	capellaschema.AddAttr(dataAttrs, "oldest_snapshot", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getComputedProjectSnapshotAttrs(),
	})

	return dataAttrs
}

func ProjectSnapshotBackupSchema() schema.Schema {

	attrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(attrs, "organization_id", projectSnapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "project_id", projectSnapshotBackupBuilder, requiredStringWithValidator())
	capellaschema.AddAttr(attrs, "page", projectSnapshotBackupBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "per_page", projectSnapshotBackupBuilder, optionalInt64())
	capellaschema.AddAttr(attrs, "sort_by", projectSnapshotBackupBuilder, optionalString())
	capellaschema.AddAttr(attrs, "sort_direction", projectSnapshotBackupBuilder, optionalString())

	capellaschema.AddAttr(attrs, "data", projectSnapshotBackupBuilder, &schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getProjectSnapshotBackupsDataAttrs(),
		},
	})

	capellaschema.AddAttr(attrs, "cursor", projectSnapshotBackupBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getCursorAttrs(),
	})

	return schema.Schema{
		MarkdownDescription: "The snapshot backups data source retrieves snapshot backups associated with a project",
		Attributes:          attrs,
	}
}
