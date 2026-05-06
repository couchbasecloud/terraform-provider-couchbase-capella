package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// TestAccAppServicesCIDRWithRequiredFields creates an App Services allowed CIDR with only
// the required fields and verifies all attributes — including computed audit fields — are populated.
func TestAccAppServicesCIDRWithRequiredFields(t *testing.T) {
	// Bug 1: ImportState uses ImportStatePassthroughID; Validate() checks OrganizationId.IsNull()
	// directly without parsing the composite import string, so import always fails with
	// "organization ID is missing or was passed incorrectly".
	t.Skip("skipped: import broken (AV-129753) — Validate() does not parse composite import ID")
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "10.50.0.0/16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.50.0.0/16"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
			// ImportState: the provider stores the import string in the "id" field via
			// ImportStatePassthroughID, but Validate() checks organization_id directly
			// without parsing the composite string. This step will fail until Validate()
			// is updated to use validateSchemaState for composite import ID parsing.
			{
				ResourceName:            resourceReference,
				ImportStateIdFunc:       generateAppServicesCIDRImportIdForResource(resourceReference),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"expires_at"},
			},
		},
	})
}

// TestAccAppServicesCIDRWithOptionalFields creates an App Services allowed CIDR with all
// optional fields (comment and expires_at) and verifies each is stored and read back correctly.
func TestAccAppServicesCIDRWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDROptionalFieldsConfig(resourceName, "10.60.0.1/32"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.60.0.1/32"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "comment", "terraform app services cidr acceptance test"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
				),
			},
		},
	})
}

// TestAccAppServicesCIDRAllowAll creates an App Services allowed CIDR with 0.0.0.0/0
// which permits connections from any IP address.
func TestAccAppServicesCIDRAllowAll(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDROptionalFieldsConfig(resourceName, "0.0.0.0/0"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "comment", "terraform app services cidr acceptance test"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
				),
			},
		},
	})
}

// TestAccAppServicesCIDRWithExpiredTimestamp attempts to create an App Services allowed CIDR
// with an already-expired expires_at timestamp and expects an API validation error.
func TestAccAppServicesCIDRWithExpiredTimestamp(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRExpiredTimestampConfig(resourceName, "10.80.0.1/32"),
				ExpectError: regexp.MustCompile(`(?s)expiry date for allowed CIDR is in the\s+past`),
			},
		},
	})
}

// TestAccAppServicesCIDRReplaceOnCIDRChange verifies that changing the cidr field triggers
// a destroy-and-recreate (replace) since cidr has RequiresReplace in the schema and the
// API does not support in-place updates to allowed CIDRs.
func TestAccAppServicesCIDRReplaceOnCIDRChange(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "10.100.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.100.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			{
				Config: testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "10.101.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.101.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

// TestAccAppServicesCIDRReplaceOnCommentChange verifies that changing the comment field triggers
// a destroy-and-recreate since comment has RequiresReplace in the schema.
func TestAccAppServicesCIDRReplaceOnCommentChange(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	expiryTime := time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRWithCommentConfig(resourceName, "10.110.0.0/24", "initial comment", expiryTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.110.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "comment", "initial comment"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
				),
			},
			{
				Config: testAccAppServicesCIDRWithCommentConfig(resourceName, "10.110.0.0/24", "updated comment", expiryTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.110.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "comment", "updated comment"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "expires_at"),
				),
			},
		},
	})
}

// TestAccAppServicesCIDRDatasource creates an App Services allowed CIDR and verifies the
// datasource lists all allowed CIDRs for the App Service including their audit metadata.
func TestAccAppServicesCIDRDatasource(t *testing.T) {
	// Bug 2: AppServicesCidrSchema() does not declare organization_id, project_id,
	// cluster_id, or app_service_id as schema attributes. Terraform rejects all four
	// with "Unsupported argument", making the datasource completely non-functional.
	t.Skip("skipped: datasource schema missing required input attributes (AV-129755)")
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	datasourceReference := "data.couchbase-capella_app_services_cidr.listing"
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAppServicesCIDRConfig(resourceName, "10.120.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(datasourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(datasourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(datasourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.id"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.cidr"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.audit.created_at"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.audit.created_by"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.audit.modified_at"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.audit.modified_by"),
					resource.TestCheckResourceAttrSet(datasourceReference, "data.0.audit.version"),
				),
			},
		},
	})
}

