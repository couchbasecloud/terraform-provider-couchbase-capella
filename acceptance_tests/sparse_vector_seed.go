package acceptance_tests

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

var (
	sparseVectorSeedOnce sync.Once
	sparseVectorSeedErr  error
)

// ensureSparseVectorSeedDoc writes one document with a sparse vector into the
// global collection. Sparse vector indexes fail to build on an empty keyspace,
// and the provider has no document resource, so the doc goes through the
// Query Service with a temporary allowlist entry and database credential.
func ensureSparseVectorSeedDoc(ctx context.Context) error {
	sparseVectorSeedOnce.Do(func() {
		sparseVectorSeedErr = seedSparseVectorDoc(ctx, globalClient)
	})
	return sparseVectorSeedErr
}

func seedSparseVectorDoc(ctx context.Context, client *api.Client) error {
	clusterUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, globalClusterId)

	allowlistId, err := createSparseVectorAllowlist(ctx, client, clusterUrl)
	if err != nil {
		return fmt.Errorf("allow-listing test client: %w", err)
	}
	if allowlistId != "" {
		defer deleteSparseVectorFixture(ctx, client, clusterUrl+"/allowedcidrs/"+allowlistId)
	}

	credName := randomStringWithPrefix("tf_acc_sv_seed_")
	cfg := api.EndpointCfg{Url: clusterUrl + "/users", Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, api.CreateDatabaseCredentialRequest{
		Name:   credName,
		Access: []api.Access{{Privileges: []string{"data_reader", "data_writer"}}},
	}, globalToken, nil)
	if err != nil {
		return fmt.Errorf("creating seed credential: %w", err)
	}
	var cred api.CreateDatabaseCredentialResponse
	if err = json.Unmarshal(response.Body, &cred); err != nil {
		return err
	}
	defer deleteSparseVectorFixture(ctx, client, clusterUrl+"/users/"+cred.Id.String())

	cfg = api.EndpointCfg{Url: clusterUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err = client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return err
	}
	var cluster clusterapi.GetClusterResponse
	if err = json.Unmarshal(response.Body, &cluster); err != nil {
		return err
	}
	queryUrl := fmt.Sprintf("https://%s:18093/query/service", strings.TrimPrefix(cluster.ConnectionString, "couchbases://"))

	cfg = api.EndpointCfg{Url: clusterUrl + "/certificates", Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err = client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return fmt.Errorf("fetching cluster certificate: %w", err)
	}
	var cert api.GetCertificateResponse
	if err = json.Unmarshal(response.Body, &cert); err != nil {
		return err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(cert.Certificate)) {
		return errors.New("cluster certificate is not valid PEM")
	}
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}},
	}

	statement := fmt.Sprintf(
		"UPSERT INTO `%s`.`%s`.`%s` (KEY, VALUE) VALUES (\"tf_acc_sv_seed_doc\", {\"type\": \"sv_seed\", \"sv\": [[3,17,42],[0.8,0.5,0.2]]})",
		globalBucketName, globalScopeName, globalCollectionName)

	deadline := time.Now().Add(5 * time.Minute)
	for {
		err = runSparseVectorQuery(ctx, httpClient, queryUrl, credName, cred.Password, statement)
		if err == nil {
			log.Print("seeded sparse vector document")
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("seeding sparse vector document: %w", err)
		}
		log.Printf("sparse vector seed not ready yet, retrying: %v", err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
		}
	}
}

// createSparseVectorAllowlist returns "" when the client CIDR is already
// allow-listed, so an entry that the suite did not create is left alone.
func createSparseVectorAllowlist(ctx context.Context, client *api.Client, clusterUrl string) (string, error) {
	cidr, err := sparseVectorClientCidr(ctx)
	if err != nil {
		return "", err
	}

	listCfg := api.EndpointCfg{Url: clusterUrl + "/allowedcidrs", Method: http.MethodGet, SuccessStatus: http.StatusOK}
	entries, err := api.GetPaginated[[]api.GetAllowListResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Cidr == cidr {
			return "", nil
		}
	}

	createCfg := api.EndpointCfg{Url: clusterUrl + "/allowedcidrs", Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, createCfg, api.CreateAllowListRequest{
		Cidr:      cidr,
		Comment:   "tf acc sparse vector seed",
		ExpiresAt: time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
	}, globalToken, nil)
	if err != nil {
		return "", err
	}
	var created api.CreateAllowListResponse
	if err = json.Unmarshal(response.Body, &created); err != nil {
		return "", err
	}
	return created.Id.String(), nil
}

func sparseVectorClientCidr(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("determining public IP: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)) + "/32", nil
}

func runSparseVectorQuery(ctx context.Context, httpClient *http.Client, queryUrl, user, password, statement string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, queryUrl, strings.NewReader(url.Values{"statement": {statement}}.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result struct {
		Errors []struct {
			Msg string `json:"msg"`
		} `json:"errors"`
	}
	_ = json.Unmarshal(body, &result)
	if len(result.Errors) > 0 {
		return errors.New(result.Errors[0].Msg)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func deleteSparseVectorFixture(ctx context.Context, client *api.Client, resourceUrl string) {
	cfg := api.EndpointCfg{Url: resourceUrl, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	if _, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil); err != nil {
		log.Printf("deleting sparse vector seed fixture %s: %v", resourceUrl, err)
	}
}
