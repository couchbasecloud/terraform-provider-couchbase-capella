package acceptance_tests

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
)

var (

	// these global variables are set by env vars.
	globalHost  string
	globalToken string
	Username    string
	Password    string
	globalOrgId string

	// these global variables are set by setup().
	globalProjectId  string
	globalClusterId  string
	globalBucketName = "default"
	globalBucketId   string

	// this global variable is set in TestMain.
	globalProviderBlock string

	// globalProtoV6ProviderFactory are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	globalProtoV6ProviderFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
	}
)

const (
	timeout = 60 * time.Second
)
