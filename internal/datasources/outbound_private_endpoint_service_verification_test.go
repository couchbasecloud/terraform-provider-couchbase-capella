package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	serviceVerificationOrgID      = "5ae82d34-5da8-4313-af5f-90ea4142a4c6"
	serviceVerificationProjectID  = "6be82d34-5da8-4313-af5f-90ea4142a4c6"
	serviceVerificationClusterID  = "7ce82d34-5da8-4313-af5f-90ea4142a4c6"
	serviceVerificationDataSource = "data.couchbase-capella_outbound_private_endpoint_service_verification.test"
)

type serviceVerificationTestProvider struct {
	data *providerschema.Data
}

func (p *serviceVerificationTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "couchbase-capella"
}

func (p *serviceVerificationTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschemafw.Schema{}
}

func (p *serviceVerificationTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.DataSourceData = p.data
}

func (p *serviceVerificationTestProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *serviceVerificationTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewOutboundPrivateEndpointServiceVerification}
}

// serviceVerificationFactories serves the verification route from a fake that checks the request
// carries serviceName, answers with code and body, and fails any other request.
func serviceVerificationFactories(t *testing.T, serviceName string, code int, body string) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	wantPath := fmt.Sprintf("/v4/organizations/%s/projects/%s/clusters/%s/outboundPrivateEndpoints/serviceVerification",
		serviceVerificationOrgID, serviceVerificationProjectID, serviceVerificationClusterID)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != wantPath {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		raw, _ := io.ReadAll(r.Body)
		var req map[string]any
		if err := json.Unmarshal(raw, &req); err != nil || req["serviceName"] != serviceName || len(req) != 1 {
			t.Errorf("unexpected request body %s", raw)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	data := &providerschema.Data{ClientV1: &api.Client{Client: srv.Client()}, HostURL: srv.URL}
	return map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(&serviceVerificationTestProvider{data: data}),
	}
}

