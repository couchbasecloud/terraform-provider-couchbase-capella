package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/couchbase/tools-common/types/ptr"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	bucketapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/bucket"
	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
)

// The anchor is one twoWay replication between a bucket on the global cluster and a bucket on the
// DM cluster. It exists from the first replication test until cleanup(). While it exists, Capella
// keeps the network link between the two clusters, in both directions. Without it, each test
// deletes the last replication: both clusters go into the Peering state, Capella blocks a new link
// for up to 10 minutes, and the next create must build the link again. It is twoWay so that a test
// that replicates from the DM cluster, such as a twoWay test, does not have to build the reverse
// link either.
var (
	replicationAnchorOnce sync.Once
	replicationAnchorErr  error
)

// waitForReplicationClusters waits until both source and target clusters are Healthy.
func waitForReplicationClusters() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	client := api.NewClient(timeout)
	if err := clusterWait(ctx, client, false); err != nil {
		return err
	}
	return dmClusterWait(ctx, client, false)
}

// requireReplicationAnchor makes sure that the DM cluster and the anchor replication exist.
func requireReplicationAnchor(t *testing.T) {
	t.Helper()
	requireDMCluster(t)
	if err := ensureReplicationAnchor(); err != nil {
		t.Fatalf("ensureReplicationAnchor: %v", err)
	}
}

func ensureReplicationAnchor() error {
	replicationAnchorOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
		defer cancel()
		client := api.NewClient(timeout)

		replicationAnchorErr = createReplicationAnchor(ctx, client)
	})
	return replicationAnchorErr
}

func createReplicationAnchor(ctx context.Context, client *api.Client) error {
	var err error
	replicationAnchorSourceBucketId, err = createReplicationAnchorBucket(ctx, client, globalClusterId, randomStringWithPrefix(replicationAnchorSourceBucketPrefix))
	if err != nil {
		return err
	}
	// Set as soon as the first part exists, so cleanup() deletes whatever was made.
	replicationAnchorCreated = true
	replicationAnchorTargetBucketId, err = createReplicationAnchorBucket(ctx, client, dmClusterId, randomStringWithPrefix(replicationAnchorTargetBucketPrefix))
	if err != nil {
		return err
	}

	// A new bucket starts a rebalance, and the replication API rejects a cluster that is not Healthy.
	if err = waitForReplicationClusters(); err != nil {
		return err
	}

	jobId, err := postReplicationAnchor(ctx, client)
	if err != nil {
		return err
	}
	replicationAnchorId, err = waitForReplicationAnchorJob(ctx, client, jobId)
	if err != nil {
		return err
	}
	log.Printf("replication anchor created - %s", replicationAnchorId)

	// The create builds the network link in both directions. It queues a deploy for each direction,
	// so a cluster can go Peering, Healthy, then Peering again. Check twice, a short time apart.
	if err = waitForReplicationClusters(); err != nil {
		return err
	}
	time.Sleep(2 * time.Minute)
	return waitForReplicationClusters()
}

func createReplicationAnchorBucket(ctx context.Context, client *api.Client, clusterId, name string) (string, error) {
	request := bucketapi.CreateBucketRequest{
		Name:                     name,
		StorageBackend:           ptr.To("couchstore"),
		MemoryAllocationInMb:     ptr.To(int64(100)),
		BucketConflictResolution: ptr.To("seqno"),
	}
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets", globalHost, globalOrgId, globalProjectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, request, globalToken, nil)
	if err != nil {
		return "", fmt.Errorf("creating replication anchor bucket %s: %w", name, err)
	}

	bucket := bucketapi.GetBucketResponse{}
	if err = json.Unmarshal(response.Body, &bucket); err != nil {
		return "", err
	}
	log.Printf("replication anchor bucket created: %s (%s)", name, bucket.Id)
	return bucket.Id, nil
}

