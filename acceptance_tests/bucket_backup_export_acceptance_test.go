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
		// ImportStateVerify is the point: ImportState alone passes even if import returns nothing.
		{
			ResourceName:      resourceReference,
			ImportStateIdFunc: generateBucketBackupExportImportIdForResource(resourceReference),
			ImportState:       true,
			ImportStateVerify: true,
			// These five move on their own between create and import; the IDs and created_at still match.
			ImportStateVerifyIgnore: []string{
				"status",
				"size_in_bytes",
				"sha256_checksum",
				"expiration",
				"backup_download_url",
			},
		},
	}

	// RequiresReplace on all five inputs is the only reason Update, which always errors, is unreachable.
	steps = append(steps, bucketBackupExportRequiresReplaceSteps(
		bucketResourceName, backupResourceName, resourceName, resourceReference,
	)...)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

// bucketBackupExportRequiresReplaceSteps: one plan-only step per input asserting Replace, then a restore step expecting an empty plan.
func bucketBackupExportRequiresReplaceSteps(bucketName, backupName, exportName, resourceReference string) []resource.TestStep {
	defaults := defaultBucketBackupExportInputs(bucketName, backupName)

	// Well formed but different from the fixtures, so this is a value change not a rejection.
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
				// PostApplyPostRefresh, not PreApply: the library rejects PreApply with PlanOnly.
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

// TestAccBucketBackupExportResourceDriftRemoval reaches the Read 404 branch that drops the resource from state.
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
				// Read removes the resource, so Terraform reports the import produced nothing, not the 404.
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

// bucketBackupExportInputs holds the five required inputs as HCL expressions so a step can vary one.
type bucketBackupExportInputs struct {
	organizationId string
	projectId      string
	clusterId      string
	bucketId       string
	backupId       string
}

// defaultBucketBackupExportInputs returns what the happy-path configuration uses.
func defaultBucketBackupExportInputs(bucketName, backupName string) bucketBackupExportInputs {
	return bucketBackupExportInputs{
		organizationId: strconv.Quote(globalOrgId),
		projectId:      strconv.Quote(globalProjectId),
		clusterId:      strconv.Quote(globalClusterId),
		bucketId:       "couchbase-capella_bucket." + bucketName + ".id",
		backupId:       "couchbase-capella_backup." + backupName + ".id",
	}
}

// testAccBucketBackupExportResourceConfigInputs is the same trio, with the export's inputs as expressions.
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

// retrieveBucketBackupExportFromServer fetches one export directly from the V4 API, bypassing the provider.
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

// bucketBackupExportAttributes returns one address's state attributes, erroring when it is absent.
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
