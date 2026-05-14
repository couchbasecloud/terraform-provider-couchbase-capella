package acceptance_tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// this is the entry point for acceptance tests.
// it setups a common project, cluster and app service.
// once acceptance tests are completed, these resources are destroyed.
func TestMain(m *testing.M) {
	if err := getEnvVars(); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	globalProviderBlock = `
variable "host" {
  description = "The globalHost URL of Couchbase Cloud."
}

variable "auth_token" {
  description = "Authentication API Key"
  sensitive   = true
}

provider "couchbase-capella" {
  host                 = var.host
  authentication_token = var.auth_token
}
`

	var code int
	ctx := context.Background()
	client := api.NewClient(timeout)
	globalClient = client

	err := setup(ctx, client)
	if err != nil {
		log.Print(err)
		code = 1
	} else {
		code = m.Run()
	}

	if err = cleanup(ctx, client); err != nil {
		log.Print(err)
		code = 1
	}

	os.Exit(code)
}

func setup(ctx context.Context, client *api.Client) error {
	// Create project only if not provided via env var
	if globalProjectId == "" {
		created, err := createProject(ctx, client)
		if err != nil {
			return err
		}
		globalProjectCreated = created
	} else {
		log.Printf("Using existing project: %s", globalProjectId)
	}

	// Create cluster only if not provided via env var
	if globalClusterId == "" {
		if err := createCluster(ctx, client); err != nil {
			return err
		}
		globalClusterCreated = true
		if err := clusterWait(ctx, client, false); err != nil {
			return err
		}
	} else {
		log.Printf("Using existing cluster: %s", globalClusterId)
	}

	// Create bucket only if not provided via env var
	if globalBucketId == "" {
		if err := createBucket(ctx, client); err != nil {
			return err
		}
		if err := bucketWait(ctx, client); err != nil {
			return err
		}
		// Bucket creation triggers a cluster rebalance; wait for the cluster
		// to return to Healthy before creating dependent resources, otherwise
		// createAppService races and fails with 412 "cluster is rebalancing".
		if err := clusterWait(ctx, client, false); err != nil {
			return err
		}
	} else {
		bucketName, err := resolveBucketNameById(ctx, client, globalBucketId)
		if err != nil {
			return err
		}
		globalBucketName = bucketName
		log.Printf("Using existing bucket: %s (%s)", globalBucketName, globalBucketId)
	}

	// Create app service only if not provided via env var
	if globalAppServiceId == "" {
		if err := createAppService(ctx, client); err != nil {
			return err
		}
		globalAppServiceCreated = true
		if err := appServiceWait(ctx, client, false); err != nil {
			return err
		}
	} else {
		log.Printf("Using existing app service: %s", globalAppServiceId)
	}

	appEndpointCreated, err := createAppEndpoint(ctx, client, globalAppEndpointName, globalBucketName)
	if err != nil {
		return err
	}
	globalAppEndpointCreated = appEndpointCreated
	if err := appEndpointWait(ctx, client, globalAppEndpointName); err != nil {
		return err
	}

	return nil
}

func cleanup(ctx context.Context, client *api.Client) error {
	if err := cleanupAppEndpointTestEnvironment(ctx, client); err != nil {
		return err
	}

	if globalAppEndpointCreated {
		if err := deleteFixtureEndpoint(ctx, client, globalAppEndpointName); err != nil {
			return err
		}
	}

	if globalAppServiceCreated {
		if err := destroyAppService(ctx, client); err != nil {
			return err
		}

		if err := appServiceWait(ctx, client, true); err != nil {
			return err
		}
	}

	if globalClusterCreated {
		if err := destroyCluster(ctx, client); err != nil {
			return err
		}

		if err := clusterWait(ctx, client, true); err != nil {
			return err
		}
	}

	if globalProjectCreated {
		if err := destroyProject(ctx, client); err != nil {
			return err
		}
	}

	return nil
}
