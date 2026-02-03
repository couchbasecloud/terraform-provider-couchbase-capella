package acceptance_tests

import (
	"context"
	"fmt"
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

	globalProviderBlock = fmt.Sprint(`
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
`)

	var code int
	ctx := context.Background()
	client := api.NewClient(timeout)

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
		if err := createProject(ctx, client); err != nil {
			return err
		}
	} else {
		log.Printf("Using existing project: %s", globalProjectId)
	}

	// Create cluster only if not provided via env var
	if globalClusterId == "" {
		if err := createCluster(ctx, client); err != nil {
			return err
		}
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
	} else {
		log.Printf("Using existing bucket: %s", globalBucketId)
	}

	// Create app service only if not provided via env var
	if globalAppServiceId == "" {
		if err := createAppService(ctx, client); err != nil {
			return err
		}
		if err := appServiceWait(ctx, client, false); err != nil {
			return err
		}
	} else {
		log.Printf("Using existing app service: %s", globalAppServiceId)
	}

	if err := createAppEndpoint(ctx, client); err != nil {
		return err
	}
	if err := appEndpointWait(ctx, client); err != nil {
		return err
	}

	return nil
}

func cleanup(ctx context.Context, client *api.Client) error {
	// Only destroy app service if it was created by setup (not provided via env var)
	if globalAppServiceId != "" && os.Getenv("TF_VAR_app_service_id") == "" {
		if err := destroyAppService(ctx, client); err != nil {
			return err
		}

		if err := appServiceWait(ctx, client, true); err != nil {
			return err
		}
	}

	// Only destroy cluster if it was created by setup (not provided via env var)
	if globalClusterId != "" && os.Getenv("TF_VAR_cluster_id") == "" {
		if err := destroyCluster(ctx, client); err != nil {
			return err
		}

		if err := clusterWait(ctx, client, true); err != nil {
			return err
		}
	}

	// Only destroy project if it was created by setup (not provided via env var)
	if globalProjectId != "" && os.Getenv("TF_VAR_project_id") == "" {
		if err := destroyProject(ctx, client); err != nil {
			return err
		}
	}

	return nil
}
