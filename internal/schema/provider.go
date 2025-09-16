package schema

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/apigen"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Config maps provider schema data to a Go type.
type Config struct {
	Host                types.String `tfsdk:"host"`
	AuthenticationToken types.String `tfsdk:"authentication_token"`
}

// Data is provider-defined data, clients, etc. that is passed
// to data sources or resources in the provider that implement the Configure method.
type Data struct {
	ClientV1 *api.Client
	ClientV2 *apigen.ClientWithResponses
	HostURL  string
	Token    string
}
