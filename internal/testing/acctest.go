package testing

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	apiRequestTimeout = 60 * time.Second
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
}

// TestAccPreCheck You can add code here to run prior to any test case execution, for
// example assertions about the appropriate environment variables being set
// are common to see in a pre-check function.
func TestAccPreCheck(t *testing.T) {
	if os.Getenv("TF_VAR_host") == "" {
		t.Fatalf(errors.ErrTFVarHostIsNotSet.Error())
	}
	if os.Getenv("TF_VAR_auth_token") == "" {
		t.Fatalf(errors.ErrTFVARAuthTokenIsNotSet.Error())
	}
	if os.Getenv("TF_VAR_organization_id") == "" {
		t.Fatalf(errors.ErrTFVAROrganizationIdIsNotSet.Error())
	}
}

// TestClient returns a common Capella client setup needed for the
// sweeper functions.
func TestClient() (*providerschema.Data, error) {
	host := os.Getenv("TF_VAR_host")
	authenticationToken := os.Getenv("TF_VAR_auth_token")

	if host == "" {
		return nil, errors.ErrTFVarHostIsNotSet
	}
	if authenticationToken == "" {
		return nil, errors.ErrTFVARAuthTokenIsNotSet
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
	result := uuid.New().String()
	return result
}

func TestAccWait(duration time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		time.Sleep(duration)
		return nil
	}
}

func GetCIDR(provider string) (string, error) {
	jwt, err := GetJWT()
	orgId := os.Getenv("TF_VAR_organization_id")
	hostName := os.Getenv("TF_VAR_host")
	hostName = strings.Replace(hostName, "cloudapi", "api", 1)
	url := fmt.Sprintf(
		"%s/v2/organizations/%s/clusters/deployment-options?provider=%s",
		hostName,
		orgId,
		provider,
	)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+jwt)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	cidr := struct {
		SuggestedCidr string `json:"suggestedCidr"`
	}{}
	json.Unmarshal(body, &cidr)
	return cidr.SuggestedCidr, nil
}

func GetJWT() (string, error) {

	hostName := os.Getenv("TF_VAR_host")
	hostName = strings.Replace(hostName, "cloudapi", "api", 1)
	url := hostName + "/sessions"
	// Username and password for Basic Authentication
	username := os.Getenv("CAPELLA_USERNAME")
	password := os.Getenv("CAPELLA_PASSWORD")

	// Create a Basic Authentication token
	authToken := createBasicAuthToken(username, password)

	// Create a new HTTP client
	client := &http.Client{}

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Add Basic Authentication header to the request
	request.Header.Add("Authorization", "Basic "+authToken)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	res := struct {
		Jwt string `json:"jwt"`
	}{}
	json.Unmarshal(body, &res)
	return res.Jwt, nil
}

func createBasicAuthToken(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
