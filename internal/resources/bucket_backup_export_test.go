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
// Terraform state. A stub serves the responses so both the pending and the completed shape can be
// exercised in milliseconds, without waiting on the backup infrastructure for the second one.

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
