package acceptance_tests

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	scopeapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/scope"
)

var (
	sparseVectorOnce sync.Once
	sparseVectorErr  error
)

// requireSparseVectorCluster skips the test if no pre-provisioned SPARSE
// VECTOR-capable cluster is configured, then lazily provisions the shared
// bucket/scope/collection, client network access, a database credential, and
// a training document on first use. Sparse vector indexes need >=1 document
// to train on and the provider has no data-plane document resource, so the
// doc is seeded directly via the Query Service REST API.
func requireSparseVectorCluster(t *testing.T) {
	t.Helper()
	if sparseVectorClusterId == "" {
		t.Skip("TF_VAR_sparse_vector_cluster_id not set; sparse vector index tests need a pre-provisioned 8.5.0+ cluster (automated 8.5.0 cluster creation is blocked by AV-145087)")
	}
	sparseVectorOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		client := api.NewClient(timeout)
		sparseVectorErr = ensureSparseVectorFixture(ctx, client)
	})
	if sparseVectorErr != nil {
		t.Fatalf("ensureSparseVectorFixture: %v", sparseVectorErr)
	}
}

func ensureSparseVectorFixture(ctx context.Context, client *api.Client) error {
	if err := resolveSparseVectorBucket(ctx, client); err != nil {
		return fmt.Errorf("resolving sparse vector bucket: %w", err)
	}
	if err := sparseVectorBucketWait(ctx, client); err != nil {
		return fmt.Errorf("waiting for sparse vector bucket: %w", err)
	}
	if err := ensureSparseVectorScope(ctx, client); err != nil {
		return fmt.Errorf("creating sparse vector scope: %w", err)
	}
	if err := ensureSparseVectorCollection(ctx, client); err != nil {
		return fmt.Errorf("creating sparse vector collection: %w", err)
	}
	if err := ensureSparseVectorAllowlist(ctx, client); err != nil {
		return fmt.Errorf("allow-listing sparse vector test client: %w", err)
	}
	if err := createSparseVectorCredential(ctx, client); err != nil {
		return fmt.Errorf("creating sparse vector database credential: %w", err)
	}
	if err := seedSparseVectorDoc(ctx, client); err != nil {
		return fmt.Errorf("seeding sparse vector training document: %w", err)
	}
	return nil
}

func resolveSparseVectorBucket(ctx context.Context, client *api.Client) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, sparseVectorClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return err
	}
	for _, bucket := range buckets {
		if bucket.Name == sparseVectorBucketName {
			sparseVectorBucketId = bucket.Id
			log.Printf("Discovered existing sparse vector bucket: %s (%s)", sparseVectorBucketName, sparseVectorBucketId)
			return nil
		}
	}

	createCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, createCfg, bucketapi.CreateBucketRequest{
		Name:           sparseVectorBucketName,
		StorageBackend: ptr.To("magma"),
	}, globalToken, nil)
	if err != nil {
		return err
	}

	var bucketResp bucketapi.GetBucketResponse
	if err = json.Unmarshal(response.Body, &bucketResp); err != nil {
		return err
	}
	sparseVectorBucketId = bucketResp.Id
	sparseVectorBucketCreated = true
	log.Printf("Created sparse vector bucket: %s (%s)", sparseVectorBucketName, sparseVectorBucketId)
	return nil
}

func sparseVectorBucketWait(ctx context.Context, client *api.Client) error {
	if !sparseVectorBucketCreated {
		return nil
	}

	const maxWaitTime = 5 * time.Minute
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeoutWaitingForBucket
		case <-ticker.C:
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
				globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorBucketId)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
			if err == nil {
				log.Print("sparse vector bucket ready")
				return nil
			}
			var apiError *api.Error
			if !errors.As(err, &apiError) {
				return err
			}
			if apiError.HttpStatusCode != http.StatusNotFound {
				return err
			}
		}
	}
}

