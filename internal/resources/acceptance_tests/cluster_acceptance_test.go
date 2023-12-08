package acceptance_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TODO: AV-66543 - Make Acceptance test concurrent
// Concurrent Project creation is failing.

// TestAccClusterResourceWithOnlyReqFieldAWS is a Terraform acceptance test that covers the lifecycle of a cluster resource
// creation, retrieval, and update. It focuses on a cluster with only the required fields specified and uses an AWS cloud provider.
//
// The test configures and verifies the following aspects of the cluster resource:
// - Creation of the cluster with minimal required fields.
// - Retrieval and verification of the cluster attributes.
// - Import state testing for the created cluster.
// - Update of the cluster by modifying various fields.
// - Verification of the updated cluster attributes.
func TestAccClusterResourceWithOnlyReqFieldAWS(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.250.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(1 * time.Second)
				},
				Config: testAccClusterResourceConfigWithOnlyReqField(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "io2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),

					//When the cluster is created for the first time, the ETag of the created cluster is 5
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 5"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update number of nodes, compute type, disk size and type, cluster name, support plan, time zone and description from empty string,
			// and Read testing
			{
				Config: testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnlyAndIfMatch(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster Update"),
					resource.TestCheckResourceAttr(resourceReference, "description", "Cluster Updated."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "8"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "32"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "IST"),
				),
			},
			{
				Config: testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnly(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster Update 2"),
					resource.TestCheckResourceAttr(resourceReference, "description", "Cluster Updated."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "8"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "32"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "IST"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// TestAccClusterResourceWithOptionalFieldAWS is a Terraform acceptance test that covers the lifecycle of a cluster resource
// creation, retrieval, and import testing. It focuses on a cluster with both required and optional fields specified and uses
// an AWS cloud provider.
//
// The test configures and verifies the following aspects of the cluster resource:
// - Creation of the cluster with required and optional fields.
// - Retrieval and verification of the cluster attributes.
// - Import state testing for the created cluster.
func TestAccClusterResourceWithOptionalFieldAWS(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.251.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(1 * time.Minute)
				},
				Config: testAccClusterResourceConfigWithAllField(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

// TestAccClusterResourceAzure is a Terraform acceptance test that covers the lifecycle of a cluster resource
// creation, retrieval, and update for an Azure cloud provider.
//
// The test configures and verifies the following aspects of the cluster resource:
// - Creation of the cluster with various fields
// - Retrieval and verification of the cluster attributes.
// - Import state testing for the created cluster.
// - Multiple updates to the cluster, including changes to disk types and horizontal scaling.
func TestAccClusterResourceAzure(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.252.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(2 * time.Minute)
				},
				Config: testAccClusterResourceConfigAzure(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "azure"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "eastus"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "1024"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "Ultra"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "2000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},

			{
				Config: testAccClusterResourceConfigAzureUpdateToDiskTypeP6(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "azure"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "eastus"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "32"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "64"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "P6"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "240"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},

			{
				Config: testAccClusterResourceConfigAzureUpdateToDiskP30(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "azure"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "eastus"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "1024"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "P30"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "5000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},

			{
				Config: testAccClusterResourceConfigAzureUpdateToDiskUltraAndHorizontalScaling(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "azure"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "eastus"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "1024"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "Ultra"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "2000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "64"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "P6"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.iops", "240"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "ET"),
				),
			},
		},
	})
}

// TestAccClusterResourceGCP is a Terraform acceptance test that covers the lifecycle of a cluster resource
// creation, retrieval, and update for a GCP (Google Cloud Platform) cloud provider.
//
// The test configures and verifies the following aspects of the cluster resource:
// - Creation of the cluster with various fields
// - Retrieval and verification of the cluster attributes.
// - Import state testing for the created cluster.
// - An update to the cluster, including changes to horizontal scaling.
func TestAccClusterResourceGCP(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.253.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(3 * time.Minute)
				},
				Config: testAccClusterResourceConfigGCP(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "gcp"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "64"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "pd-ssd"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccClusterResourceConfigGCPUpdateWithHorizontalScaling(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster"),
					resource.TestCheckResourceAttr(resourceReference, "description", "My first test cluster for multiple services."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "gcp"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.1"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "8"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "pd-ssd"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "52"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "pd-ssd"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "ET"),
				),
			},
		},
	})
}

// TestAccClusterResourceWithOptionalFieldAWSInvalidScenario is a Terraform acceptance test that covers an invalid scenario
// during the creation of a cluster resource with optional fields for an AWS (Amazon Web Services) cloud provider.
//
// The test aims to validate that an error is correctly returned when an invalid disk type ("gp2") is provided in the cluster configuration.
// This scenario is expected to fail with the error message matching the regular expression "The disk type provided, gp2, is not valid".
func TestAccClusterResourceWithOptionalFieldAWSInvalidScenario(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.251.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(4 * time.Minute)
				},
				Config:      testAccClusterResourceConfigWithAllFieldInvalidScenario(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("gp2, is not valid"),
			},
		},
	})
}

