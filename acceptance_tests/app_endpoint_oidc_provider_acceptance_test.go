package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppEndpointOidcProvider exercises CRUD and import for the
// couchbase-capella_app_endpoint_oidc_provider resource using the common
// pre-created endpoint (globalAppEndpointName).
//
// Runs sequentially (resource.Test, not ParallelTest) to avoid a race
// condition with TestAccAppEndpointDefaultOidcProvider: both resources write
// OIDC providers to the same endpoint, and concurrent create/destroy cycles
// can cause transient HTTP 500 responses from the API.
func TestAccAppEndpointOidcProvider(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_oidc_")
	resourceReference := "couchbase-capella_app_endpoint_oidc_provider." + resourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcProviderConfig(resourceName, globalAppEndpointName, "https://accounts.google.com", "example-client-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "client_id", "example-client-id"),
					resource.TestCheckResourceAttrSet(resourceReference, "provider_id"),
				),
			},
			{
				Config: testAccOidcProviderFullConfig(resourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "register", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_prefix", "google_"),
					resource.TestCheckResourceAttr(resourceReference, "username_claim", "email"),
					resource.TestCheckResourceAttr(resourceReference, "roles_claim", "roles"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateOidcProviderImportId(resourceReference),
				ImportState:       true,
				// ImportStateVerify omitted: resource has no "id" attribute;
				// composite IDs are not compatible with the default verifier.
			},
		},
	})
}

// TestAccAppEndpointDefaultOidcProvider exercises CRUD and import for the
// couchbase-capella_app_endpoint_default_oidc_provider resource.
// It creates an OIDC provider first (via couchbase-capella_app_endpoint_oidc_provider)
// to obtain a provider_id, then marks it as default.
//
// Runs sequentially after TestAccAppEndpointOidcProvider to avoid concurrent
// OIDC writes to the same endpoint causing API 500 errors.
func TestAccAppEndpointDefaultOidcProvider(t *testing.T) {
	oidcResourceName := randomStringWithPrefix("tf_acc_oidc_")
	defaultResourceName := randomStringWithPrefix("tf_acc_default_oidc_")
	defaultResourceReference := "couchbase-capella_app_endpoint_default_oidc_provider." + defaultResourceName

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, globalAppEndpointName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(defaultResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(defaultResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(defaultResourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_endpoint_name", globalAppEndpointName),
					resource.TestCheckResourceAttrSet(defaultResourceReference, "provider_id"),
				),
			},
			{
				ResourceName:      defaultResourceReference,
				ImportStateIdFunc: generateDefaultOidcProviderImportId(defaultResourceReference),
				ImportState:       true,
				// ImportStateVerify omitted: composite import IDs are not
				// compatible with the default state verifier.
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccOidcProviderConfig(resourceName, endpointName, issuer, clientId string) string {
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
		globalClusterId,
		globalAppServiceId,
		endpointName,
		issuer,
		clientId,
	)
}

func testAccOidcProviderFullConfig(resourceName, endpointName string) string {
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
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		endpointName,
	)
}

func testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, endpointName string) string {
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
		globalClusterId,
		globalAppServiceId,
		endpointName,
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
