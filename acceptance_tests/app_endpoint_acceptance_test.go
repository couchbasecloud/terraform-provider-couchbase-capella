package acceptance_tests

import (
	"fmt"
	re "regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpoint(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_app_endpoint_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, bucket, "syncFnXattr", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "syncFnXattr"),
				),
			},
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, bucket, "new_xattr", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "new_xattr"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointImportId(resourceReference),
				ImportState:       true,
			},
		},
	})
}

func TestAccAppEndpointInexistentCollection(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("bkt_")
	cfg := fmt.Sprintf(`
	%[1]s
	
	resource "couchbase-capella_bucket" "%[2]s_bucket" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		name           = "%[7]s"
	}
	
	resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		app_service_id  = "%[6]s"
		bucket          = "%[7]s"
		name            = "%[8]s"
		scopes = {
			"_default" = {
			  collections = {
				"INVALID_COLLLECTION" = {}
			  }
			}
		}
		depends_on = [couchbase-capella_bucket.%[2]s_bucket]
	}`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucket,
		epName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      cfg,
				ExpectError: re.MustCompile("Collection Not Found"),
			},
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName, endpointName, bucketName, userXattr string, deltaSync bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name           = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"
	user_xattr_key  = "%[9]s"
	delta_sync_enabled = %[10]t
	cors = {
		origin = ["*"]
	}
	oidc = [
		{
			issuer   = "https://accounts.google.com"
			client_id = "example_client_id"
		}
	]
	
	scopes = {
		"_default" = {
		  collections = {
			"_default" = {}
		  }
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
		userXattr,
		deltaSync,
	)
}

// TestAccAppEndpointNoCors verifies that creating an app endpoint without a
// cors block does not produce perpetual state drift on subsequent plans.
func TestAccAppEndpointNoCors(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_ep_nocors_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_ep_nocors_")
	bucket := randomStringWithPrefix("tf_acc_ep_nocors_bkt_")

	cfg := testAccAppEndpointNoCorsConfig(resourceName, epName, bucket)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckNoResourceAttr(resourceReference, "cors"),
				),
			},
			// Re-apply the same config; expect no changes (no perpetual drift).
			{
				Config:             cfg,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccAppEndpointNoCorsConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name           = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"
	scopes = {
		"_default" = {
		  collections = {
			"_default" = {}
		  }
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

func generateAppEndpointImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		// Import uses the endpoint name
		return fmt.Sprintf("app_endpoint_name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["name"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}

// ── S4: cors.disabled=false without origin — API 422 "App Endpoint CORS Origin is empty" ──
// Provider schema marks origin as Optional but the API requires it when a cors
// block is present and disabled=false.
func TestAccAppEndpointCorsDisabledFalseNoOrigin(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointCorsDisabledFalseResourceConfig(resourceName, epName, bucket),
				ExpectError: re.MustCompile("CORS Origin is empty"), 
				// Remove when the bug is fixed, as the config will be valid and should apply successfully. https://jira.issues.couchbase.com/browse/AV-128217
			},
		},
	})
}

// ── S6: Full cors config (all fields) — happy path (also covers I2 import) ───
func TestAccAppEndpointCorsFullConfig(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.disabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "cors.origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "*"),
					resource.TestCheckResourceAttr(resourceReference, "cors.login_origin.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "cors.headers.#", "2"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppEndpointImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"state"},
			},
		},
	})
}

// ── S19: Multiple OIDC providers — API 500 Internal Server Error ──────────────
// Server-side "OIDC discovery config validation failed" when ≥2 providers are
// supplied. Single-provider creation (S17/S18) works correctly.
func TestAccAppEndpointMultipleOIDC(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointMultipleOIDCResourceConfig(resourceName, epName, bucket),
				ExpectError: re.MustCompile("500"),
				// Remove when the bug is fixed https://jira.issues.couchbase.com/browse/AV-128222 
			},
		},
	})
}

