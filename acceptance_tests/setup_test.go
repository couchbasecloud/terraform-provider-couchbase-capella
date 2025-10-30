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
	if err := createProject(ctx, client); err != nil {
		return err
	}
	if err := createCluster(ctx, client); err != nil {
		return err
	}
	if err := clusterWait(ctx, client, false); err != nil {
		return err
	}
	if err := createBucket(ctx, client); err != nil {
		return err
	}
	if err := bucketWait(ctx, client); err != nil {
		return err
	}
	if err := createAppService(ctx, client); err != nil {
		return err
	}
	if err := appServiceWait(ctx, client, false); err != nil {
		return err
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
	if globalAppServiceId != "" {
		if err := destroyAppService(ctx, client); err != nil {
			return err
		}

		if err := appServiceWait(ctx, client, true); err != nil {
			return err
		}
	}

	if globalClusterId != "" {
		if err := destroyCluster(ctx, client); err != nil {
			return err
		}

		if err := clusterWait(ctx, client, true); err != nil {
			return err
		}
	}

	if globalProjectId != "" {
		if err := destroyProject(ctx, client); err != nil {
			return err
		}
	}

	return nil
}
