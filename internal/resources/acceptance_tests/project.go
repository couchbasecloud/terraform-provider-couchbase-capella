package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

func CreateProject(ctx context.Context, client *api.Client) error {
	projectRequest := api.CreateProjectRequest{
		Name: "tf_acc_test_project_common",
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects", Host, OrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(
		ctx,
		cfg,
		projectRequest,
		Token,
		nil,
	)
	if err != nil {
		return err
	}

	projectResponse := api.GetProjectResponse{}
	if err = json.Unmarshal(response.Body, &projectResponse); err != nil {
		return err
	}

	ProjectId = projectResponse.Id.String()

	return nil
}

func DestroyProject(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", Host, OrgId, ProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		Token,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
