package datasources

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschemafw "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	opeCapacityOrgID      = "5ae82d34-5da8-4313-af5f-90ea4142a4c6"
	opeCapacityProjectID  = "8b1c2d3e-4f5a-4b6c-8d7e-9f0a1b2c3d4e"
	opeCapacityClusterID  = "0c1d2e3f-4a5b-4c6d-8e7f-a0b1c2d3e4f5"
	opeCapacityDataSource = "data.couchbase-capella_outbound_private_endpoint_capacity.test"
)

type opeCapacityTestProvider struct {
	data *providerschema.Data
}

func (p *opeCapacityTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "couchbase-capella"
}

func (p *opeCapacityTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschemafw.Schema{}
}

func (p *opeCapacityTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.data
}

func (p *opeCapacityTestProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *opeCapacityTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewOutboundPrivateEndpointCapacity}
}

// opeCapacityFactories serves the capacity route from a fake that answers with code and body,
// and fails any other request.
func opeCapacityFactories(t *testing.T, code int, body string) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	wantPath := fmt.Sprintf("/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/capacity",
		opeCapacityOrgID, opeCapacityProjectID, opeCapacityClusterID)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != wantPath {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	data := &providerschema.Data{ClientV1: &api.Client{Client: srv.Client()}, HostURL: srv.URL}
	return map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(&opeCapacityTestProvider{data: data}),
	}
}

func TestOutboundPrivateEndpointCapacity(t *testing.T) {
	config := fmt.Sprintf(`
data "couchbase-capella_outbound_private_endpoint_capacity" "test" {
  organization_id = %q
  project_id      = %q
  cluster_id      = %q
}
`, opeCapacityOrgID, opeCapacityProjectID, opeCapacityClusterID)

	// Terraform wraps diagnostics, so words are matched across any whitespace.
	const offHint = `The\s+capacity\s+can\s+only\s+be\s+read\s+while\s+the\s+cluster\s+is\s+on`

	tests := []struct {
		name        string
		code        int
		body        string
		check       tfresource.TestCheckFunc
		expectError *regexp.Regexp
	}{
		{
			name: "some used",
			code: http.StatusOK,
			body: `{"clusterEndpointLimit": {"used": 1, "limit": 3}}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "organization_id", opeCapacityOrgID),
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "project_id", opeCapacityProjectID),
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "cluster_id", opeCapacityClusterID),
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "cluster_endpoint_limit.used", "1"),
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "cluster_endpoint_limit.limit", "3"),
			),
		},
		{
			// A cluster with no endpoints reads 0 rather than null.
			name: "none used, raised limit",
			code: http.StatusOK,
			body: `{"clusterEndpointLimit": {"used": 0, "limit": 10}}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "cluster_endpoint_limit.used", "0"),
				tfresource.TestCheckResourceAttr(opeCapacityDataSource, "cluster_endpoint_limit.limit", "10"),
			),
		},
		{
			name:        "cluster turned off",
			code:        http.StatusUnprocessableEntity,
			body:        `{"code": 4025, "hint": "The cluster is turned off.", "httpStatusCode": 422, "message": "The cluster is turned off."}`,
			expectError: regexp.MustCompile(`(?s)Error Reading Outbound Private Endpoint Capacity.*` + opeCapacityClusterID + `.*"httpStatusCode":422}\s+` + offHint),
		},
		{
			name: "not found",
			code: http.StatusNotFound,
			body: `{"code": 4025, "hint": "The cluster was not found.", "httpStatusCode": 404, "message": "The cluster was not found."}`,
			// The cluster-off hint is only for a 422.
			expectError: regexp.MustCompile(`(?s)Could not read the outbound private endpoint capacity of cluster\s+` + opeCapacityClusterID + `.*"httpStatusCode":404}\s*$`),
		},
		{
			name:        "malformed body",
			code:        http.StatusOK,
			body:        `{"clusterEndpointLimit": {"used": "one"}}`,
			expectError: regexp.MustCompile(`(?s)Could not unmarshal the outbound private endpoint capacity of cluster\s+` + opeCapacityClusterID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: opeCapacityFactories(t, tt.code, tt.body),
				Steps: []tfresource.TestStep{{
					Config:      config,
					Check:       tt.check,
					ExpectError: tt.expectError,
				}},
			})
		})
	}
}