// TestAccClusterResourceForGCPWithIOPSFieldPopulatedInvalidScenario is a Terraform acceptance test that focuses on
// testing a failure scenario where the IOPS field is set while creating a GCP (Google Cloud Platform) cluster resource.
// The aim of this test is to ensure that the creation of the cluster fails as expected.
func TestAccClusterResourceForGCPWithIOPSFieldPopulatedInvalidScenario(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.248.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					//time.Sleep(4 * time.Minute)
				},
				Config:      testAccClusterResourceForGCPWithIOPSFieldPopulatedInvalidScenarioConfig(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				ExpectError: regexp.MustCompile("Could not create cluster, unexpected error: iops for gcp cluster cannot be"),
			},
		},
	})
}

// TestAccClusterResourceWithConfigurationTypeFieldAdded is a Terraform acceptance test that validates
// the creation of a cluster resource with the addition of the "configuration_type" field set to "singleNode"
// for an AWS (Amazon Web Services) cloud provider.
//
// This test ensures that a cluster resource can be successfully created with the specified configuration type.
func TestAccClusterResourceWithConfigurationTypeFieldAdded(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.249.250.0/23"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(5 * time.Minute)
				},
				Config: testAccClusterResourceConfigWithConfigurationTypeFieldAdded(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "singleNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "8"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "single"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),
				),
			},
			//// ImportState testing
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateClusterImportIdForResource(resourceReference),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

// TestAccClusterResourceNotFound is a Terraform acceptance test that simulates the scenario where a cluster is created
// from Terraform, but it is deleted by a REST API call and the deletion is successful. Then, updating the cluster via Terraform
// should not cause any issues and should create a new cluster with the updated configuration.
//
// This test ensures that Terraform can handle the scenario where the original cluster no longer exists and can
// create a new cluster with the specified configuration when updating.
func TestAccClusterResourceNotFound(t *testing.T) {
	resourceName := "acc_cluster_" + acctest.GenerateRandomResourceName()
	resourceReference := "couchbase-capella_cluster." + resourceName
	projectResourceName := "acc_project_" + acctest.GenerateRandomResourceName()
	projectResourceReference := "couchbase-capella_project." + projectResourceName
	cidr := "10.254.250.0/23"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: func() {
					time.Sleep(6 * time.Minute)
				},
				Config: testAccClusterResourceConfigWithOnlyReqField(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", resourceName),
					resource.TestCheckResourceAttr(resourceReference, "description", ""),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "50"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "io2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3000"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.1", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.2", "query"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "developer pro"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "PT"),

					//When the cluster is created for the first time, the ETag of the created cluster is 5
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 5"),

					//Delete the cluster from the server and wait until the deletion is successful.
					testAccDeleteClusterResource(resourceReference),
				),
				ExpectNonEmptyPlan: true,
				RefreshState:       false,
			},
			// Here, we are attempting an update, but since the original cluster has been deleted, it should create a new cluster.
			{
				Config: testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnly(acctest.Cfg, resourceName, projectResourceName, projectResourceReference, cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccExistsClusterResource(resourceReference),
					resource.TestCheckResourceAttr(resourceReference, "name", "Terraform Acceptance Test Cluster Update 2"),
					resource.TestCheckResourceAttr(resourceReference, "description", "Cluster Updated."),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.type", "aws"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceReference, "cloud_provider.cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "couchbase_server.version", "7.2"),
					resource.TestCheckResourceAttr(resourceReference, "configuration_type", "multiNode"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.cpu", "8"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.compute.ram", "32"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.num_of_nodes", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.#", "2"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.0", "index"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.1.services.1", "query"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.cpu", "4"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.compute.ram", "16"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.storage", "51"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.type", "gp3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.node.disk.iops", "3001"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.num_of_nodes", "3"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.#", "1"),
					resource.TestCheckResourceAttr(resourceReference, "service_groups.0.services.0", "data"),
					resource.TestCheckResourceAttr(resourceReference, "availability.type", "multi"),
					resource.TestCheckResourceAttr(resourceReference, "support.plan", "enterprise"),
					resource.TestCheckResourceAttr(resourceReference, "support.timezone", "IST"),

					//When the cluster is created for the first time, the ETag of the created cluster is 5
					resource.TestCheckResourceAttr(resourceReference, "etag", "Version: 5"),
				),
			},
		},
	})
}

