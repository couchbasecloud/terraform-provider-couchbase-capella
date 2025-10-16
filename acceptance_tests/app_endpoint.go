package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
)

func createAppEndpoint(ctx context.Context, client *api.Client) error {
	appEndpointRequest := app_endpoints.AppEndpointRequest{
		Name:             globalAppEndpointName,
		Bucket:           globalBucketName,
		UserXattrKey:     "syncFnXattr",
		DeltaSyncEnabled: true,
		Scopes: app_endpoints.Scopes{
			"_default": app_endpoints.Scope{
				Collections: map[string]app_endpoints.Collection{
					"_default": app_endpoints.Collection{
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

	_, err := client.ExecuteWithRetry(
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

func appEndpointWait(ctx context.Context, client *api.Client, destroy bool) error {
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
			resourceNotFound, errMsg := api.CheckResourceNotFoundError(err)
			if destroy && resourceNotFound {
				log.Print("app endpoint destroyed")
				return nil
			}
			if destroy && !resourceNotFound {
				return errors.New(errMsg)
			}

			return err
		}

		if !destroy {
			var appEndpointResponse app_endpoints.GetAppEndpointResponse
			if err = json.Unmarshal(response.Body, &appEndpointResponse); err != nil {
				return fmt.Errorf("Error unmarshalling app endpoint response: %v", err)
			}

			if appEndpointResponse.State == "Online" || appEndpointResponse.State == "Offline" {
				log.Print("app endpoint created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for app endpoint to be created or destroyed")
}

func destroyAppEndpoint(ctx context.Context, client *api.Client) error {
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
		Method:        http.MethodDelete,
		SuccessStatus: http.StatusAccepted,
	}

	_, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
