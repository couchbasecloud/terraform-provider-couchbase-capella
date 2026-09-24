package datasources

import (
	"encoding/json"
	"testing"

	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Unit coverage for mapBucketBackupExport, the data source's half of the export mapping. It is a
// pure function, so every response shape is reachable here - including the ones the live API only
// produces after the job completes, and the ones it may never produce at all.

func mapTestExport(t *testing.T, body string) providerschema.BucketBackupExportData {
	t.Helper()

	var resp backupapi.GetBucketBackupExportResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshalling test payload: %v", err)
	}
	return mapBucketBackupExport(&resp, providerschema.BucketBackupExportData{})
}

// TestMapBucketBackupExportArchiveFields pins the four optional archive attributes in both
// directions, mirroring TestRefreshBucketBackupExportArchiveFields on the resource side. The two
// helpers are separate code that must agree, so they need separate guards.
func TestMapBucketBackupExportArchiveFields(t *testing.T) {
	pending := mapTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample",`+
		`"status":"pending","createdAt":"2026-09-20T10:00:00Z"}`)

	for name, got := range map[string]bool{
		"size_in_bytes":       pending.SizeInBytes.IsNull(),
		"sha256_checksum":     pending.Sha256Checksum.IsNull(),
		"expiration":          pending.Expiration.IsNull(),
		"backup_download_url": pending.BackupDownloadURL.IsNull(),
	} {
		if !got {
			t.Errorf("pending export: %s should be null before the archive exists", name)
		}
	}

	complete := mapTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample",`+
		`"status":"complete","createdAt":"2026-09-20T10:00:00Z","sizeInBytes":1234,`+
		`"sha256Checksum":"abc123","expiration":"2026-09-20T22:00:00Z",`+
		`"backupDownloadURL":"https://example.invalid/archive.zip?sig=1"}`)

	if got := complete.SizeInBytes.ValueInt64(); got != 1234 {
		t.Errorf("size_in_bytes = %d, want 1234", got)
	}
	if got := complete.Sha256Checksum.ValueString(); got != "abc123" {
		t.Errorf("sha256_checksum = %q, want %q", got, "abc123")
	}
	if got := complete.Expiration.ValueString(); got != "2026-09-20T22:00:00Z" {
		t.Errorf("expiration = %q, want %q", got, "2026-09-20T22:00:00Z")
	}
	if complete.BackupDownloadURL.IsNull() {
		t.Error("backup_download_url should be set once the export is complete")
	}
}

// TestMapBucketBackupExportOmittedCreatedAt asserts created_at is null when the response carries
// no createdAt.
//
// AV-144898: this FAILS today, for the same reason and with the same one-line fix as
// TestRefreshBucketBackupExportOmittedCreatedAt on the resource side. CreatedAt is a non-pointer
// time.Time (internal/api/backup/export.go:45), so an omitted createdAt unmarshals to the Go zero
// time and internal/datasources/bucket_backup_export.go:153 formats it unconditionally, yielding
// "0001-01-01T00:00:00Z".
//
// The assertion is written for the FIXED behaviour so it becomes a regression guard once the
// helper guards the zero time with IsZero().
func TestMapBucketBackupExportOmittedCreatedAt(t *testing.T) {
	t.Skip("AV-144898: an omitted createdAt maps to 0001-01-01T00:00:00Z instead of null; unskip once mapBucketBackupExport guards the zero time with IsZero()")

	got := mapTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample","status":"pending"}`)

	if !got.CreatedAt.IsNull() {
		t.Errorf("created_at = %q, want null when the API omits createdAt", got.CreatedAt.ValueString())
	}
}
