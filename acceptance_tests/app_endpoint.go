package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/resources"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/app_endpoints"
)

func createAppEndpoint(ctx context.Context, client *api.Client, name, bucket string) error {
	const maxWait = 5 * time.Minute
	const retryInterval = 30 * time.Second

	checkUrl := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		name,
	)
	checkCfg := api.EndpointCfg{
		Url:           checkUrl,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	// Retry the existence check: 403 and 5xx are transient right after app
	// service creation (the endpoint API warms up after the service is deployed).
	var err error
	checkDeadline := time.Now().Add(maxWait)
	for {
		_, err = client.ExecuteWithRetry(ctx, checkCfg, nil, globalToken, nil)
		if err == nil {
			log.Printf("App endpoint '%s' already exists", name)
			return nil
		}
		// The app endpoint API returns 403 (instead of 404) when the endpoint
		// does not exist yet, so treat it the same as Not Found.
		var apiErr *api.Error
		if notFound, _ := api.CheckResourceNotFoundError(err); notFound {
			break
		} else if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusForbidden {
			break
		} else if permErr := permanentAPIError(err); permErr != nil {
			return fmt.Errorf("failed to check whether app endpoint %s exists: %w", name, permErr)
		}
		if time.Now().After(checkDeadline) {
			return fmt.Errorf("timeout waiting for app endpoint check for %s: %w", name, err)
		}
		log.Printf("transient error checking app endpoint %s, retrying in %s: %v", name, retryInterval, err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done while checking app endpoint: %w", ctx.Err())
		case <-time.After(retryInterval):
		}
	}

	appEndpointRequest := app_endpoints.AppEndpointRequest{
		Name:             name,
		Bucket:           bucket,
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

	postUrl := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
	)
	postCfg := api.EndpointCfg{
		Url:           postUrl,
		Method:        http.MethodPost,
		SuccessStatus: http.StatusCreated,
	}

	deadline := time.Now().Add(maxWait)
	for {
		_, err = client.ExecuteWithRetry(ctx, postCfg, appEndpointRequest, globalToken, nil)
		if err == nil {
			return nil
		}
		if permErr := permanentAPIError(err); permErr != nil {
			return permErr
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout creating app endpoint %s: %w", name, err)
		}
		log.Printf("transient error creating app endpoint %s, retrying in %s: %v", name, retryInterval, err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done while creating app endpoint: %w", ctx.Err())
		case <-time.After(retryInterval):
		}
	}
}

// permanentAPIError returns a non-nil error if err is a permanent 4xx API
// failure (excluding 429 Too Many Requests). Returns nil for transient errors
// so callers can continue retrying.
func permanentAPIError(err error) error {
	var apiErr *api.Error
	if errors.As(err, &apiErr) {
		if apiErr.HttpStatusCode != 0 &&
			apiErr.HttpStatusCode != http.StatusTooManyRequests &&
			apiErr.HttpStatusCode < 500 {
			return fmt.Errorf("permanent API error (HTTP %d): %s", apiErr.HttpStatusCode, apiErr.CompleteError())
		}
	}
	return nil
}

func appEndpointWait(ctx context.Context, client *api.Client, name string) error {
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
			name,
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
			// 403 and 404 are transient right after endpoint creation (API warms
			// up after the app service is deployed); retry them like 5xx.
			var apiErr *api.Error
			if errors.As(err, &apiErr) &&
				(apiErr.HttpStatusCode == http.StatusForbidden || apiErr.HttpStatusCode == http.StatusNotFound) {
				log.Printf("transient error fetching app endpoint (HTTP %d), retrying: %v", apiErr.HttpStatusCode, err)
			} else if permErr := permanentAPIError(err); permErr != nil {
				return permErr
			} else {
				log.Printf("transient error fetching app endpoint, retrying: %v", err)
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("context done while waiting to retry: %w", ctx.Err())
			case <-time.After(checkInterval):
			}
			continue
		}

		var appEndpointResponse app_endpoints.GetAppEndpointResponse
		if err = json.Unmarshal(response.Body, &appEndpointResponse); err != nil {
			return fmt.Errorf("error unmarshalling app endpoint response: %w", err)
		}

		if appEndpointResponse.State == resources.AppEndpointStateOnline || appEndpointResponse.State == resources.AppEndpointStateOffline {
			log.Printf("app endpoint state %s", appEndpointResponse.State)
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
