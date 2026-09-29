package acceptance_tests

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Coverage for what only exists once the export completes: archive attributes, download URL, 14062.

// archiveDownloadTimeout bounds the fetch so a stalled URL fails this step, not the whole suite.
const archiveDownloadTimeout = 5 * time.Minute

// Statuses the export job reports. Only "complete" publishes an archive.
const (
	bucketBackupExportStatusComplete = "complete"
	bucketBackupExportStatusFailed   = "failed"
	bucketBackupExportStatusExpired  = "expired"
)

// bucketBackupExportRef carries the IDs a later PreConfig needs, since PreConfig cannot read state.
type bucketBackupExportRef struct {
	organizationId string
	projectId      string
	clusterId      string
	bucketId       string
	backupId       string
	exportId       string
}

// TestAccBucketBackupExportDownloadArchive downloads the archive and checks size and checksum.
func TestAccBucketBackupExportDownloadArchive(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_dl_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_dl_backup_")
	exportResourceName := randomStringWithPrefix("tf_acc_export_dl_export_")
	dsName := randomStringWithPrefix("tf_acc_export_dl_ds_")

	exportReference := "couchbase-capella_bucket_backup_export." + exportResourceName
	dsReference := "data.couchbase-capella_bucket_backup_export." + dsName

	var export bucketBackupExportRef

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Still pending here; this step creates the export and records its IDs for the next PreConfig.
			{
				Config: testAccBucketBackupExportDownloadConfig(bucketResourceName, backupResourceName, exportResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketBackupExportResource(t, exportReference),
					captureBucketBackupExportRef(exportReference, &export),
				),
			},
			// Same config once the job has finished, so both mapping helpers take their archive branches.
			{
				PreConfig: func() { waitForBucketBackupExportComplete(t, export) },
				Config:    testAccBucketBackupExportDownloadConfig(bucketResourceName, backupResourceName, exportResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(exportReference, "status", bucketBackupExportStatusComplete),
					resource.TestCheckResourceAttrSet(exportReference, "size_in_bytes"),
					resource.TestCheckResourceAttrSet(exportReference, "sha256_checksum"),
					resource.TestCheckResourceAttrSet(exportReference, "expiration"),
					resource.TestCheckResourceAttrSet(exportReference, "backup_download_url"),

					// The data source must describe the same completed export.
					resource.TestCheckResourceAttr(dsReference, "status", bucketBackupExportStatusComplete),
					resource.TestCheckResourceAttrPair(dsReference, "size_in_bytes", exportReference, "size_in_bytes"),
					resource.TestCheckResourceAttrPair(dsReference, "sha256_checksum", exportReference, "sha256_checksum"),
					resource.TestCheckResourceAttrPair(dsReference, "expiration", exportReference, "expiration"),
					// A fresh URL is minted per read, so the two must not match; assert the data source's works.
					resource.TestCheckResourceAttrSet(dsReference, "backup_download_url"),

					testAccDownloadBucketBackupExportArchive(exportReference),
					testAccDownloadBucketBackupExportArchive(dsReference),
				),
			},
		},
	})
}

// TestAccBucketBackupExportDuplicateCompletedCycle pins 14062; the older test accepts either code.
func TestAccBucketBackupExportDuplicateCompletedCycle(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_dupc_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_dupc_backup_")
	firstName := randomStringWithPrefix("tf_acc_export_dupc_first_")
	secondName := randomStringWithPrefix("tf_acc_export_dupc_second_")

	firstReference := "couchbase-capella_bucket_backup_export." + firstName

	var first bucketBackupExportRef

	// The code pins which rejection this is; the prose pins the message support will see quoted.
	duplicateCompleted := regexp.MustCompile(
		`(?s)Error creating backup export job.*` +
			terraformDiagnosticPattern("An export for this backup cycle has already completed.") +
			`.*"code":[\s│]*14062`,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfig(bucketResourceName, backupResourceName, firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketBackupExportResource(t, firstReference),
					captureBucketBackupExportRef(firstReference, &first),
				),
			},
			// Only once the first export has completed does a second one get 14062 rather than 14061.
			{
				PreConfig:   func() { waitForBucketBackupExportComplete(t, first) },
				Config:      testAccBucketBackupExportDuplicateConfig(bucketResourceName, backupResourceName, firstName, secondName),
				ExpectError: duplicateCompleted,
			},
		},
	})
}

