package testing

import (
	"math/rand"
	"os"

	"testing"
	"time"

	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/errors"
	"terraform-provider-capella/internal/provider"
	providerschema "terraform-provider-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
