package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// cluster is created with enterprise plan as some features require this.
func createCluster(ctx context.Context, client *api.Client) error {
	// TODO: AV-96938: generate CIDR dynamically
	cidr := "10.200.250.0/23"

	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}

	_ = node.FromDiskAWS(diskAws)

	clusterRequest := clusterapi.CreateClusterRequest{
		Name: "tf_acc_test_cluster_common",
		Availability: clusterapi.Availability{
			Type: "multi",
		},
		CloudProvider: clusterapi.CloudProvider{
			Cidr:   cidr,
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
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		clusterRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	clusterResponse := clusterapi.GetClusterResponse{}
	if err = json.Unmarshal(response.Body, &clusterResponse); err != nil {
		return err
	}

	globalClusterId = clusterResponse.Id.String()

	return nil
}

func destroyCluster(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, globalClusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	_, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func clusterWait(ctx context.Context, client *api.Client, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s",
			globalHost,
			globalOrgId,
			globalProjectId,
			globalClusterId,
		)

		cfg := api.EndpointCfg{
			Url:           url,
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}

		response, err := client.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			globalToken,
			nil,
		)
		if err != nil {
			resourceNotFound, errMsg := api.CheckResourceNotFoundError(err)
			if destroy && resourceNotFound {
				log.Print("cluster destroyed")
				return nil
			}
			if destroy && !resourceNotFound {
				return errors.New(errMsg)
			}

			return err
		}

		if !destroy {
			clusterResp := clusterapi.GetClusterResponse{}
			err = json.Unmarshal(response.Body, &clusterResp)
			if err != nil {
				return err
			}

			if clusterResp.CurrentState == clusterapi.Healthy {
				log.Print("cluster created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for cluster to be created or destroyed")
}
