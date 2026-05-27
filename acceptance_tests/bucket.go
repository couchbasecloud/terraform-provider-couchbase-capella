package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
)

// fetchBucket returns the bucket name for the given ID, or ("", notFound) if
// the bucket does not exist.
func fetchBucket(ctx context.Context, client *api.Client, bucketId string) (name string, found bool, err error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, globalClusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, apiErr := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if apiErr != nil {
		var apiErrTyped *api.Error
		if errors.As(apiErr, &apiErrTyped) && apiErrTyped.HttpStatusCode == http.StatusNotFound {
			return "", false, nil
		}
		return "", false, apiErr
	}
	var bucket bucketapi.GetBucketResponse
	if err = json.Unmarshal(response.Body, &bucket); err != nil {
		return "", false, err
	}
	return bucket.Name, true, nil
}

func discoverFirstBucket(ctx context.Context, client *api.Client) (string, string, error) {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return "", "", err
	}
	for _, bucket := range buckets {
		if bucket.Name == globalBucketName {
			return bucket.Id, bucket.Name, nil
		}
	}
	if len(buckets) > 0 {
		return buckets[0].Id, buckets[0].Name, nil
	}
	return "", "", nil
}

func createBucket(ctx context.Context, client *api.Client) error {
	bucketRequest := bucketapi.CreateBucketRequest{
		Name: globalBucketName,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, globalClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		bucketRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	bucketResponse := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucketResponse); err != nil {
		return err
	}

	globalBucketId = bucketResponse.Id
	return nil
}

func resolveBucket(ctx context.Context, client *api.Client) error {
	if globalBucketId != "" {
		name, found, err := fetchBucket(ctx, client, globalBucketId)
		if err != nil {
			return err
		}
		if found {
			globalBucketName = name
			log.Printf("Using existing bucket: %s (%s)", name, globalBucketId)
			return nil
		}
		log.Printf("TF_VAR_bucket_id=%s does not exist on cluster; discovering or creating one", globalBucketId)
		globalBucketId = ""
	}

	discoveredId, discoveredName, err := discoverFirstBucket(ctx, client)
	if err != nil {
		return err
	}
	if discoveredId != "" {
		globalBucketId = discoveredId
		globalBucketName = discoveredName
		log.Printf("Discovered existing bucket: %s (%s)", discoveredName, discoveredId)
		return nil
	}

	if err := createBucket(ctx, client); err != nil {
		return err
	}
	// Mark the bucket as created BEFORE waiting for it to become healthy. The
	// POST has already succeeded so the bucket exists remotely; if bucketWait
	// times out, cleanup must still delete it to avoid leaking the bucket
	// across flaky runs.
	globalBucketCreated = true
	if err := bucketWait(ctx, client); err != nil {
		return err
	}
	log.Printf("Created bucket: %s (%s)", globalBucketName, globalBucketId)
	return nil
}

func destroyBucket(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, globalClusterId, globalBucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	log.Printf("bucket destroyed: %s", globalBucketId)
	return nil
}

func createDMBucket(ctx context.Context, client *api.Client) error {
	bucketRequest := bucketapi.CreateBucketRequest{
		Name: dmBucketName,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, dmClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		bucketRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	bucketResponse := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucketResponse); err != nil {
		return err
	}

	dmBucketId = bucketResponse.Id
	return nil
}

func destroyDMBucket(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, dmClusterId, dmBucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	log.Printf("DM bucket destroyed: %s", dmBucketId)
	return nil
}

func resolveDMBucket(ctx context.Context, client *api.Client) error {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, dmClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return err
	}
	for _, bucket := range buckets {
		if bucket.Name == dmBucketName {
			dmBucketId = bucket.Id
			log.Printf("Discovered existing DM bucket: %s (%s)", dmBucketName, dmBucketId)
			return nil
		}
	}
	if len(buckets) > 0 {
		dmBucketId = buckets[0].Id
		dmBucketName = buckets[0].Name
		log.Printf("Discovered existing DM bucket: %s (%s)", dmBucketName, dmBucketId)
		return nil
	}

	if err := createDMBucket(ctx, client); err != nil {
		return err
	}
	dmBucketCreated = true
	if err := dmBucketWait(ctx, client); err != nil {
		return err
	}
	log.Printf("Created DM bucket: %s (%s)", dmBucketName, dmBucketId)
	return nil
}

func dmBucketWait(ctx context.Context, client *api.Client) error {
	const maxWaitTime = 5 * time.Minute

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeoutWaitingForBucket
		case <-ticker.C:
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, dmClusterId, dmBucketId)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
			if err == nil {
				log.Print("DM bucket created")
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

func resolveBucketNameById(ctx context.Context, client *api.Client, bucketID string) (string, error) {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return "", err
	}
	for _, bucket := range buckets {
		if bucket.Id == bucketID {
			return bucket.Name, nil
		}
	}

	return "", fmt.Errorf("bucket with ID %s not found", bucketID)
}

func bucketWait(ctx context.Context, client *api.Client) error {
	const maxWaitTime = 5 * time.Minute

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, maxWaitTime)
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeoutWaitingForBucket
		case <-ticker.C:
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, globalClusterId, globalBucketId)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			_, err := client.ExecuteWithRetry(
				ctx,
				cfg,
				nil,
				globalToken,
				nil,
			)
			if err == nil {
				log.Print("bucket created")
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
