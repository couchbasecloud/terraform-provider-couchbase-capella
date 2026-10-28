package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointOidcProvider exercises CRUD and import for the
// couchbase-capella_app_endpoint_oidc_provider resource.
//
// Uses a pre-created, dedicated endpoint (ensureOIDCEndpoint) so it can run in
// parallel with other tests without competing for creation slots.
func TestAccAppEndpointOidcProvider(t *testing.T) {
	ensureOIDCEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_oidc_")
	resourceReference := "couchbase-capella_app_endpoint_oidc_provider." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcProviderConfig(resourceName, globalOIDCEndpointName, "https://accounts.google.com", "example-client-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalOIDCEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "client_id", "example-client-id"),
					resource.TestCheckResourceAttrSet(resourceReference, "provider_id"),
				),
			},
			{
				Config: testAccOidcProviderFullConfig(resourceName, globalOIDCEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "register", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_prefix", "google_"),
					resource.TestCheckResourceAttr(resourceReference, "username_claim", "email"),
					resource.TestCheckResourceAttr(resourceReference, "roles_claim", "roles"),
					resource.TestCheckResourceAttr(
						resourceReference, "discovery_url", "https://accounts.google.com/.well-known/openid-configuration",
					),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateOidcProviderImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_endpoint_name",
			},
		},
	})
}

// TestAccAppEndpointDefaultOidcProvider exercises CRUD and import for the
// couchbase-capella_app_endpoint_default_oidc_provider resource.
// It creates an OIDC provider first (via couchbase-capella_app_endpoint_oidc_provider)
// to obtain a provider_id, then marks it as default.
//
// Uses its own pre-created endpoint (ensureDefaultOIDCEndpoint) so it can run
// in parallel with TestAccAppEndpointOidcProvider without state conflicts.
func TestAccAppEndpointDefaultOidcProvider(t *testing.T) {
	ensureDefaultOIDCEndpoint(t)

	oidcResourceName := randomStringWithPrefix("tf_acc_oidc_")
	defaultResourceName := randomStringWithPrefix("tf_acc_default_oidc_")
	defaultResourceReference := "couchbase-capella_app_endpoint_default_oidc_provider." + defaultResourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, globalDefaultOIDCEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(defaultResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(defaultResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(defaultResourceReference, "cluster_id", appEndpointClusterId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_service_id", appEndpointAppServiceId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_endpoint_name", globalDefaultOIDCEndpointName),
					resource.TestCheckResourceAttrSet(defaultResourceReference, "provider_id"),
				),
			},
			{
				ResourceName:                         defaultResourceReference,
				ImportStateIdFunc:                    generateDefaultOidcProviderImportId(defaultResourceReference),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "app_endpoint_name",
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccOidcProviderConfig(resourceName, epName, issuer, clientId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	issuer            = "%[8]s"
	client_id         = "%[9]s"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		epName,
		issuer,
		clientId,
	)
}

func testAccOidcProviderFullConfig(resourceName, epName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	issuer            = "https://accounts.google.com"
	client_id         = "example-client-id"
	register          = true
	user_prefix       = "google_"
	username_claim    = "email"
	roles_claim       = "roles"
	discovery_url     = "https://accounts.google.com/.well-known/openid-configuration"
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		epName,
	)
}

func testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, epName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	issuer            = "https://accounts.google.com"
	client_id         = "example-client-id"
}

