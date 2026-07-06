package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	eventingapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/eventingfunction"
)

// eventingFunctionCode is a minimal valid eventing function handler reused across tests.
const eventingFunctionCode = `function OnUpdate(doc, meta) {\n  log(\"updated\", meta.id);\n}\n`

// eventingFunctionCodeV1/V2 are valid handlers differing only by a logged marker (v1/v2) to prove a code update.
const (
	eventingFunctionCodeV1 = `function OnUpdate(doc, meta) {\n  log(\"v1\", meta.id);\n}\n`
	eventingFunctionCodeV2 = `function OnUpdate(doc, meta) {\n  log(\"v2\", meta.id);\n}\n`
)

// eventingFunctionCodeInvalid is broken JS: created undeployed but fails compilation on deploy.
const eventingFunctionCodeInvalid = `function OnUpdate(doc, meta) {\n  log(\"broken\", meta.id\n`

// eventingCodeContains checks the state's code contains marker (the API may normalise code, so exact match is avoided).
func eventingCodeContains(marker string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if !strings.Contains(value, marker) {
			return fmt.Errorf("expected code to contain %q, got %q", marker, value)
		}
		return nil
	}
}

func TestAccEventingFunctionResourceRequiredOnly(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_req_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(funcReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(funcReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttrSet(funcReference, "code"),
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					// settings are optional+computed: server defaults are populated even when omitted.
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(funcReference, "settings.feed_boundary", "from_now"),
					resource.TestCheckResourceAttr(funcReference, "settings.max_timer_context_size", "1024"),
					resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", "true"),
					resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", "false"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
			{
				// Updates are only applied while undeployed/paused, so the function stays undeployed here.
				Config: testAccEventingFunctionResourceConfigUpdatable(funcName, 3, 90),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "3"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "90"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigRequiredOnly(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }
}
`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, funcName, eventingFunctionCode, globalBucketName, globalMetadataBucketName, globalScopeName, globalCollectionName)
}

func testAccEventingFunctionResourceConfigUpdatable(funcName string, workerCount, scriptTimeout int) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    worker_count   = %[11]d
    script_timeout = %[12]d
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		workerCount,
		scriptTimeout)
}

// TestAccEventingFunctionResourceFullPayload creates a function with every optional field (deploy + all bindings) and verifies each round-trips.
func TestAccEventingFunctionResourceFullPayload(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_full_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create and Read: every optional attribute populated, deployed on create.
			{
				Config: testAccEventingFunctionResourceConfigFullPayload(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(funcReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(funcReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "description", "Eventing function created with every optional field populated."),
					resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
					resource.TestCheckResourceAttrSet(funcReference, "code"),

					// explicit keyspaces on both source and metadata storage.
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.collection", globalCollectionName),

					// every settings field explicitly set (feed_boundary differs from the default of "from_now").
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "2"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(funcReference, "settings.feed_boundary", "everything"),
					resource.TestCheckResourceAttr(funcReference, "settings.max_timer_context_size", "1024"),
					resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", "true"),
					resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", "false"),

					// bucket binding.
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "dst_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.permission", "readWrite"),

					// URL binding with basic auth. The password is sensitive and not asserted.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "myEndpoint"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.url", "https://example.com/api"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "false"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.username", "svc_user"),

					// constant binding.
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.0.alias", "maxRetries"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.0.value", "3"),
				),
			},
			// ImportState without ImportStateVerify: the sensitive URL password is omitted by the API and would mismatch.
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccEventingFunctionResourceConfigFullPayload(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  description     = "Eventing function created with every optional field populated."
  state           = "deployed"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    worker_count           = 2
    script_timeout         = 60
    sql_consistency        = "none"
    language_compatibility = "7.2.0"
    feed_boundary          = "everything"
    max_timer_context_size = 1024
    allow_sync_documents   = true
    cursor_aware           = false
  }

  bindings = {
    buckets = [
      {
        alias      = "dst_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      }
    ]

    urls = [
      {
        alias                    = "myEndpoint"
        url                      = "https://example.com/api"
        allow_cookies            = true
        validate_tls_certificate = false
        authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_password"
        }
      }
    ]

    constants = [
      {
        alias = "maxRetries"
        value = "3"
      }
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceUpdateCode (scenario 03): a deployed function's code change requires moving to undeployed/paused in the same apply.
func TestAccEventingFunctionResourceUpdateCode(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_upd_code_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Create the function deployed with the v1 handler.
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV1, "deployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
					resource.TestCheckResourceAttrWith(funcReference, "code", eventingCodeContains("v1")),
				),
			},
			// Changing the code while the function stays deployed is rejected by the provider.
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV2, "deployed"),
				ExpectError: regexp.MustCompile("Cannot change eventing function while deployed"),
			},
			// Changing the code AND undeploying in the same apply succeeds and stores the v2 handler.
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV2, "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					resource.TestCheckResourceAttrWith(funcReference, "code", eventingCodeContains("v2")),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceEmptyCode: an empty code string is rejected at plan time by the LengthAtLeast(1) validator.
func TestAccEventingFunctionResourceEmptyCode(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_empty_code_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, "", "undeployed"),
				ExpectError: regexp.MustCompile("string length must be at least 1"),
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidCode: invalid JS is created undeployed but fails on deploy (ERR_HANDLER_COMPILATION).
func TestAccEventingFunctionResourceInvalidCode(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_code_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeInvalid, "deployed"),
				ExpectError: regexp.MustCompile("Error setting state of eventing function after create|Error creating eventing function"),
			},
		},
	})
}

// TestAccEventingFunctionResourceCodeOmitted (TC-CR-01/09): omitting code makes the server supply a default boilerplate handler.
func TestAccEventingFunctionResourceCodeOmitted(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_no_code_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigNoCode(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					// code was omitted; the server returns a default comment-stub boilerplate.
					resource.TestCheckResourceAttrSet(funcReference, "code"),
					resource.TestCheckResourceAttrWith(funcReference, "code", eventingCodeContains("Add your mutation logic here")),
					resource.TestCheckResourceAttrWith(funcReference, "code", eventingCodeContains("Add your delete handling logic here")),
				),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigCodeState builds a function config with an explicit code and state.
func testAccEventingFunctionResourceConfigCodeState(funcName, code, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[11]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		code,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		state)
}

// testAccEventingFunctionResourceConfigNoCode builds a function config with the code attribute omitted.
func testAccEventingFunctionResourceConfigNoCode(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"

  event_source = {
    bucket     = "%[6]s"
    scope      = "%[8]s"
    collection = "%[9]s"
  }

  event_metadata_storage = {
    bucket     = "%[7]s"
    scope      = "%[8]s"
    collection = "%[9]s"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceAppServicesCompat (scenario 04): allow_sync_documents=false (compat on) + cursor_aware=true round-trip; other settings keep defaults.
func TestAccEventingFunctionResourceAppServicesCompat(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_appsvc_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigAppServicesCompat(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),

					// App Services compatibility ON, and cursor_aware flipped from its default.
					resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", "false"),
					resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", "true"),

					// remaining settings keep their server defaults (guards against these flags disturbing others).
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(funcReference, "settings.feed_boundary", "from_now"),
					resource.TestCheckResourceAttr(funcReference, "settings.max_timer_context_size", "1024"),

					// bucket binding into the App Services-linked keyspace.
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "app_svc_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.scope", globalScopeName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.collection", globalCollectionName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.permission", "readWrite"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccEventingFunctionResourceAppServicesCompatAdvanced (TC-CR-06): allow_sync_documents=false + cursor_aware=false persist independently.
func TestAccEventingFunctionResourceAppServicesCompatAdvanced(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_appsvc_adv_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigSettings(
					funcName,
					"allow_sync_documents = false\n    cursor_aware         = false",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", "false"),
					resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", "false"),
					// remaining settings keep their server defaults.
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigAppServicesCompat(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    # App Services compatibility ENABLED: skip App Services / mobile-sync documents.
    allow_sync_documents = false
    cursor_aware         = true
  }

  bindings = {
    buckets = [
      {
        alias      = "app_svc_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      }
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// Scenario 05: each URL auth type (none/basic/bearer/digest) round-trips; secrets carry forward from the plan; functions stay undeployed.

func TestAccEventingFunctionResourceURLAuthNone(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_none_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := `authentication = {
          type = "none"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "noneEp", "https://example.com/none", false, true, authBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "noneEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.url", "https://example.com/none"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "false"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "none"),
				),
			},
		},
	})
}

