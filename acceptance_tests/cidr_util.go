package acceptance_tests

import (
	cryptorand "crypto/rand"
	"fmt"
	"net/netip"
	"strings"
	"sync"
	"testing"

	"github.com/couchbase/tools-common/types/ptr"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

const (
	// maxCIDRGenerationAttempts is the maximum number of random CIDRs to generate
	// when trying to find one that doesn't conflict with existing clusters.
	maxCIDRGenerationAttempts = 10

	// maxCIDRConflictRetries is the maximum number of times to retry a test when
	// a CIDR conflict occurs due to a race condition (another cluster claiming the
	// CIDR between generation and deployment).
	maxCIDRConflictRetries = 3

	// cidrConflictErrorSubstring is the error message substring returned by the API
	// when a cluster deployment fails due to another cluster already using the same
	// CIDR.
	// For example, the existing cluster using CIDR 10.0.0.0/24, then attempting to
	// make a new cluster using 10.0.0.0/24 would result in a conflict error.
	cidrConflictErrorSubstring = "ensure you are passing a unique CIDR block"

	// cidrOverlapErrorSubstring is the error message substring returned by the API
	// when a cluster deployment fails due to another cluster using a CIDR that
	// overlaps with the new cluster's CIDR.
	// For example, the existing cluster using 10.0.0.0/24, then attempting to
	// make a new cluster using 10.0.0.0/23 would result in an overlap error.
	cidrOverlapErrorSubstring = "overlaps with existing resource with CIDR"

	// paginationPerPage is the number of items to request per page when paginating
	// through API results.
	paginationPerPage = 100
)

// runTestWithUniqueCIDR runs a parallel test that requires a unique /23 CIDR block,
// retrying with a freshly generated CIDR if a race-condition conflict is detected at
// deploy time. This handles the case where another cluster uses the same CIDR between
// generation and deployment.
//
// The testCase function receives a unique CIDR string and must return a fully
// configured resource.TestCase.
func runTestWithUniqueCIDR(t *testing.T, parallel bool, testCase func(cidr string) resource.TestCase) {
	t.Helper()
	if parallel {
		t.Parallel()
	}

	for attempt := 0; attempt < maxCIDRConflictRetries; attempt++ {
		cidr := generateUniqueCIDR(t)

		cidrConflict := false
		t.Run(fmt.Sprintf("attempt-%d", attempt+1), func(t *testing.T) {
			tc := testCase(cidr)

			originalErrorCheck := tc.ErrorCheck
			tc.ErrorCheck = func(err error) error {
				if strings.Contains(err.Error(), cidrConflictErrorSubstring) ||
					strings.Contains(err.Error(), cidrOverlapErrorSubstring) {
					cidrConflict = true
					t.Skip("CIDR conflict detected, skipping")
					return nil
				}

				if originalErrorCheck != nil {
					return originalErrorCheck(err)
				}

				return err
			}

			t.Logf("Running test with CIDR %s", cidr)
			resource.Test(t, tc)
		})

		if !cidrConflict {
			return
		}

		inUseCIDRs.invalidate()
		t.Logf("CIDR conflict detected on attempt %d/%d, invalidating CIDR cache", attempt+1, maxCIDRConflictRetries)
	}

	t.Fatalf("Exhausted all %d attempts to deploy with a unique CIDR", maxCIDRConflictRetries)
}

// cidrCache is a thread-safe, lazily-populated cache of CIDRs that are currently
// in-use by the global tenant.
//
// This allows multiple parallel tests to share the same list of in-use CIDRs without
// each making their own API calls to fetch them, while still ensuring that when a
// CIDR conflict is detected when attempting a deployment (caused by a race between
// CIDR generation and deployment), the cache is invalidated so the next test run
// will fetch a fresh, up-to-date list of in-use CIDRs.
type cidrCache struct {
	cidrs []netip.Prefix
	mu    sync.Mutex
	// valid tracks if the cache has been successfully populated
	valid bool
}

// getOrFetch returns the cached CIDR list, fetching from the API if the cache
// is not yet populated or has been invalidated.
// Callers will be blocked if another goroutine is fetching due to the cache not being valid.
func (c *cidrCache) getOrFetch(t *testing.T) ([]netip.Prefix, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cache is populated and valid, can instantly return a copy of the cached CIDRs
	if c.valid {
		out := make([]netip.Prefix, len(c.cidrs))
		copy(out, c.cidrs)
		return out, nil
	}

	// Cache is not valid, need to fetch
	cidrs, err := fetchInUseCIDRs(t)
	if err == nil {
		c.cidrs = cidrs
		c.valid = true
	}

	out := make([]netip.Prefix, len(cidrs))
	copy(out, cidrs)
	return out, err
}

// invalidate marks the cache as stale so the next call to getOrFetch will
// re-fetch from the API. This is called when a CIDR conflict is detected
// when deploying, indicating the cache is out of date.
func (c *cidrCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.valid = false
}

// add appends a CIDR to the cache so that concurrent tests immediately see it
// as in-use so that two tests don't use the same CIDR.
func (c *cidrCache) add(prefix netip.Prefix) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cidrs = append(c.cidrs, prefix)
}

