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
	opesOrgID      = "5ae82d34-5da8-4313-af5f-90ea4142a4c6"
	opesProjectID  = "8b6c4f1d-2a3e-4b5c-9d7e-1f2a3b4c5d6e"
	opesClusterID  = "c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f"
	opesDataSource = "data.couchbase-capella_outbound_private_endpoints.test"
)

type outboundPrivateEndpointsTestProvider struct {
	data *providerschema.Data
}

func (p *outboundPrivateEndpointsTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "couchbase-capella"
}

func (p *outboundPrivateEndpointsTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschemafw.Schema{}
}

func (p *outboundPrivateEndpointsTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.data
}

func (p *outboundPrivateEndpointsTestProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *outboundPrivateEndpointsTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewOutboundPrivateEndpoints}
}

// outboundPrivateEndpointsFactories serves the list route from a fake that answers with code and
// body, and fails any other request.
func outboundPrivateEndpointsFactories(t *testing.T, code int, body string) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	wantPath := fmt.Sprintf("/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints", opesOrgID, opesProjectID, opesClusterID)
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
		"couchbase-capella": providerserver.NewProtocol6WithError(&outboundPrivateEndpointsTestProvider{data: data}),
	}
}

func TestOutboundPrivateEndpoints(t *testing.T) {
	config := fmt.Sprintf(`
data "couchbase-capella_outbound_private_endpoints" "test" {
  organization_id = %q
  project_id      = %q
  cluster_id      = %q
}
`, opesOrgID, opesProjectID, opesClusterID)

	// awsLinked is the spec's AWS list example.
	const awsLinked = `{
		"id": "3f2e1d0c-9b8a-4765-8432-10fedcba9876",
		"name": "cng-prod",
		"cloudProvider": {"type": "aws", "region": "us-east-1"},
		"endpointType": "interface",
		"serviceName": "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0",
		"providerEndpointId": "vpce-0a1b2c3d4e5f67890",
		"status": "linked",
		"providerStatus": "available",
		"hostname": "cng.acme.example.com",
		"isCustomHostname": true,
		"privateDnsEnabled": true,
		"ipAddresses": ["10.0.4.21"],
		"providerConfig": {"aws": {
			"regionalDnsName": "vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com",
			"ipAddresses": [{"ipAddress": "10.0.4.21", "availabilityZone": "us-east-1a"}],
			"availabilityZones": ["us-east-1a"]
		}},
		"audit": {"createdBy": "u1", "createdAt": "2026-09-30T10:00:00Z", "modifiedBy": "u2", "modifiedAt": "2026-09-30T11:20:05Z", "version": 4}
	}`

	// gcpFailed has no hostname, provider endpoint, addresses or provider config yet.
	const gcpFailed = `{
		"id": "7a6b5c4d-3e2f-4a1b-9c8d-7e6f5a4b3c2d",
		"name": "cng-gcp",
		"cloudProvider": {"type": "gcp", "region": "us-east1"},
		"endpointType": "interface",
		"serviceName": "projects/p/regions/us-east1/serviceAttachments/s",
		"status": "failed",
		"failureReason": "serviceNotFound",
		"isCustomHostname": false,
		"privateDnsEnabled": false,
		"audit": {"createdBy": "u1", "createdAt": "2026-09-29T10:00:00Z", "modifiedBy": "u1", "modifiedAt": "2026-09-29T10:00:00Z", "version": 1}
	}`

	tests := []struct {
		name        string
		code        int
		body        string
		check       tfresource.TestCheckFunc
		expectError *regexp.Regexp
	}{
		{
			name: "every field in api order",
			code: http.StatusOK,
			body: `{"data": [` + awsLinked + `,` + gcpFailed + `]}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(opesDataSource, "organization_id", opesOrgID),
				tfresource.TestCheckResourceAttr(opesDataSource, "project_id", opesProjectID),
				tfresource.TestCheckResourceAttr(opesDataSource, "cluster_id", opesClusterID),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.#", "2"),

				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.id", "3f2e1d0c-9b8a-4765-8432-10fedcba9876"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.organization_id", opesOrgID),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.project_id", opesProjectID),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.cluster_id", opesClusterID),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.name", "cng-prod"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.cloud_provider.type", "aws"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.cloud_provider.region", "us-east-1"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.endpoint_type", "interface"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.service_name", "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_endpoint_id", "vpce-0a1b2c3d4e5f67890"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.status", "linked"),
				tfresource.TestCheckNoResourceAttr(opesDataSource, "data.0.failure_reason"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_status", "available"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.hostname", "cng.acme.example.com"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.is_custom_hostname", "true"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.private_dns_enabled", "true"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.ip_addresses.#", "1"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.ip_addresses.0", "10.0.4.21"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.regional_dns_name",
					"vpce-0a1b2c3d4e5f67890-abcd1234.vpce-svc-0123456789abcdef0.us-east-1.vpce.amazonaws.com"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.ip_addresses.#", "1"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.ip_addresses.0.ip_address", "10.0.4.21"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.ip_addresses.0.availability_zone", "us-east-1a"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.availability_zones.#", "1"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.provider_config.aws.availability_zones.0", "us-east-1a"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.audit.created_by", "u1"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.audit.created_at", "2026-09-30 10:00:00 +0000 UTC"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.audit.modified_by", "u2"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.audit.modified_at", "2026-09-30 11:20:05 +0000 UTC"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.0.audit.version", "4"),

				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.id", "7a6b5c4d-3e2f-4a1b-9c8d-7e6f5a4b3c2d"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.cloud_provider.type", "gcp"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.status", "failed"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.failure_reason", "serviceNotFound"),
				tfresource.TestCheckNoResourceAttr(opesDataSource, "data.1.hostname"),
				tfresource.TestCheckNoResourceAttr(opesDataSource, "data.1.provider_endpoint_id"),
				tfresource.TestCheckNoResourceAttr(opesDataSource, "data.1.provider_status"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.is_custom_hostname", "false"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.private_dns_enabled", "false"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.ip_addresses.#", "0"),
				tfresource.TestCheckNoResourceAttr(opesDataSource, "data.1.provider_config.aws.regional_dns_name"),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.1.audit.version", "1"),
			),
		},
		{
			name: "no endpoints",
			code: http.StatusOK,
			body: `{"data": []}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(opesDataSource, "cluster_id", opesClusterID),
				tfresource.TestCheckResourceAttr(opesDataSource, "data.#", "0"),
			),
		},
		{
			name:        "api error",
			code:        http.StatusNotFound,
			body:        `{"code": 4025, "hint": "The requested cluster was not found.", "httpStatusCode": 404, "message": "Cluster not found."}`,
			expectError: regexp.MustCompile(`(?s)Error Reading Outbound Private Endpoints.*Could not read the outbound private endpoints of\s+cluster\s+` + opesClusterID),
		},
		{
			name:        "malformed body",
			code:        http.StatusOK,
			body:        `{"data": [{"id": 123}]}`,
			expectError: regexp.MustCompile(`(?s)Error Reading Outbound Private Endpoints.*Could not unmarshal the outbound private\s+endpoints of cluster\s+` + opesClusterID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: outboundPrivateEndpointsFactories(t, tt.code, tt.body),
				Steps: []tfresource.TestStep{{
					Config:      config,
					Check:       tt.check,
					ExpectError: tt.expectError,
				}},
			})
		})
	}
}
