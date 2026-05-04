package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
)

func createAppEndpoint(ctx context.Context, client *api.Client) error {
	// First, check if app endpoint already exists
	checkUrl := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		globalAppEndpointName,
	)
	checkCfg := api.EndpointCfg{
		Url:           checkUrl,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	_, err := client.ExecuteWithRetry(ctx, checkCfg, nil, globalToken, nil)
	if err == nil {
		log.Printf("App endpoint '%s' already exists", globalAppEndpointName)
		return nil
	}

	// App endpoint doesn't exist, create it
	appEndpointRequest := app_endpoints.AppEndpointRequest{
		Name:             globalAppEndpointName,
		Bucket:           globalBucketName,
		DeltaSyncEnabled: true,
		Scopes: app_endpoints.Scopes{
			"_default": app_endpoints.Scope{
				Collections: map[string]app_endpoints.Collection{
					"_default": {
						AccessControlFunction: "function (doc, oldDoc, meta) {return true;}",
						ImportFilter:          "function (doc) {return true;}",
					},
				},
			},
		},
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

	_, err = client.ExecuteWithRetry(
		ctx,
		cfg,
		appEndpointRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func appEndpointWait(ctx context.Context, client *api.Client) error {
	const maxWaitTime = 10 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
			globalHost,
			globalOrgId,
			globalProjectId,
			globalClusterId,
			globalAppServiceId,
			globalAppEndpointName,
		)

		cfg := api.EndpointCfg{
			Url:           url,
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}

		response, err := client.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			globalToken,
			nil,
		)
		if err != nil {
			if resourceNotFound, errMsg := api.CheckResourceNotFoundError(err); resourceNotFound {
				return fmt.Errorf("app endpoint not found: %s", errMsg)
			}
			// Transient errors (e.g. 500s) should not fall through to
			// json.Unmarshal on a nil response body, which would panic.
			// Log and retry instead.
			log.Printf("error fetching app endpoint, retrying: %v", err)
			select {
			case <-ctx.Done():
				return fmt.Errorf("context done while waiting to retry: %w", ctx.Err())
			case <-time.After(checkInterval):
			}
			continue
		}

		var appEndpointResponse app_endpoints.GetAppEndpointResponse
		if err = json.Unmarshal(response.Body, &appEndpointResponse); err != nil {
			return fmt.Errorf("Error unmarshalling app endpoint response: %v", err)
		}

		if appEndpointResponse.State == resources.AppEndpointStateOnline || appEndpointResponse.State == resources.AppEndpointStateOffline {
			log.Print(fmt.Sprintf("app endpoint state %s", appEndpointResponse.State))
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for app endpoint: %w", ctx.Err())
		case <-time.After(checkInterval):
		}
	}

	return fmt.Errorf("timeout waiting for app endpoint to be created or destroyed")
}
