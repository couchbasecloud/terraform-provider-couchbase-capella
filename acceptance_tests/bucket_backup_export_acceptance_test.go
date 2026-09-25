package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	backupapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/backup"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAccBucketBackupExportResource creates a bucket, backs it up and exports that backup.
//
// Everything runs on a bucket created per test run rather than the shared globalBucketId, for the
// same reason TestAccBackupResource does it: Capella serialises manual backups per bucket, and only
// one active export can exist per backup cycle, so a leaked export from a previously-killed run on
// a shared bucket would fail this test with a 409.
//
// The export is still pending when the apply returns - create deliberately does not wait for the
// archive - so the assertions cover the attributes the create and its single refresh populate, not
// the archive attributes that only appear once the export completes.
func TestAccBucketBackupExportResource(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_bucket_")
	bucketResourceReference := "couchbase-capella_bucket." + bucketResourceName

	backupResourceName := randomStringWithPrefix("tf_acc_export_backup_")
	backupResourceReference := "couchbase-capella_backup." + backupResourceName

	resourceName := randomStringWithPrefix("tf_acc_bucket_backup_export_")
	resourceReference := "couchbase-capella_bucket_backup_export." + resourceName

	steps := []resource.TestStep{
		// Create and Read testing.
		{
			Config: testAccBucketBackupExportResourceConfig(bucketResourceName, backupResourceName, resourceName),
			Check: resource.ComposeAggregateTestCheckFunc(
				testAccExistsBucketBackupExportResource(t, resourceReference),
				resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
				resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
				resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
				resource.TestCheckResourceAttrPair(resourceReference, "bucket_id", bucketResourceReference, "id"),
				resource.TestCheckResourceAttrPair(resourceReference, "backup_id", backupResourceReference, "id"),
				resource.TestCheckResourceAttrPair(resourceReference, "bucket_name", bucketResourceReference, "name"),
				// The export belongs to the cycle of the backup it was created from.
				resource.TestCheckResourceAttrPair(resourceReference, "cycle_id", backupResourceReference, "cycle_id"),
				resource.TestCheckResourceAttrSet(resourceReference, "id"),
				resource.TestCheckResourceAttrSet(resourceReference, "created_at"),
				// Status is whatever the job has reached by the time of the refresh.
				resource.TestCheckResourceAttrSet(resourceReference, "status"),
			),
		},
		// ImportState testing. The resource does not support update, so Create / Read /
		// Import / Delete is the whole lifecycle.
		//
		// ImportStateVerify is what makes this step worth having: ImportState alone only
		// proves the six-part composite ID parses, and would pass just as happily against
		// an import that returned every attribute empty.
		{
			ResourceName:      resourceReference,
			ImportStateIdFunc: generateBucketBackupExportImportIdForResource(resourceReference),
			ImportState:       true,
			ImportStateVerify: true,
			// The job advances on its own between the create refresh and this import, and
			// Capella mints a fresh pre-signed URL on every read, so these five cannot be
			// compared. What stays verified is the part import can actually get wrong:
			// that the composite ID repopulates every ID, plus cycle_id, bucket_name and
			// created_at.
			ImportStateVerifyIgnore: []string{
				"status",
				"size_in_bytes",
				"sha256_checksum",
				"expiration",
				"backup_download_url",
			},
		},
	}

	// Each of the five inputs carries RequiresReplace, which is the only reason Update is
	// unreachable - and Update returns an unconditional error, so if a modifier were dropped,
	// every change to this resource would start failing at apply instead. These plan-only
	// steps assert the replacement is still planned. Nothing is applied, so no further exports
	// are created and no further Capella resources are touched.
	steps = append(steps, bucketBackupExportRequiresReplaceSteps(
		bucketResourceName, backupResourceName, resourceName, resourceReference,
	)...)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

// bucketBackupExportRequiresReplaceSteps returns one plan-only step per required input, each
// changing exactly that input and asserting Terraform plans a replacement, followed by a step
// that restores the original configuration and asserts the plan is then empty.
//
// The restore step earns its place twice over: it guards against a perpetual diff, and it
// leaves the working directory holding a configuration whose IDs are real, so the test case's
// own destroy runs against the fixtures it created.
//
// PlanOnly steps skip the pre-apply plan entirely, so ConfigPlanChecks.PreApply would never
// run here - the checks have to hang off PostApplyPostRefresh.
func bucketBackupExportRequiresReplaceSteps(bucketName, backupName, exportName, resourceReference string) []resource.TestStep {
	defaults := defaultBucketBackupExportInputs(bucketName, backupName)

	// Well formed but different from the fixtures, so the change is a value change rather than
	// a validation failure.
	const (
		otherUUID     = "99999999-9999-4999-8999-999999999999"
		otherBucketId = "dGZfYWNjX290aGVyX2J1Y2tldA=="
	)

	changed := func(mutate func(*bucketBackupExportInputs)) resource.TestStep {
		inputs := defaults
		mutate(&inputs)

		return resource.TestStep{
			Config:             testAccBucketBackupExportResourceConfigInputs(bucketName, backupName, exportName, inputs),
			PlanOnly:           true,
			ExpectNonEmptyPlan: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPostRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
				},
			},
		}
	}

	return []resource.TestStep{
		changed(func(in *bucketBackupExportInputs) { in.organizationId = strconv.Quote(otherUUID) }),
		changed(func(in *bucketBackupExportInputs) { in.projectId = strconv.Quote(otherUUID) }),
		changed(func(in *bucketBackupExportInputs) { in.clusterId = strconv.Quote(otherUUID) }),
		changed(func(in *bucketBackupExportInputs) { in.bucketId = strconv.Quote(otherBucketId) }),
		changed(func(in *bucketBackupExportInputs) { in.backupId = strconv.Quote(otherUUID) }),
		{
			Config:   testAccBucketBackupExportResourceConfigInputs(bucketName, backupName, exportName, defaults),
			PlanOnly: true,
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PostApplyPostRefresh: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionNoop),
				},
			},
		},
	}
}

