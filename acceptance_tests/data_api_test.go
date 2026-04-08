package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAccDataAPIDataSource(t *testing.T) {
	resourceName := "data_api_ds"
	resourceReference := "data.couchbase-capella_data_api." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAPIDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttrSet(resourceReference, "state"),
				),
			},
		},
	})
}

func testAccDataAPIDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
	%[1]s

	data "couchbase-capella_data_api" "%[2]s" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
	}
	`, globalProviderBlock, resourceName, globalOrgId, globalProjectId, globalClusterId)
}

func retrieveDataAPIFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s/dataAPI", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.ClientV1.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}

	dataAPIResp := api.GetDataAPIStatusResponse{}
	err = json.Unmarshal(response.Body, &dataAPIResp)
	if err != nil {
		return err
	}

	return nil
}
