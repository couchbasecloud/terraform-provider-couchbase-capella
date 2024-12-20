package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	"net/http"
)

func createBucket(ctx context.Context, client *api.Client) error {
	bucketRequest := bucketapi.CreateBucketRequest{
		Name: "default",
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

	Bucket = bucketResponse.Name
	return nil
}
