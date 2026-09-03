package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseCredentialWithReqFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithReqFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
				),
			},
		},
	})
}

func TestAccDatabaseCredentialWithOptionalFields(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "password", "Secret12$#"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.0", "data_writer"),
				),
			},
		},
	})
}

// databaseCredentialFixtureCollections are created (idempotently) in the
// shared bucket's _default scope so TestAccDatabaseCredentialWithBucketScopeAccess
// has real, multi-element collections to reference without depending on any
// externally pre-loaded sample bucket.
var databaseCredentialFixtureCollections = []string{"tf_acc_dbcred_col_a", "tf_acc_dbcred_col_b"}

// TestAccDatabaseCredentialWithBucketScopeAccess covers the nested
// access.resources.buckets.scopes.collections shape from AV-139841/CBSE-23356
// (List/ListNestedAttribute instead of Set/SetNestedAttribute), verifying the
// resource applies cleanly against real bucket/scope/collection names and
// produces no diff on a subsequent plan.
func TestAccDatabaseCredentialWithBucketScopeAccess(t *testing.T) {
	ensureDatabaseCredentialFixtureCollections(t)

	resourceName := randomStringWithPrefix("tf_acc_database_credential_")
	resourceReference := "couchbase-capella_database_credential." + resourceName
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAddDatabaseCredWithBucketScopeAccessConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.privileges.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_reader"),
					resource.TestCheckTypeSetElemAttr(resourceReference, "access.0.privileges.*", "data_writer"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.name", globalBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.name", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.0.scopes.0.collections.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.1.name", globalMetadataBucketName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.1.scopes.0.name", globalScopeName),
					resource.TestCheckResourceAttr(resourceReference, "access.0.resources.buckets.1.scopes.0.collections.0", globalScopeName),
				),
			},
		},
	})
}

func testAccAddDatabaseCredWithBucketScopeAccessConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access = [
				{
					privileges = ["data_reader", "data_writer"]
					resources = {
						buckets = [
							{
								name = "%[6]s"
								scopes = [
									{
										name        = "%[7]s"
										collections = ["%[7]s", "%[8]s", "%[9]s"]
									},
								]
							},
							{
								name = "%[10]s"
								scopes = [
									{
										name        = "%[7]s"
										collections = ["%[7]s"]
									},
								]
							},
						]
					}
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName,
		globalBucketName, globalScopeName,
		databaseCredentialFixtureCollections[0], databaseCredentialFixtureCollections[1],
		globalMetadataBucketName)
}

// ensureDatabaseCredentialFixtureCollections creates databaseCredentialFixtureCollections
// in the shared bucket's _default scope if they don't already exist, then waits
// until each is listable. Capella rejects a database credential access grant
// referencing a collection that isn't visible yet, so callers must wait for this
// before applying config that references them.
func ensureDatabaseCredentialFixtureCollections(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	collectionsUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections",
		globalHost, globalOrgId, globalProjectId, globalClusterId, globalBucketId, globalScopeName)

	for _, name := range databaseCredentialFixtureCollections {
		if err := ensureDatabaseCredentialFixtureCollection(ctx, collectionsUrl, name); err != nil {
			t.Fatalf("failed to provision test collection %s: %v", name, err)
		}
	}
}

func ensureDatabaseCredentialFixtureCollection(ctx context.Context, collectionsUrl, name string) error {
	exists, err := databaseCredentialFixtureCollectionExists(ctx, collectionsUrl, name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	createCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	if _, err = globalClient.ExecuteWithRetry(ctx, createCfg, api.CreateCollectionRequest{Name: name}, globalToken, nil); err != nil {
		var apiErr *api.Error
		if !errors.As(err, &apiErr) || apiErr.HttpStatusCode != http.StatusConflict {
			return fmt.Errorf("creating fixture collection %q: %w", name, err)
		}
	}

	const maxWait = 2 * time.Minute
	deadline := time.Now().Add(maxWait)
	for {
		exists, err = databaseCredentialFixtureCollectionExists(ctx, collectionsUrl, name)
		if err == nil && exists {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for fixture collection %q to become listable", name)
		}
		time.Sleep(5 * time.Second)
	}
}

func databaseCredentialFixtureCollectionExists(ctx context.Context, collectionsUrl, name string) (bool, error) {
	listCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := globalClient.ExecuteWithRetry(ctx, listCfg, nil, globalToken, nil)
	if err != nil {
		return false, fmt.Errorf("listing collections for %q: %w", name, err)
	}

	var collections api.GetCollectionsResponse
	if err = json.Unmarshal(response.Body, &collections); err != nil {
		return false, fmt.Errorf("unmarshalling collections list for %q: %w", name, err)
	}
	for _, c := range collections.Data {
		if c.Name != nil && *c.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func testAccAddDatabaseCredWithReqFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s

		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}

func testAccAddDatabaseCredWithOptionalFieldsConfig(resourceName string) string {
	return fmt.Sprintf(
		`
		%[1]s
		resource "couchbase-capella_database_credential" "%[5]s" {
			name            = "%[5]s"
			organization_id = "%[2]s"
			project_id      = "%[3]s"
			cluster_id      = "%[4]s"
			password        = "Secret12$#"
			access = [
				{
					privileges = ["data_writer"]
				},
			]
		}
		`, globalProviderBlock, globalOrgId, globalProjectId, globalClusterId, resourceName)
}
