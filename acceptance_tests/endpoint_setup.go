package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
)

// fixtCollectionState holds a once-and-error pair for a single pre-created
// collection in the shared app-endpoint bucket.
type fixtCollectionState struct {
	once sync.Once
	err  error
}

var (
	fixtCollectionMu     sync.Mutex
	fixtCollectionStates = map[string]*fixtCollectionState{}
)

// ensureFixtureCollection creates the named collection in the shared
// app-endpoint bucket's _default scope exactly once per test process. App
// endpoint resource tests use one collection each instead of a dedicated
// bucket: Capella allows only one endpoint per bucket/scope/collection, and a
// cluster's bucket cap (1 per 0.2 cores) is far smaller than its collection
// cap. Sharing the bucket keeps parallel tests well under the bucket limit that
// previously caused flaky "maximum number of buckets reached" failures.
//
// The collection lives in the shared bucket and is torn down with it at the end
// of the suite, so no per-test cleanup is registered. Safe to call from
// parallel tests.
func ensureFixtureCollection(t *testing.T, name string) {
	t.Helper()
	ensureAppEndpointTestEnvironment(t)
	fixtCollectionMu.Lock()
	state, ok := fixtCollectionStates[name]
	if !ok {
		state = &fixtCollectionState{}
		fixtCollectionStates[name] = state
	}
	// Unlock immediately — the mutex only guards the map. Collection creation is
	// serialised per-name by state.once, so parallel tests for different
	// collections are not blocked behind each other.
	fixtCollectionMu.Unlock()

	state.once.Do(func() {
		ctx := context.Background()
		state.err = createFixtureCollection(ctx, globalClient, appEndpointBucketId, globalScopeName, name)
	})
	if state.err != nil {
		t.Fatalf("failed to provision test collection %s: %v", name, state.err)
	}
}

// fixtEndpointState holds a once-and-error pair for a single pre-created endpoint.
type fixtEndpointState struct {
	once            sync.Once
	err             error
	bucketCreated   bool
	endpointCreated bool
}

var (
	fixtEndpointMu     sync.Mutex
	fixtEndpointStates = map[string]*fixtEndpointState{}
)

// ensureFixtureEndpoint creates the named bucket and app endpoint exactly once
// per test process. Safe to call from parallel tests.
func ensureFixtureEndpoint(t *testing.T, endpointName, bucketName, description string) {
	t.Helper()
	ensureAppEndpointTestEnvironment(t)
	fixtEndpointMu.Lock()
	state, ok := fixtEndpointStates[endpointName]
	if !ok {
		state = &fixtEndpointState{}
		fixtEndpointStates[endpointName] = state
	}
	fixtEndpointMu.Unlock()

	state.once.Do(func() {
		ctx := context.Background()
		bucketCreated, err := createFixtureBucket(ctx, globalClient, bucketName)
		if err != nil {
			state.err = err
			return
		}
		state.bucketCreated = bucketCreated
		endpointCreated, err := createAppEndpointForAppService(ctx, globalClient, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, endpointName, bucketName)
		if err != nil {
			state.err = err
			return
		}
		state.endpointCreated = endpointCreated
		state.err = appEndpointWaitForAppService(ctx, globalClient, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, endpointName)
	})

	if state.bucketCreated {
		t.Cleanup(func() {
			if err := deleteAppEndpointFixtureBucket(context.Background(), globalClient, globalProjectId, appEndpointClusterId, bucketName); err != nil {
				t.Logf("warning: failed to delete %s fixture bucket %q: %v", description, bucketName, err)
			}
		})
	}
	if state.endpointCreated {
		t.Cleanup(func() {
			if err := deleteAppEndpointFixtureEndpoint(context.Background(), globalClient, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, endpointName); err != nil {
				t.Logf("warning: failed to delete %s fixture endpoint %q: %v", description, endpointName, err)
			}
		})
	}
	if state.err != nil {
		t.Fatalf("failed to provision %s test endpoint: %v", description, state.err)
	}
}