// testAccClusterResourceConfigWithOnlyReqField generates a Terraform configuration string for testing an acceptance test
// scenario for creating a cluster with only required fields.
func testAccClusterResourceConfigWithOnlyReqField(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 50
          type    = "io2"
          iops    = 3000
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

// testAccClusterResourceConfigWithAllField generates a Terraform configuration string for testing an acceptance test
// scenario for creating a cluster with all possible fields.
func testAccClusterResourceConfigWithAllField(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
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
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    },
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
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigWithAllFieldInvalidScenario generates a Terraform configuration string for testing an
// acceptance test scenario where a cluster is created with all possible fields, but with an invalid disk type.
func testAccClusterResourceConfigWithAllFieldInvalidScenario(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
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
          storage = 50
          type    = "gp2"
          iops    = 3000
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    },
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
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigAzure generates a Terraform configuration string for testing an
// acceptance test scenario where a cluster is created with all possible fields using Azure as the cloud provider.
func testAccClusterResourceConfigAzure(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
  }
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 1024,
          type    = "Ultra"
          iops    = 2000
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

// testAccClusterResourceConfigWithConfigurationTypeFieldAdded generates a Terraform configuration string for testing an
// acceptance test scenario where a cluster is created with the "configuration_type" field set to "singleNode".
func testAccClusterResourceConfigWithConfigurationTypeFieldAdded(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  name            = "%[2]s"
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  configuration_type = "singleNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 2
          ram = 8
        }
        disk = {
          storage = 50
          type    = "gp3"
          iops    = 3000
        }
      }
      num_of_nodes = 1
      services     = ["data", "index", "query"]
    }
  ]
  availability = {
    "type" : "single"
  }
  support = {
    plan     = "developer pro"
    timezone = "PT"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigAzureUpdateToDiskTypeP6 generates a Terraform configuration string for testing an acceptance test scenario
// where a cluster resource is updated to change the disk type to "P6".
func testAccClusterResourceConfigAzureUpdateToDiskTypeP6(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 32
        }
        disk = {
          type    = "P6"
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

// testAccClusterResourceConfigAzureUpdateToDiskP30 generates a Terraform configuration string for testing an acceptance test scenario
// where a cluster resource is updated to change the disk type to "P30".
func testAccClusterResourceConfigAzureUpdateToDiskP30(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type    = "P30"
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

// testAccClusterResourceConfigAzureUpdateToDiskUltraAndHorizontalScaling generates a Terraform configuration string for testing an acceptance test scenario
// where a cluster resource is updated to:
// - Change one of the disk types to "Ultra"
// - Horizontal Scaling
// - Change the plan to "enterprise".
// - Change the timezone to "ET".
func testAccClusterResourceConfigAzureUpdateToDiskUltraAndHorizontalScaling(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  couchbase_server = {
    version = "7.1"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 1024,
          type    = "Ultra"
          iops    = 2000
        }
      }
      num_of_nodes = 3
      services     = ["data"]
    },
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          type    = "P6"
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    }
  ]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "ET"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigGCP generates a Terraform configuration string for testing an acceptance test scenario
