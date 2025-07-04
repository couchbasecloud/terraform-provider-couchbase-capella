package acceptance_tests

import (
	"context"
	"fmt"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
	"net/http"
)

func createAppEndpoint(ctx context.Context, client *api.Client) error {
	globalAppEndpoint = "tf_acc_test_app_endpoint_common"

	appServiceRequest := app_endpoints.CreateAppEndpointRequest{
		Bucket:           globalBucketName,
		Name:             globalAppEndpoint,
		UserXattrKey:     nil,
		DeltaSyncEnabled: false,
		Scopes:           nil,
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
	)

	cfg := api.EndpointCfg{
		Url:           url,
		Method:        http.MethodPost,
		SuccessStatus: http.StatusCreated,
	}

	_, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		appServiceRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
