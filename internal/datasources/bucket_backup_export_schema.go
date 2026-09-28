package datasources

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var bucketBackupExportBuilder = capellaschema.NewSchemaBuilder("bucketBackupExport", "GetBackupExportResponse")

func BucketBackupExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	capellaschema.AddAttr(attrs, "organization_id", bucketBackupExportBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "project_id", bucketBackupExportBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "cluster_id", bucketBackupExportBuilder, requiredUUIDString())
	capellaschema.AddAttr(attrs, "bucket_id", bucketBackupExportBuilder, requiredStringWithValidator())

	backupId := requiredUUIDString()
	backupId.MarkdownDescription = capellaschema.ExportBackupIdDescription
	capellaschema.AddAttr(attrs, "backup_id", bucketBackupExportBuilder, backupId)

	exportId := requiredUUIDString()
	exportId.MarkdownDescription = capellaschema.ExportIdDescription
	capellaschema.AddAttr(attrs, "id", bucketBackupExportBuilder, exportId)

	cycleId := computedString()
	cycleId.MarkdownDescription = capellaschema.ExportCycleIdDescription
	capellaschema.AddAttr(attrs, "cycle_id", bucketBackupExportBuilder, cycleId)

	bucketName := computedString()
	bucketName.MarkdownDescription = capellaschema.ExportBucketNameDescription
	capellaschema.AddAttr(attrs, "bucket_name", bucketBackupExportBuilder, bucketName)

	status := computedString()
	status.MarkdownDescription = capellaschema.ExportStatusDescription
	capellaschema.AddAttr(attrs, "status", bucketBackupExportBuilder, status)

	createdAt := computedString()
	createdAt.MarkdownDescription = capellaschema.ExportCreatedAtDescription
	capellaschema.AddAttr(attrs, "created_at", bucketBackupExportBuilder, createdAt)

	sizeInBytes := computedInt64()
	sizeInBytes.MarkdownDescription = capellaschema.ExportSizeInBytesDescription
	capellaschema.AddAttr(attrs, "size_in_bytes", bucketBackupExportBuilder, sizeInBytes)

	checksum := computedString()
	checksum.MarkdownDescription = capellaschema.ExportChecksumDescription
	capellaschema.AddAttr(attrs, "sha256_checksum", bucketBackupExportBuilder, checksum)

	expiration := computedString()
	expiration.MarkdownDescription = capellaschema.ExportExpirationDescription
	capellaschema.AddAttr(attrs, "expiration", bucketBackupExportBuilder, expiration)

	// The URL carries short lived cloud storage credentials, so keep it out of plan output.
	downloadURL := computedString()
	downloadURL.Sensitive = true
	downloadURL.MarkdownDescription = capellaschema.ExportDownloadURLDescription
	capellaschema.AddAttr(attrs, "backup_download_url", bucketBackupExportBuilder, downloadURL)

	return schema.Schema{
		MarkdownDescription: "The bucket backup export data source retrieves one backup export job. " +
			"There is no endpoint to list exports, so the export is addressed by its `id`. " +
			"Use it to fetch a fresh `backup_download_url` from a configuration other than the one that created the export.",
		Attributes: attrs,
	}
}
