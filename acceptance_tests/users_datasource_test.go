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
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// usersDatasourcePerPage bounds the datasource read to one HTTP call.
const usersDatasourcePerPage = 100

// usersProbeBudget is how long a single GET /users?page=1&perPage=1 may take
// before we conclude the tenant's /users endpoint is too slow for this test
// (tracked by AV-131649).
const usersProbeBudget = 30 * time.Second

// probeUsersEndpoint times one /users?page=1&perPage=1 call. Returns the
// elapsed time on success; on transport/HTTP error returns the error.
func probeUsersEndpoint(ctx context.Context) (time.Duration, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/users?page=1&perPage=1", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	start := time.Now()
	_, err := api.NewClient(timeout).ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	return time.Since(start), err
}

func TestAccDatasourceUsers(t *testing.T) {
	// AV-131649: on tenants where /users is unreasonably slow, the apply step
	// will exhaust the provider's 5-min HTTP client timeout. Skip cleanly with
	// a clear message rather than burning the full test budget.
	if elapsed, err := probeUsersEndpoint(context.Background()); err != nil {
		t.Skipf("skipping: /users probe failed (%v) — see AV-131649", err)
	} else if elapsed > usersProbeBudget {
		t.Skipf("skipping: /users probe took %s (> %s budget) — slow backend, see AV-131649", elapsed, usersProbeBudget)
	}

	resourceName := randomStringWithPrefix("tf_acc_users_")
	dsName := randomStringWithPrefix("tf_acc_users_ds_")
	resourceReference := "couchbase-capella_user." + resourceName
	dsReference := "data.couchbase-capella_users." + dsName

	// Random per-run email — avoids 422 (code 8010) "email already registered"
	// when a previous panicked run leaves the invited user undeleted in the org.
	username := resourceName
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
					// data.0.id (not data.#) — data.# is "0" on empty lists.
					resource.TestCheckResourceAttrSet(dsReference, "data.0.id"),
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

func TestAccDatasourceUsersEmptyOrganization(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_users_ds_empty_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_users" "%[2]s" {
  organization_id = ""
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value.*Attribute organization_id string length must be at least 1, got: 0`),
			},
		},
	})
}

func TestAccDatasourceUsersInvalidSortDirection(t *testing.T) {
	dsName := randomStringWithPrefix("tf_acc_users_ds_sort_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%[1]s

data "couchbase-capella_users" "%[2]s" {
	organization_id = "00000000-0000-0000-0000-000000000000"
  sort_direction  = "sideways"
}
`, globalProviderBlock, dsName),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Value Match.*sideways.*asc.*desc`),
			},
		},
	})
}

// perPage <= 0 omits per_page from the datasource (walk all pages).
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

const maxUserLookupPages = 100

// testAccUserListContains paginates /users directly for the resource's id,
// covering the case where the just-created user falls outside the
// datasource's per_page slice.
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
