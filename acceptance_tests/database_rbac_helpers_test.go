package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// apiErrorPattern builds an ExpectError regex for a diagnostic the V4 API produced.
//
// Terraform hard-wraps diagnostic text, so a phrase that is contiguous in the JSON
// response arrives split across lines with leading indentation: "is not a valid
// privilege" is rendered as "is not a valid\n        privilege". Matching a literal
// phrase therefore fails for reasons that have nothing to do with the behaviour under
// test. Every space in summary and phrase is matched as arbitrary whitespace instead.
//
// Always pass a phrase. Matching on statusCode alone passes for the wrong reason:
// this endpoint answers several unrelated problems with a 422, so a name that is too
// long or a malformed body satisfies a status-only pattern just as well as the case
// the test means to pin.
//
// statusCode is a regex fragment, so both "422" and "(403|404)" are valid.
func apiErrorPattern(summary, phrase, statusCode string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(
		`(?s)%s.*%s.*"httpStatusCode":%s`,
		whitespaceTolerant(summary), whitespaceTolerant(phrase), statusCode,
	))
}

// whitespaceTolerant quotes each word and joins them so any run of whitespace,
// including a line break Terraform inserted, matches between them.
func whitespaceTolerant(phrase string) string {
	words := strings.Fields(phrase)
	for i, word := range words {
		words[i] = regexp.QuoteMeta(word)
	}
	return strings.Join(words, `\s+`)
}

// Shared helpers for the fine-grained RBAC acceptance tests: database roles,
// database credentials and the data sources that list them.
//
// Naming constraint: Capella rejects a database role name longer than 32 characters
// with a 422. randomStringWithPrefix appends 10 characters, so any prefix used as a
// role name must be at most 22. Overshooting is easy to miss because the resulting
// 422 satisfies a loosely written ExpectError and the test passes for the wrong
// reason - match on the message, not just the status code. Database credential names
// are not subject to this limit.

func databaseRoleURL(organizationId, projectId, clusterId, roleId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/roles/%s",
		globalHost, organizationId, projectId, clusterId, roleId,
	)
}

func databaseCredentialURL(organizationId, projectId, clusterId, credentialId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s",
		globalHost, organizationId, projectId, clusterId, credentialId,
	)
}

// assertGoneFromServer confirms url no longer resolves. Capella answers a GET for a
// deleted role or credential with 404. A 403 is accepted too, because the API returns
// forbidden rather than not-found once the enclosing cluster is no longer readable,
// which is indistinguishable from deletion for the purposes of this check.
func assertGoneFromServer(url, label string) error {
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	_, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
	if err == nil {
		return fmt.Errorf("%s still exists after destroy", label)
	}
	if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
		return nil
	}
	if api.IsForbiddenError(err) {
		return nil
	}
	return fmt.Errorf("unexpected error verifying destroy of %s: %w", label, err)
}

// testAccCheckDatabaseRoleDestroy verifies every database role left in state was
// actually removed from Capella, not merely dropped from the state file.
func testAccCheckDatabaseRoleDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "couchbase-capella_database_role" {
			continue
		}
		attrs := rs.Primary.Attributes
		url := databaseRoleURL(attrs["organization_id"], attrs["project_id"], attrs["cluster_id"], attrs["id"])
		if err := assertGoneFromServer(url, fmt.Sprintf("database role %q (%s)", attrs["name"], name)); err != nil {
			return err
		}
	}
	return nil
}

// testAccCheckDatabaseCredentialDestroy verifies every database credential left in
// state was actually removed from Capella.
func testAccCheckDatabaseCredentialDestroy(s *terraform.State) error {
	for name, rs := range s.RootModule().Resources {
		if rs.Type != "couchbase-capella_database_credential" {
			continue
		}
		attrs := rs.Primary.Attributes
		url := databaseCredentialURL(attrs["organization_id"], attrs["project_id"], attrs["cluster_id"], attrs["id"])
		if err := assertGoneFromServer(url, fmt.Sprintf("database credential %q (%s)", attrs["name"], name)); err != nil {
			return err
		}
	}
	return nil
}

// testAccCheckDatabaseRBACDestroy covers configs that create roles and credentials together.
func testAccCheckDatabaseRBACDestroy(s *terraform.State) error {
	if err := testAccCheckDatabaseRoleDestroy(s); err != nil {
		return err
	}
	return testAccCheckDatabaseCredentialDestroy(s)
}

// testAccDatabaseRoleDeleteOOB deletes the role behind resourceReference directly
// through the V4 API, simulating removal outside Terraform so the next refresh
// exercises the not-found branch of Read.
func testAccDatabaseRoleDeleteOOB(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := databaseRoleAttrsFromState(s, resourceReference)
		if attrs == nil {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		url := databaseRoleURL(attrs["organization_id"], attrs["project_id"], attrs["cluster_id"], attrs["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("out-of-band delete of %s failed: %w", resourceReference, err)
		}
		return nil
	}
}

// testAccCheckResourceIDChanged asserts the resource was replaced rather than updated
// in place between two steps, which is how requires-replace behaviour is verified.
// It reads the id at check time and compares it with the id captured by a prior
// testAccCheckCaptureResourceID against the same holder.
func testAccCheckCaptureResourceID(resourceReference string, holder *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := databaseRoleAttrsFromState(s, resourceReference)
		if attrs == nil {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		*holder = attrs["id"]
		return nil
	}
}

func testAccCheckResourceIDChanged(resourceReference string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := databaseRoleAttrsFromState(s, resourceReference)
		if attrs == nil {
			return fmt.Errorf("resource %s not found in state", resourceReference)
		}
		if *previous == "" {
			return fmt.Errorf("no prior id captured for %s", resourceReference)
		}
		if attrs["id"] == *previous {
			return fmt.Errorf("expected %s to be replaced, but id stayed %s", resourceReference, *previous)
		}
		return nil
	}
}
