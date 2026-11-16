package acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppServicePublicAccess restricts public access of the shared app service, imports it, then removes the
// restriction. The last step leaves public access unrestricted, because destroy does not change public access.
func TestAccAppServicePublicAccess(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_service_public_access_")
	resourceReference := "couchbase-capella_app_service_public_access." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicePublicAccessConfig(resourceName, "restricted"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "organization_id", globalOrgId),
					resource.TestCheckResourceAttr(resourceReference, "project_id", globalProjectId),
					resource.TestCheckResourceAttr(resourceReference, "cluster_id", globalClusterId),
					resource.TestCheckResourceAttr(resourceReference, "app_service_id", globalAppServiceId),
					resource.TestCheckResourceAttr(resourceReference, "public_access", "restricted"),
					resource.TestCheckResourceAttr(resourceReference, "public_access_status", "restricted"),
				),
			},
			{
				ResourceName:                         resourceReference,
				ImportStateIdFunc:                    generateAppServicePublicAccessImportId(resourceReference),
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "app_service_id",
				ImportStateVerify:                    true,
			},
			{
				Config: testAccAppServicePublicAccessConfig(resourceName, "unrestricted"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "public_access", "unrestricted"),
					resource.TestCheckResourceAttr(resourceReference, "public_access_status", "unrestricted"),
				),
			},
			{
				Config:             testAccAppServicePublicAccessConfig(resourceName, "unrestricted"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TODO: enable once the public access endpoints are public, as the schema builder then infers the public_access
// validator from the OpenAPI spec.
// func TestAccAppServicePublicAccessInvalidValue(t *testing.T) {
// 	resourceName := randomStringWithPrefix("tf_acc_app_service_public_access_")
//
// 	resource.ParallelTest(t, resource.TestCase{
// 		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      testAccAppServicePublicAccessConfig(resourceName, "restricting"),
// 				ExpectError: regexp.MustCompile(`(?s)Attribute public_access value must be one of.*restricted.*unrestricted`),
// 			},
// 		},
// 	})
// }

func testAccAppServicePublicAccessConfig(resourceName, publicAccess string) string {
	return fmt.Sprintf(`
	%[1]s

	resource "couchbase-capella_app_service_public_access" "%[2]s" {
	  organization_id = "%[3]s"
	  project_id      = "%[4]s"
	  cluster_id      = "%[5]s"
	  app_service_id  = "%[6]s"
	  public_access   = "%[7]s"
	}
	`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		publicAccess,
	)
}

func generateAppServicePublicAccessImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[resourceReference]; ok {
				rawState = v.Primary.Attributes
			}
		}
		return fmt.Sprintf(
			"app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["app_service_id"],
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}