// generateRandomCIDR generates a random /23 CIDR block in the 10.0.0.0/8 private range
// to avoid conflicts with existing clusters in the organization.
// Format: 10.X.Y.0/23 where X is 0-255 and Y is an even number (0, 2, 4, ..., 254)
func generateRandomCIDR() string {
	// Use crypto/rand for cryptographically secure randomness
	buf := make([]byte, 2)
	if _, err := cryptorand.Read(buf); err != nil {
		// If crypto/rand fails, this indicates a serious system issue.
		// For test utilities, it's appropriate to panic rather than continue with weak randomness.
		panic(fmt.Sprintf("failed to generate random CIDR: crypto/rand.Read failed: %v", err))
	}

	// Second octet: 0-255
	secondOctet := int(buf[0])

	// Third octet: must be even for /23 CIDR (0, 2, 4, ..., 254)
	thirdOctet := int(buf[1]) & 0xFE // Clear the last bit to make it even

	return fmt.Sprintf("10.%d.%d.0/23", secondOctet, thirdOctet)
}

// generateUniqueCIDR generates a random /23 CIDR block that does not conflict
// with any existing cluster CIDRs across all projects in the organisation.
// Successfully generated CIDRs are added to the cache so parallel tests
// immediately see them as in-use.
func generateUniqueCIDR(t *testing.T) string {
	t.Helper()

	existing, err := inUseCIDRs.getOrFetch(t)
	if err != nil {
		t.Logf("Error fetching existing CIDRs, assuming no existing CIDRs are being used: %v", err)
	}

	for attempt := 0; attempt < maxCIDRGenerationAttempts; attempt++ {
		candidate, err := netip.ParsePrefix(generateRandomCIDR())
		if err != nil {
			t.Fatalf("failed to parse generated CIDR: %v", err)
		}

		if !overlapsAny(candidate, existing) {
			inUseCIDRs.add(candidate)
			return candidate.String()
		}
	}

	t.Fatalf("failed to generate a non-conflicting CIDR after %d attempts", maxCIDRGenerationAttempts)
	return "" // unreachable, but required by the compiler
}

// fetchInUseCIDRs retrieves all CIDR blocks from clusters across every project in
// the organisation. This ensures that generated CIDRs don't conflict with clusters
// in any project, not just the test project.
func fetchInUseCIDRs(t *testing.T) ([]netip.Prefix, error) {
	t.Helper()

	client := newTestClient(t)

	orgUUID, err := uuid.Parse(globalOrgId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse organization_id: %v", err)
	}

	projects, err := fetchAllProjects(t, client.ClientV2, orgUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %v", err)
	}

	var cidrs []netip.Prefix
	for _, project := range projects {
		clusters, err := fetchAllClusters(t, client.ClientV2, orgUUID, project.Id)
		if err != nil {
			t.Logf("Failed to fetch clusters for project %s: %v", project.Id, err)
			continue
		}

		for _, cluster := range clusters {
			if cluster.CloudProvider.Cidr == nil || *cluster.CloudProvider.Cidr == "" {
				continue
			}

			prefix, err := netip.ParsePrefix(*cluster.CloudProvider.Cidr)
			if err != nil {
				t.Logf("Could not parse CIDR %q from cluster %s: %v", *cluster.CloudProvider.Cidr, cluster.Id, err)
				continue
			}
			cidrs = append(cidrs, prefix)
		}
	}

	return cidrs, nil
}

// fetchAllProjects returns every project in the organisation, paginating through
// all pages of results.
func fetchAllProjects(t *testing.T, client *apigen.ClientWithResponses, orgID uuid.UUID) ([]apigen.GetProjectResponse, error) {
	t.Helper()

	var projects []apigen.GetProjectResponse
	page := 1

	for {
		resp, err := client.ListProjectsWithResponse(t.Context(), orgID, &apigen.ListProjectsParams{
			Page:    &page,
			PerPage: ptr.To(paginationPerPage),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list projects (page %d): %v", page, err)
		}
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("unexpected response when listing projects (page %d): status %s", page, resp.Status())
		}

		projects = append(projects, resp.JSON200.Data...)

		if resp.JSON200.Cursor.Pages.Next == nil || *resp.JSON200.Cursor.Pages.Next == 0 {
			break
		}
		page = *resp.JSON200.Cursor.Pages.Next
	}

	return projects, nil
}

// fetchAllClusters returns every cluster in a project, paginating through all
// pages of results.
func fetchAllClusters(t *testing.T, client *apigen.ClientWithResponses, orgID, projectID uuid.UUID) ([]apigen.GetClusterResponse, error) {
	t.Helper()

	var clusters []apigen.GetClusterResponse
	page := 1

	for {
		resp, err := client.ListClustersWithResponse(t.Context(), orgID, projectID, &apigen.ListClustersParams{
			Page:    &page,
			PerPage: ptr.To(paginationPerPage),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list clusters for project %s (page %d): %v", projectID, page, err)
		}
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("unexpected response when listing clusters for project %s (page %d): status %s", projectID, page, resp.Status())
		}

		clusters = append(clusters, resp.JSON200.Data...)

		if resp.JSON200.Cursor.Pages.Next == nil || *resp.JSON200.Cursor.Pages.Next == 0 {
			break
		}
		page = *resp.JSON200.Cursor.Pages.Next
	}

	return clusters, nil
}

// overlapsAny returns true if the candidate prefix overlaps with any of the
// existing prefixes.
func overlapsAny(candidate netip.Prefix, existing []netip.Prefix) bool {
	for _, e := range existing {
		if candidate.Overlaps(e) {
			return true
		}
	}
	return false
}
