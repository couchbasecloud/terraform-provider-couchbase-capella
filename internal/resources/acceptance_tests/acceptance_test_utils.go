package acceptance_tests

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	acctest "terraform-provider-capella/internal/testing"

	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-capella/internal/api"
	clusterapi "terraform-provider-capella/internal/api/cluster"
	"terraform-provider-capella/internal/provider"
	providerschema "terraform-provider-capella/internal/schema"
)

var (
	apiRequestTimeout = 60 * time.Second
)

const (
	// charSetAlpha is the alphabetical character set for use with
	// RandStringFromCharSet.
	charSetAlpha = "abcdefghijklmnopqrstuvwxyz"

	// Length of the resource name we wish to generate.
	resourceNameLength = 10
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// TestAccPreCheck You can add code here to run prior to any test case execution, for
// example assertions about the appropriate environment variables being set
// are common to see in a pre-check function.
func TestAccPreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_host") == "" {
		t.Fatalf("host not set")
	}
	if os.Getenv("TF_VAR_auth_token") == "" {
		t.Fatalf("auth token not set")
	}
	if os.Getenv("TF_VAR_organization_id") == "" {
		t.Fatalf("organization id not set")
	}
}

// TestClient returns a common Capella client setup needed for the
// sweeper functions.
func TestClient() (*providerschema.Data, error) {
	host := os.Getenv("TF_VAR_host")
	authenticationToken := os.Getenv("TF_VAR_auth_token")

	if host == "" {
		return nil, fmt.Errorf("empty host name")
	}
	if authenticationToken == "" {
		return nil, fmt.Errorf("authentication token not set")
	}

	// Create a new capella client using the configuration values
	providerData := &providerschema.Data{
		HostURL: host,
		Token:   authenticationToken,
		Client:  api.NewClient(apiRequestTimeout),
	}
	return providerData, nil
}

// GenerateRandomResourceName builds a unique-ish resource identifier to use in
// tests.
func GenerateRandomResourceName() string {
	result := make([]byte, resourceNameLength)
	for i := 0; i < resourceNameLength; i++ {
		result[i] = charSetAlpha[randIntRange(0, len(charSetAlpha))]
	}
	return string(result)
}

// randIntRange returns a random integer between min (inclusive) and max
// (exclusive).
func randIntRange(min int, max int) int {
	return rand.Intn(max-min) + min
}

// retrieveClusterFromServer checks cluster exists in server.
func retrieveClusterFromServer(data *providerschema.Data, organizationId, projectId, clusterId string) (*clusterapi.GetClusterResponse, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", data.HostURL, organizationId, projectId, clusterId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := data.Client.Execute(
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

func testAccDeleteAllowIP(clusterResourceReference, projectResourceReference, allowIPResoureceReference string) resource.TestCheckFunc {
	log.Println("deleting the ip")
	return func(s *terraform.State) error {
		var clusterState, projectState, allowListState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[clusterResourceReference]; ok {
					clusterState = v.Primary.Attributes
				}
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
				if v, ok := m.Resources[allowIPResoureceReference]; ok {
					allowListState = v.Primary.Attributes
				}
			}
		}
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s//allowedcidrs/%s", host, orgid, projectState["id"], clusterState["id"], allowListState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		_, err = data.Client.Execute(
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccDeleteProject(projectResourceReference string) resource.TestCheckFunc {
	log.Println("Deleting the project")
	return func(s *terraform.State) error {
		var projectState map[string]string
		for _, m := range s.Modules {
			if len(m.Resources) > 0 {
				if v, ok := m.Resources[projectResourceReference]; ok {
					projectState = v.Primary.Attributes
				}
			}
		}
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s", host, orgid, projectState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
		_, err = data.Client.Execute(
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}
		return nil
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
		data, err := TestClient()
		if err != nil {
			return err
		}
		host := os.Getenv("TF_VAR_host")
		orgid := os.Getenv("TF_VAR_organization_id")
		authToken := os.Getenv("TF_VAR_auth_token")
		url := fmt.Sprintf("%s/v4/organizations/%s/projects/%s/clusters/%s", host, orgid, projectState["id"], clusterState["id"])
		cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusAccepted}
		_, err = data.Client.Execute(
			cfg,
			nil,
			authToken,
			nil,
		)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccWait(duration time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		time.Sleep(duration)
		return nil
	}
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