resource "couchbase-capella_app_endpoint_default_oidc_provider" "%[8]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[6]s"
	app_endpoint_name = "%[7]s"
	provider_id       = couchbase-capella_app_endpoint_oidc_provider.%[2]s.provider_id
	depends_on        = [couchbase-capella_app_endpoint_oidc_provider.%[2]s]
}
`,
		globalProviderBlock,
		oidcResourceName,
		globalOrgId,
		globalProjectId,
		appEndpointClusterId,
		appEndpointAppServiceId,
		epName,
		defaultResourceName,
	)
}

func generateOidcProviderImportId(resourceReference string) resource.ImportStateIdFunc {
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
		return fmt.Sprintf(
			"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s,provider_id=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["app_service_id"],
			rawState["app_endpoint_name"],
			rawState["provider_id"],
		), nil
	}
}

func generateDefaultOidcProviderImportId(resourceReference string) resource.ImportStateIdFunc {
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
		return fmt.Sprintf(
			"organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s",
			rawState["organization_id"],
			rawState["project_id"],
			rawState["cluster_id"],
			rawState["app_service_id"],
			rawState["app_endpoint_name"],
		), nil
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Regression probes
// ─────────────────────────────────────────────────────────────────────────────

// TestAccAppEndpointOidcProviderIdPopulated_AV_149878 covers AV-149878: the empty-provider_id
// branch of AppEndpointOidcProvider.Create (internal/resources/app_endpoint_oidc.go:120).
//
// PR #813 (AV-145497) hardened the case where the create 201 body cannot be unmarshalled:
// the apply now fails loudly and tells the operator the provider was orphaned. The sibling
// case is untouched. A body that IS valid JSON but carries no `providerId` unmarshals
// cleanly into "", because api.AppEndpointOIDCProviderResponse.ProviderID is a plain string
// tagged `omitempty`. Nothing checks it before it is written to state, so:
//
//  1. the apply SUCCEEDS, emitting only a warning;
//  2. state is written with provider_id = "";
//  3. the next refresh hard-fails in Validate() with
//     "could not validate resource state, please contact Couchbase Capella Support:
//     terraform resource was empty: ... providerId";
//  4. destroy reads the 404 on .../oidcProviders/ as "already gone", drops the resource
//     from state and leaves the real provider orphaned in Capella - silently.
//
// Steps 1-4 were reproduced offline against a stub API. What is NOT established is whether
// Capella ever returns such a body; that cannot be forced from a Terraform configuration,
// which is why this is a probe rather than a regression guard. Run it against a live org:
// a pass says the create response always carries providerId and the branch is unreachable
// in practice, and the test then stands as a guard against the API regressing. A failure
// says it is reachable, and the step that fails tells you which of 1-4 you are in.
//
// Two things here that TestAccAppEndpointOidcProvider above does not already cover:
//
//   - The explicit non-empty check names the defect. TestCheckResourceAttrSet also rejects
//     "" (helper/resource/testing.go:1159), so the existing test would catch step 2 - but it
//     reports "Attribute 'provider_id' expected to be set", which reads as a missing
//     attribute rather than an orphaned provider, and would send the reader looking in the
//     wrong place.
//   - The RefreshState step reaches step 3. The existing test aborts at its Check before any
//     refresh runs, so the Validate() hard-error is currently unreachable by any test.
//
// No ExpectEmptyPlan here on purpose: a perpetual diff from some unrelated attribute would
// fail this test for a reason that has nothing to do with provider_id, and muddy the live
// repro this exists to serve.
//
// Uses a client_id distinct from the test above so both can share globalOIDCEndpointName
// without colliding on an issuer/client_id duplicate.
func TestAccAppEndpointOidcProviderIdPopulated_AV_149878(t *testing.T) {
	ensureOIDCEndpoint(t)

	resourceName := randomStringWithPrefix("tf_acc_oidc_id_")
	resourceReference := "couchbase-capella_app_endpoint_oidc_provider." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcProviderConfig(
					resourceName, globalOIDCEndpointName,
					"https://accounts.google.com", "oidc-provider-id-probe-client-id",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceReference, "provider_id", requireNonEmptyProviderId),
				),
			},
			{
				// Proves the refresh does not hard-fail in Validate(). Reached only because
				// the step above wrote a usable provider_id.
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceReference, "provider_id", requireNonEmptyProviderId),
				),
			},
		},
	})
}

// requireNonEmptyProviderId fails with the consequence rather than the symptom, so a live
// failure points at the orphan instead of at a missing attribute.
func requireNonEmptyProviderId(value string) error {
	if value == "" {
		return fmt.Errorf(
			"provider_id is empty: the create response carried no providerId, so the OIDC " +
				"provider exists in Capella but is untracked. The next refresh will fail in " +
				"Validate() and destroy will silently drop it from state, orphaning it. " +
				"See internal/resources/app_endpoint_oidc.go:120",
		)
	}
	return nil
}