// createFixtureBucket creates the named bucket if it does not already exist,
// then waits until the bucket is available. Returns true if the bucket was
// newly created by this call, false if it already existed.
func createFixtureBucket(ctx context.Context, client *api.Client, name string) (bool, error) {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return false, err
	}
	for _, b := range buckets {
		if b.Name == name {
			log.Printf("fixture bucket %q already exists", name)
			return false, waitForFixtureBucket(ctx, client, b.Id)
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, bucketapi.CreateBucketRequest{Name: name}, globalToken, nil)
	if err != nil {
		return false, err
	}

	var bucketResp bucketapi.GetBucketResponse
	if err = json.Unmarshal(response.Body, &bucketResp); err != nil {
		return false, err
	}

	return true, waitForFixtureBucket(ctx, client, bucketResp.Id)
}

// createFixtureCollection creates the named collection in the given scope of the
// shared bucket if it does not already exist, then waits until it is listable.
func createFixtureCollection(ctx context.Context, client *api.Client, bucketID, scope, name string) error {
	collectionsUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections",
		globalHost, globalOrgId, globalProjectId, appEndpointClusterId, bucketID, scope)

	exists, err := fixtureCollectionExists(ctx, client, collectionsUrl, name)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("fixture collection %q already exists", name)
		return nil
	}

	createCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	if _, err = client.ExecuteWithRetry(ctx, createCfg, api.CreateCollectionRequest{Name: name}, globalToken, nil); err != nil {
		return fmt.Errorf("creating fixture collection %q: %w", name, err)
	}

	return waitForFixtureCollection(ctx, client, collectionsUrl, name)
}

// fixtureCollectionExists reports whether a collection with the given name is
// present at the collections list endpoint.
func fixtureCollectionExists(ctx context.Context, client *api.Client, collectionsUrl, name string) (bool, error) {
	listCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := client.ExecuteWithRetry(ctx, listCfg, nil, globalToken, nil)
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

// waitForFixtureCollection polls until the collection is listable. App endpoint
// creation rejects a collection that is not yet visible ("Collection Not
// Found"), so fixtures wait for it before pointing endpoints at it.
func waitForFixtureCollection(ctx context.Context, client *api.Client, collectionsUrl, name string) error {
	const maxWait = 2 * time.Minute
	deadline := time.Now().Add(maxWait)

	for {
		exists, err := fixtureCollectionExists(ctx, client, collectionsUrl, name)
		if err == nil && exists {
			return nil
		}
		if time.Now().After(deadline) {
			if err != nil {
				return fmt.Errorf("timeout waiting for fixture collection %q: %w", name, err)
			}
			return fmt.Errorf("timeout waiting for fixture collection %q to become listable", name)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done waiting for fixture collection %q: %w", name, ctx.Err())
		case <-time.After(5 * time.Second):
		}
	}
}

// deleteBucketWithRetry retries the bucket DELETE on 412 (App Service not yet
// reachable) until success or the 5-minute deadline is reached.
func deleteBucketWithRetry(ctx context.Context, client *api.Client, cfg api.EndpointCfg, name string) error {
	const maxWait = 5 * time.Minute
	const retryInterval = 30 * time.Second
	deadline := time.Now().Add(maxWait)

	for {
		_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err == nil {
			return nil
		}
		var apiErr *api.Error
		if !errors.As(err, &apiErr) || apiErr.HttpStatusCode != http.StatusPreconditionFailed {
			return fmt.Errorf("deleting fixture bucket %q: %w", name, err)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout deleting fixture bucket %q: %w", name, err)
		}
		log.Printf("App Service not yet reachable, retrying bucket deletion %q in %s", name, retryInterval)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done while deleting bucket %q: %w", name, ctx.Err())
		case <-time.After(retryInterval):
		}
	}
}