func TestAccEventingFunctionResourceURLAuthBasic(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_basic_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := `authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "basicEp", "https://example.com/basic", true, false, authBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "basicEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.url", "https://example.com/basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "false"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.username", "svc_user"),
					// secret is carried forward from the plan, so it is present in state.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "svc_pass"),
				),
			},
		},
	})
}

func TestAccEventingFunctionResourceURLAuthBearer(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_bearer_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := `authentication = {
          type         = "bearer"
          bearer_token = "tok_abc123"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "bearerEp", "https://example.com/bearer", false, true, authBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "bearerEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.url", "https://example.com/bearer"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "false"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "bearer"),
					// secret is carried forward from the plan, so it is present in state.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.bearer_token", "tok_abc123"),
				),
			},
		},
	})
}

func TestAccEventingFunctionResourceURLAuthDigest(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_digest_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := `authentication = {
          type     = "digest"
          username = "svc_user"
          password = "svc_pass"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "digestEp", "https://example.com/digest", false, true, authBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "digestEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.url", "https://example.com/digest"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "false"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "digest"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.username", "svc_user"),
					// secret is carried forward from the plan, so it is present in state.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "svc_pass"),
				),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigURLAuth builds a function with one URL binding whose auth block the caller supplies.
func testAccEventingFunctionResourceConfigURLAuth(funcName, alias, url string, allowCookies, validateTLS bool, authBlock string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    urls = [
      {
        alias                    = "%[11]s"
        url                      = "%[12]s"
        allow_cookies            = %[13]t
        validate_tls_certificate = %[14]t
        %[15]s
      },
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		alias,
		url,
		allowCookies,
		validateTLS,
		authBlock)
}

// TestAccEventingFunctionResourceMultipleBindings (scenario 06): 3 bucket + 3 URL + 3 constant bindings round-trip in order.
func TestAccEventingFunctionResourceMultipleBindings(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_multi_bind_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigMultipleBindings(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),

					// 3 bucket bindings, in order.
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "3"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "src_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.permission", "read"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.1.alias", "work_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.1.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.1.permission", "readWrite"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.2.alias", "archive_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.2.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.2.permission", "readWrite"),

					// 3 URL bindings (one per auth flavour), in order.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.#", "3"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "apiNone"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "none"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.alias", "apiBasic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.authentication.type", "basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.authentication.username", "svc_user"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.authentication.password", "svc_pass"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.2.alias", "apiBearer"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.2.authentication.type", "bearer"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.2.authentication.bearer_token", "tok_abc123"),

					// 3 constant bindings, in order.
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.#", "3"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.0.alias", "maxRetries"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.0.value", "3"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.1.alias", "region"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.1.value", "ap-south-1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.2.alias", "featureFlag"),
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.2.value", "true"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigMultipleBindings(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    buckets = [
      {
        alias      = "src_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "read"
      },
      {
        alias      = "work_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      },
      {
        alias      = "archive_col"
        bucket     = "%[8]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      },
    ]

    urls = [
      {
        alias                    = "apiNone"
        url                      = "https://example.com/none"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type = "none"
        }
      },
      {
        alias                    = "apiBasic"
        url                      = "https://example.com/basic"
        allow_cookies            = true
        validate_tls_certificate = false
        authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }
      },
      {
        alias                    = "apiBearer"
        url                      = "https://example.com/bearer"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type         = "bearer"
          bearer_token = "tok_abc123"
        }
      },
    ]

    constants = [
      {
        alias = "maxRetries"
        value = "3"
      },
      {
        alias = "region"
        value = "ap-south-1"
      },
      {
        alias = "featureFlag"
        value = "true"
      },
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceOmittedScopeCollection (scenario 08): omitting scope/collection on both keyspaces computes to _default.
func TestAccEventingFunctionResourceOmittedScopeCollection(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_omit_ks_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigOmittedScopeCollection(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),

					// event_source: bucket set explicitly, scope/collection computed to _default.
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "event_source.collection", "_default"),

					// event_metadata_storage: bucket set explicitly, scope/collection computed to _default.
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.collection", "_default"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccEventingFunctionResourceConfigOmittedScopeCollection(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  # Only the bucket is set; scope and collection are left to compute.
  event_source = {
    bucket = "%[7]s"
  }

  # Only the bucket is set; scope and collection are omitted.
  event_metadata_storage = {
    bucket = "%[8]s"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName)
}

// TestAccEventingFunctionResourceBucketBindingOmittedScope (scenario 09): a bucket binding omitting scope/collection computes to _default.
func TestAccEventingFunctionResourceBucketBindingOmittedScope(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bind_omit_ks_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigBucketBindingOmittedScope(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "dst_col"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.permission", "readWrite"),
					// scope/collection omitted in config, computed to _default by the server.
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.collection", "_default"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func testAccEventingFunctionResourceConfigBucketBindingOmittedScope(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    buckets = [
      {
        alias      = "dst_col"
        bucket     = "%[7]s"
        permission = "readWrite"
        # scope and collection intentionally omitted.
      },
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceURLValidateTLS (scenario 10): two none-auth URL bindings differing only on validate_tls_certificate round-trip.
func TestAccEventingFunctionResourceURLValidateTLS(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_tls_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLValidateTLS(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.#", "2"),

					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "tlsOnEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "none"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.validate_tls_certificate", "true"),

					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.alias", "tlsOffEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.authentication.type", "none"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.validate_tls_certificate", "false"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigURLValidateTLS(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    urls = [
      {
        alias                    = "tlsOnEp"
        url                      = "https://example.com/tls-on"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type = "none"
        }
      },
      {
        alias                    = "tlsOffEp"
        url                      = "https://example.com/tls-off"
        allow_cookies            = false
        validate_tls_certificate = false
        authentication = {
          type = "none"
        }
      },
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceURLAllowCookies (scenario 11): explicit allow_cookies=true persists; omitted computes to false.
func TestAccEventingFunctionResourceURLAllowCookies(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_url_cookies_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAllowCookies(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.#", "2"),

					// explicit value persists.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "cookiesOnEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "true"),

					// omitted -> computed default of false.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.alias", "cookiesDefEp"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.1.allow_cookies", "false"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigURLAllowCookies(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    urls = [
      {
        alias                    = "cookiesOnEp"
        url                      = "https://example.com/cookies-on"
        allow_cookies            = true
        validate_tls_certificate = true
        authentication = {
          type = "none"
        }
      },
      {
        alias                    = "cookiesDefEp"
        url                      = "https://example.com/cookies-default"
        validate_tls_certificate = true
        # allow_cookies omitted -> expect computed default of false.
        authentication = {
          type = "none"
        }
      },
    ]
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccDatasourceEventingFunctionFullPayload (scenario 12): the get data source returns the full payload (description, settings, all bindings) with export=false.
func TestAccDatasourceEventingFunctionFullPayload(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_ds_evt_full_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName
	dataSourceReference := "data.couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionDataSourceConfigFull(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", funcName),
					resource.TestCheckResourceAttr(dataSourceReference, "description", "Full payload read via data source."),
					resource.TestCheckResourceAttr(dataSourceReference, "export", "false"),
					resource.TestCheckResourceAttr(dataSourceReference, "status", "undeployed"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "code", funcReference, "code"),

					// keyspaces.
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.bucket", globalMetadataBucketName),

					// non-default settings must deserialise correctly through the data source.
					resource.TestCheckResourceAttr(dataSourceReference, "settings.worker_count", "2"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.script_timeout", "90"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.feed_boundary", "everything"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.allow_sync_documents", "false"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.cursor_aware", "true"),

					// all three binding types surfaced by the data source.
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.buckets.#", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.buckets.0.alias", "dst_col"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.buckets.0.permission", "readWrite"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.#", "2"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.0.alias", "apiBasic"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.0.authentication.type", "basic"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.0.authentication.username", "svc_user"),
					// the API redacts secrets to ***** and the data source surfaces that verbatim (no plan carry-forward).
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.0.authentication.password", "*****"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.1.alias", "apiBearer"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.1.authentication.type", "bearer"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.urls.1.authentication.bearer_token", "*****"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.constants.#", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.constants.0.alias", "maxRetries"),
					resource.TestCheckResourceAttr(dataSourceReference, "bindings.constants.0.value", "3"),
				),
			},
		},
	})
}

// testAccEventingFunctionDataSourceConfigFull creates a rich function and a get data source (export=false) reading it.
func testAccEventingFunctionDataSourceConfigFull(funcName string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  description     = "Full payload read via data source."
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    worker_count         = 2
    script_timeout       = 90
    feed_boundary        = "everything"
    allow_sync_documents = false
    cursor_aware         = true
  }

  bindings = {
    buckets = [
      {
        alias      = "dst_col"
        bucket     = "%[7]s"
        scope      = "%[9]s"
        collection = "%[10]s"
        permission = "readWrite"
      }
    ]

    urls = [
      {
        alias                    = "apiBasic"
        url                      = "https://example.com/basic"
        allow_cookies            = true
        validate_tls_certificate = false
        authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }
      },
      {
        alias                    = "apiBearer"
        url                      = "https://example.com/bearer"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type         = "bearer"
          bearer_token = "tok_abc123"
        }
      }
    ]

    constants = [
      {
        alias = "maxRetries"
        value = "3"
      }
    ]
  }
}

data "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = couchbase-capella_eventing_function.%[5]s.name
  export          = false
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName)
}

// TestAccEventingFunctionResourceImportUndeployed (scenario 15): imports an undeployed function with ImportStateVerify for a clean round-trip.
func TestAccEventingFunctionResourceImportUndeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_import_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "name", funcName),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// eventingAppCodeSizePadBytes pads the handler past the eventing appcode size limit (raise if a cluster's limit is higher).
const eventingAppCodeSizePadBytes = 8 * 1024 * 1024

// Scenario 16 error cases: each ExpectError matches a deterministic provider summary; alternations cover create-time vs deploy-time rejection.

// TestAccEventingFunctionResourceAppCodeSize (16.02, ERR_APPCODE_SIZE): an oversized handler is rejected by the eventing service.
func TestAccEventingFunctionResourceAppCodeSize(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_appcode_size_fn_")

	// A valid handler padded with a large comment so the failure is size, not compilation.
	bigCode := `function OnUpdate(doc, meta) {\n  /* ` + strings.Repeat("a", eventingAppCodeSizePadBytes) + ` */\n  log(\"x\", meta.id);\n}\n`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, bigCode, "deployed"),
				ExpectError: regexp.MustCompile("Error creating eventing function|Error setting state of eventing function after create"),
			},
		},
	})
}

