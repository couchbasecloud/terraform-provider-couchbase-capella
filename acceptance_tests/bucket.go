package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	"net/http"
	"time"
)

func CreateBucket(ctx context.Context, client *api.Client) error {
	bucketRequest := bucketapi.CreateBucketRequest{
		Name: BucketName,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", Host, OrgId, ProjectId, ClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		bucketRequest,
		Token,
		nil,
	)
	if err != nil {
		return err
	}

	bucketResponse := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucketResponse); err != nil {
		return err
	}

	BucketId = bucketResponse.Id
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
			url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", Host, OrgId, ProjectId, ClusterId, BucketId)
			cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
			_, err := client.ExecuteWithRetry(
				ctx,
				cfg,
				nil,
				Token,
				nil,
			)
			if err == nil {
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
