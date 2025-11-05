package acceptance_tests

import (
	"fmt"
	re "regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAppEndpoint(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_endpoint_")
	resourceReference := "couchbase-capella_app_endpoint." + resourceName
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("tf_acc_app_endpoint_bucket_")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, bucket, "syncFnXattr", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "bucket", bucket),
					resource.TestCheckResourceAttr(resourceReference, "name", epName),
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "true"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "syncFnXattr"),
				),
			},
			{
				Config: testAccAppEndpointResourceConfig(resourceName, epName, bucket, "new_xattr", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "delta_sync_enabled", "false"),
					resource.TestCheckResourceAttr(resourceReference, "user_xattr_key", "new_xattr"),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppEndpointImportId(resourceReference),
				ImportState:       true,
			},
		},
	})
}

func TestAccAppEndpointInexistentCollection(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_endpoint_")
	epName := randomStringWithPrefix("tf_acc_endpoint_")
	bucket := randomStringWithPrefix("bkt_")
	cfg := fmt.Sprintf(`
	%[1]s
	
	resource "couchbase-capella_bucket" "%[2]s_bucket" {
		organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		name           = "%[7]s"
	}
	
	resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
		project_id      = "%[4]s"
		cluster_id      = "%[5]s"
		app_service_id  = "%[6]s"
		bucket          = "%[7]s"
		name            = "%[8]s"
		scopes = {
			"_default" = {
			  collections = {
				"INVALID_COLLLECTION" = {}
			  }
			}
		}
		depends_on = [couchbase-capella_bucket.%[2]s_bucket]
	}`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucket,
		epName)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config:      cfg,
				ExpectError: re.MustCompile("Collection Not Found"),
			},
		},
	})
}

func testAccAppEndpointResourceConfig(resourceName, endpointName, bucketName, userXattr string, deltaSync bool) string {
	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_bucket" "%[2]s_bucket" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	name           = "%[7]s"
}

resource "couchbase-capella_app_endpoint" "%[2]s" {
	organization_id = "%[3]s"
	project_id      = "%[4]s"
	cluster_id      = "%[5]s"
	app_service_id  = "%[6]s"
	bucket          = "%[7]s"
	name            = "%[8]s"
	user_xattr_key  = "%[9]s"
	delta_sync_enabled = %[10]t
	cors = {
		origin = ["*"]
	}
	oidc = [
		{
			issuer   = "https://accounts.google.com"
			client_id = "example_client_id"
		}
	]
	
	scopes = {
		"_default" = {
		  collections = {
			"_default" = {}
		  }
		}
	}
	depends_on = [couchbase-capella_bucket.%[2]s_bucket]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		bucketName,
		endpointName,
		userXattr,
		deltaSync,
	)
}

func generateAppEndpointImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[resourceReference]; ok {
					rawState = v.Primary.Attributes
				}
			}
		}
		// Import uses the endpoint name
		return fmt.Sprintf("app_endpoint_name=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s", rawState["name"], rawState["app_service_id"], rawState["cluster_id"], rawState["project_id"], rawState["organization_id"]), nil
	}
}