func deleteAppEndpointFixtureBucket(ctx context.Context, client *api.Client, projectID, clusterID, name string) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, projectID, clusterID)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return fmt.Errorf("listing buckets for deletion of %q: %w", name, err)
	}

	for _, b := range buckets {
		if b.Name == name {
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
				globalHost, globalOrgId, projectID, clusterID, b.Id)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
			if err = deleteBucketWithRetry(ctx, client, cfg, name); err != nil {
				return err
			}
			log.Printf("deleted fixture bucket %q", name)
			return nil
		}
	}

	log.Printf("fixture bucket %q not found, skipping deletion", name)
	return nil
}

func waitForAppEndpointFixtureEndpointDeletion(ctx context.Context, client *api.Client, projectID, clusterID, appServiceID, name string) error {
	const maxWait = 5 * time.Minute
	const retryInterval = 10 * time.Second
	deadline := time.Now().Add(maxWait)

	for {
		endpointURL := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
			globalHost,
			globalOrgId,
			projectID,
			clusterID,
			appServiceID,
			url.PathEscape(name),
		)
		cfg := api.EndpointCfg{Url: endpointURL, Method: http.MethodGet, SuccessStatus: http.StatusOK}

		_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err != nil {
			var apiErr *api.Error
			if errors.As(err, &apiErr) &&
				(apiErr.HttpStatusCode == http.StatusNotFound || apiErr.HttpStatusCode == http.StatusForbidden) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for fixture endpoint %q deletion", name)
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done waiting for fixture endpoint %q deletion: %w", name, ctx.Err())
		case <-time.After(retryInterval):
		}
	}
}

// deleteFixtureEndpoint deletes the named app endpoint and waits until the API
// reports it as gone. It is a no-op if the endpoint does not exist.
func deleteFixtureEndpoint(ctx context.Context, client *api.Client, name string) error {
	return deleteAppEndpointFixtureEndpoint(ctx, client, globalProjectId, globalClusterId, globalAppServiceId, name)
}

func deleteAppEndpointFixtureEndpoint(ctx context.Context, client *api.Client, projectID, clusterID, appServiceID, name string) error {
	endpointURL := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		globalHost,
		globalOrgId,
		projectID,
		clusterID,
		appServiceID,
		url.PathEscape(name),
	)
	cfg := api.EndpointCfg{Url: endpointURL, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}

	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) &&
			(apiErr.HttpStatusCode == http.StatusNotFound || apiErr.HttpStatusCode == http.StatusForbidden) {
			log.Printf("fixture endpoint %q not found, skipping deletion", name)
			return nil
		}
		return fmt.Errorf("deleting fixture endpoint %q: %w", name, err)
	}

	if err = waitForAppEndpointFixtureEndpointDeletion(ctx, client, projectID, clusterID, appServiceID, name); err != nil {
		return err
	}
	log.Printf("deleted fixture endpoint %q", name)
	return nil
}

func waitForFixtureBucket(ctx context.Context, client *api.Client, bucketID string) error {
	const maxWait = 5 * time.Minute
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
			globalHost, globalOrgId, globalProjectId, appEndpointClusterId, bucketID)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done waiting for bucket: %w", ctx.Err())
		case <-time.After(30 * time.Second):
		}
	}

	return fmt.Errorf("timeout waiting for fixture bucket %s", bucketID)
}

func ensureACFEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalACFEndpointName, globalACFBucketName, "ACF")
}

func ensureIFEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalIFEndpointName, globalIFBucketName, "ImportFilter")
}

func ensureCORSEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalCORSEndpointName, globalCORSBucketName, "CORS")
}

func ensureCORSOriginOnlyEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalCORSOriginOnlyEndpointName, globalCORSOriginOnlyBucketName, "CORS origin-only")
}

func ensureOIDCEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalOIDCEndpointName, globalOIDCBucketName, "OIDC")
}

func ensureDefaultOIDCEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, globalDefaultOIDCEndpointName, globalDefaultOIDCBucketName, "default OIDC")
}

func ensureActivationEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, appEndpointActivationEndpointName, appEndpointActivationBucketName, "activation status")
}

func ensureLoggingEndpoint(t *testing.T) {
	t.Helper()
	ensureFixtureEndpoint(t, appEndpointLoggingEndpointName, appEndpointLoggingBucketName, "logging config")
}