func ensureSparseVectorScope(ctx context.Context, client *api.Client) error {
	scopesUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes",
		globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorBucketId)

	listCfg := api.EndpointCfg{Url: scopesUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := client.ExecuteWithRetry(ctx, listCfg, nil, globalToken, nil)
	if err != nil {
		return err
	}
	var scopesResp scopeapi.GetScopesResponse
	if err := json.Unmarshal(response.Body, &scopesResp); err != nil {
		return err
	}
	for _, s := range scopesResp.Scopes {
		if s.Name != nil && *s.Name == sparseVectorScopeName {
			log.Printf("sparse vector scope %q already exists", sparseVectorScopeName)
			return nil
		}
	}

	createCfg := api.EndpointCfg{Url: scopesUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	_, err = client.ExecuteWithRetry(ctx, createCfg, scopeapi.CreateScopeRequest{Name: sparseVectorScopeName}, globalToken, nil)
	return err
}

func ensureSparseVectorCollection(ctx context.Context, client *api.Client) error {
	collectionsUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s/scopes/%s/collections",
		globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorBucketId, sparseVectorScopeName)

	listCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := client.ExecuteWithRetry(ctx, listCfg, nil, globalToken, nil)
	if err != nil {
		return err
	}
	var collectionsResp api.GetCollectionsResponse
	if err := json.Unmarshal(response.Body, &collectionsResp); err != nil {
		return err
	}
	for _, c := range collectionsResp.Data {
		if c.Name != nil && *c.Name == sparseVectorCollectionName {
			log.Printf("sparse vector collection %q already exists", sparseVectorCollectionName)
			return nil
		}
	}

	createCfg := api.EndpointCfg{Url: collectionsUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	_, err = client.ExecuteWithRetry(ctx, createCfg, api.CreateCollectionRequest{Name: sparseVectorCollectionName}, globalToken, nil)
	return err
}

// ensureSparseVectorAllowlist allow-lists this test runner's public IP for
// direct Query Service access - reachable only via the cluster's data-port
// allowlist, which is separate from the management-API access every other
// helper here uses.
func ensureSparseVectorAllowlist(ctx context.Context, client *api.Client) error {
	cidr, err := sparseVectorClientCidr(ctx)
	if err != nil {
		return err
	}

	allowlistUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs",
		globalHost, globalOrgId, globalProjectId, sparseVectorClusterId)
	listCfg := api.EndpointCfg{Url: allowlistUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	entries, err := api.GetPaginated[[]api.GetAllowListResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Cidr == cidr {
			log.Printf("sparse vector test client CIDR already allow-listed: %s", cidr)
			return nil
		}
	}

	createCfg := api.EndpointCfg{Url: allowlistUrl, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, createCfg, api.CreateAllowListRequest{
		Cidr:    cidr,
		Comment: "sparse vector index acceptance tests - direct Query Service access for doc seeding",
	}, globalToken, nil)
	if err != nil {
		return err
	}
	var allowlistResp api.CreateAllowListResponse
	if err := json.Unmarshal(response.Body, &allowlistResp); err != nil {
		return err
	}
	sparseVectorAllowlistId = allowlistResp.Id.String()
	sparseVectorAllowlistCreated = true
	log.Printf("allow-listed sparse vector test client: %s", cidr)
	return nil
}

func sparseVectorClientCidr(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("determining test client public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)) + "/32", nil
}

func createSparseVectorCredential(ctx context.Context, client *api.Client) error {
	sparseVectorCredName = randomStringWithPrefix("tf_acc_sparse_vector_")

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users", globalHost, globalOrgId, globalProjectId, sparseVectorClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, api.CreateDatabaseCredentialRequest{
		Name: sparseVectorCredName,
		Access: []api.Access{
			{Privileges: []string{"data_reader", "data_writer"}},
		},
	}, globalToken, nil)
	if err != nil {
		return err
	}

	var credResp api.CreateDatabaseCredentialResponse
	if err := json.Unmarshal(response.Body, &credResp); err != nil {
		return err
	}
	sparseVectorCredId = credResp.Id.String()
	sparseVectorCredPassword = credResp.Password
	return nil
}

func fetchSparseVectorConnectionString(ctx context.Context, client *api.Client) (string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, sparseVectorClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return "", err
	}
	var clusterResp clusterapi.GetClusterResponse
	if err := json.Unmarshal(response.Body, &clusterResp); err != nil {
		return "", err
	}
	return clusterResp.ConnectionString, nil
}

// seedSparseVectorDoc inserts one training document directly via the Query
// Service REST API, since a sparse vector index needs >=1 document to train
// on and this provider has no data-plane document resource. UPSERT keeps
// this idempotent across repeated suite runs against the same persistent
// fixture collection.
func seedSparseVectorDoc(ctx context.Context, client *api.Client) error {
	connString, err := fetchSparseVectorConnectionString(ctx, client)
	if err != nil {
		return err
	}
	host := strings.TrimPrefix(connString, "couchbases://")

	statement := fmt.Sprintf(
		"UPSERT INTO `%s`.`%s`.`%s` (KEY, VALUE) VALUES (\"tf_acc_seed_doc\", {\"sv\": [[3,17,42],[0.8,0.5,0.2]]})",
		sparseVectorBucketName, sparseVectorScopeName, sparseVectorCollectionName,
	)

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			// This sandbox's Query Service cert isn't in the default trust
			// store; acceptable since this only seeds throwaway test data.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("https://%s:18093/query/service", host),
		strings.NewReader("statement="+statement))
	if err != nil {
		return err
	}
	req.SetBasicAuth(sparseVectorCredName, sparseVectorCredPassword)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Errors []struct {
			Msg string `json:"msg"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil && len(result.Errors) > 0 {
		return errors.New(result.Errors[0].Msg)
	}

	log.Print("seeded sparse vector training document")
	return nil
}

func destroySparseVectorFixture(ctx context.Context, client *api.Client) error {
	if sparseVectorCredId != "" {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/users/%s",
			globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorCredId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("deleting sparse vector database credential: %w", err)
		}
	}

	if sparseVectorAllowlistCreated {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/allowedcidrs/%s",
			globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorAllowlistId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("deleting sparse vector allowlist entry: %w", err)
		}
	}

	if sparseVectorBucketCreated {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s",
			globalHost, globalOrgId, globalProjectId, sparseVectorClusterId, sparseVectorBucketId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		if _, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil); err != nil {
			return fmt.Errorf("deleting sparse vector bucket: %w", err)
		}
	}

	return nil
}
