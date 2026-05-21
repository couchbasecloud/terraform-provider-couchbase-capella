package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// usersDatasourcePerPage caps the page size the test asks the
// couchbase-capella_users datasource for. Bounds the datasource read's work
// to a single HTTP call regardless of how many users the tenant has.
const usersDatasourcePerPage = 100

func TestAccDatasourceUsers(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_users_")
	dsName := randomStringWithPrefix("tf_acc_users_ds_")
	resourceReference := "couchbase-capella_user." + resourceName
	dsReference := "data.couchbase-capella_users." + dsName

	// Match the pattern in user_acceptance_test.go — the username/email pair is
	// a fixture in the test tenant and the invite flow needs a deterministic
	// value. The resource handle uses a randomised name so parallel tests do
	// not collide on terraform state.
	username := "terraform_acceptance_test_ds"
	email := username + "@couchbase.com"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccUsersDataSourceConfig(resourceName, dsName, username, email, usersDatasourcePerPage),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", username),
					resource.TestCheckResourceAttr(resourceReference, "email", email),
					resource.TestCheckResourceAttrSet(resourceReference, "id"),

					resource.TestCheckResourceAttr(dsReference, "organization_id", globalOrgId),
					// data.# is set even when the list is empty (Terraform
					// stores the count "0"), so TestCheckResourceAttrSet would
					// pass on an empty response. Assert at least one element
					// is present via data.0.id to actually prove the datasource
					// returned users.
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
					// Confirm the just-created user is present somewhere in the
					// org. The datasource itself returns at most per_page users
					// (the just-created one may sit on a later page) — so the
					// custom check walks /users directly looking for our id.
					// Bounded by maxUserLookupPages so a runaway loop can't
					// blow the test budget.
					testAccUserListContains(resourceReference),
				),
			},
		},
	})
}

func TestAccDatasourceUsersInvalidOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_users_ds_invalid_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_users" "%[2]s" {
  organization_id = "00000000-0000-0000-0000-000000000000"
}
`, globalProviderBlock, dsName),
				// Read() wraps API errors with "Error Reading Capella Users"; a
				// bogus org id is rejected by /v4/organizations/.../users with a
				// 403/404. Require both so this test only passes for that exact
				// failure mode, not for any unrelated diagnostic.
				ExpectError: regexp.MustCompile(`(?s)Error Reading Capella Users.*"httpStatusCode":(403|404)`),
			},
		},
	})
}

// testAccUsersDataSourceConfig emits HCL for the user resource + users
// datasource. A perPage of 0 (or negative) means "do not constrain the
// datasource" — callers that need every page to find the just-created user
// via data.* membership assertions (e.g. the membership test) should pass 0.
// Callers that bound the datasource read to a single page (and verify
// membership via a separate paginated lookup) should pass a positive value.
func testAccUsersDataSourceConfig(resourceName, dsName, username, email string, perPage int) string {
	perPageLine := ""
	if perPage > 0 {
		perPageLine = fmt.Sprintf("  per_page        = %d\n", perPage)
	}
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_user" "%[3]s" {
  organization_id = "%[2]s"

  name  = "%[5]s"
  email = "%[6]s"

  organization_roles = [
    "organizationMember"
  ]
}

data "couchbase-capella_users" "%[4]s" {
  organization_id = "%[2]s"
%[7]s
  depends_on = [couchbase-capella_user.%[3]s]
}
`, globalProviderBlock, globalOrgId, resourceName, dsName, username, email, perPageLine)
}

// maxUserLookupPages caps how many pages testAccUserListContains will scan
// looking for the just-created user. At perPage=100 this covers the first
// 10,000 users — generous for any realistic tenant.
const maxUserLookupPages = 100

// testAccUserListContains looks up the resource's id from state, then walks
// /v4/organizations/{org}/users in pages of 100 (independent of the datasource
// under test) until it finds that id or exhausts the lookup budget. This is
// the only way to assert "our user is in the org" without relying on
// server-side email/name filtering on the /users endpoint (filed as a
// follow-up against AV-131648).
func testAccUserListContains(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceReference]
		if !ok {
			return fmt.Errorf("resource %q not found in state", resourceReference)
		}
		wantID := rs.Primary.Attributes["id"]
		if wantID == "" {
			return fmt.Errorf("resource %q has empty id", resourceReference)
		}

		client := api.NewClient(timeout)
		ctx := context.Background()
		for page := 1; page <= maxUserLookupPages; page++ {
			q := url.Values{}
			q.Set("page", strconv.Itoa(page))
			q.Set("perPage", strconv.Itoa(usersDatasourcePerPage))
			pageURL := fmt.Sprintf("%s/v4/organizations/%s/users?%s", globalHost, globalOrgId, q.Encode())
			cfg := api.EndpointCfg{Url: pageURL, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			resp, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
			if err != nil {
				return fmt.Errorf("lookup page %d: %w", page, err)
			}
			var body struct {
				Data   []api.GetUserResponse `json:"data"`
				Cursor api.Cursor            `json:"cursor"`
			}
			if err := json.Unmarshal(resp.Body, &body); err != nil {
				return fmt.Errorf("unmarshal page %d: %w", page, err)
			}
			for _, u := range body.Data {
				if u.Id.String() == wantID {
					return nil
				}
			}
			if body.Cursor.Pages.Next == 0 {
				return fmt.Errorf("user id %q not found in org %q after %d page(s) (total items=%d)",
					wantID, globalOrgId, page, body.Cursor.Pages.TotalItems)
			}
		}
		return fmt.Errorf("user id %q not found in first %d pages of org %q",
			wantID, maxUserLookupPages, globalOrgId)
	}
}
