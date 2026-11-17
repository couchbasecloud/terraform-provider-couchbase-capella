package acceptance_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
)

// TestAccAppServicePublicAccess restricts public access of the shared app service, imports it, removes the
// restriction, then restricts it again. Destroy runs while public access is restricted, so CheckDestroy can verify
// that destroy does not change public access. Cleanup restores the public access that the app service had before.
func TestAccAppServicePublicAccess(t *testing.T) {
	restoreAppServicePublicAccess(t)

	resourceName := randomStringWithPrefix("tf_acc_app_service_public_access_")
	resourceReference := "couchbase-capella_app_service_public_access." + resourceName

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: globalProtoV6ProviderFactory,
		CheckDestroy:             testAccCheckAppServicePublicAccess(appservice.PublicAccessRestricted),
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
				Config: testAccAppServicePublicAccessConfig(resourceName, "restricted"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceReference, "public_access", "restricted"),
					resource.TestCheckResourceAttr(resourceReference, "public_access_status", "restricted"),
				),
			},
			{
				Config:             testAccAppServicePublicAccessConfig(resourceName, "restricted"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// restoreAppServicePublicAccess registers a cleanup that restores the current public access of the shared app
// service, so that the test leaves the fixture as it found it even if the test fails.
func restoreAppServicePublicAccess(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	initial, err := getAppServicePublicAccess(ctx)
	if err != nil {
		t.Fatalf("failed to read app service public access: %v", err)
	}
	if initial != appservice.PublicAccessRestricted && initial != appservice.PublicAccessUnrestricted {
		t.Fatalf("app service public access is %q, expected %q or %q",
			initial, appservice.PublicAccessRestricted, appservice.PublicAccessUnrestricted)
	}

	t.Cleanup(func() {
		if err := setAppServicePublicAccess(ctx, initial); err != nil {
			t.Errorf("failed to restore app service public access to %q: %v", initial, err)
		}
	})
}

// setAppServicePublicAccess waits for any public access change in progress to complete, requests publicAccess if
// the app service does not have it, then waits until the app service reaches it.
func setAppServicePublicAccess(ctx context.Context, publicAccess appservice.PublicAccess) error {
	const maxWaitTime = 30 * time.Minute
	const checkInterval = 10 * time.Second

	requested := false
	for deadline := time.Now().Add(maxWaitTime); time.Now().Before(deadline); time.Sleep(checkInterval) {
		current, err := getAppServicePublicAccess(ctx)
		switch {
		case err != nil:
			return err
		case current == publicAccess:
			return nil
		case current.IsFailed():
			return fmt.Errorf("public access change failed with status %q", current)
		case !requested && (current == appservice.PublicAccessRestricted || current == appservice.PublicAccessUnrestricted):
			if err := requestAppServicePublicAccess(ctx, publicAccess); err != nil {
				return err
			}
			requested = true
		}
	}

	return fmt.Errorf("timeout waiting for app service public access to be %q", publicAccess)
}

func requestAppServicePublicAccess(ctx context.Context, publicAccess appservice.PublicAccess) error {
	method := http.MethodDelete
	if publicAccess == appservice.PublicAccessRestricted {
		method = http.MethodPost
	}

	url := fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/publicAccessRestriction",
		globalHost,
		globalOrgId,
		globalProjectId,
		globalClusterId,
		globalAppServiceId,
	)
	cfg := api.EndpointCfg{Url: url, Method: method, SuccessStatus: http.StatusAccepted}

	_, err := globalClient.ExecuteWithRetry(ctx, cfg, nil, globalToken, nil)
	return err
}

func getAppServicePublicAccess(ctx context.Context) (appservice.PublicAccess, error) {
	appServiceResponse, err := getAppService(ctx, globalClient)
	if err != nil {
		return "", err
	}
	if appServiceResponse.PublicAccess == nil {
		return "", errors.New("app service response has no public access")
	}

	return *appServiceResponse.PublicAccess, nil
}

// testAccCheckAppServicePublicAccess verifies through the API that the shared app service has publicAccess.
func testAccCheckAppServicePublicAccess(publicAccess appservice.PublicAccess) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		current, err := getAppServicePublicAccess(context.Background())
		if err != nil {
			return fmt.Errorf("failed to read app service public access: %w", err)
		}
		if current != publicAccess {
			return fmt.Errorf("expected app service public access %q, got %q", publicAccess, current)
		}

		return nil
	}
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
