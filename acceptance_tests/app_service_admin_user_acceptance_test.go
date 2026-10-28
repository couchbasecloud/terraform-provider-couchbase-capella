package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

const (
	testAdminUserPassword    = "Password123!"
	testAdminUserNewPassword = "NewPassword456!"
)

// TestAccAppServiceAdminUserResource covers the create/read/import/update path of an
// admin user scoped to a named App Endpoint, and the transition from an endpoint list
// to access_all_endpoints - the two arms of the update API's oneOf body.
func TestAccAppServiceAdminUserResource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			// Create and Read, granting access to the shared App Endpoint.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					fmt.Sprintf(`endpoints = ["%s"]`, globalAppEndpointName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttr(resourceReference, "password", testAdminUserPassword),
					// endpoints is a Set, so assert membership rather than position.
					resource.TestCheckResourceAttr(resourceReference, "endpoints.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "endpoints.*", globalAppEndpointName),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
					// access_all_endpoints is deliberately not asserted here: the config omits it,
					// and what the provider stores for an omitted value is the subject of
					// TestAccAppServiceAdminUserExplicitAccessAllEndpointsFalse.
				),
			},
			// ImportState. ImportStateVerify is what actually proves the round-trip;
			// password is ignored because the API never returns it.
			//
			// This step only round-trips exactly while the granted bucket holds a single
			// App Endpoint. Capella grants admin access per bucket, so a GET returns
			// every App Endpoint on the bucket; on import there is no prior set to
			// filter that expansion against, so state takes the expanded list - which is
			// what TestAccAppServiceAdminUserBucketLevelEndpointExpansion asserts.
			//
			// That test adds a second endpoint to this same bucket, so it is deliberately
			// serial and finishes its destroy before this parallel test resumes. If the
			// shared fixture ever gains a permanent second App Endpoint on the bucket,
			// this step starts failing on endpoints.
			{
				ResourceName:            resourceReference,
				ImportStateIdFunc:       generateAppServiceAdminUserImportIdForResource(resourceReference),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
			// Update: swap the endpoint list for access_all_endpoints, which sends the
			// other arm of the update oneOf. Assert the whole shape, not just the field
			// that changed, so a reset of another attribute is caught.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					`access_all_endpoints = true`,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access_all_endpoints", "true"),
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttr(resourceReference, "password", testAdminUserPassword),
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					// endpoints must go away entirely, not become an empty set.
					resource.TestCheckNoResourceAttr(resourceReference, "endpoints.#"),
				),
			},
			// Update back the other way. This is not symmetric with the step above: the
			// read now has to filter the bucket-level expansion out of the API response
			// against the planned set, so it exercises reconcileEndpoints on an update
			// rather than just swapping oneOf arms.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					fmt.Sprintf(`endpoints = ["%s"]`, globalAppEndpointName),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "endpoints.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "endpoints.*", globalAppEndpointName),
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.modified_at"),
				),
			},
		},
	})
}

// TestAccAppServiceAdminUserResourceAccessAllEndpoints creates an admin user with
// access to every App Endpoint, the one shape that never populates the endpoints set.
func TestAccAppServiceAdminUserResourceAccessAllEndpoints(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_all_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_all_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					`access_all_endpoints = true`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access_all_endpoints", "true"),
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_at"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.created_by"),
					resource.TestCheckResourceAttrSet(resourceReference, "audit.version"),
					resource.TestCheckNoResourceAttr(resourceReference, "endpoints.#"),
				),
			},
			{
				ResourceName:            resourceReference,
				ImportStateIdFunc:       generateAppServiceAdminUserImportIdForResource(resourceReference),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// TestAccAppServiceAdminUserExplicitAccessAllEndpointsFalse asserts that writing
// access_all_endpoints = false alongside an endpoint list applies cleanly and is
// stored as false.
//
// ValidateConfig accepts that configuration, but the post-create read maps
// access_all_endpoints to null unless the API returned true, so Terraform aborts the
// apply with "Provider produced inconsistent result after apply" and taints the
// resource. Every later apply then destroys and recreates the admin user and fails
// again, so the configuration can never converge.
func TestAccAppServiceAdminUserExplicitAccessAllEndpointsFalse(t *testing.T) {
	skipKnownBug(t, "AV-149567: access_all_endpoints = false is stored as null, so apply fails with an inconsistent result and taints the resource; remove this guard once the read maps the API value and the attribute is Optional+Computed.")

	resourceName := randomStringWithPrefix("tf_acc_admin_user_false_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_false_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					fmt.Sprintf("access_all_endpoints = false\n  endpoints            = [\"%s\"]", globalAppEndpointName),
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access_all_endpoints", "false"),
					resource.TestCheckResourceAttr(resourceReference, "endpoints.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "endpoints.*", globalAppEndpointName),
				),
			},
		},
	})
}

