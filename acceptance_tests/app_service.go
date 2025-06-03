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
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
)

func createAppService(ctx context.Context, client *api.Client) error {
	var n int64 = 2
	appServiceRequest := appservice.CreateAppServiceRequest{
		Name: "tf_acc_test_app_service_common",
		Compute: appservice.AppServiceCompute{
			Cpu: 2,
			Ram: 4,
		},
		Nodes: &n,
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices", 
		globalHost, 
		globalOrgId, 
		globalProjectId, 
		globalClusterId,
	)

	cfg := api.EndpointCfg{
		Url: url, 
		Method: http.MethodPost, 
		SuccessStatus: http.StatusCreated,
	}

	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		appServiceRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	appServiceResponse := appservice.CreateAppServiceResponse{}
	if err = json.Unmarshal(response.Body, &appServiceResponse); err != nil {
		return err
	}

	globalAppServiceId = appServiceResponse.Id.String()

	return nil
}

func appServiceWait(ctx context.Context, client *api.Client, destroy bool) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 1 * time.Minute

	deadline := time.Now().Add(maxWaitTime)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf(
			"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s",
			globalHost,
			globalOrgId,
			globalProjectId,
			globalClusterId,
			globalAppServiceId)

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
				log.Print("app service destroyed")
				return nil
			}
			if destroy && !resourceNotFound {
				return errors.New(errMsg)
			}
			
			return err
		}


		if !destroy {
			var appServiceResponse appservice.GetAppServiceResponse
			if err = json.Unmarshal(response.Body, &appServiceResponse); err != nil {
				return fmt.Errorf("Error unmarshalling app service response: %v", err)
			}

			if appServiceResponse.CurrentState == appservice.Healthy {
				log.Print("app service created")
				return nil
			}
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for app service to be created or destroyed")
}

func destroyAppService(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId)

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