// TestAccBucketBackupExportResourceDriftRemoval covers the branch Read takes when Capella no
// longer has the export: CheckResourceNotFoundError recognises the 404 and the resource is
// dropped from state instead of the refresh failing.
//
// That is the same branch a practitioner reaches roughly seven days after creating an export,
// when Capella drops the record and the next plan proposes a fresh export. The timeline cannot
// be reproduced and there is no endpoint to delete an export early, so importing an ID that was
// never real is the only way to execute the code path.
//
// It is worth executing. The 404 is only recognised when the response body carries a non-zero
// "code": without one, ExecuteWithRetry returns a plain error rather than an *api.Error,
// CheckResourceNotFoundError answers false, and every refresh after the record expires would
// hard-fail rather than removing the resource.
//
// The IDs are the pair TestAccDatasourceBucketBackupExportUnknownExport already establishes
// return 404/5019 together.
func TestAccBucketBackupExportResourceDriftRemoval(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_drift_")
	resourceReference := "couchbase-capella_bucket_backup_export." + resourceName

	const (
		missingBackupId = "00000000-0000-0000-0000-000000000000"
		missingExportId = "11111111-1111-1111-1111-111111111111"
	)

	importId := fmt.Sprintf(
		"id=%s,backup_id=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
		missingExportId, missingBackupId, globalBucketId, globalClusterId, globalProjectId, globalOrgId,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId, missingBackupId,
				),
				ResourceName:  resourceReference,
				ImportState:   true,
				ImportStateId: importId,
				// Read removes the resource, so Terraform reports that the import produced
				// nothing rather than surfacing the 404 itself. Seeing the provider's own
				// "Error Reading Capella Bucket Backup Export" here instead would mean the
				// 404 was not recognised and drift is no longer handled.
				ExpectError: regexp.MustCompile(`(?s)` + terraformDiagnosticPattern("Cannot import non-existent remote object")),
			},
		},
	})
}

