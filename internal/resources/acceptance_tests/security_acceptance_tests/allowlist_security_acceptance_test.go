package security_acceptance_tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	acctest "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAllowListOrgOwner(t *testing.T) {
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// IP with required fields
			{
				PreConfig: func() { testAccCreateOrgAPI("organizationOwner") },
				Config:    testAccAddIpWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_req", "cidr", "10.1.1.0/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
				),
			},
		},
	})
}

func TestAccAllowListOrgMember(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("organizationMember") },
				Config:      testAccAddIpWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListProjCreator(t *testing.T) {
	tempId := os.Getenv("TF_VAR_auth_token")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig:   func() { testAccCreateOrgAPI("projectCreator") },
				Config:      testAccAddIpWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListProjOwner(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig: func() { testAccCreateProjAPI("organizationMember", projId, "projectOwner") },
				Config:    testAccAddIpWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_req", "cidr", "10.1.1.0/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
				),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListProjManager(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig: func() { testAccCreateProjAPI("organizationMember", projId, "projectManager") },
				Config:    testAccAddIpWithReqFields(&testCfg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("capella_allowlist.add_allowlist_req", "cidr", "10.1.1.0/32"),
					resource.TestCheckResourceAttrSet("capella_allowlist.add_allowlist_req", "id"),
				),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListProjViewer(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectViewer") },
				Config:      testAccAddIpWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListDatabaseReaderWriter(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReaderWriter") },
				Config:      testAccAddIpWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func TestAccAllowListDatabaseReader(t *testing.T) {

	tempId := os.Getenv("TF_VAR_auth_token")
	projId := os.Getenv("TF_VAR_project_id")
	testCfg := acctest.Cfg
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//IP with required fields
			{
				PreConfig:   func() { testAccCreateProjAPI("organizationMember", projId, "projectDataReader") },
				Config:      testAccAddIpWithReqFields(&testCfg),
				ExpectError: regexp.MustCompile("Access Denied"),
			},
			{
				// Dummy Test Step to set the Auth token to its original value
				PreConfig: func() { testSetAuthToken(tempId) },
				Config:    testCfg,
			},
		},
	})
}

func testSetAuthToken(tempId string) resource.TestCheckFunc {
	fmt.Println("The test set auth token is here")
	os.Setenv("TF_VAR_auth_token", tempId)
	return func(state *terraform.State) error {
		return nil
	}
}

func testAccAddIpWithReqFields(cfg *string) string {

	*cfg = fmt.Sprintf(`
%[1]s

output "add_allowlist_req"{
  value = capella_allowlist.add_allowlist_req
}

resource "capella_allowlist" "add_allowlist_req" {
  organization_id = var.organization_id
  project_id      = var.project_id
  cluster_id      = var.cluster_id
  cidr            = "10.1.1.0/32"
}

`, *cfg)
	return *cfg
}
