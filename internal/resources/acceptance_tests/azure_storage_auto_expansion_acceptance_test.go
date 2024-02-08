package acceptance_tests

import (
	"fmt"
	"regexp"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterResourceAzureDiskAutoExpansion(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("azure")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansion(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "true"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},

			//Update the autoexpansion
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOff(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
			//Turn it on
			{
				Config: testAccClusterConfigAzureDiskAutoExpansion(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "true"),
				),
			},
			//Update with invalid field
			{
				Config:      testAccClusterConfigAzureDiskAutoExpansionInvalidfield(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Error: Invalid reference"),
			},
		},
	})
}

func TestAccClusterResourceAzureAutoDiskExpansionDefault(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("azure")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionDefault(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
		},
	})

}

func TestAccClusterResourceAzureAutoDiskExpansionOff(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("azure")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansionOff(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "false"),
				),
			},
		},
	})

}

func TestAccClusterResourceAzureAutoDiskExpansionInvalidConfig(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	fmt.Println(resourceReference)
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("azure")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{

		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccClusterConfigAzureDiskAutoExpansionInvalidfield(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Error: Invalid reference"),
			},
		},
	})

}

func TestAccClusterResourceAzureDiskAutoExpansionDeleteCluster(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr, err := acctest.GetCIDR("azure")
	fmt.Println(cidr)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAzureDiskAutoExpansion(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.autoexpansion", "true"),
					testAccDeleteClusterResource(resourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
		},
	})
}

func testAccClusterConfigAzureDiskAutoExpansion(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
  			type          = "P6"
			autoexpansion = true
		}
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

func testAccClusterConfigAzureDiskAutoExpansionDefault(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
  			type          = "P6"
		}
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

func testAccClusterConfigAzureDiskAutoExpansionOff(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster" "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
    type   = "azure"
    region = "eastus"
    cidr   = "%[5]s"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
  			type          = "P6"
			autoexpansion = false
		}
      }
      num_of_nodes = 3
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

func testAccClusterConfigAzureDiskAutoExpansionInvalidfield(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
   organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster" "%[2]s" {
 organization_id = var.organization_id
 project_id      = %[4]s.id
 name            = "Terraform Acceptance Test Cluster"
 description     = "My first test cluster for multiple services."
 cloud_provider = {
   type   = "azure"
   region = "eastus"
   cidr   = "%[5]s"
 }
 service_groups = [
   {
     node = {
       compute = {
         cpu = 4
         ram = 16
       }
       disk = {
 			type          = "P6"
			autoexpansion = invalid
		}
     }
     num_of_nodes = 3
     services     = ["data", "index", "query"]
   }
 ]
 availability = {
   "type" : "multi"
 }
 support = {
   plan     = "developer pro"
   timezone = "PT"
 }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}
