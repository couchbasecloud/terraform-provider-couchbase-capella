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
	cloudAccountsOrgID      = "5ae82d34-5da8-4313-af5f-90ea4142a4c6"
	cloudAccountsDataSource = "data.couchbase-capella_cloud_accounts.test"
)

type cloudAccountsTestProvider struct {
	data *providerschema.Data
}

func (p *cloudAccountsTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "couchbase-capella"
}

func (p *cloudAccountsTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschemafw.Schema{}
}

func (p *cloudAccountsTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.data
}

func (p *cloudAccountsTestProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *cloudAccountsTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewCloudAccounts}
}

// cloudAccountsFactories serves the cloud accounts route from a fake that answers with code and
// body, and fails any other request.
func cloudAccountsFactories(t *testing.T, code int, body string) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	wantPath := fmt.Sprintf("/v4/organizations/%s/cloudAccounts", cloudAccountsOrgID)
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
		"couchbase-capella": providerserver.NewProtocol6WithError(&cloudAccountsTestProvider{data: data}),
	}
}

func TestCloudAccounts(t *testing.T) {
	config := fmt.Sprintf(`
data "couchbase-capella_cloud_accounts" "test" {
  organization_id = %q
}
`, cloudAccountsOrgID)

	tests := []struct {
		name        string
		code        int
		body        string
		check       tfresource.TestCheckFunc
		expectError *regexp.Regexp
	}{
		{
			name: "every cloud",
			code: http.StatusOK,
			body: `{
				"aws-capella-account": "123456789012",
				"gcp-capella-project": "cb-dp-7a9b1c2d",
				"azure-capella-subscription": "capella-azure-sub",
				"azure-capella-subscription-id": "0f1e2d3c-4b5a-4978-8a9b-0c1d2e3f4a5b"
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "organization_id", cloudAccountsOrgID),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "aws_capella_account", "123456789012"),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "gcp_capella_project", "cb-dp-7a9b1c2d"),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "azure_capella_subscription", "capella-azure-sub"),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "azure_capella_subscription_id", "0f1e2d3c-4b5a-4978-8a9b-0c1d2e3f4a5b"),
			),
		},
		{
			// The API leaves the subscription ID out when the organization has no Azure account.
			name: "no azure account",
			code: http.StatusOK,
			body: `{
				"aws-capella-account": "123456789012",
				"gcp-capella-project": "cb-dp-7a9b1c2d",
				"azure-capella-subscription": ""
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "aws_capella_account", "123456789012"),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "gcp_capella_project", "cb-dp-7a9b1c2d"),
				tfresource.TestCheckResourceAttr(cloudAccountsDataSource, "azure_capella_subscription", ""),
				tfresource.TestCheckNoResourceAttr(cloudAccountsDataSource, "azure_capella_subscription_id"),
			),
		},
		{
			name:        "api error",
			code:        http.StatusForbidden,
			body:        `{"code": 1002, "hint": "Access denied.", "httpStatusCode": 403, "message": "Access denied."}`,
			expectError: regexp.MustCompile(`(?s)Error Reading Cloud Accounts.*Could not read the cloud accounts of organization\s+` + cloudAccountsOrgID),
		},
		{
			name:        "malformed body",
			code:        http.StatusOK,
			body:        `{"aws-capella-account": 123}`,
			expectError: regexp.MustCompile(`(?s)Error Reading Cloud Accounts.*Could not unmarshal the cloud accounts of organization\s+` + cloudAccountsOrgID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: cloudAccountsFactories(t, tt.code, tt.body),
				Steps: []tfresource.TestStep{{
					Config:      config,
					Check:       tt.check,
					ExpectError: tt.expectError,
				}},
			})
		})
	}
}