// TestAccAppServiceAdminUserPasswordRotation asserts that changing the password
// actually changes it on the App Service.
//
// The v4 update API has no password field - its body is a oneOf of
// {accessAllEndpoints} or {endpoints} - and password carries no RequiresReplace, so
// today a password change is planned as an in-place update, sends nothing, and reports
// success while the admin user keeps its old password. Replacement is the only way to
// apply a new password, so the plan must show one.
func TestAccAppServiceAdminUserPasswordRotation(t *testing.T) {
	skipKnownBug(t, "AV-149572: a password change reports a successful in-place update but leaves the password unchanged on the App Service; remove this guard once password carries RequiresReplace.")

	resourceName := randomStringWithPrefix("tf_acc_admin_user_pw_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_pw_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					`access_all_endpoints = true`,
				),
				Check: resource.TestCheckResourceAttr(resourceReference, "password", testAdminUserPassword),
			},
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserNewPassword,
					`access_all_endpoints = true`,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttr(resourceReference, "password", testAdminUserNewPassword),
			},
		},
	})
}

// TestAccAppServiceAdminUserResourceInvalidAccess covers the ValidateConfig rules that
// keep an invalid access body from ever reaching the API. Each case is rejected at plan
// time, so none of them provisions anything.
func TestAccAppServiceAdminUserResourceInvalidAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_bad_")
	userName := randomStringWithPrefix("tf_acc_user_bad_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// Neither access_all_endpoints nor endpoints.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword, ``,
				),
				ExpectError: regexp.MustCompile(`(?s)Invalid App Service Admin User Configuration.*Exactly one of access_all_endpoints = true or a non-empty endpoints set`),
			},
			// Both access_all_endpoints = true and a non-empty endpoints set.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					fmt.Sprintf("access_all_endpoints = true\n  endpoints            = [\"%s\"]", globalAppEndpointName),
				),
				ExpectError: regexp.MustCompile(`(?s)Invalid App Service Admin User Configuration.*Exactly one of access_all_endpoints = true or a non-empty endpoints set`),
			},
			// An explicitly empty endpoints set, which the read would store as null.
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword, `endpoints = []`,
				),
				ExpectError: regexp.MustCompile(`(?s)Invalid App Service Admin User Configuration.*endpoints must not be empty`),
			},
		},
	})
}

// TestAccDatasourceAppServiceAdminUser reads back a single admin user created by the
// resource and checks the data source reports the same user.
func TestAccDatasourceAppServiceAdminUser(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_ds_")
	dataSourceReference := "data.couchbase-capella_app_service_admin_user.ds_" + resourceName
	scopedReference := "data.couchbase-capella_app_service_admin_user.scoped_ds_" + resourceName
	userName := randomStringWithPrefix("tf_acc_user_ds_")
	scopedUserName := randomStringWithPrefix("tf_acc_user_ds_scoped_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAppServiceAdminUserDestroy(t, userName),
			testAccCheckAppServiceAdminUserDestroy(t, scopedUserName),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserDataSourceConfig(resourceName, userName, scopedUserName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(dataSourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(dataSourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(dataSourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(dataSourceReference, "name", userName),
					resource.TestCheckResourceAttr(dataSourceReference, "access_all_endpoints", "true"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "id"),
					resource.TestCheckResourceAttrSet(dataSourceReference, "audit.created_at"),
					// An endpoint-scoped user exercises the other mapping branch. Note the
					// data source reports what the API holds - access_all_endpoints really
					// is false, and endpoints is the bucket-expanded list - where the
					// resource stores the narrowed set it was configured with. Membership
					// rather than a count, because the expansion depends on how many App
					// Endpoints share the bucket.
					resource.TestCheckResourceAttr(scopedReference, "name", scopedUserName),
					resource.TestCheckResourceAttr(scopedReference, "access_all_endpoints", "false"),
					resource.TestCheckTypeSetElemAttr(scopedReference, "endpoints.*", globalAppEndpointName),
					resource.TestCheckResourceAttrSet(scopedReference, "audit.created_at"),
				),
			},
		},
	})
}

