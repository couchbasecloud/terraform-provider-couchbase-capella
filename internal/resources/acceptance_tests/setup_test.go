package acceptance_tests

import (
	"context"
	"fmt"

	"log"
	"os"
	"testing"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

func TestMain(m *testing.M) {
	if err := GetEnvVars(); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	ProviderBlock = fmt.Sprint(`
variable "host" {
  description = "The Host URL of Couchbase Cloud."
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

	ctx := context.Background()
	client := api.NewClient(Timeout)

	if err := setup(ctx, client); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	code := m.Run()

	if err := cleanup(ctx, client); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	os.Exit(code)
}

func setup(ctx context.Context, client *api.Client) error {
	if err := CreateProject(ctx, client); err != nil {
		return err
	}
	if err := CreateCluster(ctx, client); err != nil {
		return err
	}
	if err := Wait(ctx, client, false); err != nil {
		return err
	}
	if err := CreateBucket(ctx, client); err != nil {
		return err
	}

	return nil
}

func cleanup(ctx context.Context, client *api.Client) error {
	if err := DestroyCluster(ctx, client); err != nil {
		return err
	}
	if err := Wait(ctx, client, true); err != nil {
		return err
	}
	if err := DestroyProject(ctx, client); err != nil {
		return err
	}

	return nil
}
