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
	globalOrgId string

	// these global variables are set by setup().
	globalProjectId       string
	globalClusterId       string
	globalBucketName      = "default"
	globalScopeName       = "_default"
	globalCollectionName  = "_default"
	globalBucketId        string
	globalAppServiceId    string
	globalAppEndpointName = "tf_acc_test_app_endpoint_common"

	// this global variable is set in TestMain.
	globalProviderBlock string

	// globalProtoV6ProviderFactory are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	globalProtoV6ProviderFactory = map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
	}

	// inUseCIDRs caches the CIDR prefixes currently in use by clusters across
	// all projects in the global tenant to ensure acceptance tests use unique
	// CIDR prefixes for cluster creation.
	inUseCIDRs cidrCache
)

const (
	timeout = 60 * time.Second
)
