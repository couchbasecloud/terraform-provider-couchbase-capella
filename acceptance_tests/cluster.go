package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

// cluster is created with enterprise plan as some features require this.
func createCluster(ctx context.Context, client *api.Client) error {
	id, err := createClusterNamed(ctx, client, "tf_acc_test_cluster_common")
	if err != nil {
		return err
	}
	globalClusterId = id
	return nil
}

func createSnapshotCluster(ctx context.Context, client *api.Client) error {
	id, err := createClusterNamed(ctx, client, "tf_acc_test_cluster_snapshot")
	if err != nil {
		return err
	}
	globalSnapshotClusterId = id
	globalSnapshotClusterCreated = true
	return nil
}

func createClusterNamed(ctx context.Context, client *api.Client, name string) (string, error) {
	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}

	_ = node.FromDiskAWS(diskAws)

	clusterRequest := clusterapi.CreateClusterRequest{
		Name: name,
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
		return "", err
	}

	clusterResponse := clusterapi.GetClusterResponse{}
	if err = json.Unmarshal(response.Body, &clusterResponse); err != nil {
		return "", err
	}

	return clusterResponse.Id.String(), nil
}

func destroyCluster(ctx context.Context, client *api.Client) error {
	return destroyClusterByID(ctx, client, globalClusterId)
}

func destroySnapshotCluster(ctx context.Context, client *api.Client) error {
	return destroyClusterByID(ctx, client, globalSnapshotClusterId)
}

func destroyClusterByID(ctx context.Context, client *api.Client, clusterID string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, clusterID)
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
	return clusterWaitByID(ctx, client, globalClusterId, destroy)
}

func snapshotClusterWait(ctx context.Context, client *api.Client, destroy bool) error {
	return clusterWaitByID(ctx, client, globalSnapshotClusterId, destroy)
}

func clusterWaitByID(ctx context.Context, client *api.Client, clusterID string, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s",
			globalHost,
			globalOrgId,
			globalProjectId,
			clusterID,
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

// createClustersConcurrently provisions the primary and snapshot clusters in
// parallel so the additional ~10 min cluster build doesn't extend total CI
// runtime serially. Either side is skipped when its corresponding
// TF_VAR_cluster_id / TF_VAR_snapshot_cluster_id env var is set.
func createClustersConcurrently(ctx context.Context, client *api.Client) error {
	var wg sync.WaitGroup
	var primaryErr, snapshotErr error

	if globalClusterId == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := createCluster(ctx, client); err != nil {
				primaryErr = fmt.Errorf("create primary cluster: %w", err)
				return
			}
			if err := clusterWait(ctx, client, false); err != nil {
				primaryErr = fmt.Errorf("wait primary cluster: %w", err)
			}
		}()
	} else {
		log.Printf("Using existing cluster: %s", globalClusterId)
	}

	if globalSnapshotClusterId == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := createSnapshotCluster(ctx, client); err != nil {
				snapshotErr = fmt.Errorf("create snapshot cluster: %w", err)
				return
			}
			if err := snapshotClusterWait(ctx, client, false); err != nil {
				snapshotErr = fmt.Errorf("wait snapshot cluster: %w", err)
			}
		}()
	} else {
		log.Printf("Using existing snapshot cluster: %s", globalSnapshotClusterId)
	}

	wg.Wait()

	switch {
	case primaryErr != nil && snapshotErr != nil:
		return fmt.Errorf("%w; %w", primaryErr, snapshotErr)
	case primaryErr != nil:
		return primaryErr
	case snapshotErr != nil:
		return snapshotErr
	}
	return nil
}

// destroyClustersConcurrently tears down both clusters in parallel. Each side
// is skipped when its env var was provided (pre-existing cluster).
func destroyClustersConcurrently(ctx context.Context, client *api.Client) error {
	var wg sync.WaitGroup
	var primaryErr, snapshotErr error

	if globalClusterId != "" && os.Getenv("TF_VAR_cluster_id") == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := destroyCluster(ctx, client); err != nil {
				primaryErr = fmt.Errorf("destroy primary cluster: %w", err)
				return
			}
			if err := clusterWait(ctx, client, true); err != nil {
				primaryErr = fmt.Errorf("wait primary cluster destroy: %w", err)
			}
		}()
	}

	if globalSnapshotClusterCreated {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := destroySnapshotCluster(ctx, client); err != nil {
				snapshotErr = fmt.Errorf("destroy snapshot cluster: %w", err)
				return
			}
			if err := snapshotClusterWait(ctx, client, true); err != nil {
				snapshotErr = fmt.Errorf("wait snapshot cluster destroy: %w", err)
			}
		}()
	}

	wg.Wait()

	switch {
	case primaryErr != nil && snapshotErr != nil:
		return fmt.Errorf("%w; %w", primaryErr, snapshotErr)
	case primaryErr != nil:
		return primaryErr
	case snapshotErr != nil:
		return snapshotErr
	}
	return nil
}
