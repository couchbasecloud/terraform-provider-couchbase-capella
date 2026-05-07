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

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

var appEndpointEnvironmentOnce struct {
	sync.Once
	err error
}

func ensureAppEndpointTestEnvironment(t *testing.T) {
	t.Helper()

	appEndpointEnvironmentOnce.Do(func() {
		ctx := context.Background()
		appEndpointEnvironmentOnce.err = setupAppEndpointTestEnvironment(ctx, globalClient)
	})
	if appEndpointEnvironmentOnce.err != nil {
		t.Fatalf("failed to provision app endpoint test environment: %v", appEndpointEnvironmentOnce.err)
	}
}

func setupAppEndpointTestEnvironment(ctx context.Context, client *api.Client) error {
	suffix := time.Now().UTC().Format("20060102150405")
	appEndpointClusterName = "tf_acc_app_endpoint_cluster_" + suffix
	appEndpointAppServiceName = "tf_acc_app_endpoint_app_service_" + suffix

	clusterID, err := createAppEndpointTestCluster(ctx, client)
	if err != nil {
		return err
	}
	appEndpointClusterId = clusterID
	appEndpointClusterCreated = true
	if err = waitForAppEndpointTestCluster(ctx, client, false); err != nil {
		return err
	}

	bucketID, bucketCreated, err := createAppEndpointTestBucket(ctx, client, appEndpointBucketName)
	if err != nil {
		return err
	}
	appEndpointBucketId = bucketID
	appEndpointBucketCreated = bucketCreated
	if err = waitForAppEndpointTestBucket(ctx, client, appEndpointBucketId); err != nil {
		return err
	}

	appServiceID, err := createAppEndpointTestAppService(ctx, client)
	if err != nil {
		return err
	}
	appEndpointAppServiceId = appServiceID
	appEndpointAppServiceCreated = true
	if err = waitForAppEndpointTestAppService(ctx, client, false); err != nil {
		return err
	}

	endpointCreated, err := createAppEndpointForAppService(ctx, client, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, appEndpointCommonEndpointName, appEndpointBucketName)
	if err != nil {
		return err
	}
	appEndpointCreated = endpointCreated
	return appEndpointWaitForAppService(ctx, client, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, appEndpointCommonEndpointName)
}

func cleanupAppEndpointTestEnvironment(ctx context.Context, client *api.Client) error {
	if appEndpointCreated {
		if err := deleteAppEndpointFixtureEndpoint(ctx, client, globalProjectId, appEndpointClusterId, appEndpointAppServiceId, appEndpointCommonEndpointName); err != nil {
			return err
		}
	}

	if appEndpointBucketCreated {
		if err := deleteAppEndpointFixtureBucket(ctx, client, globalProjectId, appEndpointClusterId, appEndpointBucketName); err != nil {
			return err
		}
	}

	if appEndpointAppServiceCreated {
		if err := destroyAppEndpointTestAppService(ctx, client); err != nil {
			return err
		}
		if err := waitForAppEndpointTestAppService(ctx, client, true); err != nil {
			return err
		}
	}

	if appEndpointClusterCreated {
		if err := destroyAppEndpointTestCluster(ctx, client); err != nil {
			return err
		}
		if err := waitForAppEndpointTestCluster(ctx, client, true); err != nil {
			return err
		}
	}

	return nil
}

func createAppEndpointTestCluster(ctx context.Context, client *api.Client) (string, error) {
	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}
	_ = node.FromDiskAWS(diskAws)

	clusterRequest := clusterapi.CreateClusterRequest{
		Name: appEndpointClusterName,
		Availability: clusterapi.Availability{
			Type: "multi",
		},
		CloudProvider: clusterapi.CloudProvider{
			Region: "us-east-1",
			Type:   "aws",
		},
		ServiceGroups: []clusterapi.ServiceGroup{
			{
				Node: &clusterapi.Node{
					Compute: clusterapi.Compute{
						Cpu: 4,
						Ram: 16,
					},
					Disk: node.Disk,
				},
				Services: &[]clusterapi.Service{
					clusterapi.Service("data"),
					clusterapi.Service("index"),
					clusterapi.Service("query")},
				NumOfNodes: ptr.To(3),
			},
		},
		Support: clusterapi.Support{
			Plan:     "enterprise",
			Timezone: "PT",
		},
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	response, err := client.ExecuteWithRetry(ctx, cfg, clusterRequest, globalToken, nil)
	if err != nil {
		return "", err
	}

	clusterResponse := clusterapi.GetClusterResponse{}
	if err = json.Unmarshal(response.Body, &clusterResponse); err != nil {
		return "", err
	}

	return clusterResponse.Id.String(), nil
}

