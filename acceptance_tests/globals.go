package acceptance_tests

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"time"
)

var (
	ProviderBlock string
	Host          string
	Token         string

	Username string
	Password string

	OrgId     string
	ProjectId string
	ClusterId string

	BucketName = "default"
	BucketId   string

	// TestAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(provider.New()()),
	}
)

const (
	Timeout = 60 * time.Second
)