// TestAccEventingFunctionResourceBucketMissing (16.03, ERR_BUCKET_MISSING): the source bucket does not exist on the cluster.
func TestAccEventingFunctionResourceBucketMissing(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bkt_missing_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigSourceKeyspace(funcName, "no_such_bucket_xyz", "_default", "_default", "deployed"),
				ExpectError: regexp.MustCompile("Error creating eventing function|Error setting state of eventing function after create"),
			},
		},
	})
}

// TestAccEventingFunctionResourceCollectionMissing (16.04, ERR_COLLECTION_MISSING): the source scope/collection do not exist.
func TestAccEventingFunctionResourceCollectionMissing(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_col_missing_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigSourceKeyspace(funcName, globalBucketName, "no_such_scope", "no_such_collection", "deployed"),
				ExpectError: regexp.MustCompile("Error creating eventing function|Error setting state of eventing function after create"),
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidWorkerCount (16.10): worker_count=0 is rejected by the client-side range validator.
func TestAccEventingFunctionResourceInvalidWorkerCount(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_worker_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// worker_count has a client-side range validator (1-64); 0 is rejected at plan time.
				Config:      testAccEventingFunctionResourceConfigWorkerCount(funcName, 0, "undeployed"),
				ExpectError: regexp.MustCompile("value must be between 1 and 64"),
			},
		},
	})
}