func TestOutboundPrivateEndpointServiceVerification(t *testing.T) {
	const (
		awsService   = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
		gcpService   = "projects/acme-network/regions/us-central1/serviceAttachments/cng-gateway"
		azureService = "cng-gateway.4f5e6d7c-1a2b-3c4d-5e6f-7a8b9c0d1e2f.eastus.azure.privatelinkservice"
	)
	ds := serviceVerificationDataSource

	tests := []struct {
		name        string
		serviceName string
		code        int
		body        string
		check       tfresource.TestCheckFunc
		expectError *regexp.Regexp
	}{
		{
			name:        "aws ready",
			serviceName: awsService,
			code:        http.StatusOK,
			body: `{
				"cloudProvider": {"type": "aws", "region": "us-east-1"},
				"serviceName": "` + awsService + `",
				"ready": true,
				"providerConfig": {"aws": {
					"availabilityZones": ["us-east-1a", "us-east-1b", "us-east-1c"],
					"acceptanceRequired": true,
					"serviceOwner": "210987654321",
					"privateDnsName": "*.acme.example.com",
					"privateDnsNameVerification": "verified"
				}},
				"checks": [
					{"name": "format", "result": "passed"},
					{"name": "region", "result": "passed"},
					{"name": "accessGranted", "result": "passed"},
					{"name": "serviceZones", "result": "passed"},
					{"name": "notAlreadyConnected", "result": "passed"},
					{"name": "privateDnsName", "result": "passed"}
				]
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(ds, "organization_id", serviceVerificationOrgID),
				tfresource.TestCheckResourceAttr(ds, "project_id", serviceVerificationProjectID),
				tfresource.TestCheckResourceAttr(ds, "cluster_id", serviceVerificationClusterID),
				tfresource.TestCheckResourceAttr(ds, "service_name", awsService),
				tfresource.TestCheckResourceAttr(ds, "cloud_provider.type", "aws"),
				tfresource.TestCheckResourceAttr(ds, "cloud_provider.region", "us-east-1"),
				tfresource.TestCheckResourceAttr(ds, "ready", "true"),
				tfresource.TestCheckResourceAttr(ds, "checks.#", "6"),
				tfresource.TestCheckResourceAttr(ds, "checks.0.name", "format"),
				tfresource.TestCheckResourceAttr(ds, "checks.0.result", "passed"),
				tfresource.TestCheckNoResourceAttr(ds, "checks.0.message"),
				tfresource.TestCheckResourceAttr(ds, "checks.5.name", "privateDnsName"),
				tfresource.TestCheckResourceAttr(ds, "checks.5.result", "passed"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.availability_zones.#", "3"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.availability_zones.0", "us-east-1a"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.availability_zones.2", "us-east-1c"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.acceptance_required", "true"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.service_owner", "210987654321"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.private_dns_name", "*.acme.example.com"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.private_dns_name_verification", "verified"),
			),
		},
		{
			// No private DNS name: the field is absent and its check only warns.
			name:        "aws without private dns name",
			serviceName: awsService,
			code:        http.StatusOK,
			body: `{
				"cloudProvider": {"type": "aws", "region": "us-east-1"},
				"serviceName": "` + awsService + `",
				"ready": true,
				"providerConfig": {"aws": {
					"availabilityZones": ["us-east-1a"],
					"acceptanceRequired": false,
					"serviceOwner": "210987654321",
					"privateDnsNameVerification": "none"
				}},
				"checks": [
					{"name": "privateDnsName", "result": "warning", "message": "The service has no verified private DNS name."}
				]
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(ds, "ready", "true"),
				tfresource.TestCheckResourceAttr(ds, "checks.0.result", "warning"),
				tfresource.TestCheckResourceAttr(ds, "checks.0.message", "The service has no verified private DNS name."),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.acceptance_required", "false"),
				tfresource.TestCheckNoResourceAttr(ds, "provider_config.aws.private_dns_name"),
				tfresource.TestCheckResourceAttr(ds, "provider_config.aws.private_dns_name_verification", "none"),
			),
		},
		{
			// GCP reports no provider config fields yet, so the object is present with aws null.
			name:        "gcp ready",
			serviceName: gcpService,
			code:        http.StatusOK,
			body: `{
				"cloudProvider": {"type": "gcp", "region": "us-central1"},
				"serviceName": "` + gcpService + `",
				"ready": true,
				"providerConfig": {},
				"checks": [
					{"name": "format", "result": "passed"},
					{"name": "region", "result": "passed"},
					{"name": "accessGranted", "result": "notApplicable", "message": "GCP does not let Capella check the accept list. Make sure the Capella project is on it."},
					{"name": "serviceZones", "result": "notApplicable"},
					{"name": "notAlreadyConnected", "result": "passed"},
					{"name": "privateDnsName", "result": "notApplicable"}
				]
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(ds, "service_name", gcpService),
				tfresource.TestCheckResourceAttr(ds, "cloud_provider.type", "gcp"),
				tfresource.TestCheckResourceAttr(ds, "ready", "true"),
				tfresource.TestCheckResourceAttr(ds, "checks.#", "6"),
				tfresource.TestCheckResourceAttr(ds, "checks.2.name", "accessGranted"),
				tfresource.TestCheckResourceAttr(ds, "checks.2.result", "notApplicable"),
				tfresource.TestCheckResourceAttr(ds, "checks.2.message", "GCP does not let Capella check the accept list. Make sure the Capella project is on it."),
				tfresource.TestCheckNoResourceAttr(ds, "provider_config.aws"),
			),
		},
		{
			// A failed check is a 200: the read succeeds and reports ready false.
			name:        "azure not visible",
			serviceName: azureService,
			code:        http.StatusOK,
			body: `{
				"cloudProvider": {"type": "azure", "region": "eastus"},
				"serviceName": "` + azureService + `",
				"ready": false,
				"checks": [
					{"name": "format", "result": "passed"},
					{"name": "region", "result": "passed"},
					{"name": "accessGranted", "result": "failed", "message": "Add the Capella subscription to the service's visibility."},
					{"name": "serviceZones", "result": "notApplicable"},
					{"name": "notAlreadyConnected", "result": "passed"},
					{"name": "privateDnsName", "result": "notApplicable"}
				]
			}`,
			check: tfresource.ComposeAggregateTestCheckFunc(
				tfresource.TestCheckResourceAttr(ds, "cloud_provider.type", "azure"),
				tfresource.TestCheckResourceAttr(ds, "cloud_provider.region", "eastus"),
				tfresource.TestCheckResourceAttr(ds, "ready", "false"),
				tfresource.TestCheckResourceAttr(ds, "checks.2.result", "failed"),
				tfresource.TestCheckResourceAttr(ds, "checks.2.message", "Add the Capella subscription to the service's visibility."),
				tfresource.TestCheckNoResourceAttr(ds, "provider_config"),
			),
		},
		{
			name:        "api error",
			serviceName: awsService,
			code:        http.StatusUnprocessableEntity,
			body:        `{"code": 4001, "hint": "The cluster is not healthy.", "httpStatusCode": 422, "message": "The cluster is not healthy."}`,
			expectError: regexp.MustCompile(`(?s)Error Verifying Outbound Private Endpoint Service.*Could not verify service\s+` + regexp.QuoteMeta(awsService) + `\s+for\s+cluster\s+` + serviceVerificationClusterID),
		},
		{
			name:        "malformed body",
			serviceName: awsService,
			code:        http.StatusOK,
			body:        `{"ready": "yes"}`,
			expectError: regexp.MustCompile(`(?s)Error Verifying Outbound Private Endpoint Service.*Could not unmarshal the verification of service`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := fmt.Sprintf(`
data "couchbase-capella_outbound_private_endpoint_service_verification" "test" {
  organization_id = %q
  project_id      = %q
  cluster_id      = %q
  service_name    = %q
}
`, serviceVerificationOrgID, serviceVerificationProjectID, serviceVerificationClusterID, tt.serviceName)

			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: serviceVerificationFactories(t, tt.serviceName, tt.code, tt.body),
				Steps: []tfresource.TestStep{{
					Config:      config,
					Check:       tt.check,
					ExpectError: tt.expectError,
				}},
			})
		})
	}
}
