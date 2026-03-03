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

func createBucket(ctx context.Context, client *api.Client) error {
	// First, check if bucket already exists
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, globalClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	
	// Use the paginated API to get all buckets
	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err == nil {
		// Check if bucket with globalBucketName already exists
		for _, bucket := range buckets {
			if bucket.Name == globalBucketName {
				globalBucketId = bucket.Id
				log.Printf("Bucket '%s' already exists with ID: %s", globalBucketName, globalBucketId)
				return nil
			}
		}
	}

	// Bucket doesn't exist, create it
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
