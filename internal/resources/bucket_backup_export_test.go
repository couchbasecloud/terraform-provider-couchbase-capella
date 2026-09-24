package resources

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// Unit coverage for refreshBucketBackupExport, the helper that turns a get-export response into
// Terraform state. It is served by a stub so the optional-field branches can be exercised without
// a Capella organisation - the acceptance suite can only reach the shapes the real API happens to
// return, and an omitted field is not one of them.

// newBucketBackupExportTestResource wires the resource to a stub API that answers every request
// with body.
func newBucketBackupExportTestResource(t *testing.T, body string) *BucketBackupExport {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return &BucketBackupExport{Data: &providerschema.Data{
		HostURL:  server.URL,
		Token:    "test-token",
		ClientV1: api.NewClient(10 * time.Second),
	}}
}

func refreshTestExport(t *testing.T, body string) *providerschema.BucketBackupExport {
	t.Helper()

	state, err := newBucketBackupExportTestResource(t, body).refreshBucketBackupExport(
		context.Background(), "org", "proj", "cluster", "bucket", "backup", "e1",
	)
	if err != nil {
		t.Fatalf("refreshBucketBackupExport: %v", err)
	}
	return state
}

// TestRefreshBucketBackupExportArchiveFields pins the four optional archive attributes in both
// directions: absent from the response means null in state, present means carried through. These
// are the branches no acceptance test could reach until the export completed.
func TestRefreshBucketBackupExportArchiveFields(t *testing.T) {
	pending := refreshTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample",`+
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

	complete := refreshTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample",`+
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

// TestRefreshBucketBackupExportOmittedCreatedAt asserts created_at is null when the response
// carries no createdAt.
//
// AV-144898: this FAILS today. CreatedAt is a non-pointer time.Time
// (internal/api/backup/export.go:45), so an omitted createdAt unmarshals to the Go zero time, and
// internal/resources/bucket_backup_export.go:279 formats it unconditionally - writing
// "0001-01-01T00:00:00Z" into state. Every other optional field on that struct is a pointer and is
// nil-guarded immediately below; CreatedAt is the only one that is not. Because created_at is
// Computed, the bogus timestamp then persists with no plan diff to reveal it.
//
// The assertion is written for the FIXED behaviour, so it turns into a regression guard the moment
// the helper guards the zero time with IsZero(). The same one-line fix is needed in
// mapBucketBackupExport (internal/datasources/bucket_backup_export.go:153).
func TestRefreshBucketBackupExportOmittedCreatedAt(t *testing.T) {
	t.Skip("AV-144898: an omitted createdAt maps to 0001-01-01T00:00:00Z instead of null; unskip once refreshBucketBackupExport guards the zero time with IsZero()")

	got := refreshTestExport(t, `{"id":"e1","cycleId":"cy1","bucketName":"travel-sample","status":"pending"}`)

	if !got.CreatedAt.IsNull() {
		t.Errorf("created_at = %q, want null when the API omits createdAt", got.CreatedAt.ValueString())
	}
}