// TestAccDatasourceAppServiceAdminUsers lists the admin users of the shared App
// Service, with and without a name filter. The non-matching filter is the interesting
// case: it must yield an empty list rather than a null one.
func TestAccDatasourceAppServiceAdminUsers(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_users_ds_")
	userName := randomStringWithPrefix("tf_acc_user_list_")
	filtered := "data.couchbase-capella_app_service_admin_users.filtered_" + resourceName
	noMatch := "data.couchbase-capella_app_service_admin_users.nomatch_" + resourceName
	all := "data.couchbase-capella_app_service_admin_users.all_" + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUsersDataSourceConfig(resourceName, userName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(filtered, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(filtered, "app_service_id", globalAppServiceId),
					// The filter is applied provider-side, so exactly one user must come back.
					resource.TestCheckResourceAttr(filtered, "data.#", "1"),
					resource.TestCheckResourceAttr(filtered, "data.0.name", userName),
					resource.TestCheckResourceAttr(filtered, "data.0.access_all_endpoints", "true"),
					resource.TestCheckResourceAttrSet(filtered, "data.0.id"),
					resource.TestCheckResourceAttrSet(filtered, "data.0.audit.created_at"),
					// An empty list, not a null one: "0" and "absent" are the two states
					// being distinguished here, so asserting the count is the point.
					resource.TestCheckResourceAttr(noMatch, "data.#", "0"),
					// The unfiltered listing is the common usage and takes a different
					// path through Read (no filter, so no provider-side matching). The
					// suite runs in parallel, so assert membership rather than a count.
					testCheckAdminUsersListContains(all, userName),
				),
			},
		},
	})
}

