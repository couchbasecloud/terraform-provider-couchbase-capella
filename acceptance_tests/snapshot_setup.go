package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
)

var (
	snapshotClusterOnce    sync.Once
	snapshotClusterID      string
	snapshotClusterCreated bool
	snapshotClusterErr     error
)

// ensureSnapshotCluster lazily provisions a dedicated cluster used by the
// cloud_snapshot_* acceptance tests so that the long-running snapshot/restore
// operations don't disrupt the primary cluster shared by other tests. The
// cluster is built on first call (typically from the first snapshot test that
// runs) and reused by subsequent calls. cleanup() destroys it via
// destroySnapshotClusterIfCreated.
func ensureSnapshotCluster() (string, error) {
	snapshotClusterOnce.Do(func() {
		// 45-minute timeout covers cluster creation (up to 30 min) plus a safety margin.
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
		defer cancel()
		client := api.NewClient(timeout)
		id, created, err := createSnapshotClusterAndWait(ctx, client)
		if err != nil {
			snapshotClusterErr = err
			return
		}
		snapshotClusterID = id
		snapshotClusterCreated = created
		log.Printf("snapshot cluster ready: %s (created=%v)", id, created)
	})
	return snapshotClusterID, snapshotClusterErr
}

// findSnapshotClusterByName returns the ID of the first cluster in the current
// project whose name matches, or "" if none is found.
func findSnapshotClusterByName(ctx context.Context, client *api.Client, name string) (string, error) {
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

func createSnapshotClusterAndWait(ctx context.Context, client *api.Client) (string, bool, error) {
	const clusterName = "tf_acc_test_cluster_snapshot"

	// Reuse an existing cluster with this name to avoid duplicates from leaked runs.
	if existing, err := findSnapshotClusterByName(ctx, client, clusterName); err != nil {
		return "", false, fmt.Errorf("find existing snapshot cluster: %w", err)
	} else if existing != "" {
		log.Printf("reusing existing snapshot cluster: %s", existing)
		if err := waitForClusterHealthy(ctx, client, existing); err != nil {
			return "", false, fmt.Errorf("wait for existing snapshot cluster: %w", err)
		}
		return existing, false, nil
	}

	node := clusterapi.Node{}
	diskAws := clusterapi.DiskAWS{
		Type:    clusterapi.DiskAWSType("gp3"),
		Storage: 50,
		Iops:    3000,
	}
	if err := node.FromDiskAWS(diskAws); err != nil {
		return "", false, fmt.Errorf("node.FromDiskAWS: %w", err)
	}

	req := clusterapi.CreateClusterRequest{
		Name:         clusterName,
		Availability: clusterapi.Availability{Type: "multi"},
		CloudProvider: clusterapi.CloudProvider{
			Region: "us-east-1",
			Type:   "aws",
		},
		ServiceGroups: []clusterapi.ServiceGroup{
			{
				Node: &clusterapi.Node{
					Compute: clusterapi.Compute{Cpu: 4, Ram: 16},
					Disk:    node.Disk,
				},
				Services: &[]clusterapi.Service{
					clusterapi.Service("data"),
					clusterapi.Service("index"),
					clusterapi.Service("query"),
				},
				NumOfNodes: ptr.To(3),
			},
		},
		Support: clusterapi.Support{Plan: "enterprise", Timezone: "PT"},
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
	resp, err := client.ExecuteWithRetry(ctx, cfg, req, globalToken, nil)
	if err != nil {
		return "", false, fmt.Errorf("create snapshot cluster: %w", err)
	}

	var clusterResp clusterapi.GetClusterResponse
	if err = json.Unmarshal(resp.Body, &clusterResp); err != nil {
		return "", false, fmt.Errorf("decode snapshot cluster response: %w", err)
	}

	id := clusterResp.Id.String()
	if err := waitForClusterHealthy(ctx, client, id); err != nil {
		return id, false, fmt.Errorf("wait snapshot cluster healthy: %w", err)
	}
	return id, true, nil
}

// destroySnapshotClusterIfCreated tears down the snapshot cluster if
// ensureSnapshotCluster successfully created one. Safe to call when no
// snapshot cluster was provisioned (e.g. no snapshot tests ran).
func destroySnapshotClusterIfCreated(ctx context.Context, client *api.Client) error {
	if !snapshotClusterCreated || snapshotClusterID == "" {
		return nil
	}
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", globalHost, globalOrgId, globalProjectId, snapshotClusterID)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
	if _, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil); err != nil {
		return fmt.Errorf("destroy snapshot cluster: %w", err)
	}
	if err := pollSnapshotCluster(ctx, client, snapshotClusterID, true); err != nil {
		return fmt.Errorf("wait snapshot cluster destroyed: %w", err)
	}
	return nil
}

func waitForClusterHealthy(ctx context.Context, client *api.Client, clusterID string) error {
	return pollSnapshotCluster(ctx, client, clusterID, false)
}

func pollSnapshotCluster(ctx context.Context, client *api.Client, clusterID string, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s",
			globalHost, globalOrgId, globalProjectId, clusterID,
		)
		cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
		resp, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err != nil {
			notFound, msg := api.CheckResourceNotFoundError(err)
			if destroy && notFound {
				return nil
			}
			if destroy && !notFound {
				return errors.New(msg)
			}
			return err
		}
		if !destroy {
			var got clusterapi.GetClusterResponse
			if err := json.Unmarshal(resp.Body, &got); err != nil {
				return err
			}
			if got.CurrentState == clusterapi.Healthy {
				return nil
			}
		}
		time.Sleep(checkInterval)
	}
	return fmt.Errorf("timeout waiting for snapshot cluster %s (destroy=%v)", clusterID, destroy)
}
