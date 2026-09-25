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

// Coverage for the half of the bucket backup export lifecycle that only exists once the job
// has finished: the archive attributes, the download URL, and the API's refusal to export a
// cycle that already has a completed export.
//
// The rest of the export suite deliberately stops at the pending state, which leaves the
// `if ... != nil` branch of every archive attribute unexecuted in both mapping helpers. These
// tests are the only thing that reaches them.
//
// They wait on the backup infrastructure, so they are slower than the rest of the suite and
// are kept out of sanity.list. Both export a bucket that was created empty moments earlier,
// so the archive itself is small.

// Statuses the export job reports. Only "complete" publishes an archive.
const (
	bucketBackupExportStatusComplete = "complete"
	bucketBackupExportStatusFailed   = "failed"
	bucketBackupExportStatusExpired  = "expired"
)

// bucketBackupExportRef is the set of IDs needed to poll one export job. A step's PreConfig
// cannot read Terraform state, so the preceding step's Check records them here instead.
type bucketBackupExportRef struct {
	organizationId string
	projectId      string
	clusterId      string
	bucketId       string
	backupId       string
	exportId       string
}

// TestAccBucketBackupExportDownloadArchive proves the exported backup is actually retrievable:
// it waits for the job to complete, then downloads the archive from the pre-signed URL and
// checks the bytes against the size and checksum Capella reported.
//
// It also covers the completed-export read path through both the resource and the data source,
// which nothing else in the suite does.
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
			// The export is still pending here. This step exists to create it and record the
			// IDs the next step's PreConfig needs.
			{
				Config: testAccBucketBackupExportDownloadConfig(bucketResourceName, backupResourceName, exportResourceName, dsName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsBucketBackupExportResource(t, exportReference),
					captureBucketBackupExportRef(exportReference, &export),
				),
			},
			// The same configuration, applied again once the job has finished. The resource is
			// refreshed and the data source re-read, so both mapping helpers take their
			// archive branches.
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
					// The URL is the one attribute that must NOT match: Capella mints a fresh
					// pre-signed URL on every read, so a pair assertion here would be wrong.
					// Assert instead that the data source's own URL works.
					resource.TestCheckResourceAttrSet(dsReference, "backup_download_url"),

					testAccDownloadBucketBackupExportArchive(exportReference),
					testAccDownloadBucketBackupExportArchive(dsReference),
				),
			},
		},
	})
}

// TestAccBucketBackupExportDuplicateCompletedCycle asserts the 14062 rejection specifically -
// the cycle already has a *completed* export.
//
// TestAccBucketBackupExportResourceDuplicateCycle accepts either 14061 or 14062, and because
// create does not wait, the second export there almost always races in while the first is
// still pending and gets 14061. So without this test the completed-export rejection is never
// exercised at all.
func TestAccBucketBackupExportDuplicateCompletedCycle(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_dupc_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_dupc_backup_")
	firstName := randomStringWithPrefix("tf_acc_export_dupc_first_")
	secondName := randomStringWithPrefix("tf_acc_export_dupc_second_")

	firstReference := "couchbase-capella_bucket_backup_export." + firstName

	var first bucketBackupExportRef

	// Both halves are asserted on purpose. The code pins which rejection this is; the prose
	// pins the message support will actually see quoted in a ticket.
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
			// Only once the first export has completed does a second one get 14062 rather
			// than 14061.
			{
				PreConfig:   func() { waitForBucketBackupExportComplete(t, first) },
				Config:      testAccBucketBackupExportDuplicateConfig(bucketResourceName, backupResourceName, firstName, secondName),
				ExpectError: duplicateCompleted,
			},
		},
	})
}

// testAccBucketBackupExportDownloadConfig declares a bucket, a backup, an export and a data
// source reading that export back, so one apply exercises both read paths.
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

// captureBucketBackupExportRef records the IDs of an applied export so a later step can poll
// it. Check functions run before the next step's PreConfig, which is what makes this work.
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

// waitForBucketBackupExportComplete polls until the export job publishes an archive, and fails
// the test rather than returning if it reaches any other terminal state - a test that carried
// on past a failed export would assert against null archive attributes and be confusing.
//
// Export duration scales with backup size. These tests export a bucket created empty moments
// earlier, so the wait is normally short; the ceiling is generous only to absorb a busy
// backup fleet.
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

// testAccDownloadBucketBackupExportArchive downloads the archive from the pre-signed URL held
// in state and checks it against the size and checksum Capella reported for it.
//
// Asserting that backup_download_url is merely set would pass against a URL that 403s or
// serves an error document. This is the only assertion in the suite that proves an exported
// backup can actually be retrieved.
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

		// The URL is pre-signed, so it is fetched without the Capella auth header.
		resp, err := http.Get(downloadURL) // #nosec G107 -- the URL is minted by the Capella API
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

// terraformDiagnosticPattern turns a plain sentence into a regex that still matches once the
// Terraform CLI has wrapped the diagnostic across lines and prefixed each continuation with
// "│". Asserting on a message longer than the diagnostic box without this is a coin flip that
// depends on where the wrap happens to land.
func terraformDiagnosticPattern(sentence string) string {
	words := strings.Fields(sentence)
	for i, word := range words {
		words[i] = regexp.QuoteMeta(word)
	}
	return strings.Join(words, `[\s│]+`)
}