// ── S20: cors with specific (non-wildcard) origins — happy path ───────────────
func TestAccAppEndpointCorsSpecificOrigins(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.origin.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "https://app.example.com"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "https://admin.example.com"),
					resource.TestCheckResourceAttr(resourceReference, "cors.login_origin.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "600"),
				),
			},
		},
	})
}

// ── S21: cors.max_age=0 — silent state drift (omitempty drops 0 from CREATE request) ──
// Bug: apply succeeds without error, but state shows max_age=0 while the API
// stores its default value (the API omits maxAge from GET responses when at
// default). Users cannot set max_age=0 — the value is always silently replaced.
func TestAccAppEndpointCorsMaxAgeZeroSilentDrift(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "0"),
					// https://jira.issues.couchbase.com/browse/AV-128218
				),
			},
		},
	})
}

// ── S18: OIDC with all optional fields — happy path ───────────────────────────
func TestAccAppEndpointOIDCFullFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointOIDCAllFieldsResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.client_id", "example-client-id"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.register", "true"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.user_prefix", "google_"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.username_claim", "email"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.roles_claim", "roles"),
				),
			},
		},
	})
}

// ── S22: OIDC with discovery_url — happy path ─────────────────────────────────
func TestAccAppEndpointOIDCDiscoveryURL(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointOIDCDiscoveryURLResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.discovery_url", "https://accounts.google.com/.well-known/openid-configuration"),
				),
			},
		},
	})
}

// ── U1: Expand cors from minimal (origin only) to full (all fields) — happy path ──
func TestAccAppEndpointUpdateCorsExpand(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.origin.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "*"),
				),
			},
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.disabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "3600"),
					resource.TestCheckResourceAttr(resourceReference, "cors.login_origin.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "cors.headers.#", "2"),
				),
			},
		},
	})
}

// ── U2: Remove cors block via update — API 422 "App Endpoint CORS Origin is empty" ──
// Once cors is set, the API rejects any PUT that omits the cors body entirely.
// cors is effectively write-once via Terraform.
func TestAccAppEndpointUpdateRemoveCors(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, bucket),
			},
			{
				Config:      testAccAppEndpointNoCorsResourceConfig(resourceName, epName, bucket),
				ExpectError: re.MustCompile("CORS Origin is empty"),
				// Remove once the bug is fixed https://jira.issues.couchbase.com/browse/AV-128229, https://jira.issues.couchbase.com/browse/AV-128217
			},
		},
	})
}

// ── U3: cors.disabled false → true — API 409 "CORS cannot be disabled, config not empty" ──
// The API rejects disabling CORS when other cors fields (origin etc.) are also set.
func TestAccAppEndpointUpdateCorsDisableToggle(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, bucket),
			},
			{
				Config:      testAccAppEndpointCorsDisabledTrueResourceConfig(resourceName, epName, bucket),
				ExpectError: re.MustCompile("CORS cannot be disabled"),
				// Remove once the bug is fixed https://jira.issues.couchbase.com/browse/AV-128229
			},
		},
	})
}

// ── U5: cors origin wildcard → specific URLs — happy path ────────────────────
func TestAccAppEndpointUpdateCorsOriginWildcardToSpecific(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "*"),
				),
			},
			{
				Config: testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.origin.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "https://app.example.com"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "https://admin.example.com"),
				),
			},
		},
	})
}

// ── U7: Add OIDC block to existing endpoint — happy path ─────────────────────
func TestAccAppEndpointUpdateAddOIDC(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, bucket),
			},
			{
				Config: testAccAppEndpointWithOIDCResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.client_id", "example-client-id"),
					resource.TestCheckResourceAttrSet(resourceReference, "oidc.0.provider_id"),
				),
			},
		},
	})
}

