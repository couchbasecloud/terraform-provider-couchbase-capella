package acceptance_tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

const projectName = "tf_acc_test_project_common"

// createProject creates a new test project. On quota error (403) it falls back
// to findProjectByName and reuses an existing project so that concurrent CI
// runs sharing one Capella account do not all fail when the project limit is
// reached. Returns (created=true) only when a new project was actually created.
func createProject(ctx context.Context, client *api.Client) (bool, error) {
	projectRequest := api.CreateProjectRequest{
		Name: projectName,
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/projects", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := client.ExecuteWithRetry(ctx, cfg, projectRequest, globalToken, nil)
	if err != nil {
		// On quota exhaustion, try to reuse an existing project rather than
		// failing the whole test suite immediately.
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.HttpStatusCode == http.StatusForbidden {
			id, findErr := findProjectByName(ctx, client, projectName)
			if findErr == nil && id != "" {
				log.Printf("project quota hit; reusing existing project %s", id)
				globalProjectId = id
				return false, nil
			}
		}
		return false, err
	}

	projectResponse := api.GetProjectResponse{}
	if err = json.Unmarshal(response.Body, &projectResponse); err != nil {
		return false, err
	}

	log.Print("project created")
	globalProjectId = projectResponse.Id.String()
	return true, nil
}

// findProjectByName returns the ID of the first project in the org whose name
// matches, or "" if none is found.
func findProjectByName(ctx context.Context, client *api.Client, name string) (string, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects", globalHost, globalOrgId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	projects, err := api.GetPaginated[[]api.GetProjectResponse](ctx, client, globalToken, cfg, api.SortById)
	if err != nil {
		return "", err
	}
	for _, p := range projects {
		if p.Name == name {
			return p.Id.String(), nil
		}
	}
	return "", nil
}

func destroyProject(ctx context.Context, client *api.Client) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", globalHost, globalOrgId, globalProjectId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err := client.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	if err != nil {
		return err
	}
	log.Print("project destroyed")
	return nil
}
