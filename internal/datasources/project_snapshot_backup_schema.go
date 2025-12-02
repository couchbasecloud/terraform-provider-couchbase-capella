package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var projectSnapshotBackupBuilder = capellaschema.NewSchemaBuilder("projectSnapshotBackup", "GetProjectLevelCloudSnapshotBackupResponse")

// Separate builder for nested snapshot fields which reference GetCloudSnapshotBackupResponse
var nestedSnapshotBuilder = capellaschema.NewSchemaBuilder("nestedSnapshot", "GetCloudSnapshotBackupResponse")

// Builders for pagination cursor fields
var hrefsBuilder = capellaschema.NewSchemaBuilder("hrefs", "Hrefs")
var pagesBuilder = capellaschema.NewSchemaBuilder("pages", "Pages")

func gethrefsAttrs() map[string]schema.Attribute {
	hrefsAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(hrefsAttrs, "first", hrefsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "last", hrefsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "next", hrefsBuilder, computedString())
	capellaschema.AddAttr(hrefsAttrs, "previous", hrefsBuilder, computedString())

	return hrefsAttrs
}

func getPagesAttrs() map[string]schema.Attribute {
	pagesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(pagesAttrs, "last", pagesBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "next", pagesBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "page", pagesBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "per_page", pagesBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "previous", pagesBuilder, computedInt64())
	capellaschema.AddAttr(pagesAttrs, "total_items", pagesBuilder, computedInt64())

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
	serverBuilder := capellaschema.NewSchemaBuilder("server", "CouchbaseServer")
	serverAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(serverAttrs, "version", serverBuilder, computedString())

	return serverAttrs
}

func getProgressAttrs() map[string]schema.Attribute {
	progressBuilder := capellaschema.NewSchemaBuilder("progress", "CloudSnapshotBackupProgress")
	progressAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(progressAttrs, "status", progressBuilder, computedString())
	capellaschema.AddAttr(progressAttrs, "time", progressBuilder, computedString())

	return progressAttrs
}

func getCrossRegionCopiesAttrs() map[string]schema.Attribute {
	crossRegionCopiesBuilder := capellaschema.NewSchemaBuilder("crossRegionCopies", "CloudSnapshotBackupCrossRegionCopies")
	crossRegionCopiesAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(crossRegionCopiesAttrs, "region_code", crossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "status", crossRegionCopiesBuilder, computedString())
	capellaschema.AddAttr(crossRegionCopiesAttrs, "time", crossRegionCopiesBuilder, computedString())

	return crossRegionCopiesAttrs
}

func getCmekAttrs() map[string]schema.Attribute {
	cmekBuilder := capellaschema.NewSchemaBuilder("cmek", "ClusterCMEKConfig")
	cmekAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(cmekAttrs, "id", cmekBuilder, computedString())
	capellaschema.AddAttr(cmekAttrs, "provider_id", cmekBuilder, computedString())

	return cmekAttrs
}

func getComputedProjectSnapshotAttrs() map[string]schema.Attribute {
	computedProjectSnapshotAttrs := make(map[string]schema.Attribute)
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cluster_id", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "created_at", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "expiration", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "id", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "project_id", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "app_service", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "retention", nestedSnapshotBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "database_size", nestedSnapshotBuilder, computedInt64())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "organization_id", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "type", nestedSnapshotBuilder, computedString())
	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cross_region_copies", nestedSnapshotBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCrossRegionCopiesAttrs(),
		},
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "progress", nestedSnapshotBuilder, &schema.SingleNestedAttribute{
		Computed:   true,
		Attributes: getProgressAttrs(),
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "cmek", nestedSnapshotBuilder, &schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: getCmekAttrs(),
		},
	})

	capellaschema.AddAttr(computedProjectSnapshotAttrs, "server", nestedSnapshotBuilder, &schema.SingleNestedAttribute{
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