// testAccBucketBackupExportDownloadConfig declares bucket, backup, export and a data source.
func testAccBucketBackupExportDownloadConfig(bucketName, backupName, exportName, dsName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
}

resource "couchbase-capella_backup" "%[6]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
}

resource "couchbase-capella_bucket_backup_export" "%[7]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
}

data "couchbase-capella_bucket_backup_export" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
  id              = couchbase-capella_bucket_backup_export.%[7]s.id
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, exportName, dsName)
}

// captureBucketBackupExportRef records an export's IDs; Checks run before the next PreConfig.
func captureBucketBackupExportRef(resourceReference string, out *bucketBackupExportRef) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := bucketBackupExportAttributes(s, resourceReference)
		if err != nil {
			return err
		}

		*out = bucketBackupExportRef{
			organizationId: attrs["organization_id"],
			projectId:      attrs["project_id"],
			clusterId:      attrs["cluster_id"],
			bucketId:       attrs["bucket_id"],
			backupId:       attrs["backup_id"],
			exportId:       attrs["id"],
		}
		return nil
	}
}

// waitForBucketBackupExportComplete polls until the archive exists, failing on any other status.
func waitForBucketBackupExportComplete(t *testing.T, ref bucketBackupExportRef) {
	t.Helper()

	const (
		maxWaitTime   = 30 * time.Minute
		checkInterval = 15 * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), maxWaitTime)
	defer cancel()

	data := newTestClient(t)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		export, err := retrieveBucketBackupExportFromServer(
			data, ref.organizationId, ref.projectId, ref.clusterId, ref.bucketId, ref.backupId, ref.exportId,
		)
		if err != nil {
			t.Fatalf("polling backup export %s: %v", ref.exportId, err)
		}

		switch export.Status {
		case bucketBackupExportStatusComplete:
			return
		case bucketBackupExportStatusFailed, bucketBackupExportStatusExpired:
			t.Fatalf("backup export %s reached terminal status %q without publishing an archive", ref.exportId, export.Status)
		}

		select {
		case <-ctx.Done():
			t.Fatalf("timed out after %s waiting for backup export %s to complete, last status %q",
				maxWaitTime, ref.exportId, export.Status)
		case <-ticker.C:
		}
	}
}

// testAccDownloadBucketBackupExportArchive fetches the archive and verifies its size and checksum.
func testAccDownloadBucketBackupExportArchive(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := bucketBackupExportAttributes(s, resourceReference)
		if err != nil {
			return err
		}

		downloadURL := attrs["backup_download_url"]
		if downloadURL == "" {
			return fmt.Errorf("%s has no backup_download_url in state", resourceReference)
		}

		// Pre-signed, so no auth header; its own client because http.DefaultClient has no timeout.
		client := &http.Client{Timeout: archiveDownloadTimeout}
		resp, err := client.Get(downloadURL) // #nosec G107 -- the URL is minted by the Capella API
		if err != nil {
			return fmt.Errorf("downloading export archive for %s: %w", resourceReference, err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
			return fmt.Errorf("downloading export archive for %s: HTTP %d, body: %s", resourceReference, resp.StatusCode, body)
		}

		// Hash while streaming: the archive never needs to be held in memory whole.
		checksum := sha256.New()
		downloaded, err := io.Copy(checksum, resp.Body)
		if err != nil {
			return fmt.Errorf("reading export archive for %s: %w", resourceReference, err)
		}

		if want := attrs["size_in_bytes"]; want != strconv.FormatInt(downloaded, 10) {
			return fmt.Errorf("%s: downloaded %d bytes but size_in_bytes is %s", resourceReference, downloaded, want)
		}
		if got, want := hex.EncodeToString(checksum.Sum(nil)), attrs["sha256_checksum"]; !strings.EqualFold(got, want) {
			return fmt.Errorf("%s: downloaded archive hashes to %s but sha256_checksum is %s", resourceReference, got, want)
		}
		return nil
	}
}

// terraformDiagnosticPattern tolerates the line wrapping and "│" the Terraform CLI inserts.
func terraformDiagnosticPattern(sentence string) string {
	words := strings.Fields(sentence)
	for i, word := range words {
		words[i] = regexp.QuoteMeta(word)
	}
	return strings.Join(words, `[\s│]+`)
}
