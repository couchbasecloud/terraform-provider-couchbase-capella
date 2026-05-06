package acceptance_tests

import (
	"context"
	"encoding/json"
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
	fixtBucketMu.Unlock()

	state.once.Do(func() {
		ctx := context.Background()
		state.err = createFixtureBucket(ctx, globalClient, name)
	})
	if state.err != nil {
		t.Fatalf("failed to provision test bucket %s: %v", name, state.err)
	}
}

var (
	acfEndpointOnce sync.Once
	acfEndpointErr  error

	ifEndpointOnce sync.Once
	ifEndpointErr  error

	corsEndpointOnce sync.Once
	corsEndpointErr  error

	corsOriginOnlyEndpointOnce sync.Once
	corsOriginOnlyEndpointErr  error

	oidcEndpointOnce sync.Once
	oidcEndpointErr  error

	defaultOIDCEndpointOnce sync.Once
	defaultOIDCEndpointErr  error
)

// createFixtureBucket creates the named bucket if it does not already exist,
// then waits until the bucket is available. Each fixture endpoint needs its own
// bucket because Capella only permits one endpoint per bucket/scope/collection.
func createFixtureBucket(ctx context.Context, client *api.Client, name string) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets",
		globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err == nil {
		for _, b := range buckets {
			if b.Name == name {
				log.Printf("fixture bucket %q already exists", name)
				return waitForFixtureBucket(ctx, client, b.Id)
			}
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
	acfEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalACFBucketName); err != nil {
			acfEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalACFEndpointName, globalACFBucketName); err != nil {
			acfEndpointErr = err
			return
		}
		acfEndpointErr = appEndpointWait(ctx, globalClient, globalACFEndpointName)
	})
	if acfEndpointErr != nil {
		t.Fatalf("failed to provision ACF test endpoint: %v", acfEndpointErr)
	}
}

func ensureIFEndpoint(t *testing.T) {
	t.Helper()
	ifEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalIFBucketName); err != nil {
			ifEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalIFEndpointName, globalIFBucketName); err != nil {
			ifEndpointErr = err
			return
		}
		ifEndpointErr = appEndpointWait(ctx, globalClient, globalIFEndpointName)
	})
	if ifEndpointErr != nil {
		t.Fatalf("failed to provision ImportFilter test endpoint: %v", ifEndpointErr)
	}
}

func ensureCORSEndpoint(t *testing.T) {
	t.Helper()
	corsEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalCORSBucketName); err != nil {
			corsEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalCORSEndpointName, globalCORSBucketName); err != nil {
			corsEndpointErr = err
			return
		}
		corsEndpointErr = appEndpointWait(ctx, globalClient, globalCORSEndpointName)
	})
	if corsEndpointErr != nil {
		t.Fatalf("failed to provision CORS test endpoint: %v", corsEndpointErr)
	}
}

func ensureCORSOriginOnlyEndpoint(t *testing.T) {
	t.Helper()
	corsOriginOnlyEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalCORSOriginOnlyBucketName); err != nil {
			corsOriginOnlyEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalCORSOriginOnlyEndpointName, globalCORSOriginOnlyBucketName); err != nil {
			corsOriginOnlyEndpointErr = err
			return
		}
		corsOriginOnlyEndpointErr = appEndpointWait(ctx, globalClient, globalCORSOriginOnlyEndpointName)
	})
	if corsOriginOnlyEndpointErr != nil {
		t.Fatalf("failed to provision CORS origin-only test endpoint: %v", corsOriginOnlyEndpointErr)
	}
}

func ensureOIDCEndpoint(t *testing.T) {
	t.Helper()
	oidcEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalOIDCBucketName); err != nil {
			oidcEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalOIDCEndpointName, globalOIDCBucketName); err != nil {
			oidcEndpointErr = err
			return
		}
		oidcEndpointErr = appEndpointWait(ctx, globalClient, globalOIDCEndpointName)
	})
	if oidcEndpointErr != nil {
		t.Fatalf("failed to provision OIDC test endpoint: %v", oidcEndpointErr)
	}
}

func ensureDefaultOIDCEndpoint(t *testing.T) {
	t.Helper()
	defaultOIDCEndpointOnce.Do(func() {
		ctx := context.Background()
		if err := createFixtureBucket(ctx, globalClient, globalDefaultOIDCBucketName); err != nil {
			defaultOIDCEndpointErr = err
			return
		}
		if err := createAppEndpoint(ctx, globalClient, globalDefaultOIDCEndpointName, globalDefaultOIDCBucketName); err != nil {
			defaultOIDCEndpointErr = err
			return
		}
		defaultOIDCEndpointErr = appEndpointWait(ctx, globalClient, globalDefaultOIDCEndpointName)
	})
	if defaultOIDCEndpointErr != nil {
		t.Fatalf("failed to provision default OIDC test endpoint: %v", defaultOIDCEndpointErr)
	}
}
