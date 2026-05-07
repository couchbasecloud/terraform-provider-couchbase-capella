package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
)

func bucketExists(ctx context.Context, client *api.Client, bucketId string) (bool, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, globalClusterId, bucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err == nil {
		return true, nil
	}
	if apiErr, ok := err.(*api.Error); ok && apiErr.HttpStatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, err
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
	globalBucketCreated = true
	return nil
}

func destroyBucket(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, globalClusterId, globalBucketId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		if apiErr, ok := err.(*api.Error); ok && apiErr.HttpStatusCode == http.StatusNotFound {
			return nil
		}
		return err
	}
	log.Printf("bucket destroyed: %s", globalBucketId)
	return nil
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

			apiError, ok := err.(*api.Error)
			if ok {
				if apiError.HttpStatusCode != http.StatusNotFound {
					return err
				}
			} else {
				return err
			}
		}
	}
}
