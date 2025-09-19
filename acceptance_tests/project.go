package acceptance_tests

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

func createProject(ctx context.Context, client *apigen.ClientWithResponses) error {
	req := apigen.CreateProjectRequest{
		Name: "tf_acc_test_project_common",
	}

	orgUUID, _ := uuid.Parse(globalOrgId)
	resp, err := client.PostProjectWithResponse(ctx, orgUUID, req)
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		return fmt.Errorf("unexpected status: %s", resp.Status())
	}
	globalProjectId = resp.JSON201.Id.String()
	return nil
}

func destroyProject(ctx context.Context, client *apigen.ClientWithResponses) error {
	orgUUID, _ := uuid.Parse(globalOrgId)
	projUUID, _ := uuid.Parse(globalProjectId)
	_, err := client.DeleteProjectByIDWithResponse(ctx, orgUUID, projUUID)
	return err
}