// TestAccDatasourceAppServiceAdminUsersInvalidFilter checks every way the filter block
// can be wrong: half-specified, filtering on an unsupported attribute, and an empty
// values set. All three are rejected at plan time, so none of them lists anything.
func TestAccDatasourceAppServiceAdminUsersInvalidFilter(t *testing.T) {
	withFilter := func(filterBody string) string {
		return fmt.Sprintf(`
	%[1]s

	data "couchbase-capella_app_service_admin_users" "bad_filter" {
	  organization_id = "%[2]s"
	  project_id      = "%[3]s"
	  cluster_id      = "%[4]s"
	  app_service_id  = "%[5]s"

	  filter {
	    %[6]s
	  }
	}
	`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, filterBody)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			// values without name.
			{
				Config:      withFilter(`values = ["something"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Filters Configuration.*Both 'name' and 'values' in filter block must be configured`),
			},
			// name without values.
			{
				Config:      withFilter(`name = "name"`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Filters Configuration.*Both 'name' and 'values' in filter block must be configured`),
			},
			// Only filtering by admin user name is supported.
			{
				Config: withFilter(`name   = "id"
	    values = ["something"]`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*Attribute filter.name value must be one of: \["name"\], got: "id"`),
			},
			// An empty values set would filter everything out silently.
			{
				Config: withFilter(`name   = "name"
	    values = []`),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute filter.values set must contain at least 1 elements, got: 0`),
			},
		},
	})
}

// testAccAppServiceAdminUserResourceConfig builds an admin user resource. accessBlock
// carries the access attributes verbatim so each test can set access_all_endpoints,
// endpoints, both or neither.
func testAccAppServiceAdminUserResourceConfig(resourceName, userName, password, accessBlock string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_admin_user" "%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	  name            = "%[7]s"
	  password        = "%[8]s"
	  %[9]s
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		globalAppServiceId, userName, password, accessBlock)
}

func testAccAppServiceAdminUserDataSourceConfig(resourceName, userName, scopedUserName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_admin_user" "%[2]s" {
	  organization_id      = "%[3]s"
	  project_id           = "%[4]s"
	  cluster_id           = "%[5]s"
	  app_service_id       = "%[6]s"
	  name                 = "%[7]s"
	  password             = "%[8]s"
	  access_all_endpoints = true
	}

	resource "couchbase-capella_app_service_admin_user" "scoped_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	  name            = "%[9]s"
	  password        = "%[8]s"
	  endpoints       = ["%[10]s"]
	}

	data "couchbase-capella_app_service_admin_user" "ds_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	  id              = couchbase-capella_app_service_admin_user.%[2]s.id
	}

	data "couchbase-capella_app_service_admin_user" "scoped_ds_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	  id              = couchbase-capella_app_service_admin_user.scoped_%[2]s.id
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		globalAppServiceId, userName, testAdminUserPassword, scopedUserName, globalAppEndpointName)
}

func testAccAppServiceAdminUsersDataSourceConfig(resourceName, userName string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_admin_user" "%[2]s" {
	  organization_id      = "%[3]s"
	  project_id           = "%[4]s"
	  cluster_id           = "%[5]s"
	  app_service_id       = "%[6]s"
	  name                 = "%[7]s"
	  password             = "%[8]s"
	  access_all_endpoints = true
	}

	data "couchbase-capella_app_service_admin_users" "filtered_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"

	  filter {
	    name   = "name"
	    values = ["%[7]s"]
	  }

	  depends_on = [couchbase-capella_app_service_admin_user.%[2]s]
	}

	data "couchbase-capella_app_service_admin_users" "all_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"

	  depends_on = [couchbase-capella_app_service_admin_user.%[2]s]
	}

	data "couchbase-capella_app_service_admin_users" "nomatch_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"

	  filter {
	    name   = "name"
	    values = ["tf_acc_no_such_admin_user"]
	  }

	  depends_on = [couchbase-capella_app_service_admin_user.%[2]s]
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		globalAppServiceId, userName, testAdminUserPassword)
}

// generateAppServiceAdminUserImportIdForResource builds the composite import key the
// resource's Validate() expects. The order of the pairs does not matter, but the set
// of keys must be exactly these five.
func generateAppServiceAdminUserImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		if rawState == nil {
			return "", fmt.Errorf("resource %s not found in state", resourceReference)
		}
		return fmt.Sprintf(
			"id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"], rawState["app_service_id"], rawState["cluster_id"],
			rawState["project_id"], rawState["organization_id"],
		), nil
	}
}

// testAccCheckAppServiceAdminUserDestroy asserts no admin user of the shared App
// Service still carries userName. It matches on name rather than id because the state
// handed to CheckDestroy no longer holds the destroyed resource.
func testAccCheckAppServiceAdminUserDestroy(t *testing.T, userName string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		data := newTestClient(t)

		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers",
			data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId,
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

		// The list endpoint has no sortBy parameter.
		users, err := api.GetPaginated[[]apigen.AppServiceAdminUser](
			context.Background(), data.ClientV1, data.Token, cfg, "",
		)
		if err != nil {
			return fmt.Errorf("failed to list App Service admin users after destroy: %w", err)
		}

		for _, user := range users {
			if user.Name != nil && *user.Name == userName {
				return fmt.Errorf("App Service admin user %q still exists after destroy", userName)
			}
		}
		return nil
	}
}

// TestAccAppServiceAdminUserNameForcesReplacement checks that renaming an admin user
// destroys and recreates it.
//
// The v4 update API cannot rename a user - its body is a oneOf of {accessAllEndpoints}
// or {endpoints} - so name carries RequiresReplace. Were that modifier ever dropped,
// Terraform would route a rename to Update, which would send no name, report success
// and leave the user called something else: the same silent no-op as AV-149572. The
// plan check is what guards it; the state assertions alone would not, because Update
// echoes the planned name into state either way.
func TestAccAppServiceAdminUserNameForcesReplacement(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_rename_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_rename_")
	renamed := randomStringWithPrefix("tf_acc_user_renamed_")

	var firstId string

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckAppServiceAdminUserDestroy(t, userName),
			testAccCheckAppServiceAdminUserDestroy(t, renamed),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword, `access_all_endpoints = true`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					captureAdminUserAttr(resourceReference, "id", &firstId),
				),
			},
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, renamed, testAdminUserPassword, `access_all_endpoints = true`,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", renamed),
					resource.TestCheckResourceAttr(resourceReference, "access_all_endpoints", "true"),
					// A replacement means a genuinely new user, so the id must change.
					// Without this, an Update that merely echoed the new name would pass.
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceReference]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceReference)
						}
						if rs.Primary.Attributes["id"] == firstId {
							return fmt.Errorf("expected a new admin user id after rename, still %s", firstId)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAppServiceAdminUserInvalidIdentifiers checks that malformed identifiers and an
// empty name are rejected at plan time rather than at apply.
//
// Deferring this class to apply is a recurring complaint on this provider
// (AV-132142 / AV-132143), and it is the cheapest possible test: no step here
// provisions anything.
func TestAccAppServiceAdminUserInvalidIdentifiers(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_ids_")
	userName := randomStringWithPrefix("tf_acc_user_ids_")

	badUUID := func(field string) string {
		cfg := map[string]string{
			"organization_id": globalOrgId,
			"project_id":      globalProjectId,
			"cluster_id":      globalClusterId,
			"app_service_id":  globalAppServiceId,
		}
		cfg[field] = "not-a-uuid"
		return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_admin_user" "%[2]s" {
	  organization_id      = "%[3]s"
	  project_id           = "%[4]s"
	  cluster_id           = "%[5]s"
	  app_service_id       = "%[6]s"
	  name                 = "%[7]s"
	  password             = "%[8]s"
	  access_all_endpoints = true
	}
	`, globalProviderBlock, resourceName, cfg["organization_id"], cfg["project_id"],
			cfg["cluster_id"], cfg["app_service_id"], userName, testAdminUserPassword)
	}

	steps := []resource.TestStep{}
	for _, field := range []string{"organization_id", "project_id", "cluster_id", "app_service_id"} {
		steps = append(steps, resource.TestStep{
			Config:      badUUID(field),
			ExpectError: regexp.MustCompile(fmt.Sprintf(`(?s)Invalid Attribute Value Match.*Attribute %s must be a valid UUID, got: not-a-uuid`, field)),
		})
	}
	// An empty name must fail the length validator rather than reach the API.
	steps = append(steps, resource.TestStep{
		Config: testAccAppServiceAdminUserResourceConfig(
			resourceName, "", testAdminUserPassword, `access_all_endpoints = true`,
		),
		ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Length.*Attribute name string length must be at least 1, got: 0`),
	})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps:                    steps,
	})
}

// TestAccAppServiceAdminUserInvalidEndpoint checks that granting access to an App
// Endpoint that does not exist fails the apply instead of silently creating a user with
// no access. The provider cannot catch this at plan time - endpoint names are not known
// to the schema - so the guard is that the error surfaces at all.
func TestAccAppServiceAdminUserInvalidEndpoint(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_badep_")
	userName := randomStringWithPrefix("tf_acc_user_badep_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: testAccAppServiceAdminUserResourceConfig(
					resourceName, userName, testAdminUserPassword,
					`endpoints = ["tf_acc_no_such_app_endpoint"]`,
				),
				// This matches any create failure, not specifically the rejected
				// endpoint, so in principle an environment problem - a token without
				// create permission, say - would also satisfy it. Narrowing it needs the
				// API's actual message for this case, which has not been captured yet;
				// pinning a status code would be guesswork, because the App Endpoint
				// APIs return 403 rather than 404 for things that do not exist.
				//
				// It is still worth having: the assertion that matters is that a
				// nonexistent endpoint does not silently produce a working admin user,
				// and a token that could not create would fail every other test here
				// rather than passing quietly.
				ExpectError: regexp.MustCompile(`(?s)Error Creating App Service Admin User`),
			},
		},
	})
}

// TestAccAppServiceAdminUserDeletedExternally checks that an admin user deleted outside
// Terraform is dropped from state by the next Read and recreated, rather than erroring.
//
// Customers delete things in the UI (CBSE-22525), so Read must treat a 404 as
// "remove from state". Without this test nothing exercises that branch.
func TestAccAppServiceAdminUserDeletedExternally(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_extdel_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_extdel_")

	var userId string

	cfg := testAccAppServiceAdminUserResourceConfig(
		resourceName, userName, testAdminUserPassword, `access_all_endpoints = true`,
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					captureAdminUserAttr(resourceReference, "id", &userId),
				),
			},
			// Delete the admin user behind Terraform's back, then re-apply the identical
			// configuration. Read must remove it from state and Create must put it back.
			{
				PreConfig: func() {
					if err := deleteAppServiceAdminUserOutOfBand(t, userId); err != nil {
						t.Fatalf("failed to delete admin user %s out of band: %v", userId, err)
					}
					// The delete is accepted with 202, so it may not have taken effect
					// yet. Re-applying before it has lets Read still find the user, which
					// would skip the remove-from-state branch this test exists to cover
					// and fail the id assertion below intermittently.
					if err := waitForAppServiceAdminUserGone(t, userId); err != nil {
						t.Fatalf("admin user %s was not deleted: %v", userId, err)
					}
				},
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttr(resourceReference, "access_all_endpoints", "true"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceReference]
						if !ok {
							return fmt.Errorf("resource %s not found in state", resourceReference)
						}
						if rs.Primary.Attributes["id"] == userId {
							return fmt.Errorf("expected the admin user to be recreated with a new id, still %s", userId)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccAppServiceAdminUserBucketLevelEndpointExpansion covers the one thing the rest
// of this file cannot: what happens when the granted bucket holds more than one App
// Endpoint.
//
// Capella grants admin access per bucket, so a GET returns every App Endpoint on the
// bucket, not just the configured one. reconcileEndpoints filters that expansion back
// out against the prior set. With only the shared fixture endpoint on the bucket there
// is nothing to expand, so that filter never executes - this test provisions a second
// App Endpoint on the same bucket to make it run.
//
// Two behaviours are pinned:
//
//   - Create filters the expansion, so state holds only the configured endpoint and the
//     step's own post-apply plan comes back empty. Without the filter this is a
//     perpetual diff.
//   - Import has no prior set to filter against, so it takes the full expanded list.
//     That is surprising but self-healing - one apply narrows state back - and pinning
//     it means a future change to reconcileEndpoints cannot quietly make it permanent.
//
// The second App Endpoint needs its own collection, because an App Endpoint maps a
// keyspace and two endpoints cannot share one. Both live on the shared bucket, which is
// what makes the admin grant expand across them.
func TestAccAppServiceAdminUserBucketLevelEndpointExpansion(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_admin_user_expand_")
	resourceReference := "couchbase-capella_app_service_admin_user." + resourceName
	userName := randomStringWithPrefix("tf_acc_user_expand_")
	collectionName := randomStringWithPrefix("tf_acc_col_expand_")
	secondEndpoint := randomStringWithPrefix("tf_acc_ep_expand_")

	cfg := testAccAppServiceAdminUserExpansionConfig(resourceName, userName, collectionName, secondEndpoint)

	// Deliberately NOT parallel. This test adds a second App Endpoint to the shared
	// bucket, and TestAccAppServiceAdminUserResource's import step asserts a single
	// endpoint round-trips. A test that does not call t.Parallel() runs to completion -
	// including its destroy - before any paused parallel test resumes, so the extra
	// endpoint never exists while the lifecycle test is observing the bucket. Adding
	// t.Parallel() here reintroduces that race; isolate this onto its own bucket first
	// if it ever needs to run concurrently.
	//
	// The CRUD slot is belt-and-braces: App Endpoint CRUD is capped at
	// appEndpointCRUDConcurrency to avoid the transient 500s seen on concurrent destroy.
	defer acquireAppEndpointCRUDSlot()()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServiceAdminUserDestroy(t, userName),
		Steps: []resource.TestStep{
			// The admin user is granted only the shared endpoint, but the API reports
			// both. Exactly one endpoint in state is the whole point of the filter.
			{
				Config: cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", userName),
					resource.TestCheckResourceAttr(resourceReference, "endpoints.#", "1"),
					// Together these pin the filter: exactly one element, and it is the
					// configured endpoint. An unfiltered read would put the second
					// endpoint here too, because it shares the bucket.
					resource.TestCheckTypeSetElemAttr(resourceReference, "endpoints.*", globalAppEndpointName),
				),
			},
			// Import takes the expanded list. ImportStateVerify would fail by design, so
			// assert the expansion directly instead of comparing against prior state.
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServiceAdminUserImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					attrs := states[0].Attributes
					count := attrs["endpoints.#"]
					if count != "2" {
						return fmt.Errorf(
							"expected import to take the bucket-expanded list of 2 endpoints, got endpoints.#=%q; "+
								"if reconcileEndpoints now filters on import this test should assert 1 instead",
							count,
						)
					}
					found := map[string]bool{}
					for k, v := range attrs {
						if strings.HasPrefix(k, "endpoints.") && k != "endpoints.#" {
							found[v] = true
						}
					}
					for _, want := range []string{globalAppEndpointName, secondEndpoint} {
						if !found[want] {
							return fmt.Errorf("imported endpoints %v is missing %q", found, want)
						}
					}
					return nil
				},
			},
		},
	})
}

// testAccAppServiceAdminUserExpansionConfig builds a second App Endpoint on the shared
// bucket, plus the collection it needs, and an admin user granted only the shared
// endpoint. cors and oidc are omitted deliberately - both are optional, and oidc carries
// its own ordering problems that would muddy a failure here.
func testAccAppServiceAdminUserExpansionConfig(resourceName, userName, collectionName, secondEndpoint string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_collection" "col_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  bucket_id       = "%[6]s"
	  scope_name      = "%[7]s"
	  collection_name = "%[8]s"
	}

	resource "couchbase-capella_app_endpoint" "ep_%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[9]s"
	  bucket          = "%[10]s"
	  name            = "%[11]s"

	  scopes = {
	    "%[7]s" = {
	      collections = {
	        "%[8]s" = {}
	      }
	    }
	  }

	  depends_on = [couchbase-capella_collection.col_%[2]s]
	}

	resource "couchbase-capella_app_service_admin_user" "%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[9]s"
	  name            = "%[12]s"
	  password        = "%[13]s"
	  endpoints       = ["%[14]s"]

	  # The grant is bucket-level, so the second endpoint must already exist for the
	  # read to return an expanded list.
	  depends_on = [couchbase-capella_app_endpoint.ep_%[2]s]
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId,
		globalBucketId, globalScopeName, collectionName, globalAppServiceId, globalBucketName,
		secondEndpoint, userName, testAdminUserPassword, globalAppEndpointName)
}

// TestAccDatasourceAppServiceAdminUserNotFound checks that reading an admin user that
// does not exist reports an error rather than returning an empty object.
//
// The id is one that genuinely existed and was then deleted, rather than an invented
// string. A made-up id is rejected on its *format*, so the data source would error for
// the wrong reason and the test would pass even if genuine not-found handling regressed.
// Creating and deleting a real user is the only way to get an id that is unambiguously
// well-formed and unambiguously absent.
func TestAccDatasourceAppServiceAdminUserNotFound(t *testing.T) {
	userName := randomStringWithPrefix("tf_acc_user_gone_")

	deletedId, err := createAppServiceAdminUserOutOfBand(t, userName)
	if err != nil {
		t.Fatalf("failed to create the admin user to delete: %v", err)
	}
	if err := deleteAppServiceAdminUserOutOfBand(t, deletedId); err != nil {
		t.Fatalf("failed to delete admin user %s: %v", deletedId, err)
	}
	if err := waitForAppServiceAdminUserGone(t, deletedId); err != nil {
		t.Fatalf("admin user %s was not deleted: %v", deletedId, err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
	%[1]s

	data "couchbase-capella_app_service_admin_user" "missing" {
	  organization_id = "%[2]s"
	  project_id      = "%[3]s"
	  cluster_id      = "%[4]s"
	  app_service_id  = "%[5]s"
	  id              = "%[6]s"
	}
	`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, deletedId),
				// The status code is left broad on purpose: the App Endpoint APIs return
				// 403 rather than 404 for things that do not exist, so pinning one code
				// would be guesswork. The id above is what makes this assertion mean
				// something - it can only be a not-found.
				ExpectError: regexp.MustCompile(`(?s)Error Reading App Service Admin User.*"httpStatusCode":(400|403|404)`),
			},
		},
	})
}