// ── U8: Remove OIDC block via update — "Provider produced inconsistent result after apply" ──
// The provider omits oidc from the PUT request when removed from config, but the
// API does not remove the OIDC provider. refreshAppEndpoint re-populates state.Oidc
// from the GET response, and Terraform detects the plan/state mismatch.
func TestAccAppEndpointUpdateRemoveOIDC(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointWithOIDCResourceConfig(resourceName, epName, bucket),
			},
			{
				Config:      testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, bucket),
				ExpectError: re.MustCompile("Provider produced inconsistent result after apply"),
				// Remove once the bug is fixed https://jira.issues.couchbase.com/browse/AV-128167
			},
		},
	})
}

// ── U9: Update access_control_function body — happy path ─────────────────────
func TestAccAppEndpointUpdateACF(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	initialACF := "function(doc, oldDoc, meta) { channel(doc.channels); }"
	updatedACF := "function(doc, oldDoc, meta) { channel(doc.type); }"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointACFResourceConfig(resourceName, epName, bucket, initialACF),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
				),
			},
			{
				Config: testAccAppEndpointACFResourceConfig(resourceName, epName, bucket, updatedACF),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
				),
			},
		},
	})
}

// ── U10: cors.max_age non-zero → zero — "Provider produced inconsistent result after apply" ──
// Root cause: MaxAge int64 `json:"maxAge,omitempty"` — zero is omitted from PUT
// request. API retains old value (3600). refreshAppEndpoint reads 3600 back.
// Terraform detects plan=0 vs state=3600 and raises inconsistency error.
func TestAccAppEndpointUpdateCorsMaxAgeToZero(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, bucket, 3600),
			},
			{
				Config:      testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, bucket, 0),
				ExpectError: re.MustCompile("Provider produced inconsistent result after apply"),
				// https://jira.issues.couchbase.com/browse/AV-128218
			},
		},
	})
}

// ── U11: cors.max_age zero → non-zero (drift recovery) — happy path ───────────
// After Bug S21 silently stores 0 in state (API has its default), updating to a
// non-zero value recovers the drift. The non-zero value is not omitted by
// omitempty and is applied correctly by the API.
func TestAccAppEndpointUpdateCorsMaxAgeFromZero(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_bucket_")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, epName, bucket),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "0"),
				),
			},
			{
				Config: testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, bucket, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "3600"),
				),
			},
		},
	})
}

// ── Scenario B: Four parallel creates — exercises missing resource state timing bug ──
// Terraform's default parallelism issues all four POST+GET cycles concurrently,
// maximising the chance of hitting the window where an immediate GET lands before
// the endpoint is stable. If any endpoint returns empty state while still visible
// in Capella, a second apply will see HTTP 412 (already exists) on re-create.
// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

