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
			"refresh the resource until its status is complete to obtain `backup_download_url`, a pre-signed URL valid for one hour.\n\n" +
			"An export is a server side job that Capella cannot cancel or delete, so its lifecycle differs from most resources. " +
			"Destroy only removes the export from state.\n\n" +
			"**The archive expires 12 hours after the export completes.** `status` becomes `expired` and `backup_download_url` becomes null. " +
			"Both are computed, so `terraform plan` proposes no change when this happens; check `status` or `expiration` directly. " +
			"To get a new archive, run `terraform apply -replace` on this resource.\n\n" +
			"**Replacing fails until then.** Only one active export can exist per backup cycle, so `-replace` and `taint` fail with " +
			"error 14061 while the export is pending or processing, and 14062 while it is complete and not yet expired.\n\n" +
			"**The export is re-created about 7 days after it completes.** Capella then drops the export record, so the resource is removed " +
			"from state and the next `terraform apply` starts a new export. Configurations applied on a schedule or from CI will keep " +
			"re-exporting, so remove the resource from the configuration once you have the download.\n\n" +
			"**A failed export stays in state** until you replace it, since Capella keeps failed export records indefinitely.",
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
