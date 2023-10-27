package testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"terraform-provider-capella/internal/api"
	"terraform-provider-capella/internal/provider"
	providerschema "terraform-provider-capella/internal/schema"
	"time"
)

var (
	apiRequestTimeout = 60 * time.Second
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"capella": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// SharedClient returns a common Capella client setup needed for the
// sweeper functions.
func SharedClient(host, authenticationToken string) (*providerschema.Data, error) {
	// Check if 'host' is empty and set it from the environment variable if needed.
	if host == "" {
		host = os.Getenv("TF_VAR_host")
	}

	// Check if 'host' is still empty after trying the environment variable.
	if host == "" {
		return nil, fmt.Errorf("TF_VAR_host is not set")
	}

	// Check if 'authenticationToken' is empty and set it from the environment variable if needed.
	if authenticationToken == "" {
		authenticationToken = os.Getenv("TF_VAR_auth_token")
	}

	// Check if 'authenticationToken' is still empty after trying the environment variable.
	if authenticationToken == "" {
		return nil, fmt.Errorf("TF_VAR_auth_token is not set")
	}

	// Create a new capella client using the configuration values
	providerData := &providerschema.Data{
		HostURL: host,
		Token:   authenticationToken,
		Client:  api.NewClient(apiRequestTimeout),
	}
	return providerData, nil
}
