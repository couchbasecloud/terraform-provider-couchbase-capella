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

	// Resolve bucket: validate env var, else discover, else create.
	if globalBucketId != "" {
		exists, err := bucketExists(ctx, client, globalBucketId)
		if err != nil {
			return err
		}
		if exists {
			log.Printf("Using existing bucket: %s", globalBucketId)
		} else {
			log.Printf("TF_VAR_bucket_id=%s does not exist on cluster; discovering or creating one", globalBucketId)
			globalBucketId = ""
		}
	}
	if globalBucketId == "" {
		discoveredId, discoveredName, err := discoverFirstBucket(ctx, client)
		if err != nil {
			return err
		}
		if discoveredId != "" {
			globalBucketId = discoveredId
			globalBucketName = discoveredName
			log.Printf("Discovered existing bucket: %s (%s)", discoveredName, discoveredId)
		} else {
			if err := createBucket(ctx, client); err != nil {
				return err
			}
			if err := bucketWait(ctx, client); err != nil {
				return err
			}
			log.Printf("Created bucket: %s (%s)", globalBucketName, globalBucketId)
		}
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

	// Only destroy bucket if it was created by setup.
	if globalBucketCreated && globalBucketId != "" {
		if err := destroyBucket(ctx, client); err != nil {
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