// captureAdminUserAttr stores an attribute of a resource in state so a later step can
// compare against it - used to prove a replacement produced a genuinely new user.
func captureAdminUserAttr(resourceReference, attribute string, dest *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		*dest = rs.Primary.Attributes[attribute]
		return nil
	}
}

// testCheckAdminUsersListContains asserts the listing data source returned a user with
// this name. The suite runs in parallel so the list also holds other tests' users, which
// rules out an exact count assertion.
func testCheckAdminUsersListContains(dataSourceReference, userName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceReference]
		if !ok {
			return fmt.Errorf("data source %s not found in state", dataSourceReference)
		}
		count, err := strconv.Atoi(rs.Primary.Attributes["data.#"])
		if err != nil {
			return fmt.Errorf("data source %s has no data.# attribute: %w", dataSourceReference, err)
		}
		for i := 0; i < count; i++ {
			if rs.Primary.Attributes[fmt.Sprintf("data.%d.name", i)] == userName {
				return nil
			}
		}
		return fmt.Errorf("admin user %q is missing from the %d users returned by %s", userName, count, dataSourceReference)
	}
}

// deleteAppServiceAdminUserOutOfBand removes an admin user through the API, simulating a
// deletion made in the Capella UI.
func deleteAppServiceAdminUserOutOfBand(t *testing.T, userId string) error {
	t.Helper()
	data := newTestClient(t)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
	return err
}

