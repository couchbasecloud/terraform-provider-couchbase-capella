package acceptance_tests

import (
	"context"
	"sync"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

var (
	dmClusterOnce sync.Once
	dmClusterErr  error
)

// ensureDMCluster lazily provisions the dedicated data-management cluster (and
// its bucket) used by the data_management_* acceptance tests. It is built on
// first call from the first DM test that runs and reused thereafter, so runs
// that exercise no DM tests don't pay the cluster-creation cost. cleanup()
// tears it down via the dmClusterCreated/dmBucketCreated flags.
func ensureDMCluster() error {
	dmClusterOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
		defer cancel()
		client := api.NewClient(timeout)

		if err := setupDMCluster(ctx, client); err != nil {
			dmClusterErr = err
			return
		}
		if err := resolveDMBucket(ctx, client); err != nil {
			dmClusterErr = err
			return
		}
		// Bucket creation triggers a rebalance; wait for the cluster to return
		// to Healthy before dependent resources are created.
		if err := dmClusterWait(ctx, client, false); err != nil {
			dmClusterErr = err
			return
		}
	})
	return dmClusterErr
}