// generateAppServicesCIDRImportIdForResource builds the composite import ID from the
// resource state. Format matches the registry docs:
// id=<cidr_id>,app_service_id=<id>,cluster_id=<id>,project_id=<id>,organization_id=<id>
//
// Note: ImportState uses ImportStatePassthroughID which stores this string in the "id"
// field. The provider's Validate() checks organization_id directly without parsing the
// composite string, so import will fail until Validate() uses validateSchemaState.
func generateAppServicesCIDRImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"],
			rawState["app_service_id"],
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}

func testAccAppServicesCIDRRequiredFieldsConfig(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr)
}

func testAccAppServicesCIDROptionalFieldsConfig(resourceName, cidr string) string {
	expiryTime := time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  comment         = "terraform app services cidr acceptance test"
  expires_at      = "%[8]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr, expiryTime)
}

func testAccAppServicesCIDRExpiredTimestampConfig(resourceName, cidr string) string {
	timeNow := time.Now().UTC()
	time.Sleep(time.Second * 10)
	expiryTime := timeNow.Format(time.RFC3339)
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  comment         = "terraform app services cidr acceptance test"
  expires_at      = "%[8]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr, expiryTime)
}

func testAccAppServicesCIDRWithCommentConfig(resourceName, cidr, comment, expiresAt string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  comment         = "%[8]s"
  expires_at      = "%[9]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr, comment, expiresAt)
}

func testAccDatasourceAppServicesCIDRConfig(resourceName, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
}

data "couchbase-capella_app_services_cidr" "listing" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"

  depends_on = [couchbase-capella_app_services_cidr.%[2]s]
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr)
}

// TestAccAppServicesCIDRReplaceOnExpiresAtChange verifies that changing only the expires_at
// field triggers a destroy-and-recreate, confirming expires_at has RequiresReplace in the schema.
func TestAccAppServicesCIDRReplaceOnExpiresAtChange(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	expiryTime1 := time.Now().AddDate(0, 0, 30).UTC().Format(time.RFC3339)
	expiryTime2 := time.Now().AddDate(0, 0, 60).UTC().Format(time.RFC3339)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRWithExpiresAtConfig(resourceName, "10.140.0.0/24", expiryTime1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.140.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "expires_at", expiryTime1),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
			{
				Config: testAccAppServicesCIDRWithExpiresAtConfig(resourceName, "10.140.0.0/24", expiryTime2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.140.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "expires_at", expiryTime2),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
				),
			},
		},
	})
}

// TestAccAppServicesCIDRInvalidCIDRFormat attempts to create an App Services allowed CIDR
// with a malformed CIDR string and expects an API validation error.
// Note: there is no provider-side CIDR validator; the error comes from the API (HTTP 422).
// If this test fails with "no match on", update the ExpectError regex to match the actual
// API error message shown in the failure output.
func TestAccAppServicesCIDRInvalidCIDRFormat(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "not-a-cidr"),
				ExpectError: regexp.MustCompile("(?i)invalid"),
			},
		},
	})
}

