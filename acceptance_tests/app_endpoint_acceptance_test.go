package acceptance_tests

import (
	"fmt"
	re "regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpoint(t *testing.T) {
	ensureFixtureCollection(t, globalEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, globalEPCollectionName, "syncFnXattr", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "bucket", appEndpointBucketName),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "syncFnXattr"),
				),
			},
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, globalEPCollectionName, "new_xattr", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
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
	ensureAppEndpointTestEnvironment(t)

	resourceName := randomStringWithPrefix("tf_acc_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	cfg := fmt.Sprintf(`
	%[1]s

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
	}`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		appEndpointBucketName,
		epName)

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
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

func testAccAppEndpointResourceConfig(resourceName, endpointName, collectionName, userXattr string, deltaSync bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id    = "%[3]s"
	project_id         = "%[4]s"
	cluster_id         = "%[5]s"
	app_service_id     = "%[6]s"
	bucket             = "`+appEndpointBucketName+`"
	name               = "%[8]s"
	user_xattr_key     = "%[9]s"
	delta_sync_enabled = %[10]t
	cors = {
		origin = ["*"]
	}
	oidc = [
		{
			issuer    = "https://accounts.google.com"
			client_id = "example_client_id"
		}
	]

	scopes = {
		"_default" = {
		  collections = {
			"%[7]s" = {}
		  }
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
		userXattr,
		deltaSync,
	)
}

// TestAccAppEndpointNoCors verifies that creating an app endpoint without a
// cors block does not produce perpetual state drift on subsequent plans.
func TestAccAppEndpointNoCors(t *testing.T) {
	ensureFixtureCollection(t, globalNoCorsEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_ep_nocors_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_ep_nocors_")

	cfg := testAccAppEndpointNoCorsConfig(resourceName, epName, globalNoCorsEPCollectionName)

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
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

func testAccAppEndpointNoCorsConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"
	scopes = {
		"_default" = {
		  collections = {
			"%[7]s" = {}
		  }
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

func generateAppEndpointImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[resourceReference]; ok {
				rawState = v.Primary.Attributes
				break
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf("app_endpoint_name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["name"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}

func testAccAppEndpointComputedAttrs(resourceReference string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceReference, "state"),
		resource.TestCheckResourceAttrSet(resourceReference, "admin_url"),
		resource.TestCheckResourceAttrSet(resourceReference, "public_url"),
		resource.TestCheckResourceAttrSet(resourceReference, "metrics_url"),
	)
}

// ── AV-128217: cors.disabled=false without origin must fail provider validation ──
// The API rejects empty/missing origin when a cors block is present.
// Provider schema must reject this before issuing the API request.
func TestAccAppEndpointCorsDisabledFalseNoOrigin(t *testing.T) {
	ensureFixtureCollection(t, globalCorsDisabledFalseEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppEndpointCorsDisabledFalseResourceConfig(resourceName, epName, globalCorsDisabledFalseEPCollectionName),
				ExpectError: re.MustCompile(`(?s).*cors\.origin.*at least 1 value.*`),
			},
			{
				Config:      testAccAppEndpointCorsEmptyOriginResourceConfig(resourceName, epName, globalCorsDisabledFalseEPCollectionName),
				ExpectError: re.MustCompile(`(?s).*cors\.origin.*at least 1 value.*`),
			},
			{
				Config: testAccAppEndpointCorsDisabledTrueNoOriginResourceConfig(resourceName, epName, globalCorsDisabledFalseEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "cors.disabled", "true"),
				),
			},
		},
	})
}

// ── S6: Full cors config (all fields) — happy path (also covers I2 import) ───
func TestAccAppEndpointCorsFullConfig(t *testing.T) {
	ensureFixtureCollection(t, globalCorsFullEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, globalCorsFullEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
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
	t.Skip("AV-128222: multiple OIDC providers cause a 500 Internal Server Error; unskip once the bug is fixed")
	ensureFixtureCollection(t, globalMultipleOIDCEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointMultipleOIDCResourceConfig(resourceName, epName, globalMultipleOIDCEPCollectionName),
			},
		},
	})
}

// ── S20: cors with specific (non-wildcard) origins — happy path ───────────────
func TestAccAppEndpointCorsSpecificOrigins(t *testing.T) {
	ensureFixtureCollection(t, globalCorsSpecificEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, epName, globalCorsSpecificEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
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
	t.Skip("AV-128218: cors.max_age=0 is silently omitted by the API request model; unskip once zero values round-trip correctly")
	ensureFixtureCollection(t, globalCorsMaxAge0EPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, epName, globalCorsMaxAge0EPCollectionName),
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
	ensureFixtureCollection(t, globalOIDCFullEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointOIDCAllFieldsResourceConfig(resourceName, epName, globalOIDCFullEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.client_id", "example-client-id"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.register", "true"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.user_prefix", "google_"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.username_claim", "email"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.roles_claim", "roles"),
					resource.TestCheckResourceAttrSet(resourceReference, "oidc.0.provider_id"),
					resource.TestCheckResourceAttrSet(resourceReference, "oidc.0.is_default"),
				),
			},
		},
	})
}

// ── S22: OIDC with discovery_url — happy path ─────────────────────────────────
func TestAccAppEndpointOIDCDiscoveryURL(t *testing.T) {
	ensureFixtureCollection(t, globalOIDCDiscEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointOIDCDiscoveryURLResourceConfig(resourceName, epName, globalOIDCDiscEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "oidc.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "oidc.0.discovery_url", "https://accounts.google.com/.well-known/openid-configuration"),
					resource.TestCheckResourceAttrSet(resourceReference, "oidc.0.provider_id"),
				),
			},
		},
	})
}

// ── U1: Expand cors from minimal (origin only) to full (all fields) — happy path ──
func TestAccAppEndpointUpdateCorsExpand(t *testing.T) {
	ensureFixtureCollection(t, globalCorsExpandEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, globalCorsExpandEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "bucket", appEndpointBucketName),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "*"),
				),
			},
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, globalCorsExpandEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccAppEndpointComputedAttrs(resourceReference),
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
	t.Skip("AV-128229 / AV-128217: removing the cors block after it is set should succeed once the bugs are fixed")
	ensureFixtureCollection(t, globalRemoveCorsEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, epName, globalRemoveCorsEPCollectionName),
			},
			{
				Config: testAccAppEndpointNoCorsResourceConfig(resourceName, epName, globalRemoveCorsEPCollectionName),
			},
		},
	})
}

// ── U3: cors.disabled false → true — API 409 "CORS cannot be disabled, config not empty" ──
// The API rejects disabling CORS when other cors fields (origin etc.) are also set.
func TestAccAppEndpointUpdateCorsDisableToggle(t *testing.T) {
	t.Skip("AV-128229: toggling cors.disabled=true while other CORS fields are set should succeed once the bug is fixed")
	ensureFixtureCollection(t, globalCorsDisableToggleEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, globalCorsDisableToggleEPCollectionName),
			},
			{
				Config: testAccAppEndpointCorsDisabledTrueResourceConfig(resourceName, epName, globalCorsDisableToggleEPCollectionName),
			},
		},
	})
}

// ── U5: cors origin wildcard → specific URLs — happy path ────────────────────
func TestAccAppEndpointUpdateCorsOriginWildcardToSpecific(t *testing.T) {
	ensureFixtureCollection(t, globalCorsWildEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, globalCorsWildEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceReference, "cors.origin.*", "*"),
				),
			},
			{
				Config: testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, epName, globalCorsWildEPCollectionName),
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
	ensureFixtureCollection(t, globalAddOIDCEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, globalAddOIDCEPCollectionName),
			},
			{
				Config: testAccAppEndpointWithOIDCResourceConfig(resourceName, epName, globalAddOIDCEPCollectionName),
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
	t.Skip("AV-128167: removing the oidc block should succeed once the bug is fixed")
	ensureFixtureCollection(t, globalRemoveOIDCEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointWithOIDCResourceConfig(resourceName, epName, globalRemoveOIDCEPCollectionName),
			},
			{
				Config: testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, epName, globalRemoveOIDCEPCollectionName),
			},
		},
	})
}

// ── U9: Update access_control_function body — happy path ─────────────────────
func TestAccAppEndpointUpdateACF(t *testing.T) {
	ensureFixtureCollection(t, globalACFUpdateEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	initialACF := "function(doc, oldDoc, meta) { channel(doc.channels); }"
	updatedACF := "function(doc, oldDoc, meta) { channel(doc.type); }"

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointACFResourceConfig(resourceName, epName, globalACFUpdateEPCollectionName, initialACF),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
				),
			},
			{
				Config: testAccAppEndpointACFResourceConfig(resourceName, epName, globalACFUpdateEPCollectionName, updatedACF),
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
	t.Skip("AV-128218: updating cors.max_age to 0 currently produces inconsistent state; unskip once zero values round-trip correctly")
	ensureFixtureCollection(t, globalCorsMaxAgeZeroEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, globalCorsMaxAgeZeroEPCollectionName, 3600),
			},
			{
				Config:      testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, globalCorsMaxAgeZeroEPCollectionName, 0),
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
	ensureFixtureCollection(t, globalCorsMaxAgeFromZeroEPCollectionName)

	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")

	t.Parallel()
	defer acquireAppEndpointCRUDSlot()()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, epName, globalCorsMaxAgeFromZeroEPCollectionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "0"),
				),
			},
			{
				Config: testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, epName, globalCorsMaxAgeFromZeroEPCollectionName, 3600),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cors.max_age", "3600"),
				),
			},
		},
	})
}

// ── Scenario B: Four parallel creates — exercises missing resource state timing bug ──
// Terraform's default parallelism issues all four POST+GET cycles concurrently,
// maximising the chance of hitting the window where an endpoint returns empty state
// while still visible in Capella. If any endpoint returns empty state while still
// visible in Capella, a second apply will see HTTP 412 (already exists) on re-create.
// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

// testAccAppEndpointNoCorsResourceConfig creates an endpoint with no cors block.
// Used by: TestAccAppEndpointUpdateRemoveCors (U2 phase 2, skipped).
func testAccAppEndpointNoCorsResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsDisabledFalseResourceConfig creates an endpoint with
// cors { disabled=false } and no origin field. Used by: TestAccAppEndpointCorsDisabledFalseNoOrigin.
func testAccAppEndpointCorsDisabledFalseResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		disabled = false
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsDisabledTrueNoOriginResourceConfig creates an endpoint with
// cors { disabled=true } and no origin field. Used by: TestAccAppEndpointCorsDisabledFalseNoOrigin.
func testAccAppEndpointCorsDisabledTrueNoOriginResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		disabled = true
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsEmptyOriginResourceConfig creates an endpoint with
// cors { disabled=false, origin=[] }. Used by: TestAccAppEndpointCorsDisabledFalseNoOrigin.
func testAccAppEndpointCorsEmptyOriginResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		disabled = false
		origin   = []
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsMaxAgeZeroResourceConfig creates an endpoint with
// cors.max_age=0. Used by: TestAccAppEndpointCorsMaxAgeZeroSilentDrift (S21),
// TestAccAppEndpointUpdateCorsMaxAgeFromZero (U11 phase 1).
func testAccAppEndpointCorsMaxAgeZeroResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		disabled = false
		origin   = ["*"]
		max_age  = 0
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointMultipleOIDCResourceConfig creates an endpoint with two OIDC
// providers. Used by: TestAccAppEndpointMultipleOIDC (S19, skipped).
func testAccAppEndpointMultipleOIDCResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsAllFieldsResourceConfig creates an endpoint with all
// cors fields set. Used by: TestAccAppEndpointCorsFullConfig (S6),
// TestAccAppEndpointUpdateCorsExpand (U1 phase 2), TestAccAppEndpointUpdateRemoveCors (U2 phase 1, skipped).
func testAccAppEndpointCorsAllFieldsResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsSpecificOriginsResourceConfig creates an endpoint with
// specific (non-wildcard) cors origins. Used by: TestAccAppEndpointCorsSpecificOrigins (S20),
// TestAccAppEndpointUpdateCorsOriginWildcardToSpecific (U5 phase 2).
func testAccAppEndpointCorsSpecificOriginsResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointOIDCAllFieldsResourceConfig creates an endpoint with an OIDC
// provider using all optional fields. Used by: TestAccAppEndpointOIDCFullFields (S18).
func testAccAppEndpointOIDCAllFieldsResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointOIDCDiscoveryURLResourceConfig creates an endpoint with an OIDC
// provider that includes a discovery_url. Used by: TestAccAppEndpointOIDCDiscoveryURL (S22).
func testAccAppEndpointOIDCDiscoveryURLResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsOriginOnlyResourceConfig creates an endpoint with a
// minimal cors block (origin only). Used by: TestAccAppEndpointUpdateCorsExpand (U1 phase 1),
// TestAccAppEndpointUpdateCorsDisableToggle (U3 phase 1, skipped),
// TestAccAppEndpointUpdateCorsOriginWildcardToSpecific (U5 phase 1),
// TestAccAppEndpointUpdateAddOIDC (U7 phase 1), TestAccAppEndpointUpdateRemoveOIDC (U8 phase 2, skipped).
func testAccAppEndpointCorsOriginOnlyResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsDisabledTrueResourceConfig creates an endpoint with
// cors.disabled=true. Used by: TestAccAppEndpointUpdateCorsDisableToggle (U3 phase 2, skipped).
func testAccAppEndpointCorsDisabledTrueResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		disabled = true
		origin   = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointCorsMaxAgeResourceConfig creates an endpoint with a
// parameterised cors.max_age. Used by: TestAccAppEndpointUpdateCorsMaxAgeToZero (U10),
// TestAccAppEndpointUpdateCorsMaxAgeFromZero (U11 phase 2).
func testAccAppEndpointCorsMaxAgeResourceConfig(resourceName, endpointName, collectionName string, maxAge int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		origin  = ["*"]
		max_age = %[9]d
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
		maxAge,
	)
}

// testAccAppEndpointWithOIDCResourceConfig creates an endpoint with cors and a
// minimal OIDC provider. Used by: TestAccAppEndpointUpdateAddOIDC (U7 phase 2),
// TestAccAppEndpointUpdateRemoveOIDC (U8 phase 1, skipped).
func testAccAppEndpointWithOIDCResourceConfig(resourceName, endpointName, collectionName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
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
				"%[7]s" = {}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
	)
}

// testAccAppEndpointACFResourceConfig creates an endpoint with a parameterised
// access_control_function body. Used by: TestAccAppEndpointUpdateACF (U9).
func testAccAppEndpointACFResourceConfig(resourceName, endpointName, collectionName, acfBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "`+appEndpointBucketName+`"
	name            = "%[8]s"

	cors = {
		origin = ["*"]
	}

	scopes = {
		"_default" = {
			collections = {
				"%[7]s" = {
					access_control_function = "%[9]s"
				}
			}
		}
	}
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		collectionName,
		endpointName,
		acfBody,
	)
}
