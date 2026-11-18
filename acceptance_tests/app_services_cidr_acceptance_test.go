package acceptance_tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAppServicesCIDRUpgrade creates an allowed CIDR with the last provider release before access_mode existed,
// then verifies that the current provider plans no changes for the same configuration and keeps access_mode null.
func TestAccAppServicesCIDRUpgrade(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	cidr := fmt.Sprintf("172.16.%d.%d/32", rand.Intn(256), rand.Intn(256)) // #nosec G404
	config := testAccAppServicesCIDRConfig(resourceName, cidr, "")

	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"couchbase-capella": {Source: "couchbasecloud/couchbase-capella", VersionConstraint: "1.13.0"},
				},
				Config: config,
			},
			{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Config:                   config,
				PlanOnly:                 true,
			},
			{
				ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
				Config:                   config,
				Check:                    resource.TestCheckNoResourceAttr(resourceReference, "access_mode"),
			},
		},
	})
}

// TestAccAppServicesCIDRAdminAccessMode changes access_mode between omitted and admin. The API treats both the same, so
// the allowed CIDR updates in place.
func TestAccAppServicesCIDRAdminAccessMode(t *testing.T) {
	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	cidr := fmt.Sprintf("172.16.%d.%d/32", rand.Intn(256), rand.Intn(256)) // #nosec G404
	expectUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRConfig(resourceName, cidr, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceReference, "id"),
					resource.TestCheckNoResourceAttr(resourceReference, "access_mode"),
				),
			},
			{
				Config:           testAccAppServicesCIDRConfig(resourceName, cidr, "admin"),
				ConfigPlanChecks: expectUpdate,
				Check:            resource.TestCheckResourceAttr(resourceReference, "access_mode", "admin"),
			},
			{
				Config:           testAccAppServicesCIDRConfig(resourceName, cidr, ""),
				ConfigPlanChecks: expectUpdate,
				Check:            resource.TestCheckNoResourceAttr(resourceReference, "access_mode"),
			},
		},
	})
}

// TestAccAppServicesCIDRPublicAccessMode restricts public access of the shared app service, adds a public allowed CIDR
// and imports it, changes its access mode to admin and back to public, then removes the public allowed CIDR and the
// restriction in the same change. It does not run in parallel, because it changes public access of the shared app
// service. Cleanup restores the public access that the app service had before.
func TestAccAppServicesCIDRPublicAccessMode(t *testing.T) {
	restoreAppServicePublicAccess(t)

	resourceName := randomStringWithPrefix("tf_acc_app_services_cidr_")
	resourceReference := "couchbase-capella_app_services_cidr." + resourceName
	dataSourceReference := "data.couchbase-capella_app_services_cidr." + resourceName
	allDataSourceReference := dataSourceReference + "_all"
	cidr := fmt.Sprintf("198.51.100.%d/32", rand.Intn(256)) // #nosec G404
	expectReplace := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionReplace),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServicesCIDRPublicAccessModeConfig(resourceName, "restricted", cidr, "public"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "cidr", cidr),
					resource.TestCheckResourceAttr(resourceReference, "access_mode", "public"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceReference, "data.*", map[string]string{
						"cidr":        cidr,
						"access_mode": "public",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(allDataSourceReference, "data.*", map[string]string{
						"cidr":        cidr,
						"access_mode": "public",
					}),
				),
			},
			{
				ResourceName:      resourceReference,
				ImportStateIdFunc: generateAppServicesCIDRImportId(resourceReference),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:           testAccAppServicesCIDRPublicAccessModeConfig(resourceName, "restricted", cidr, "admin"),
				ConfigPlanChecks: expectReplace,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "access_mode", "admin"),
					resource.TestCheckTypeSetElemNestedAttrs(allDataSourceReference, "data.*", map[string]string{
						"cidr":        cidr,
						"access_mode": "admin",
					}),
				),
			},
			{
				Config:           testAccAppServicesCIDRPublicAccessModeConfig(resourceName, "restricted", cidr, "public"),
				ConfigPlanChecks: expectReplace,
				Check:            resource.TestCheckResourceAttr(resourceReference, "access_mode", "public"),
			},
			{
				Config: testAccAppServicesCIDRPublicAccessModeConfig(resourceName, "unrestricted", "", ""),
				Check:  resource.TestCheckResourceAttr(dataSourceReference, "data.#", "0"),
			},
		},
	})
}

// testAccAppServicesCIDRConfig adds an allowed CIDR. An empty accessMode omits access_mode.
func testAccAppServicesCIDRConfig(resourceName, cidr, accessMode string) string {
	accessModeAttr := ""
	if accessMode != "" {
		accessModeAttr = fmt.Sprintf("access_mode     = %q", accessMode)
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  cidr            = "%[7]s"
  %[8]s
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		cidr,
		accessModeAttr,
	)
}

// testAccAppServicesCIDRPublicAccessModeConfig sets public access and lists the public allowed CIDRs, and the admin and
// public allowed CIDRs. A non-empty cidr adds it as an allowed CIDR with the access mode.
func testAccAppServicesCIDRPublicAccessModeConfig(resourceName, publicAccess, cidr, accessMode string) string {
	cidrResource := ""
	dataSourceDependsOn := "couchbase-capella_app_service_public_access." + resourceName
	if cidr != "" {
		cidrResource = fmt.Sprintf(`
resource "couchbase-capella_app_services_cidr" "%[1]s" {
  organization_id = "%[2]s"
  project_id      = "%[3]s"
  cluster_id      = "%[4]s"
  app_service_id  = "%[5]s"
  cidr            = "%[6]s"
  access_mode     = "%[7]s"

  depends_on = [couchbase-capella_app_service_public_access.%[1]s]
}
`, resourceName, globalOrgId, globalProjectId, globalClusterId, globalAppServiceId, cidr, accessMode)
		dataSourceDependsOn = "couchbase-capella_app_services_cidr." + resourceName
	}

	return fmt.Sprintf(`
%[1]s

resource "couchbase-capella_app_service_public_access" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  public_access   = "%[7]s"
}
%[8]s
data "couchbase-capella_app_services_cidr" "%[2]s" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  access_mode     = ["public"]

  depends_on = [%[9]s]
}

data "couchbase-capella_app_services_cidr" "%[2]s_all" {
  organization_id = "%[3]s"
  project_id      = "%[4]s"
  cluster_id      = "%[5]s"
  app_service_id  = "%[6]s"
  access_mode     = ["admin", "public"]

  depends_on = [%[9]s]
}
`,
		globalProviderBlock,
		resourceName,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
		publicAccess,
		cidrResource,
		dataSourceDependsOn,
	)
}

func generateAppServicesCIDRImportId(resourceReference string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		var rawState map[string]string
		for _, m := range state.Modules {
			if v, ok := m.Resources[resourceReference]; ok {
				rawState = v.Primary.Attributes
			}
		}
		return fmt.Sprintf(
			"id=%s,app_service_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			rawState["id"],
			rawState["app_service_id"],
			rawState["cluster_id"],
			rawState["project_id"],
			rawState["organization_id"],
		), nil
	}
}