func waitForAppEndpointTestCluster(ctx context.Context, client *api.Client, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err != nil {
			resourceNotFound, errMsg := api.CheckResourceNotFoundError(err)
			if destroy && resourceNotFound {
				log.Print("app endpoint test cluster destroyed")
				return nil
			}
			if destroy && !resourceNotFound {
				return errors.New(errMsg)
			}

			return err
		}

		if !destroy {
			clusterResp := clusterapi.GetClusterResponse{}
			if err = json.Unmarshal(response.Body, &clusterResp); err != nil {
				return err
			}
			if clusterResp.CurrentState == clusterapi.Healthy {
				log.Print("app endpoint test cluster created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for app endpoint test cluster to be created or destroyed")
}

func destroyAppEndpointTestCluster(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	return err
}

func createAppEndpointTestBucket(ctx context.Context, client *api.Client, name string) (string, bool, error) {
	listUrl := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	listCfg := api.EndpointCfg{Url: listUrl, Method: http.MethodGet, SuccessStatus: http.StatusOK}

	buckets, err := api.GetPaginated[[]bucketapi.GetBucketResponse](ctx, client, globalToken, listCfg, api.SortById)
	if err != nil {
		return "", false, err
	}
	for _, bucket := range buckets {
		if bucket.Name == name {
			log.Printf("app endpoint test bucket %q already exists", name)
			return bucket.Id, false, nil
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, bucketapi.CreateBucketRequest{Name: name}, globalToken, nil)
	if err != nil {
		return "", false, err
	}

	bucketResponse := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucketResponse); err != nil {
		return "", false, err
	}

	return bucketResponse.Id, true, nil
}

func waitForAppEndpointTestBucket(ctx context.Context, client *api.Client, bucketID string) error {
	const maxWait = 5 * time.Minute
	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, appEndpointClusterId, bucketID)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done waiting for app endpoint test bucket: %w", ctx.Err())
		case <-time.After(30 * time.Second):
		}
	}

	return fmt.Errorf("timeout waiting for app endpoint test bucket %s", bucketID)
}

func createAppEndpointTestAppService(ctx context.Context, client *api.Client) (string, error) {
	var n int64 = 2
	appServiceRequest := appservice.CreateAppServiceRequest{
		Name: appEndpointAppServiceName,
		Compute: appservice.AppServiceCompute{
			Cpu: 2,
			Ram: 4,
		},
		Nodes: &n,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices", globalHost, globalOrgId, globalProjectId, appEndpointClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, appServiceRequest, globalToken, nil)
	if err != nil {
		return "", err
	}

	appServiceResponse := appservice.CreateAppServiceResponse{}
	if err = json.Unmarshal(response.Body, &appServiceResponse); err != nil {
		return "", err
	}

	return appServiceResponse.Id.String(), nil
}

func waitForAppEndpointTestAppService(ctx context.Context, client *api.Client, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", globalHost, globalOrgId, globalProjectId, appEndpointClusterId, appEndpointAppServiceId)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		response, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err != nil {
			resourceNotFound, errMsg := api.CheckResourceNotFoundError(err)
			if destroy && resourceNotFound {
				log.Print("app endpoint test app service destroyed")
				return nil
			}
			if destroy && !resourceNotFound {
				return errors.New(errMsg)
			}

			return err
		}

		if !destroy {
			var appServiceResponse appservice.GetAppServiceResponse
			if err = json.Unmarshal(response.Body, &appServiceResponse); err != nil {
				return fmt.Errorf("Error unmarshalling app service response: %v", err)
			}
			if appServiceResponse.CurrentState == appservice.Healthy {
				log.Print("app endpoint test app service created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for app endpoint test app service to be created or destroyed")
}

func destroyAppEndpointTestAppService(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s", globalHost, globalOrgId, globalProjectId, appEndpointClusterId, appEndpointAppServiceId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	return err
}