// postReplicationAnchor sends the create call. It sends the call again while the API reports that a cluster is not Healthy.
func postReplicationAnchor(ctx context.Context, client *api.Client) (string, error) {
	request := replicationapi.CreateReplicationRequest{
		SourceBucket: replicationAnchorSourceBucketId,
		Target: replicationapi.CreateReplicationTarget{
			Cluster: dmClusterId,
			Bucket:  replicationAnchorTargetBucketId,
		},
		Mode:      replicationapi.ModeAsync,
		Direction: ptr.To("twoWay"),
	}
	cfg := api.EndpointCfg{
		Url:           replicationapi.ClusterReplicationsURL(globalHost, globalOrgId, globalProjectId, globalClusterId),
		Method:        http.MethodPost,
		SuccessStatus: http.StatusAccepted,
	}

	for {
		response, err := client.ExecuteWithRetry(ctx, cfg, request, globalToken, nil)
		if err == nil {
			created := replicationapi.CreateReplicationAsyncResponse{}
			if err = json.Unmarshal(response.Body, &created); err != nil {
				return "", err
			}
			if created.JobId == "" {
				return "", errors.New("replication anchor create response has no job ID")
			}
			return created.JobId, nil
		}

		var apiErr *api.Error
		if !errors.As(err, &apiErr) || apiErr.HttpStatusCode != http.StatusUnprocessableEntity ||
			!strings.Contains(apiErr.Message, "must be in a healthy state") {
			return "", fmt.Errorf("creating replication anchor: %w", err)
		}
		log.Printf("replication anchor: a cluster is not healthy, retrying the create")

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("creating replication anchor: %w", ctx.Err())
		case <-time.After(30 * time.Second):
		}
	}
}

// waitForReplicationAnchorJob polls the create job until it completes. Capella retries a job that
// failed with a non-fatal error, so only a fatal error stops the poll.
func waitForReplicationAnchorJob(ctx context.Context, client *api.Client, jobId string) (string, error) {
	cfg := api.EndpointCfg{
		Url:           replicationapi.ReplicationJobURL(globalHost, globalOrgId, globalProjectId, globalClusterId, jobId),
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	for {
		response, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
		if err == nil {
			job := replicationapi.GetReplicationJobResponse{}
			if err = json.Unmarshal(response.Body, &job); err != nil {
				return "", err
			}
			switch {
			case job.State == replicationapi.JobStateComplete && job.ReplicationId != "":
				return job.ReplicationId, nil
			case job.State == replicationapi.JobStateFailed && replicationapi.IsFatalJobError(job.LastError):
				return "", fmt.Errorf("replication anchor job %s failed: %s", jobId, job.LastError)
			}
		} else {
			log.Printf("replication anchor: could not get job %s, retrying: %v", jobId, err)
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("waiting for replication anchor job %s: %w", jobId, ctx.Err())
		case <-time.After(15 * time.Second):
		}
	}
}

// destroyReplicationAnchor deletes the anchor replication (one delete removes both directions) and its buckets, skipping the parts
// that were not made. It runs before the DM and global clusters are deleted, so that clusters given by the caller are left clean.
func destroyReplicationAnchor(ctx context.Context, client *api.Client) error {
	if replicationAnchorId != "" {
		url := replicationapi.ReplicationURL(globalHost, globalOrgId, globalProjectId, globalClusterId, replicationAnchorId)
		if err := deleteReplicationAnchorResource(ctx, client, url); err != nil {
			return fmt.Errorf("deleting replication anchor: %w", err)
		}
		log.Printf("replication anchor destroyed - %s", replicationAnchorId)

		// Deleting the last replication removes the network link: both clusters go into the
		// Peering state, and a cluster cannot be deleted while it is busy.
		if err := waitForReplicationClusters(); err != nil {
			return err
		}
	}

	buckets := []struct{ clusterId, bucketId string }{
		{globalClusterId, replicationAnchorSourceBucketId},
		{dmClusterId, replicationAnchorTargetBucketId},
	}
	for _, b := range buckets {
		if b.bucketId == "" {
			continue
		}
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/buckets/%s", globalHost, globalOrgId, globalProjectId, b.clusterId, b.bucketId)
		if err := deleteReplicationAnchorResource(ctx, client, url); err != nil {
			return fmt.Errorf("deleting replication anchor bucket %s: %w", b.bucketId, err)
		}
	}
	return nil
}

// deleteReplicationAnchorResource deletes a resource and waits until GET returns 404. The
// replication delete call can return 500 while the delete continues in the background.
func deleteReplicationAnchorResource(ctx context.Context, client *api.Client, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	deleteCfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	if _, err := client.ExecuteWithRetry(ctx, deleteCfg, nil, globalToken, nil); err != nil {
		var apiErr *api.Error
		if !errors.As(err, &apiErr) || (apiErr.HttpStatusCode != http.StatusNotFound && apiErr.HttpStatusCode != http.StatusInternalServerError) {
			return err
		}
	}

	getCfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	for {
		_, err := client.ExecuteWithRetry(ctx, getCfg, nil, globalToken, nil)
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusNotFound {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting for %s to be deleted: %w", url, ctx.Err())
		case <-time.After(30 * time.Second):
		}
	}
}
