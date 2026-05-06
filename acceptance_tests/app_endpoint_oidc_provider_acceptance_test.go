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
// Creates its own bucket and app endpoint so it can run in parallel with other
// tests without competing for the shared common endpoint's OIDC state.
func TestAccAppEndpointOidcProvider(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_oidc_")
	resourceReference := "couchbase-capella_app_endpoint_oidc_provider." + resourceName
	bucketName := randomStringWithPrefix("tf_acc_oidc_bkt_")
	epName := randomStringWithPrefix("tf_acc_oidc_ep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccOidcProviderConfig(resourceName, bucketName, epName, "https://accounts.google.com", "example-client-id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "app_endpoint_name", epName),
					resource.TestCheckResourceAttr(resourceReference, "issuer", "https://accounts.google.com"),
					resource.TestCheckResourceAttr(resourceReference, "client_id", "example-client-id"),
					resource.TestCheckResourceAttrSet(resourceReference, "provider_id"),
				),
			},
			{
				Config: testAccOidcProviderFullConfig(resourceName, bucketName, epName),
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
				// ImportStateVerify cannot be used here: ImportStatePassthroughID
				// stores the full composite ID in app_endpoint_name, but Read
				// normalises it to just the endpoint name. The verifier then
				// fails to locate the resource by the original composite value.
			},
		},
	})
}

// TestAccAppEndpointDefaultOidcProvider exercises CRUD and import for the
// couchbase-capella_app_endpoint_default_oidc_provider resource.
// It creates an OIDC provider first (via couchbase-capella_app_endpoint_oidc_provider)
// to obtain a provider_id, then marks it as default.
//
// Creates its own bucket and app endpoint so it can run in parallel with
// TestAccAppEndpointOidcProvider without competing for the same endpoint.
func TestAccAppEndpointDefaultOidcProvider(t *testing.T) {
	oidcResourceName := randomStringWithPrefix("tf_acc_oidc_")
	defaultResourceName := randomStringWithPrefix("tf_acc_default_oidc_")
	defaultResourceReference := "couchbase-capella_app_endpoint_default_oidc_provider." + defaultResourceName
	bucketName := randomStringWithPrefix("tf_acc_doidc_bkt_")
	epName := randomStringWithPrefix("tf_acc_doidc_ep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, bucketName, epName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(defaultResourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(defaultResourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(defaultResourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(defaultResourceReference, "app_endpoint_name", epName),
					resource.TestCheckResourceAttrSet(defaultResourceReference, "provider_id"),
				),
			},
			{
				ResourceName:      defaultResourceReference,
				ImportStateIdFunc: generateDefaultOidcProviderImportId(defaultResourceReference),
				ImportState:       true,
				// ImportStateVerify cannot be used here: ImportStatePassthroughID
				// stores the full composite ID in app_endpoint_name, but Read
				// normalises it to just the endpoint name. The verifier then
				// fails to locate the resource by the original composite value.
			},
		},
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Config helpers
// ─────────────────────────────────────────────────────────────────────────────

func testAccOidcProviderConfig(resourceName, bucketName, epName, issuer, clientId string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[6]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s_ep" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[7]s"
	bucket          = "%[6]s"
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

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	issuer            = "%[9]s"
	client_id         = "%[10]s"
	depends_on        = [couchbase-capella_app_endpoint.%[2]s_ep]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		bucketName,
		globalAppServiceId,
		epName,
		issuer,
		clientId,
	)
}

func testAccOidcProviderFullConfig(resourceName, bucketName, epName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[6]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s_ep" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[7]s"
	bucket          = "%[6]s"
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

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	issuer            = "https://accounts.google.com"
	client_id         = "example-client-id"
	register          = true
	user_prefix       = "google_"
	username_claim    = "email"
	roles_claim       = "roles"
	depends_on        = [couchbase-capella_app_endpoint.%[2]s_ep]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		bucketName,
		globalAppServiceId,
		epName,
	)
}

func testAccDefaultOidcProviderConfig(oidcResourceName, defaultResourceName, bucketName, epName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name            = "%[6]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s_ep" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[7]s"
	bucket          = "%[6]s"
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

resource "couchbase-capella_app_endpoint_oidc_provider" "%[2]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	issuer            = "https://accounts.google.com"
	client_id         = "example-client-id"
	depends_on        = [couchbase-capella_app_endpoint.%[2]s_ep]
}

resource "couchbase-capella_app_endpoint_default_oidc_provider" "%[9]s" {
	organization_id   = "%[3]s"
	project_id        = "%[4]s"
	cluster_id        = "%[5]s"
	app_service_id    = "%[7]s"
	app_endpoint_name = "%[8]s"
	provider_id       = couchbase-capella_app_endpoint_oidc_provider.%[2]s.provider_id
	depends_on        = [couchbase-capella_app_endpoint_oidc_provider.%[2]s]
}
`,
		globalProviderBlock,
		oidcResourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		bucketName,
		globalAppServiceId,
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