// testAccAppEndpointNoCorsResourceConfig creates an endpoint with no cors block.
// Used by: TestAccAppEndpointNoCors (S1/S2), TestAccAppEndpointUpdateRemoveCors (U2 phase 2).
func testAccAppEndpointNoCorsResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsDisabledFalseResourceConfig creates an endpoint with
// cors { disabled=false } and no origin field. Used by: TestAccAppEndpointCorsDisabledFalseNoOrigin (S4).
func testAccAppEndpointCorsDisabledFalseResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		disabled = false
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsMaxAgeZeroResourceConfig creates an endpoint with
// cors.max_age=0. Used by: TestAccAppEndpointCorsMaxAgeZeroSilentDrift (S21),
// TestAccAppEndpointUpdateCorsMaxAgeFromZero (U11 phase 1).
func testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		disabled = false
		origin   = ["*"]
		max_age  = 0
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointMultipleOIDCResourceConfig creates an endpoint with two OIDC
// providers. Used by: TestAccAppEndpointMultipleOIDC (S19).
func testAccAppEndpointMultipleOIDCResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	oidc = [
		{
			issuer    = "https://accounts.google.com"
			client_id = "google-client-id"
			register  = true
		},
		{
			issuer    = "https://login.microsoftonline.com/common/v2.0"
			client_id = "azure-client-id"
		}
	]

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsAllFieldsResourceConfig creates an endpoint with all
// cors fields set. Used by: TestAccAppEndpointCorsFullConfig (S6),
// TestAccAppEndpointUpdateCorsExpand (U1 phase 2), TestAccAppEndpointUpdateRemoveCors (U2 phase 1).
func testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		disabled     = false
		origin       = ["*"]
		login_origin = ["*"]
		headers      = ["Authorization", "Content-Type"]
		max_age      = 3600
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsSpecificOriginsResourceConfig creates an endpoint with
// specific (non-wildcard) cors origins. Used by: TestAccAppEndpointCorsSpecificOrigins (S20),
// TestAccAppEndpointUpdateCorsOriginWildcardToSpecific (U5 phase 2).
func testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		disabled     = false
		origin       = ["https://app.example.com", "https://admin.example.com"]
		login_origin = ["https://app.example.com"]
		headers      = ["Authorization", "Content-Type", "X-Custom-Header"]
		max_age      = 600
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointOIDCAllFieldsResourceConfig creates an endpoint with an OIDC
// provider using all optional fields. Used by: TestAccAppEndpointOIDCFullFields (S18).
func testAccAppEndpointOIDCAllFieldsResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	oidc = [
		{
			issuer         = "https://accounts.google.com"
			client_id      = "example-client-id"
			register       = true
			user_prefix    = "google_"
			username_claim = "email"
			roles_claim    = "roles"
		}
	]

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointOIDCDiscoveryURLResourceConfig creates an endpoint with an OIDC
// provider that includes a discovery_url. Used by: TestAccAppEndpointOIDCDiscoveryURL (S22).
func testAccAppEndpointOIDCDiscoveryURLResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	oidc = [
		{
			issuer        = "https://accounts.google.com"
			client_id     = "example-client-id"
			discovery_url = "https://accounts.google.com/.well-known/openid-configuration"
		}
	]

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsOriginOnlyResourceConfig creates an endpoint with a
// minimal cors block (origin only). Used by: TestAccAppEndpointUpdateCorsExpand (U1 phase 1),
// TestAccAppEndpointUpdateCorsDisableToggle (U3 phase 1),
// TestAccAppEndpointUpdateCorsOriginWildcardToSpecific (U5 phase 1),
// TestAccAppEndpointUpdateAddOIDC (U7 phase 1), TestAccAppEndpointUpdateRemoveOIDC (U8 phase 2).
func testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsDisabledTrueResourceConfig creates an endpoint with
// cors.disabled=true. Used by: TestAccAppEndpointUpdateCorsDisableToggle (U3 phase 2).
func testAccAppEndpointCorsDisabledTrueResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		disabled = true
		origin   = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointCorsMaxAgeResourceConfig creates an endpoint with a
// parameterised cors.max_age. Used by: TestAccAppEndpointUpdateCorsMaxAgeToZero (U10),
// TestAccAppEndpointUpdateCorsMaxAgeFromZero (U11 phase 2).
func testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, endpointName, bucketName string, maxAge int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin  = ["*"]
		max_age = %[9]d
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
		maxAge,
	)
}

// testAccAppEndpointWithOIDCResourceConfig creates an endpoint with cors and a
// minimal OIDC provider. Used by: TestAccAppEndpointUpdateAddOIDC (U7 phase 2),
// TestAccAppEndpointUpdateRemoveOIDC (U8 phase 1).
func testAccAppEndpointWithOIDCResourceConfig(resourceName, endpointName, bucketName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	oidc = [
		{
			issuer    = "https://accounts.google.com"
			client_id = "example-client-id"
		}
	]

	scopes = {
		"_default" = {
			collections = {
				"_default" = {}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
	)
}

// testAccAppEndpointACFResourceConfig creates an endpoint with a parameterised
// access_control_function body. Used by: TestAccAppEndpointUpdateACF (U9).
func testAccAppEndpointACFResourceConfig(resourceName, endpointName, bucketName, acfBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"_default" = {
					access_control_function = "%[9]s"
				}
			}
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
		acfBody,
	)
}