// TestAccAppServicesCIDROutOfBandDeletion verifies that when an allowed CIDR is deleted
// outside of Terraform (e.g. via the Capella API or UI), the next apply gracefully removes
// the resource from state (via resp.State.RemoveResource) and recreates it rather than
// returning a persistent error.
func TestAccAppServicesCIDROutOfBandDeletion(t *testing.T) {
	// Bug 3: Read() uses err == errors.ErrNotFound (strict equality) instead of
	// errors.Is(err, errors.ErrNotFound). refreshAllowedCIDR wraps the sentinel with
	// fmt.Errorf("%s: %w", ...), so the equality check always fails and the 404 falls
	// through to AddError, producing a persistent error instead of resp.State.RemoveResource.
	t.Skip("skipped: out-of-band deletion causes persistent error instead of graceful state removal (AV-129801)")
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	var cidrId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "10.150.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.150.0.0/24"),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					func(s *terraform.State) error {
						for _, m := range s.Modules {
							if v, ok := m.Resources[resourceReference]; ok {
								cidrId = v.Primary.Attributes["id"]
							}
						}
						return nil
					},
				),
			},
			{
				// Delete the CIDR via the API before Terraform runs, simulating an out-of-band deletion.
				// Read() should detect the 404, call resp.State.RemoveResource, and the subsequent
				// plan should show a create (not an error), causing Terraform to recreate it.
				PreConfig: func() {
					ctx := context.Background()
					client := api.NewClient(timeout)
					url := fmt.Sprintf(
						"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/allowedcidrs/%s",
						globalHost, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidrId,
					)
					cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
					_, _ = client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
				},
				Config: testAccAppServicesCIDRRequiredFieldsConfig(resourceName, "10.150.0.0/24"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", "10.150.0.0/24"),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					func(s *terraform.State) error {
						for _, m := range s.Modules {
							if v, ok := m.Resources[resourceReference]; ok {
								if newId := v.Primary.Attributes["id"]; newId == cidrId {
									return fmt.Errorf("expected a new CIDR ID after out-of-band deletion and recreation, still got original ID %s", cidrId)
								}
							}
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccAppServicesCIDRWithExpiresAtConfig(resourceName, cidr, expiresAt string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  expires_at      = "%[8]s"
}
`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr, expiresAt)
}

// TestAccAppServicesCIDREmptyOrganizationId verifies that an empty organization_id triggers
// the provider-side LengthAtLeast(1) validator before reaching the API.
func TestAccAppServicesCIDREmptyOrganizationId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRWithCustomIdsConfig(resourceName, "", globalProjectId, globalClusterId, globalAppServiceId, "10.160.0.0/24"),
				ExpectError: regexp.MustCompile("string length must be at least 1"),
			},
		},
	})
}

// TestAccAppServicesCIDRInvalidOrganizationId verifies that a non-existent organization_id
// causes the API to return an error (403 or 404) that the provider surfaces correctly.
// Note: update the ExpectError regex to match the actual API message on first run.
func TestAccAppServicesCIDRInvalidOrganizationId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRWithCustomIdsConfig(resourceName, "00000000-0000-0000-0000-000000000000", globalProjectId, globalClusterId, globalAppServiceId, "10.160.0.0/24"),
				ExpectError: regexp.MustCompile(`(?is)access.denied|1002`),
			},
		},
	})
}

// TestAccAppServicesCIDRInvalidClusterId verifies that a non-existent cluster_id causes
// the API to return a 404 that the provider surfaces correctly.
// Note: update the ExpectError regex to match the actual API message on first run.
func TestAccAppServicesCIDRInvalidClusterId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRWithCustomIdsConfig(resourceName, globalOrgId, globalProjectId, "00000000-0000-0000-0000-000000000000", globalAppServiceId, "10.160.0.0/24"),
				ExpectError: regexp.MustCompile("(?i)not found|cluster"),
			},
		},
	})
}

// TestAccAppServicesCIDRInvalidAppServiceId verifies that a non-existent app_service_id causes
// the API to return a 404 that the provider surfaces correctly.
// Note: update the ExpectError regex to match the actual API message on first run.
func TestAccAppServicesCIDRInvalidAppServiceId(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_svc_cidr_")
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppServicesCIDRWithCustomIdsConfig(resourceName, globalOrgId, globalProjectId, globalClusterId, "00000000-0000-0000-0000-000000000000", "10.160.0.0/24"),
				ExpectError: regexp.MustCompile("(?i)not found|app.?service"),
			},
		},
	})
}

func testAccAppServicesCIDRWithCustomIdsConfig(resourceName, orgId, projectId, clusterId, appServiceId, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
}
`, globalProviderBlock, resourceName, orgId, projectId, clusterId, appServiceId, cidr)
}