// TestAccDatasourceEventingFunctionNotFound (16.07, ERR_APP_NOT_FOUND): the get data source looks up a nonexistent function.
func TestAccDatasourceEventingFunctionNotFound(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_missing_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionDataSourceConfigByName(funcName),
				ExpectError: regexp.MustCompile("Error Reading Capella Eventing Function"),
			},
		},
	})
}

// TestAccEventingFunctionResourcePauseUndeployed (16.08, ERR_APP_NOT_DEPLOYED): pausing an undeployed function is rejected by the API.
func TestAccEventingFunctionResourcePauseUndeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_pause_undep_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "undeployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
			},
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "paused"),
				ExpectError: regexp.MustCompile("Error setting state of eventing function for update"),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigSourceKeyspace builds a function with a caller-supplied source keyspace and state.
func testAccEventingFunctionResourceConfigSourceKeyspace(funcName, srcBucket, srcScope, srcCollection, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[11]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "_default"
    collection = "_default"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		srcBucket,
		globalMetadataBucketName,
		srcScope,
		srcCollection,
		state)
}

// testAccEventingFunctionResourceConfigWorkerCount builds a function with a caller-supplied worker_count and state.
func testAccEventingFunctionResourceConfigWorkerCount(funcName string, workerCount int, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[11]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    worker_count = %[12]d
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		state,
		workerCount)
}

// testAccEventingFunctionDataSourceConfigByName builds a lone get data source looking up funcName (no resource).
func testAccEventingFunctionDataSourceConfigByName(funcName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  export          = false
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName)
}

