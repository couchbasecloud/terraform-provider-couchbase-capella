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
	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}

	_ = node.FromDiskAWS(diskAws)

	clusterRequest := clusterapi.CreateClusterRequest{
		Name: globalClusterName,
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

// findClusterByName lists clusters in the current project and returns the ID of
// the first cluster matching name, or "" if none is found.
func findClusterByName(ctx context.Context, client *api.Client, name string) (string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	clusters, err := api.GetPaginated[[]clusterapi.GetClusterResponse](ctx, client, globalToken, cfg, api.SortById)
	if err != nil {
		return "", err
	}
	for _, c := range clusters {
		if c.Name == name {
			return c.Id.String(), nil
		}
	}
	return "", nil
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

func setupDMCluster(ctx context.Context, client *api.Client) error {
	if dmClusterId != "" {
		log.Printf("Using existing DM cluster: %s", dmClusterId)
		return nil
	}
	if err := createDMCluster(ctx, client); err != nil {
		var apiErr *api.Error
		if !errors.As(err, &apiErr) || apiErr.HttpStatusCode < 500 {
			return err
		}
		id, findErr := findClusterByName(ctx, client, dmClusterName)
		if findErr != nil || id == "" {
			return err
		}
		log.Printf("createDMCluster returned 5xx but cluster was found; adopting %s", id)
		dmClusterId = id
	} else {
		dmClusterCreated = true
	}
	return dmClusterWait(ctx, client, false)
}

func createDMCluster(ctx context.Context, client *api.Client) error {
	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}

	_ = node.FromDiskAWS(diskAws)

	clusterRequest := clusterapi.CreateClusterRequest{
		Name: dmClusterName,
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

	dmClusterId = clusterResponse.Id.String()

	return nil
}

func destroyDMCluster(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, dmClusterId)
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

func dmClusterWait(ctx context.Context, client *api.Client, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s",
			globalHost,
			globalOrgId,
			globalProjectId,
			dmClusterId,
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
				log.Print("DM cluster destroyed")
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
				log.Printf("DM cluster created - %s", dmClusterId)
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for DM cluster to be created or destroyed")
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
				log.Printf("cluster created - %s", globalClusterId)
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for cluster to be created or destroyed")
}
