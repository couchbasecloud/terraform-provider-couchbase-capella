package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"

	"log"
	"os"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"
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
	clientV2, _ := apigen.NewClientWithResponses(globalHost,
		apigen.WithHTTPClient(client),
		apigen.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+globalToken)
			return nil
		}),
	)

	err := setup(ctx, client, clientV2)
	if err != nil {
		log.Print(err)
		code = 1
	} else {
		code = m.Run()
	}

	if err = cleanup(ctx, client, clientV2); err != nil {
		log.Print(err)
		code = 1
	}

	os.Exit(code)
}

func setup(ctx context.Context, client *api.Client, clientV2 *apigen.ClientWithResponses) error {
	if err := createProject(ctx, clientV2); err != nil {
		return err
	}
	if err := createCluster(ctx, clientV2); err != nil { // v2
		return err
	}
	if err := clusterWait(ctx, clientV2, false); err != nil { // v2
		return err
	}
	if err := createBucket(ctx, clientV2); err != nil { // v2
		return err
	}
	if err := bucketWait(ctx, clientV2); err != nil { // v2
		return err
	}
	if err := createAppService(ctx, client); err != nil {
		return err
	}
	if err := appServiceWait(ctx, client, false); err != nil {
		return err
	}

	return nil
}

func cleanup(ctx context.Context, client *api.Client, clientV2 *apigen.ClientWithResponses) error {
	if globalAppServiceId != "" {
		if err := destroyAppService(ctx, client); err != nil {
			return err
		}

		if err := appServiceWait(ctx, client, true); err != nil {
			return err
		}
	}

	if globalClusterId != "" {
		if err := destroyCluster(ctx, clientV2); err != nil { // v2
			return err
		}

		if err := clusterWait(ctx, clientV2, true); err != nil { // v2
			return err
		}
	}

	if globalProjectId != "" {
		if err := destroyProject(ctx, clientV2); err != nil {
			return err
		}
	}

	return nil
}
