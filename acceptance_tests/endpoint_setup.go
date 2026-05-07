package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
)

// fixtBucketState holds a once-and-error pair for a single pre-created bucket.
type fixtBucketState struct {
	once sync.Once
	err  error
}

var (
	fixtBucketMu     sync.Mutex
	fixtBucketStates = map[string]*fixtBucketState{}
)

// ensureFixtureBucketByName creates the named bucket exactly once per test
// process. Safe to call from parallel tests.
func ensureFixtureBucketByName(t *testing.T, name string) {
	t.Helper()
	fixtBucketMu.Lock()
	state, ok := fixtBucketStates[name]
	if !ok {
		state = &fixtBucketState{}
		fixtBucketStates[name] = state
	}
	// Unlock immediately — the mutex only guards the map. Bucket creation is
	// serialised per-name by state.once (sync.Once), so parallel tests for
	// different buckets are not blocked behind each other.
	fixtBucketMu.Unlock()

	state.once.Do(func() {
		ctx := context.Background()
		state.err = createFixtureBucket(ctx, globalClient, name)
	})
	if state.err != nil {
		t.Fatalf("failed to provision test bucket %s: %v", name, state.err)
	}

	t.Cleanup(func() {
		if err := deleteFixtureBucket(context.Background(), globalClient, name); err != nil {
			t.Logf("warning: failed to delete fixture bucket %q: %v", name, err)
		}
	})
}

// fixtEndpointState holds a once-and-error pair for a single pre-created endpoint.
type fixtEndpointState struct {
	once sync.Once
	err  error
}

var (
	fixtEndpointMu     sync.Mutex
	fixtEndpointStates = map[string]*fixtEndpointState{}
)

// ensureFixtureEndpoint creates the named bucket and app endpoint exactly once
// per test process. Safe to call from parallel tests.
func ensureFixtureEndpoint(t *testing.T, endpointName, bucketName, description string) {
	t.Helper()
	fixtEndpointMu.Lock()
	state, ok := fixtEndpointStates[endpointName]
	if !ok {
		state = &fixtEndpointState{}
		fixtEndpointStates[endpointName] = state
	}
	fixtEndpointMu.Unlock()

	state.once.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, bucketName); err != nil {
			state.err = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, endpointName, bucketName); err != nil {
			state.err = err
			return
		}
		state.err = appEndpointWait(ctx, globalClient, endpointName)
	})
	if state.err != nil {
		t.Fatalf("failed to provision %s test endpoint: %v", description, state.err)
	}
}

// createFixtureBucket creates the named bucket if it does not already exist,
// then waits until the bucket is available. Each fixture endpoint needs its own
// bucket because Capella only permits one endpoint per bucket/scope/collection.
func createFixtureBucket(ctx context.Context, client *api.Client, name string) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return err
	}
	for _, b := range buckets {
		if b.Name == name {
			log.Printf("fixture bucket %q already exists", name)
			return waitForFixtureBucket(ctx, client, b.Id)
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, globalClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, bucketapi.CreateBucketRequest{Name: name}, globalToken, nil)
	if err != nil {
		return err
	}

	var bucketResp bucketapi.GetBucketResponse
	if err = json.Unmarshal(response.Body, &bucketResp); err != nil {
		return err
	}

	return waitForFixtureBucket(ctx, client, bucketResp.Id)
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

// deleteFixtureBucket looks up the named bucket by name and deletes it. It is a
// no-op if the bucket does not exist.
func deleteFixtureBucket(ctx context.Context, client *api.Client, name string) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return fmt.Errorf("listing buckets for deletion of %q: %w", name, err)
	}

	for _, b := range buckets {
		if b.Name == name {
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
				globalHost, globalOrgId, globalProjectId, globalClusterId, b.Id)
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

func waitForFixtureBucket(ctx context.Context, client *api.Client, bucketID string) error {
	const maxWait = 5 * time.Minute
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
			globalHost, globalOrgId, globalProjectId, globalClusterId, bucketID)
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
