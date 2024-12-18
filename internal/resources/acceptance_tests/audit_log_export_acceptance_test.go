package acceptance_tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"
)

func TestAccAuditLogExportTestCases(t *testing.T) {
	clusterResourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	clusterResourceReference := "couchbase-capella_cluster." + clusterResourceName
	projectName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectName
	cidr, err := acctest.GetCIDR("aws")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	auditLogSettingsResourceName := "acc_audit_log_settings_" + acctest.GenerateRandomResourceName()
	auditLogExportResourceName := "acc_audit_log_export_" + acctest.GenerateRandomResourceName()
	auditLogExportResourceReference := "couchbase-capella_audit_log_export." + auditLogExportResourceName

	resource.Test(
		t, resource.TestCase{
			PreCheck:                 func() { acctest.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccAuditLogExportSetup(
						acctest.Cfg, projectName, projectResourceReference, clusterResourceName,
						clusterResourceReference, cidr, auditLogSettingsResourceName,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccAuditLogExportGetCluster(clusterResourceReference),
					),
				},
				{
					Config: testAccAuditLogExportConfig(
						acctest.Cfg, projectName, projectResourceReference, clusterResourceName,
						clusterResourceReference, cidr, auditLogExportResourceName, auditLogExportResourceReference,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet(auditLogExportResourceReference, "id"),
					),
					ExpectNonEmptyPlan: true,
				},
			},
		},
	)
}

// create cluster and enable audit logs
func testAccAuditLogExportSetup(providerAndVariables, projectName, projectResourceReference, clusterResourceName, clusterResourceReference, cidr, auditLogSettingsResourceName string) string {
	config := fmt.Sprintf(
		`
%[1]s

resource "couchbase-capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "terraform audit log export test cluster"
}

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = var.organization_id
  project_id      = %[3]s.id
  name            = "terraform audit log export test cluster"
  description     = "terraform audit log export test cluster"
  configuration_type = "multiNode"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[6]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}

resource "couchbase-capella_audit_log_settings" "%[7]s" {
  organization_id   = var.organization_id
  project_id        = %[3]s.id
  cluster_id        = %[5]s.id
  audit_enabled     = true
  enabled_event_ids = [8243, 8257, 8265]
  disabled_users    = []
}

`, providerAndVariables, projectName, projectResourceReference, clusterResourceName, clusterResourceReference, cidr,
		auditLogSettingsResourceName,
	)

	return config
}

func testAccAuditLogExportGetCluster(resourceReference string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state

		var rawState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		err = getCluster(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}

func getCluster(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId,
	)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	_, err := data.Client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func testAccAuditLogExportConfig(providerAndVariables, projectName, projectResourceReference, clusterResourceName, clusterResourceReference, cidr, auditLogExportResourceName, auditLogExportResourceReference string) string {
	return fmt.Sprintf(
		`
%[1]s

resource "couchbase-capella_project" "%[2]s" {
    organization_id = var.organization_id
	name            = "%[2]s"
	description     = "terraform audit log export test cluster"
}

resource "couchbase-capella_cluster" "%[4]s" {
  organization_id = var.organization_id
  project_id      = %[3]s.id
  name            = "terraform audit log export test cluster"
  description     = "terraform audit log export test cluster"
  configuration_type = "multiNode"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[6]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "PT"
  }
}

output "%[7]s" {
  value = %[8]s
}

resource "couchbase-capella_audit_log_export" "%[7]s" {
 organization_id = var.organization_id
 project_id = %[3]s.id 
 cluster_id = %[5]s.id
 start    = "%[9]s"
 end      = "%[10]s"
}

`, providerAndVariables, projectName, projectResourceReference, clusterResourceName, clusterResourceReference, cidr,
		auditLogExportResourceName, auditLogExportResourceReference,
		time.Now().Add(-2*time.Hour).Format("2006-01-02T15:04:05-07:00"),
		time.Now().Add(-1*time.Hour).Format("2006-01-02T15:04:05-07:00"),
	)
}