// skipKnownBug skips a test that asserts behaviour a filed bug currently breaks.
//
// These tests are written against the *fixed* behaviour, so until the ticket lands a
// failure is the expected outcome and a pass means the guard can be removed. Set
// ACC_RUN_KNOWN_BUGS=1 to run them anyway and confirm the bug still reproduces against
// a live Capella organisation; leave it unset and the suite stays green.
func skipKnownBug(t *testing.T, reason string) {
	t.Helper()
	if os.Getenv("ACC_RUN_KNOWN_BUGS") == "" {
		t.Skip(reason)
	}
	t.Logf("ACC_RUN_KNOWN_BUGS is set, so this test will run: %s", reason)
	t.Log("A FAILURE here is the expected result. A PASS means the bug is fixed and the guard should be removed.")
}

// createAppServiceAdminUserOutOfBand creates an admin user through the API and returns
// its id, bypassing Terraform. Used to obtain a real id for a user that is then deleted.
func createAppServiceAdminUserOutOfBand(t *testing.T, userName string) (string, error) {
	t.Helper()
	data := newTestClient(t)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers",
		data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}

	body := map[string]any{
		"name":     userName,
		"password": testAdminUserPassword,
		"access":   map[string]any{"accessAllEndpoints": true},
	}
	response, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, body, data.Token, nil)
	if err != nil {
		return "", fmt.Errorf("creating admin user %q: %w", userName, err)
	}

	var created apigen.CreateAppServiceAdminUserResponse
	if err := json.Unmarshal(response.Body, &created); err != nil {
		return "", fmt.Errorf("unmarshalling create response for %q: %w", userName, err)
	}
	if created.Id == nil || *created.Id == "" {
		return "", fmt.Errorf("create response for %q carried no id", userName)
	}
	return *created.Id, nil
}

// waitForAppServiceAdminUserGone polls until the admin user is no longer readable.
//
// Deletion is accepted with 202, so it is not necessarily complete when the call
// returns. Anything that depends on the user being gone - a Terraform refresh expected
// to drop it from state, a data source expected to report not-found - has to wait, or it
// races the backend and fails intermittently.
//
// The App Endpoint APIs return 403 rather than 404 for resources that do not exist, so
// both are treated as gone.
func waitForAppServiceAdminUserGone(t *testing.T, userId string) error {
	t.Helper()
	data := newTestClient(t)

	const (
		maxWait       = 2 * time.Minute
		retryInterval = 5 * time.Second
	)

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/adminUsers/%s",
		data.HostURL, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, userId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	deadline := time.Now().Add(maxWait)
	for {
		_, err := data.ClientV1.ExecuteWithRetry(context.Background(), cfg, nil, data.Token, nil)
		if err != nil {
			if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
				return nil
			}
			var apiErr *api.Error
			if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusForbidden {
				return nil
			}
			return fmt.Errorf("unexpected error polling admin user %s: %w", userId, err)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("admin user %s still readable after %s", userId, maxWait)
		}
		time.Sleep(retryInterval)
	}
}
