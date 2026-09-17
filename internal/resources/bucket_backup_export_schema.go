package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	capellaschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

var bucketBackupExportBuilder = capellaschema.NewSchemaBuilder("bucketBackupExport", "GetBackupExportResponse")

func BucketBackupExportSchema() schema.Schema {
	attrs := make(map[string]schema.Attribute)

	id := stringAttribute([]string{computed, useStateForUnknown})
	id.MarkdownDescription = capellaschema.ExportIdDescription
	capellaschema.AddAttr(attrs, "id", bucketBackupExportBuilder, id)

	capellaschema.AddAttr(attrs, "organization_id", bucketBackupExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "project_id", bucketBackupExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "cluster_id", bucketBackupExportBuilder, requiredUUIDStringAttribute())
	capellaschema.AddAttr(attrs, "bucket_id", bucketBackupExportBuilder, requiredNonEmptyStringAttribute())

	backupId := requiredUUIDStringAttribute()
	backupId.MarkdownDescription = capellaschema.ExportBackupIdDescription
	capellaschema.AddAttr(attrs, "backup_id", bucketBackupExportBuilder, backupId)

	addBucketBackupExportComputedAttrs(attrs)

	return schema.Schema{
		MarkdownDescription: "This resource allows you to export a bucket backup for an operational cluster into a downloadable zip archive. " +
			"The export runs asynchronously on the backup infrastructure, so it is still pending when the apply returns; " +
			"refresh the resource until its status is complete to obtain `backup_download_url`, a pre-signed URL valid for one hour. " +
			"The archive is deleted from cloud storage at `expiration`, roughly 12 hours after the export completes, " +
			"after which the backup must be exported again.",
		Attributes: attrs,
	}
}

// addBucketBackupExportComputedAttrs adds the attributes the get endpoint returns. The archive
// details are only set once the export completes.
func addBucketBackupExportComputedAttrs(attrs map[string]schema.Attribute) {
	cycleId := stringAttribute([]string{computed})
	cycleId.MarkdownDescription = capellaschema.ExportCycleIdDescription
	capellaschema.AddAttr(attrs, "cycle_id", bucketBackupExportBuilder, cycleId)

	bucketName := stringAttribute([]string{computed})
	bucketName.MarkdownDescription = capellaschema.ExportBucketNameDescription
	capellaschema.AddAttr(attrs, "bucket_name", bucketBackupExportBuilder, bucketName)

	status := stringAttribute([]string{computed})
	status.MarkdownDescription = capellaschema.ExportStatusDescription
	capellaschema.AddAttr(attrs, "status", bucketBackupExportBuilder, status)

	createdAt := stringAttribute([]string{computed})
	createdAt.MarkdownDescription = capellaschema.ExportCreatedAtDescription
	capellaschema.AddAttr(attrs, "created_at", bucketBackupExportBuilder, createdAt)

	sizeInBytes := int64Attribute(computed)
	sizeInBytes.MarkdownDescription = capellaschema.ExportSizeInBytesDescription
	capellaschema.AddAttr(attrs, "size_in_bytes", bucketBackupExportBuilder, sizeInBytes)

	checksum := stringAttribute([]string{computed})
	checksum.MarkdownDescription = capellaschema.ExportChecksumDescription
	capellaschema.AddAttr(attrs, "sha256_checksum", bucketBackupExportBuilder, checksum)

	expiration := stringAttribute([]string{computed})
	expiration.MarkdownDescription = capellaschema.ExportExpirationDescription
	capellaschema.AddAttr(attrs, "expiration", bucketBackupExportBuilder, expiration)

	// The URL carries short lived cloud storage credentials, so keep it out of plan output.
	downloadURL := stringAttribute([]string{computed, sensitive})
	downloadURL.MarkdownDescription = capellaschema.ExportDownloadURLDescription
	capellaschema.AddAttr(attrs, "backup_download_url", bucketBackupExportBuilder, downloadURL)
}