// where a cluster resource is created with the GCP (Google Cloud Platform) cloud provider.
func testAccClusterResourceConfigGCP(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster"  "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
	type = "gcp",
	region = "us-east1",
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
  } 
  service_groups = [
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					storage = 64,
					type = "pd-ssd"
				}
			}
			num_of_nodes = 4
			services = ["data", "query", "index"]
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

// testAccClusterResourceForGCPWithIOPSFieldPopulatedInvalidScenarioConfig generates a Terraform configuration string for testing an acceptance test scenario
// where a we are trying to create cluster resource for the GCP (Google Cloud Platform) cloud provider while populating iops field.
func testAccClusterResourceForGCPWithIOPSFieldPopulatedInvalidScenarioConfig(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster"  "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
	type = "gcp",
	region = "us-east1",
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
  } 
  service_groups = [
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					storage = 64,
					type    = "pd-ssd"
					iops    = 3000
				}
			},
			num_of_nodes = 4
			services = ["data", "query", "index"]
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

// testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnly generates a Terraform configuration string for testing an acceptance test scenario
// where an existing cluster resource is updated where cluster is created with required fields only.
func testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnly(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  name            = "Terraform Acceptance Test Cluster Update 2"
  description     = "Cluster Updated."
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 8
          ram = 32
        }
        disk = {
          storage = 51
          type    = "gp3"
          iops    = 3001
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    },
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 51
          type    = "gp3"
          iops    = 3001
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
    timezone = "IST"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnlyAndIfMatch generates a Terraform configuration string for testing an acceptance test scenario
