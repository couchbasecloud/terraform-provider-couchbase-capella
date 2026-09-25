package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strconv"
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

// accessShape is an order-insensitive description of one database role access entry.
//
// access and its nested buckets, scopes and collections are Lists, so Terraform
// compares them positionally. The V4 API does not echo back the order it was sent,
// and reconcileAccess can only realign the response against prior state. Import has
// no prior state, so an imported role's entries arrive in whatever order the API
// chose and a blanket ImportStateVerify fails on any role with more than one entry -
// not because anything was lost, but because entry 0 is no longer the same entry.
// Tests describe the grants they expect with accessShape instead and compare them as
// an unordered collection. See AV-143880.
//
// The shape covers what these tests configure: at most one bucket per entry and at
// most one scope per bucket. Granting several buckets in a single entry needs a
// richer comparison than this.
type accessShape struct {
	// privileges is compared as a set, matching the schema attribute.
	privileges []string

	// bucket is empty for an entry with no resources block.
	bucket string

	// scope is empty when the grant stops at bucket level.
	scope string

	// collections is nil when the grant stops at scope level.
	collections []string
}

// key renders the shape so that two shapes describing the same grant compare equal
// whatever order their privileges and collections happen to be stored in.
func (s accessShape) key() string {
	return strings.Join([]string{
		strings.Join(sortedCopy(s.privileges), ","),
		s.bucket,
		s.scope,
		strings.Join(sortedCopy(s.collections), ","),
	}, "|")
}

func sortedCopy(values []string) []string {
	sorted := append([]string(nil), values...)
	sort.Strings(sorted)
	return sorted
}

func accessShapeKeys(shapes []accessShape) []string {
	keys := make([]string, len(shapes))
	for i, shape := range shapes {
		keys[i] = shape.key()
	}
	sort.Strings(keys)
	return keys
}

// checkImportedAccess asserts the imported access list holds exactly the wanted
// grants, ignoring the order the API returned them in. Pair it with
// ImportStateVerifyIgnore on "access" so every attribute outside the access list is
// still compared attribute by attribute.
func checkImportedAccess(want ...accessShape) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected 1 imported state, got %d", len(states))
		}
		got, err := parseAccessShapes(states[0].Attributes)
		if err != nil {
			return err
		}
		gotKeys, wantKeys := accessShapeKeys(got), accessShapeKeys(want)
		if !slices.Equal(gotKeys, wantKeys) {
			return fmt.Errorf(
				"imported access entries (order ignored) = %v, want %v",
				gotKeys, wantKeys,
			)
		}
		return nil
	}
}

// parseAccessShapes reads a database role's access list out of the flat attribute map
// Terraform keeps instance state in.
func parseAccessShapes(attrs map[string]string) ([]accessShape, error) {
	count, err := strconv.Atoi(attrs["access.#"])
	if err != nil {
		return nil, fmt.Errorf("access.# is %q, want a count: %w", attrs["access.#"], err)
	}
	shapes := make([]accessShape, count)
	for i := range shapes {
		bucket := fmt.Sprintf("access.%d.resources.buckets.0.", i)
		scope := bucket + "scopes.0."
		shapes[i] = accessShape{
			privileges:  elementsOf(attrs, fmt.Sprintf("access.%d.privileges", i)),
			bucket:      attrs[bucket+"name"],
			scope:       attrs[scope+"name"],
			collections: elementsOf(attrs, scope+"collections"),
		}
	}
	return shapes, nil
}

// elementsOf collects the elements of a flat list or set attribute. A missing count
// means the attribute is absent from state, which is reported as nil rather than as
// an error so callers can distinguish an omitted block from an empty one.
func elementsOf(attrs map[string]string, key string) []string {
	count, err := strconv.Atoi(attrs[key+".#"])
	if err != nil || count == 0 {
		return nil
	}
	elements := make([]string, count)
	for i := range elements {
		elements[i] = attrs[fmt.Sprintf("%s.%d", key, i)]
	}
	return elements
}