// TestAccBucketBackupExportResourceDuplicateCycle asserts that a second export for a cycle that
// already has one is rejected with 409. depends_on forces the two creates to be ordered, so the
// second reliably loses rather than racing the first.
func TestAccBucketBackupExportResourceDuplicateCycle(t *testing.T) {
	bucketResourceName := randomStringWithPrefix("tf_acc_export_dup_bucket_")
	backupResourceName := randomStringWithPrefix("tf_acc_export_dup_backup_")
	firstName := randomStringWithPrefix("tf_acc_export_dup_first_")
	secondName := randomStringWithPrefix("tf_acc_export_dup_second_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportDuplicateConfig(bucketResourceName, backupResourceName, firstName, secondName),
				// 14061 is "already has a pending export", 14062 "already has a complete export".
				// Which one comes back depends on how far the first export got.
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":1406[12]`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceUnknownBackup asserts that a backup that does not exist is
// reported as 404/5017 rather than being silently accepted.
func TestAccBucketBackupExportResourceUnknownBackup(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_unknown_backup_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"00000000-0000-0000-0000-000000000000",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":5017`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceUnknownCluster asserts the cluster in the path is checked before
// the backup, so a bad cluster is 404/4025 and not a backup-not-found.
func TestAccBucketBackupExportResourceUnknownCluster(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_unknown_cluster_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId,
					"00000000-0000-0000-0000-000000000000", globalBucketId,
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":4025`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceWrongProject asserts that a project which does not own the
// cluster is rejected with 422/4031, rather than leaking whether the backup exists.
func TestAccBucketBackupExportResourceWrongProject(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_wrong_project_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, "00000000-0000-0000-0000-000000000000",
					globalClusterId, globalBucketId,
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Error creating backup export job.*"code":4031`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceMalformedBackupId asserts the schema validator rejects a
// non-UUID backup_id at plan time, before any API call is made.
func TestAccBucketBackupExportResourceMalformedBackupId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_bad_backup_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, globalBucketId,
					"not-a-uuid",
				),
				ExpectError: regexp.MustCompile(`(?s)Attribute backup_id must be a valid UUID`),
			},
		},
	})
}

// TestAccBucketBackupExportResourceEmptyBucketId asserts bucket_id is rejected when empty. It is
// not a UUID - it is the base64 encoding of the bucket name - so it only carries a length
// validator.
func TestAccBucketBackupExportResourceEmptyBucketId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_export_empty_bucket_id_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketBackupExportResourceConfigWithIDs(
					resourceName, globalOrgId, globalProjectId, globalClusterId, "",
					"11111111-1111-1111-1111-111111111111",
				),
				ExpectError: regexp.MustCompile(`(?s)Attribute bucket_id string length must be at least 1, got: 0`),
			},
		},
	})
}

// testAccBucketBackupExportResourceConfig declares a fresh bucket, a backup of it and an export of
// that backup. Terraform orders the three by their references, and the backup resource waits for
// the backup to reach a final state, so the export is only created once the backup is ready.
func testAccBucketBackupExportResourceConfig(bucketName, backupName, exportName string) string {
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
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, exportName)
}

// testAccBucketBackupExportDuplicateConfig exports the same backup twice. depends_on orders the
// second create after the first so the conflict is deterministic.
func testAccBucketBackupExportDuplicateConfig(bucketName, backupName, firstName, secondName string) string {
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

resource "couchbase-capella_bucket_backup_export" "%[8]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  bucket_id       = couchbase-capella_bucket.%[5]s.id
  backup_id       = couchbase-capella_backup.%[6]s.id
  depends_on      = [couchbase-capella_bucket_backup_export.%[7]s]
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, firstName, secondName)
}

// bucketBackupExportInputs holds the export resource's five required inputs as raw HCL
// expressions, so a step can change exactly one of them and leave the rest pointing at the
// per-run fixtures.
type bucketBackupExportInputs struct {
	organizationId string
	projectId      string
	clusterId      string
	bucketId       string
	backupId       string
}

// defaultBucketBackupExportInputs returns what the happy-path configuration uses: quoted
// literals for the three Capella IDs, and references to the bucket and backup created by the
// same configuration.
func defaultBucketBackupExportInputs(bucketName, backupName string) bucketBackupExportInputs {
	return bucketBackupExportInputs{
		organizationId: strconv.Quote(globalOrgId),
		projectId:      strconv.Quote(globalProjectId),
		clusterId:      strconv.Quote(globalClusterId),
		bucketId:       "couchbase-capella_bucket." + bucketName + ".id",
		backupId:       "couchbase-capella_backup." + backupName + ".id",
	}
}

// testAccBucketBackupExportResourceConfigInputs builds the same bucket / backup / export trio
// as testAccBucketBackupExportResourceConfig, but takes the export's inputs as expressions so
// one of them can be varied. The bucket and backup always use the real IDs, so a step that
// points the export elsewhere still destroys its own fixtures cleanly.
func testAccBucketBackupExportResourceConfigInputs(bucketName, backupName, exportName string, in bucketBackupExportInputs) string {
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
  organization_id = %[8]s
  project_id      = %[9]s
  cluster_id      = %[10]s
  bucket_id       = %[11]s
  backup_id       = %[12]s
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, bucketName, backupName, exportName,
		in.organizationId, in.projectId, in.clusterId, in.bucketId, in.backupId)
}

// testAccBucketBackupExportResourceConfigWithIDs is used by the invalid-input tests, which fail
// before an export is created and so need no bucket or backup of their own.
func testAccBucketBackupExportResourceConfigWithIDs(resourceName, organizationId, projectId, clusterId, bucketId, backupId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket_backup_export" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  bucket_id       = "%[6]s"
  backup_id       = "%[7]s"
}
`, globalProviderBlock, resourceName, organizationId, projectId, clusterId, bucketId, backupId)
}

func generateBucketBackupExportImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		return fmt.Sprintf(
			"id=%s,backup_id=%s,bucket_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["backup_id"], rawState["bucket_id"],
			rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

// retrieveBucketBackupExportFromServer fetches one export directly from the V4 API, bypassing
// the provider. It returns the response so callers that need more than existence - polling for
// completion, for instance - do not have to repeat the request.
func retrieveBucketBackupExportFromServer(data *providerschema.Data, organizationId, projectId, clusterId, bucketId, backupId, exportId string) (*backupapi.GetBucketBackupExportResponse, error) {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/backups/%s/exports/%s",
		data.HostURL, organizationId, projectId, clusterId, bucketId, backupId, exportId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	if err != nil {
		return nil, err
	}

	exportResp := backupapi.GetBucketBackupExportResponse{}
	if err := json.Unmarshal(response.Body, &exportResp); err != nil {
		return nil, err
	}
	if exportResp.Id != exportId {
		return nil, errors.ErrNotFound
	}
	return &exportResp, nil
}

// bucketBackupExportAttributes returns the flat state attributes of one export resource or data
// source, failing loudly when the address is not in state. Reporting that as an error matters:
// a check that silently read a nil map would compare empty strings and pass.
func bucketBackupExportAttributes(s *terraform.State, resourceReference string) (map[string]string, error) {
	for _, m := range s.Modules {
		if v, ok := m.Resources[resourceReference]; ok && v.Primary != nil {
			return v.Primary.Attributes, nil
		}
	}
	return nil, fmt.Errorf("%s not found in state", resourceReference)
}

func testAccExistsBucketBackupExportResource(t *testing.T, resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rawState, err := bucketBackupExportAttributes(s, resourceReference)
		if err != nil {
			return err
		}
		data := newTestClient(t)
		_, err = retrieveBucketBackupExportFromServer(
			data,
			rawState["organization_id"], rawState["project_id"], rawState["cluster_id"],
			rawState["bucket_id"], rawState["backup_id"], rawState["id"],
		)
		return err
	}
}
