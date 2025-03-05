package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

func createProject(ctx context.Context, client *api.Client) error {
	projectRequest := api.CreateProjectRequest{
		Name: "tf_acc_test_project_common",
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		projectRequest,
		globalToken,
		nil,
	)
	if err != nil {
		return err
	}

	projectResponse := api.GetProjectResponse{}
	if err = json.Unmarshal(response.Body, &projectResponse); err != nil {
		return err
	}

	log.Print("project created")

	globalProjectId = projectResponse.Id.String()

	return nil
}

func destroyProject(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
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

	log.Print("project destroyed")

	return nil
}
