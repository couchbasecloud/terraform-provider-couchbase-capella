package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// TestAccCloudAccountsDataSource reads the test organization's cloud accounts. The shared test
// cluster is on AWS, so the organization always has a Capella AWS account. Every output is also
// compared with the API's own response, so the GCP and Azure outputs are checked whether or not
// the organization has those accounts.
func TestAccCloudAccountsDataSource(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_cloud_accounts_")
	resourceReference := "data.couchbase-capella_cloud_accounts." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudAccountsDataSourceConfig(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestMatchResourceAttr(resourceReference, "aws_capella_account", regexp.MustCompile(`^\d{12}$`)),
					testAccCheckCloudAccountsMatchAPI(resourceReference),
				),
			},
		},
	})
}

// testAccCheckCloudAccountsMatchAPI reads the route the data source reads and compares every
// output with it. The response is decoded by property name rather than through the provider's
// own struct, so a property name the provider gets wrong leaves an output empty and fails here.
// The subscription ID is null when the response has none.
func testAccCheckCloudAccountsMatchAPI(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		cfg := api.EndpointCfg{
			Url:           fmt.Sprintf("%s/v4/organizations/%s/cloudAccounts", globalHost, globalOrgId),
			Method:        http.MethodGet,
			SuccessStatus: http.StatusOK,
		}
		response, err := globalClient.ExecuteWithRetry(context.Background(), cfg, nil, globalToken, nil)
		if err != nil {
			return fmt.Errorf("reading the cloud accounts: %w", err)
		}

		var want map[string]any
		if err := json.Unmarshal(response.Body, &want); err != nil {
			return fmt.Errorf("decoding the cloud accounts: %w", err)
		}

		var checks []resource.TestCheckFunc
		for attr, property := range map[string]string{
			"aws_capella_account":        "aws-capella-account",
			"gcp_capella_project":        "gcp-capella-project",
			"azure_capella_subscription": "azure-capella-subscription",
		} {
			value, _ := want[property].(string)
			checks = append(checks, resource.TestCheckResourceAttr(resourceReference, attr, value))
		}

		if id, ok := want["azure-capella-subscription-id"].(string); ok {
			checks = append(checks, resource.TestCheckResourceAttr(resourceReference, "azure_capella_subscription_id", id))
		} else {
			checks = append(checks, resource.TestCheckNoResourceAttr(resourceReference, "azure_capella_subscription_id"))
		}

		return resource.ComposeAggregateTestCheckFunc(checks...)(s)
	}
}

func testAccCloudAccountsDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%[1]s

data "couchbase-capella_cloud_accounts" "%[3]s" {
  organization_id = "%[2]s"
}
`, globalProviderBlock, globalOrgId, resourceName)
}
