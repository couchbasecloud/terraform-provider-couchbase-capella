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

	ProviderBlock = fmt.Sprintf(`
provider "couchbase-capella" {
  host                 = "%s"
  authentication_token = "%s"
}`, Host, Token)

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