// where an existing cluster resource is updated where cluster is created with required fields only and if match flag is provided.
func testAccClusterResourceConfigUpdateWhenClusterCreatedWithReqFieldOnlyAndIfMatch(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
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
  name            = "Terraform Acceptance Test Cluster Update"
  description     = "Cluster Updated."
  cloud_provider = {
    type   = "aws"
    region = "us-east-1"
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  service_groups = [
    {
      node = {
        compute = {
          cpu = 8
          ram = 32
        }
        disk = {
          storage = 51
          type    = "gp3"
          iops    = 3001
        }
      }
      num_of_nodes = 2
      services     = ["index", "query"]
    },
    {
      node = {
        compute = {
          cpu = 4
          ram = 16
        }
        disk = {
          storage = 51
          type    = "gp3"
          iops    = 3001
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
    timezone = "IST"
  }
  if_match = 5
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// testAccClusterResourceConfigGCPUpdateWithHorizontalScaling generates a Terraform configuration string for testing an acceptance test scenario
// where an existing GCP (Google Cloud Platform) cluster resource is updated with horizontal scaling.
func testAccClusterResourceConfigGCPUpdateWithHorizontalScaling(cfg, resourceName, projectResourceName, projectResourceReference, cidr string) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_project" "%[3]s" {
    organization_id = var.organization_id
	name            = "acc_test_project_name"
	description     = "description"
}

resource "couchbase-capella_cluster"  "%[2]s" {
  organization_id = var.organization_id
  project_id      = %[4]s.id
  name            = "Terraform Acceptance Test Cluster"
  description     = "My first test cluster for multiple services."
  cloud_provider = {
	type = "gcp",
	region = "us-east1",
    cidr   = "%[5]s"
  }
  configuration_type = "multiNode"
  couchbase_server = {
    version = "7.1"
  } 
  service_groups = [
		{
			node = {
				compute = {
					cpu = 8
					ram = 16
				}
				disk = {
					storage = 51,
					type = "pd-ssd"
				}
			}
			num_of_nodes = 3
			services = ["data"]
		},
		{
			node = {
				compute = {
					cpu = 4
					ram = 16
				}
				disk = {
					storage = 52,
					type = "pd-ssd"
				}
			}
			num_of_nodes = 2
			services = ["query", "index"]
		}
	]
  availability = {
    "type" : "multi"
  }
  support = {
    plan     = "enterprise"
    timezone = "ET"
  }
}
`, cfg, resourceName, projectResourceName, projectResourceReference, cidr)
}

// This function takes a resource reference string and returns a resource.TestCheckFunc. The returned function, when used
// in Terraform acceptance tests, ensures the successful deletion of the specified cluster resource. It retrieves
// the resource by name from the Terraform state, initiates the deletion, checks the status of the deletion, and
// confirms that the resource no longer exists. If the resource is successfully deleted, it returns nil; otherwise,
// it returns an error.
func testAccDeleteClusterResource(resourceReference string) resource.TestCheckFunc {
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

		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		err = deleteClusterFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		fmt.Printf("delete initiated")
		err = checkClusterStatus(data, context.Background(), rawState["organization_id"], rawState["project_id"], rawState["id"])
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			return fmt.Errorf(errString)
		}
		fmt.Printf("successfully deleted")
		return nil
	}
}

// deleteClusterFromServer deletes cluster from server.
func deleteClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) error {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
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

// checkClusterStatus checks the current state of cluster.
func checkClusterStatus(data *providerschema.Data, ctx context.Context, organizationId, projectId, ClusterId string) error {
	var (
		clusterResp *clusterapi.GetClusterResponse
		err         error
	)

	// Assuming 60 minutes is the max time deployment takes, can change after discussion
	const timeout = time.Minute * 60

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	const sleep = time.Second * 3
	timer := time.NewTimer(2 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return errors.ErrClusterCreationTimeoutAfterInitiation
		case <-timer.C:
			clusterResp, err = retrieveClusterFromServer(data, organizationId, projectId, ClusterId)
			switch err {
			case nil:
				if clusterapi.IsFinalState(clusterResp.CurrentState) {
					return nil
				}
				const msg = "waiting for cluster to complete the execution"
				tflog.Info(ctx, msg)
			default:
				return err
			}
			timer.Reset(sleep)
		}
	}
}

// generateClusterImportIdForResource generates a cluster import ID based on the provided resource reference
// and the attributes in the Terraform state.
//
// This function takes a resource reference as input and returns a function of type `resource.ImportStateIdFunc`.
// The generated import ID is in the format "id=<value>,organization_id=<value>,project_id=<value>".
func generateClusterImportIdForResource(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		fmt.Printf("raw state %s", rawState)
		return fmt.Sprintf("id=%s,organization_id=%s,project_id=%s", rawState["id"], rawState["organization_id"], rawState["project_id"]), nil
	}
}

func testAccDeleteCluster(clusterResourceReference, projectResourceReference string) resource.TestCheckFunc {
	log.Println("Deleting the cluster")
	return func(s *terraform.State) error {
		var clusterState, projectState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
			}
		}
		data, err := acctest.TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", host, orgid, projectState["id"], clusterState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
		_, err = data.Client.ExecuteWithRetry(
			context.Background(),
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}

		//get the cluster state
		err = checkClusterStatus(data, context.Background(), orgid, projectState["id"], clusterState["id"])
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if !resourceNotFound {
			return fmt.Errorf(errString)
		}
		fmt.Printf("successfully deleted")
		return nil

	}
}

// retrieveClusterFromServer checks cluster exists in server.
func retrieveClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.Client.ExecuteWithRetry(
		context.Background(),
		cfg,
		nil,
		data.Token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	clusterResp := clusterapi.GetClusterResponse{}
	err = json.Unmarshal(response.Body, &clusterResp)
	if err != nil {
		return nil, err
	}
	clusterResp.Etag = response.Response.Header.Get("ETag")
	return &clusterResp, nil
}

// This function takes a resource reference string and returns a resource.TestCheckFunc. The returned function, when used
// in Terraform acceptance tests, ensures that the specified cluster resource exists in the Terraform state. It retrieves
// the resource by name from the Terraform state and checks its existence. If the resource exists, it returns nil; otherwise,
// it returns an error.
func testAccExistsClusterResource(resourceReference string) resource.TestCheckFunc {
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
		_, err = retrieveClusterFromServer(data, rawState["organization_id"], rawState["project_id"], rawState["id"])
		if err != nil {
			return err
		}
		return nil
	}
}