// TestAccEventingFunctionResourceImportURLSecret (scenario 17 / TC-IMP-Secret-01): import leaves the URL secret empty (no prior to carry forward) and never leaks it.
func TestAccEventingFunctionResourceImportURLSecret(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_import_secret_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := `authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "authEp", "https://example.com/api", false, true, authBlock),
				Check: resource.ComposeAggregateTestCheckFunc(
					// after create, the secret is present in state (carried forward from the plan).
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.username", "svc_user"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "svc_pass"),
				),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateCheck:                     testAccCheckEventingFunctionImportDropsURLSecret,
			},
		},
	})
}

// testAccCheckEventingFunctionImportDropsURLSecret asserts import keeps the non-secret auth fields but leaves the password empty.
func testAccCheckEventingFunctionImportDropsURLSecret(states []*terraform.InstanceState) error {
	for _, is := range states {
		attrs := is.Attributes
		if got := attrs["bindings.urls.0.authentication.type"]; got != "basic" {
			return fmt.Errorf("expected imported auth type %q, got %q", "basic", got)
		}
		if got := attrs["bindings.urls.0.authentication.username"]; got != "svc_user" {
			return fmt.Errorf("expected imported username %q, got %q", "svc_user", got)
		}
		if pw := attrs["bindings.urls.0.authentication.password"]; pw != "" {
			return fmt.Errorf("expected imported URL binding password to be empty (dropped), got %q", pw)
		}
	}
	return nil
}

// Scenario 18: out-of-range numeric settings are rejected by client-side range validators (like worker_count in 16.10).

// TestAccEventingFunctionResourceInvalidScriptTimeout asserts script_timeout = 0 is rejected.
func TestAccEventingFunctionResourceInvalidScriptTimeout(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_timeout_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// script_timeout has a client-side minimum validator; 0 is rejected at plan time.
				Config:      testAccEventingFunctionResourceConfigSettings(funcName, "script_timeout = 0"),
				ExpectError: regexp.MustCompile("value must be at least 1"),
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidMaxTimerContextSize: max_timer_context_size=0 is rejected by the range validator.
func TestAccEventingFunctionResourceInvalidMaxTimerContextSize(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_timer_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				// max_timer_context_size has a client-side range validator (20-20971520); 0 is rejected at plan time.
				Config:      testAccEventingFunctionResourceConfigSettings(funcName, "max_timer_context_size = 0"),
				ExpectError: regexp.MustCompile("value must be between 20"),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigSettings builds an undeployed function with a caller-supplied settings body.
func testAccEventingFunctionResourceConfigSettings(funcName, settingsBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "undeployed"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  settings = {
    %[11]s
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		settingsBody)
}

// TestAccEventingFunctionResourceUpdateDescriptionUndeployed (scenario 19): a first description update on an undeployed function succeeds with consistent state (regression guard).
func TestAccEventingFunctionResourceUpdateDescriptionUndeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_drift_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigDescription(funcName, "initial description", "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "description", "initial description"),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					// computed values resolved on create.
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
				),
			},
			{
				// First update on the undeployed function must succeed with no unknown computed values leaking.
				Config: testAccEventingFunctionResourceConfigDescription(funcName, "updated description", "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "description", "updated description"),
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					// computed values remain known/consistent after the update.
					resource.TestCheckResourceAttr(funcReference, "event_source.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "event_source.collection", "_default"),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.scope", "_default"),
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
				),
			},
		},
	})
}

func testAccEventingFunctionResourceConfigDescription(funcName, description, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  description     = "%[6]s"
  code            = "%[7]s"
  state           = "%[10]s"

  event_source           = { bucket = "%[8]s" }
  event_metadata_storage = { bucket = "%[9]s" }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		description,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		state)
}

// TestAccEventingFunctionResourceOutOfBandSettingDrift (scenario 20): an out-of-band worker_count change is detected and reconciled on the next plan.
func TestAccEventingFunctionResourceOutOfBandSettingDrift(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_oob_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigSettings(funcName, "worker_count = 2"),
				Check:  resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "2"),
			},
			{
				// Change worker_count out of band; the re-plan must show an in-place update, then reconcile to 2.
				PreConfig: testAccEventingFunctionSetWorkerCountOOB(funcName, 4),
				Config:    testAccEventingFunctionResourceConfigSettings(funcName, "worker_count = 2"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(funcReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "2"),
			},
		},
	})
}

// testAccEventingFunctionSetWorkerCountOOB is a PreConfig hook that changes worker_count via the API (panics on failure).
func testAccEventingFunctionSetWorkerCountOOB(funcName string, workerCount int64) func() {
	return func() {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s",
			globalHost, globalOrgId, globalProjectId, globalClusterId, funcName,
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodPut, SuccessStatus: http.StatusNoContent}
		body := eventingapi.UpdateEventingFunctionRequest{
			Settings: &eventingapi.Settings{WorkerCount: &workerCount},
		}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, body, globalToken, nil); err != nil {
			panic(fmt.Sprintf("out-of-band worker_count update for %q failed: %v", funcName, err))
		}
	}
}

// Update test cases: a deployed function's definition cannot be changed in place; move it to undeployed/paused in the same apply.

// TestAccEventingFunctionResourceUpdateSettings (TC-UP-Settings): updates worker_count/script_timeout/sql_consistency while undeployed; others keep defaults.
func TestAccEventingFunctionResourceUpdateSettings(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_upd_settings_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigSettings(funcName, "worker_count   = 2\n    script_timeout = 60\n    sql_consistency = \"none\""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "2"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "none"),
				),
			},
			{
				Config: testAccEventingFunctionResourceConfigSettings(funcName, "worker_count   = 5\n    script_timeout = 120\n    sql_consistency = \"request\""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "settings.worker_count", "5"),
					resource.TestCheckResourceAttr(funcReference, "settings.script_timeout", "120"),
					resource.TestCheckResourceAttr(funcReference, "settings.sql_consistency", "request"),
					// unchanged settings keep their defaults.
					resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(funcReference, "settings.feed_boundary", "from_now"),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceUpdateLanguageCompatibility (TC-UP-LangCompat): cycles language_compatibility through all four values while undeployed.
func TestAccEventingFunctionResourceUpdateLanguageCompatibility(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_langcompat_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	step := func(version string) resource.TestStep {
		return resource.TestStep{
			Config: testAccEventingFunctionResourceConfigSettings(funcName, fmt.Sprintf("language_compatibility = %q", version)),
			Check:  resource.TestCheckResourceAttr(funcReference, "settings.language_compatibility", version),
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			step("6.0.0"),
			step("6.5.0"),
			step("6.6.2"),
			step("7.2.0"),
		},
	})
}

// TestAccEventingFunctionResourceUpdateAppSvcCompatMatrix (TC-UP-AppSvcCompat): toggles allow_sync_documents/cursor_aware through all four combos while undeployed.
func TestAccEventingFunctionResourceUpdateAppSvcCompatMatrix(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_appsvc_mtx_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	step := func(allowSync, cursorAware bool) resource.TestStep {
		body := fmt.Sprintf("allow_sync_documents = %t\n    cursor_aware         = %t", allowSync, cursorAware)
		return resource.TestStep{
			Config: testAccEventingFunctionResourceConfigSettings(funcName, body),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(funcReference, "settings.allow_sync_documents", fmt.Sprintf("%t", allowSync)),
				resource.TestCheckResourceAttr(funcReference, "settings.cursor_aware", fmt.Sprintf("%t", cursorAware)),
			),
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			step(true, true),
			step(true, false),
			step(false, true),
			step(false, false),
		},
	})
}

// TestAccEventingFunctionResourceUpdateKeyspacesUndeployed (TC-UP-Source/Meta-Undeployed): swaps source/metadata buckets while undeployed.
func TestAccEventingFunctionResourceUpdateKeyspacesUndeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_upd_ks_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigKeyspaces(funcName, globalBucketName, globalMetadataBucketName, "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalMetadataBucketName),
				),
			},
			{
				Config: testAccEventingFunctionResourceConfigKeyspaces(funcName, globalMetadataBucketName, globalBucketName, "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "event_source.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(funcReference, "event_metadata_storage.bucket", globalBucketName),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceRemoveDescription (TC-UP-Optional-Omit): clearing a set description should empty it on read; the fix lives in branch AV-136448-desc, so this is skipped until that merges.
func TestAccEventingFunctionResourceRemoveDescription(t *testing.T) {
	t.Skip("description-clear fix is delivered separately (AV-136448-desc); un-skip once it merges")

	funcName := randomStringWithPrefix("tf_acc_evt_rm_desc_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigMaybeDescription(funcName, "initial description", "undeployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "description", "initial description"),
			},
			{
				Config: testAccEventingFunctionResourceConfigMaybeDescription(funcName, "", "undeployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "description", ""),
			},
		},
	})
}

// TestAccEventingFunctionResourceUpdateNoOp (TC-UP-NoOp): re-applying an unchanged config yields an empty plan.
func TestAccEventingFunctionResourceUpdateNoOp(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_noop_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
			},
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccEventingFunctionResourceUpdateWhileDeployedRejected (TC-UP-Desc/Source-Deployed/Mixed): a deployed-function change is rejected; moving to undeployed in the same apply succeeds.
func TestAccEventingFunctionResourceUpdateWhileDeployedRejected(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_upd_dep_rej_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigMaybeDescription(funcName, "initial description", "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				// Definition change while deployed is rejected.
				Config:      testAccEventingFunctionResourceConfigMaybeDescription(funcName, "changed description", "deployed"),
				ExpectError: regexp.MustCompile("Cannot change eventing function while deployed"),
			},
			{
				// Supported pattern: change the definition and move to undeployed in the same apply.
				Config: testAccEventingFunctionResourceConfigMaybeDescription(funcName, "changed description", "undeployed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
					resource.TestCheckResourceAttr(funcReference, "description", "changed description"),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceSecretRotate (TC-UP-SecretRotate): a new URL secret is carried into state.
func TestAccEventingFunctionResourceSecretRotate(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_secret_rot_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := func(password string) string {
		return fmt.Sprintf(`authentication = {
          type     = "basic"
          username = "svc_user"
          password = %q
        }`, password)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "authEp", "https://example.com/api", false, true, authBlock("secret_v1")),
				Check:  resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "secret_v1"),
			},
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "authEp", "https://example.com/api", false, true, authBlock("secret_v2")),
				Check:  resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "secret_v2"),
			},
		},
	})
}

// TestAccEventingFunctionResourceSecretPassthrough (TC-UP-SecretPassthrough): passing ***** keeps the existing secret while other URL fields change.
func TestAccEventingFunctionResourceSecretPassthrough(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_secret_pass_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	authBlock := func(password string) string {
		return fmt.Sprintf(`authentication = {
          type     = "basic"
          username = "svc_user"
          password = %q
        }`, password)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "authEp", "https://example.com/api", false, true, authBlock("secret_v1")),
				Check:  resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "basic"),
			},
			{
				// Change a non-secret field (allow_cookies) and pass ***** to preserve the secret.
				Config: testAccEventingFunctionResourceConfigURLAuth(funcName, "authEp", "https://example.com/api", true, true, authBlock("*****")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.allow_cookies", "true"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.password", "*****"),
				),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigKeyspaces builds a function with explicit source/metadata buckets and state.
func testAccEventingFunctionResourceConfigKeyspaces(funcName, srcBucket, metaBucket, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[9]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "_default"
    collection = "_default"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "_default"
    collection = "_default"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		srcBucket,
		metaBucket,
		state)
}

// testAccEventingFunctionResourceConfigMaybeDescription builds a function, including the description line only when non-empty.
func testAccEventingFunctionResourceConfigMaybeDescription(funcName, description, state string) string {
	descLine := ""
	if description != "" {
		descLine = fmt.Sprintf("\n  description     = %q", description)
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[9]s"%[10]s

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[11]s"
    collection = "%[12]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[11]s"
    collection = "%[12]s"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		state,
		descLine,
		globalScopeName,
		globalCollectionName)
}

// testAccEventingFunctionResourceConfigBindings builds a function with a caller-supplied bindings body and state.
func testAccEventingFunctionResourceConfigBindings(funcName, state, bindingsBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"
  state           = "%[11]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  bindings = {
    %[12]s
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		state,
		bindingsBody)
}

// TestAccEventingFunctionResourceUpdateCodePaused (TC-UP-Code): code is updated while paused and the function stays paused.
func TestAccEventingFunctionResourceUpdateCodePaused(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_code_paused_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV1, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				// Pause (state change only).
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV1, "paused"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
			{
				// Change code while paused: appcode is updated, the function stays paused.
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCodeV2, "paused"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "state", "paused"),
					resource.TestCheckResourceAttrWith(funcReference, "code", eventingCodeContains("v2")),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceUpdateBindingsPaused (TC-UP-BucketBinding/UrlBinding-AuthChange/ConstBinding): modifies bindings while paused.
func TestAccEventingFunctionResourceUpdateBindingsPaused(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bind_paused_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	body1 := fmt.Sprintf(`buckets = [
      {
        alias      = "b1"
        bucket     = %[1]q
        scope      = %[2]q
        collection = %[3]q
        permission = "read"
      },
      {
        alias      = "b2"
        bucket     = %[1]q
        scope      = %[2]q
        collection = %[3]q
        permission = "read"
      },
    ]

    urls = [
      {
        alias                    = "u1"
        url                      = "https://example.com/api"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type = "none"
        }
      },
    ]

    constants = [
      {
        alias = "c1"
        value = "1"
      },
    ]`, globalBucketName, globalScopeName, globalCollectionName)

	// b1 permission changed, b2 removed, b3 added; URL auth none -> basic; constant value 1 -> 2.
	body2 := fmt.Sprintf(`buckets = [
      {
        alias      = "b1"
        bucket     = %[1]q
        scope      = %[2]q
        collection = %[3]q
        permission = "readWrite"
      },
      {
        alias      = "b3"
        bucket     = %[1]q
        scope      = %[2]q
        collection = %[3]q
        permission = "readWrite"
      },
    ]

    urls = [
      {
        alias                    = "u1"
        url                      = "https://example.com/api"
        allow_cookies            = false
        validate_tls_certificate = true
        authentication = {
          type     = "basic"
          username = "svc_user"
          password = "svc_pass"
        }
      },
    ]

    constants = [
      {
        alias = "c1"
        value = "2"
      },
    ]`, globalBucketName, globalScopeName, globalCollectionName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "deployed", body1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "2"),
				),
			},
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "paused", body1),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "paused", body2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "state", "paused"),
					// bucket bindings: one modified (permission), one removed, one added.
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.#", "2"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "b1"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.permission", "readWrite"),
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.1.alias", "b3"),
					// URL binding auth changed none -> basic.
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.type", "basic"),
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.authentication.username", "svc_user"),
					// constant binding value changed.
					resource.TestCheckResourceAttr(funcReference, "bindings.constants.0.value", "2"),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceUpdateMetaPaused (TC-UP-Meta-Paused): a keyspace change while paused is rejected (422); keyspace changes require undeployed.
func TestAccEventingFunctionResourceUpdateMetaPaused(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_meta_paused_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigKeyspaces(funcName, globalBucketName, globalMetadataBucketName, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				Config: testAccEventingFunctionResourceConfigKeyspaces(funcName, globalBucketName, globalMetadataBucketName, "paused"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
			{
				// Swapping keyspaces while paused is rejected by the eventing service.
				Config:      testAccEventingFunctionResourceConfigKeyspaces(funcName, globalMetadataBucketName, globalBucketName, "paused"),
				ExpectError: regexp.MustCompile("Error updating eventing function"),
			},
		},
	})
}

// TestAccEventingFunctionResourceActivationLifecycle (TC-AS-UP-01/02/03/05 + Idemp-Deploy/Pause): walks the terminal activation states and asserts idempotent re-applies.
func TestAccEventingFunctionResourceActivationLifecycle(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_activation_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	stateStep := func(state string) resource.TestStep {
		return resource.TestStep{
			Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, state),
			Check:  resource.TestCheckResourceAttr(funcReference, "state", state),
		}
	}
	// noopStep re-applies the same state and asserts an empty plan (no activation call).
	noopStep := func(state string) resource.TestStep {
		return resource.TestStep{
			Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, state),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
			},
			Check: resource.TestCheckResourceAttr(funcReference, "state", state),
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			stateStep("undeployed"), // create
			stateStep("deployed"),   // TC-AS-UP-01: undeployed -> deploy
			noopStep("deployed"),    // TC-AS-Idemp-Deploy
			stateStep("undeployed"), // TC-AS-UP-02: deployed -> undeploy
			stateStep("deployed"),   // redeploy to set up pause
			stateStep("paused"),     // TC-AS-UP-03: deployed -> pause
			noopStep("paused"),      // TC-AS-Idemp-Pause
			stateStep("undeployed"), // TC-AS-UP-05: paused -> undeploy
		},
	})
}

// TestAccEventingFunctionResourceActivationResume (TC-AS-UP-04): resuming a paused function (state=deployed, deploy verb) ends deployed.
func TestAccEventingFunctionResourceActivationResume(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_resume_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "paused"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
			{
				// Resume: paused -> deployed.
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
		},
	})
}

// testAccCheckEventingFunctionDestroy verifies every eventing function in state was removed from the cluster.
func testAccCheckEventingFunctionDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "couchbase-capella_eventing_function" {
			continue
		}
		attrs := rs.Primary.Attributes
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s",
			globalHost, attrs["organization_id"], attrs["project_id"], attrs["cluster_id"], attrs["name"],
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		_, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
		if err == nil {
			return fmt.Errorf("eventing function %q (%s) still exists after destroy", attrs["name"], name)
		}
		if notFound, _ := api.CheckResourceNotFoundError(err); !notFound {
			return fmt.Errorf("unexpected error verifying destroy of %q: %w", attrs["name"], err)
		}
	}
	return nil
}

// testAccEventingFunctionDeleteOOB deletes the named function directly via the API (simulates out-of-band removal).
func testAccEventingFunctionDeleteOOB(funcName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/eventingFunctions/%s",
			globalHost, globalOrgId, globalProjectId, globalClusterId, funcName,
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("out-of-band delete of %q failed: %w", funcName, err)
		}
		return nil
	}
}

// TestAccEventingFunctionResourceDeleteUndeployed (TC-DEL-01): destroying an undeployed function removes it from the cluster.
func TestAccEventingFunctionResourceDeleteUndeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_del_undep_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckEventingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "undeployed"),
			},
		},
	})
}

// TestAccEventingFunctionResourceDeleteDeployed (TC-DEL-02): destroying a deployed function undeploys then deletes it.
func TestAccEventingFunctionResourceDeleteDeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_del_dep_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckEventingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
		},
	})
}

// TestAccEventingFunctionResourceDeletePaused (TC-DEL-03): destroying a paused function undeploys then deletes it.
func TestAccEventingFunctionResourceDeletePaused(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_del_paused_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckEventingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "paused"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
		},
	})
}

// TestAccEventingFunctionResourceDeleteAlreadyGone (TC-DEL-05): destroy gracefully handles a function already removed out of band (404).
func TestAccEventingFunctionResourceDeleteAlreadyGone(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_del_gone_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckEventingFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccEventingFunctionResourceConfigRequiredOnly(funcName),
				Check:              testAccEventingFunctionDeleteOOB(funcName),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccEventingFunctionResourceImportDeployed (TC-IMP-02): imports a deployed function with ImportStateVerify (status=deployed).
func TestAccEventingFunctionResourceImportDeployed(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_imp_dep_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccEventingFunctionResourceImportPaused (TC-IMP-03): imports a paused function (status=paused).
func TestAccEventingFunctionResourceImportPaused(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_imp_paused_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "deployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "deployed"),
			},
			{
				Config: testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "paused"),
				Check:  resource.TestCheckResourceAttr(funcReference, "state", "paused"),
			},
			{
				ResourceName:                         funcReference,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    generateEventingFunctionImportIdForResource(funcReference),
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

// TestAccEventingFunctionResourceImportNotFound (TC-IMP-04): a well-formed ID for a nonexistent function yields no imported state.
func TestAccEventingFunctionResourceImportNotFound(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_imp_nf_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
			},
			{
				ResourceName: funcReference,
				ImportState:  true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf(
						"function_name=tf_acc_nonexistent_fn,cluster_id=%s,project_id=%s,organization_id=%s",
						globalClusterId, globalProjectId, globalOrgId,
					), nil
				},
				ExpectError: regexp.MustCompile("Cannot import non-existent remote object"),
			},
		},
	})
}

// TestAccEventingFunctionResourceImportMalformedID (TC-IMP-05): a malformed import ID is rejected while parsing, before any API call.
func TestAccEventingFunctionResourceImportMalformedID(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_imp_bad_id_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigRequiredOnly(funcName),
			},
			{
				ResourceName:  funcReference,
				ImportState:   true,
				ImportStateId: "function_name=foo,cluster_id=bar",
				ExpectError:   regexp.MustCompile("terraform import parameters did not match"),
			},
		},
	})
}

// ─── Coverage-gap tests (review follow-up) ───────────────────────────────────

// TestAccEventingFunctionResourceCreatePausedRejected (gap 1): creating a function directly in the paused state is rejected.
func TestAccEventingFunctionResourceCreatePausedRejected(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_create_paused_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "paused"),
				ExpectError: regexp.MustCompile("Cannot create a paused eventing function"),
			},
		},
	})
}

// TestAccEventingFunctionResourceEmptyBindings (gap 2): an empty bindings={} is accepted and does not drift on re-apply.
func TestAccEventingFunctionResourceEmptyBindings(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_empty_bind_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "undeployed", ""),
				Check:  resource.TestCheckResourceAttr(funcReference, "name", funcName),
			},
			{
				// Re-applying the identical config must be a no-op (guards empty-vs-null bindings handling).
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "undeployed", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidState (gap 3): an invalid state value is rejected by the OneOf validator at plan time.
func TestAccEventingFunctionResourceInvalidState(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_state_fn_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigCodeState(funcName, eventingFunctionCode, "running"),
				ExpectError: regexp.MustCompile("value must be one of"),
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidPermission (gap 4): an invalid bucket binding permission is rejected by the OneOf validator at plan time.
func TestAccEventingFunctionResourceInvalidPermission(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_perm_fn_")

	bindingsBody := fmt.Sprintf(`buckets = [
      {
        alias      = "b1"
        bucket     = %[1]q
        scope      = %[2]q
        collection = %[3]q
        permission = "write"
      },
    ]`, globalBucketName, globalScopeName, globalCollectionName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigBindings(funcName, "undeployed", bindingsBody),
				ExpectError: regexp.MustCompile("value must be one of"),
			},
		},
	})
}

// TestAccEventingFunctionResourceInvalidAuthType (gap 5): an invalid URL auth type is rejected by the OneOf validator at plan time.
func TestAccEventingFunctionResourceInvalidAuthType(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_bad_auth_fn_")

	authBlock := `authentication = {
          type = "oauth"
        }`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccEventingFunctionResourceConfigURLAuth(funcName, "badEp", "https://example.com/api", false, true, authBlock),
				ExpectError: regexp.MustCompile("value must be one of"),
			},
		},
	})
}

// TestAccEventingFunctionResourceRenameRequiresReplace (gap 6): changing name forces resource replacement.
func TestAccEventingFunctionResourceRenameRequiresReplace(t *testing.T) {
	label := randomStringWithPrefix("tf_acc_evt_rename_")
	nameA := randomStringWithPrefix("tf_acc_evt_name_a_")
	nameB := randomStringWithPrefix("tf_acc_evt_name_b_")
	funcReference := "couchbase-capella_eventing_function." + label

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigNamed(label, nameA, "undeployed"),
				Check:  resource.TestCheckResourceAttr(funcReference, "name", nameA),
			},
			{
				Config: testAccEventingFunctionResourceConfigNamed(label, nameB, "undeployed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(funcReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttr(funcReference, "name", nameB),
			},
		},
	})
}

// TestAccEventingFunctionResourceBucketBindingOmittedPermission (gap 7): an omitted bucket binding permission is computed by the server.
func TestAccEventingFunctionResourceBucketBindingOmittedPermission(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_omit_perm_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	bindingsBody := fmt.Sprintf(`buckets = [
      {
        alias  = "dst_col"
        bucket = %[1]q
        scope  = %[2]q
        collection = %[3]q
      },
    ]`, globalBucketName, globalScopeName, globalCollectionName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "undeployed", bindingsBody),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.buckets.0.alias", "dst_col"),
					// permission omitted -> computed by the server.
					resource.TestCheckResourceAttrSet(funcReference, "bindings.buckets.0.permission"),
				),
			},
		},
	})
}

// TestAccEventingFunctionResourceURLBindingOmittedAuth (gap 8): an omitted URL authentication block is computed by the server.
func TestAccEventingFunctionResourceURLBindingOmittedAuth(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_evt_omit_auth_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName

	bindingsBody := `urls = [
      {
        alias = "noAuthEp"
        url   = "https://example.com/api"
      },
    ]`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionResourceConfigBindings(funcName, "undeployed", bindingsBody),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(funcReference, "bindings.urls.0.alias", "noAuthEp"),
					// authentication omitted -> computed by the server.
					resource.TestCheckResourceAttrSet(funcReference, "bindings.urls.0.authentication.type"),
				),
			},
		},
	})
}

// testAccEventingFunctionResourceConfigNamed builds a function with separate resource label and name, so a rename keeps the address constant.
func testAccEventingFunctionResourceConfigNamed(label, name, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[6]s"
  code            = "%[7]s"
  state           = "%[12]s"

  event_source = {
    bucket     = "%[8]s"
    scope      = "%[10]s"
    collection = "%[11]s"
  }

  event_metadata_storage = {
    bucket     = "%[9]s"
    scope      = "%[10]s"
    collection = "%[11]s"
  }
}
`, globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		label,
		name,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		state)
}

func generateEventingFunctionImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
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
			"function_name=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["name"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

func TestAccDatasourceEventingFunction(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_ds_evt_fn_")
	funcReference := "couchbase-capella_eventing_function." + funcName
	dataSourceReference := "data.couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionDataSourceConfig(funcName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", funcName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.sql_consistency", "none"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.language_compatibility", "7.2.0"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.feed_boundary", "from_now"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.max_timer_context_size", "1024"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.allow_sync_documents", "true"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.cursor_aware", "false"),
					resource.TestCheckResourceAttrPair(dataSourceReference, "code", funcReference, "code"),
					// export defaults to false, so the read-only fields (currently just status) is returned
					resource.TestCheckResourceAttr(dataSourceReference, "status", "undeployed"),
				),
			},
		},
	})
}

// TestAccDatasourceEventingFunctionExport: export=true on the get data source omits the read-only status field.
func TestAccDatasourceEventingFunctionExport(t *testing.T) {
	funcName := randomStringWithPrefix("tf_acc_ds_evt_exp_fn_")
	dataSourceReference := "data.couchbase-capella_eventing_function." + funcName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccEventingFunctionDataSourceConfig(funcName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", funcName),
					resource.TestCheckResourceAttr(dataSourceReference, "export", "true"),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.bucket", globalBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_source.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.bucket", globalMetadataBucketName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.scope", globalScopeName),
					resource.TestCheckResourceAttr(dataSourceReference, "event_metadata_storage.collection", globalCollectionName),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.worker_count", "1"),
					resource.TestCheckResourceAttr(dataSourceReference, "settings.script_timeout", "60"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "code"),
					// export omits the read-only status field from the response, so check that is omitted
					resource.TestCheckNoResourceAttr(dataSourceReference, "status"),
				),
			},
		},
	})
}

// testAccEventingFunctionDataSourceConfig creates a function and a get data source reading it.
func testAccEventingFunctionDataSourceConfig(funcName string, export bool) string {
	exportLine := ""
	if export {
		exportLine = "\n  export          = true"
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = "%[5]s"
  code            = "%[6]s"

  event_source = {
    bucket     = "%[7]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }

  event_metadata_storage = {
    bucket     = "%[8]s"
    scope      = "%[9]s"
    collection = "%[10]s"
  }
}

data "couchbase-capella_eventing_function" "%[5]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  name            = couchbase-capella_eventing_function.%[5]s.name%[11]s
}
`,
		globalProviderBlock,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		funcName,
		eventingFunctionCode,
		globalBucketName,
		globalMetadataBucketName,
		globalScopeName,
		globalCollectionName,
		exportLine,
	)
}
