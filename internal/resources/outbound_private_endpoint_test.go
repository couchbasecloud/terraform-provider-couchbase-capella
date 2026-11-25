package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschemafw "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	outboundpe "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/outbound_private_endpoint"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	opeOrgID     = "5ae82d34-5da8-4313-af5f-90ea4142a4c6"
	opeProjectID = "b02d46ce-c3b1-4a5e-8f3f-ad49112b1dc6"
	opeClusterID = "533bc6d6-1628-49f5-9384-2dd1c5d025fe"

	opeServiceA = "com.amazonaws.vpce.us-east-1.vpce-svc-0123456789abcdef0"
	opeServiceB = "com.amazonaws.vpce.us-east-1.vpce-svc-0fedcba9876543210"

	opeResource = "couchbase-capella_outbound_private_endpoint.test"
)

// fastOutboundPrivateEndpointPolling shrinks the poll and retry timing so waits resolve in milliseconds.
func fastOutboundPrivateEndpointPolling(t *testing.T) {
	t.Helper()
	poll, timeout, wait := outboundPrivateEndpointPollInterval, outboundPrivateEndpointTimeout, outboundPrivateEndpointConflictWait
	outboundPrivateEndpointPollInterval = time.Millisecond
	outboundPrivateEndpointTimeout = 5 * time.Second
	outboundPrivateEndpointConflictWait = time.Millisecond
	t.Cleanup(func() {
		outboundPrivateEndpointPollInterval = poll
		outboundPrivateEndpointTimeout = timeout
		outboundPrivateEndpointConflictWait = wait
	})
}

// fakeOutboundPrivateEndpoints stands in for the v4 outbound private endpoint routes. A created
// endpoint reads pending once, then createStatus. A deleted one reads deleting once, then 404.
type fakeOutboundPrivateEndpoints struct {
	mu sync.Mutex

	endpoints map[string]*outboundpe.GetOutboundPrivateEndpointResponse
	pending   map[string]int
	deleted   map[string]bool
	nextID    int

	// createStatus and createFailureReason are what a new endpoint reads once it leaves pending.
	createStatus        string
	createFailureReason string

	// busy answers this many requests with 409 16017 before handling one.
	busy int

	creates, renames, deletes int
	lastCreate                outboundpe.CreateOutboundPrivateEndpointRequest
}

func newFakeOutboundPrivateEndpoints() *fakeOutboundPrivateEndpoints {
	return &fakeOutboundPrivateEndpoints{
		endpoints:    map[string]*outboundpe.GetOutboundPrivateEndpointResponse{},
		pending:      map[string]int{},
		deleted:      map[string]bool{},
		createStatus: outboundpe.StatusPendingAcceptance,
	}
}

func (f *fakeOutboundPrivateEndpoints) handler(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	const prefix = "/v4/organizations/" + opeOrgID + "/projects/" + opeProjectID + "/clusters/" + opeClusterID + "/outboundPrivateEndpoints"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		writeAPIError(w, http.StatusNotFound, 404, "unexpected path "+r.URL.Path)
		return
	}
	id := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, prefix), "/")

	if r.Method != http.MethodGet && f.busy > 0 {
		f.busy--
		writeAPIError(w, http.StatusConflict, outboundpe.ErrCodeOperationInProgress, "another change is in progress")
		return
	}

	switch {
	case r.Method == http.MethodPost && id == "":
		f.create(w, r)
	case r.Method == http.MethodGet && id != "":
		f.get(w, id)
	case r.Method == http.MethodPut && id != "":
		f.rename(w, r, id)
	case r.Method == http.MethodDelete && id != "":
		f.delete(w, id)
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, 405, "unexpected request")
	}
}

func (f *fakeOutboundPrivateEndpoints) create(w http.ResponseWriter, r *http.Request) {
	var req outboundpe.CreateOutboundPrivateEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, 400, err.Error())
		return
	}
	f.creates++
	f.lastCreate = req
	f.nextID++
	id := fmt.Sprintf("3f2e1d0c-9b8a-4765-8432-%012d", f.nextID)

	hostname := id[:8] + ".assigned.example.com"
	if req.Hostname != "" {
		hostname = req.Hostname
	}
	endpointType := req.EndpointType
	if endpointType == "" {
		endpointType = "interface"
	}
	now := time.Date(2026, 11, 23, 10, 0, 0, 0, time.UTC)

	f.endpoints[id] = &outboundpe.GetOutboundPrivateEndpointResponse{
		Id:                 id,
		Name:               req.Name,
		CloudProvider:      outboundpe.CloudProvider{Type: "aws", Region: "us-east-1"},
		EndpointType:       endpointType,
		ServiceName:        req.ServiceName,
		ProviderEndpointId: strPtr("vpce-0a1b2c3d4e5f6789" + fmt.Sprint(f.nextID)),
		Status:             outboundpe.StatusPending,
		Hostname:           strPtr(hostname),
		IsCustomHostname:   req.Hostname != "",
		IpAddresses:        []string{"10.0.4.21"},
		ProviderConfig: &outboundpe.ProviderConfig{AWS: &outboundpe.AWSConfig{
			RegionalDnsName:   strPtr("vpce-1.vpce-svc-1.us-east-1.vpce.amazonaws.com"),
			IpAddresses:       []outboundpe.AWSIpAddress{{IpAddress: "10.0.4.21", AvailabilityZone: "us-east-1a"}},
			AvailabilityZones: []string{"us-east-1a"},
		}},
		Audit: api.CouchbaseAuditData{CreatedAt: now, CreatedBy: "apikey", ModifiedAt: now, ModifiedBy: "apikey", Version: 1},
	}
	f.pending[id] = 1

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(outboundpe.CreateOutboundPrivateEndpointResponse{Id: id})
}

func (f *fakeOutboundPrivateEndpoints) get(w http.ResponseWriter, id string) {
	endpoint, ok := f.endpoints[id]
	if !ok {
		writeAPIError(w, http.StatusNotFound, 16004, "not found")
		return
	}

	switch {
	case f.pending[id] > 0:
		f.pending[id]--
	case endpoint.Status == outboundpe.StatusPending:
		endpoint.Status = f.createStatus
		if f.createFailureReason != "" {
			endpoint.FailureReason = strPtr(f.createFailureReason)
		}
	case endpoint.Status == outboundpe.StatusDeleting && f.deleted[id]:
		delete(f.endpoints, id)
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(endpoint)
	if endpoint.Status == outboundpe.StatusDeleting {
		f.deleted[id] = true
	}
}

func (f *fakeOutboundPrivateEndpoints) rename(w http.ResponseWriter, r *http.Request, id string) {
	endpoint, ok := f.endpoints[id]
	if !ok {
		writeAPIError(w, http.StatusNotFound, 16004, "not found")
		return
	}
	var req outboundpe.UpdateOutboundPrivateEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, 400, err.Error())
		return
	}
	f.renames++
	endpoint.Name = req.Name
	endpoint.Audit.Version++
	w.WriteHeader(http.StatusNoContent)
}

func (f *fakeOutboundPrivateEndpoints) delete(w http.ResponseWriter, id string) {
	endpoint, ok := f.endpoints[id]
	if !ok {
		writeAPIError(w, http.StatusNotFound, 16004, "not found")
		return
	}
	f.deletes++
	endpoint.Status = outboundpe.StatusDeleting
	w.WriteHeader(http.StatusAccepted)
}

// removeAll drops every endpoint, as if they were deleted outside Terraform.
func (f *fakeOutboundPrivateEndpoints) removeAll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	clear(f.endpoints)
}

func (f *fakeOutboundPrivateEndpoints) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.endpoints)
}

func writeAPIError(w http.ResponseWriter, status, code int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(api.Error{Code: code, HttpStatusCode: status, Message: message})
}

// outboundPrivateEndpointTestProvider serves only the outbound private endpoint resource, wired
// to a fake backend, so plan and apply run through Terraform itself.
type outboundPrivateEndpointTestProvider struct {
	data *providerschema.Data
}

func (p *outboundPrivateEndpointTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "couchbase-capella"
}

func (p *outboundPrivateEndpointTestProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschemafw.Schema{}
}

func (p *outboundPrivateEndpointTestProvider) Configure(_ context.Context, _ provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.ResourceData = p.data
}

func (p *outboundPrivateEndpointTestProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{NewOutboundPrivateEndpoint}
}

func (p *outboundPrivateEndpointTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func outboundPrivateEndpointFactories(t *testing.T, fake *fakeOutboundPrivateEndpoints) map[string]func() (tfprotov6.ProviderServer, error) {
	t.Helper()
	data := newTestProviderData(t, fake.handler)
	return map[string]func() (tfprotov6.ProviderServer, error){
		"couchbase-capella": providerserver.NewProtocol6WithError(&outboundPrivateEndpointTestProvider{data: data}),
	}
}

func outboundPrivateEndpointConfig(name, serviceName, hostname string) string {
	hostnameLine := ""
	if hostname != "" {
		hostnameLine = fmt.Sprintf("  hostname        = %q\n", hostname)
	}
	return fmt.Sprintf(`
resource "couchbase-capella_outbound_private_endpoint" "test" {
  organization_id = %q
  project_id      = %q
  cluster_id      = %q
  name            = %q
  service_name    = %q
%s}
`, opeOrgID, opeProjectID, opeClusterID, name, serviceName, hostnameLine)
}

func TestOutboundPrivateEndpointPlanAndApply(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)
	fake := newFakeOutboundPrivateEndpoints()

	var firstID string
	sameID := func(s *terraform.State) error {
		if got := s.RootModule().Resources[opeResource].Primary.ID; got != firstID {
			return fmt.Errorf("id changed from %s to %s on a rename", firstID, got)
		}
		return nil
	}

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: outboundPrivateEndpointFactories(t, fake),
		CheckDestroy: func(*terraform.State) error {
			if n := fake.count(); n != 0 {
				return fmt.Errorf("%d endpoints left after destroy", n)
			}
			return nil
		},
		Steps: []tfresource.TestStep{
			{
				// Create returns once the endpoint leaves pending, with the assigned hostname.
				Config: outboundPrivateEndpointConfig("cng-prod", opeServiceA, ""),
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttr(opeResource, "status", outboundpe.StatusPendingAcceptance),
					tfresource.TestCheckResourceAttr(opeResource, "endpoint_type", "interface"),
					tfresource.TestCheckResourceAttr(opeResource, "is_custom_hostname", "false"),
					tfresource.TestCheckResourceAttr(opeResource, "hostname", "3f2e1d0c.assigned.example.com"),
					tfresource.TestCheckResourceAttr(opeResource, "cloud_provider.type", "aws"),
					tfresource.TestCheckResourceAttr(opeResource, "ip_addresses.0", "10.0.4.21"),
					tfresource.TestCheckResourceAttr(opeResource, "provider_config.aws.ip_addresses.0.availability_zone", "us-east-1a"),
					tfresource.TestCheckResourceAttr(opeResource, "audit.version", "1"),
					func(s *terraform.State) error {
						firstID = s.RootModule().Resources[opeResource].Primary.ID
						return nil
					},
				),
			},
			{
				// A rename updates in place, and the fields a rename can't change keep their values in the plan.
				Config: outboundPrivateEndpointConfig("cng-prod-east", opeServiceA, ""),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(opeResource, tfjsonpath.New("hostname"), knownvalue.StringExact("3f2e1d0c.assigned.example.com")),
						plancheck.ExpectKnownValue(opeResource, tfjsonpath.New("provider_endpoint_id"), knownvalue.StringExact("vpce-0a1b2c3d4e5f67891")),
						plancheck.ExpectKnownValue(opeResource, tfjsonpath.New("is_custom_hostname"), knownvalue.Bool(false)),
						plancheck.ExpectKnownValue(opeResource, tfjsonpath.New("cloud_provider").AtMapKey("type"), knownvalue.StringExact("aws")),
						plancheck.ExpectUnknownValue(opeResource, tfjsonpath.New("status")),
					},
				},
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttr(opeResource, "name", "cng-prod-east"),
					tfresource.TestCheckResourceAttr(opeResource, "audit.version", "2"),
					sameID,
				),
			},
			{
				// The service can't change in place.
				Config: outboundPrivateEndpointConfig("cng-prod-east", opeServiceB, ""),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: tfresource.TestCheckResourceAttr(opeResource, "service_name", opeServiceB),
			},
			{
				// Setting a hostname on an endpoint with an assigned one replaces it.
				Config: outboundPrivateEndpointConfig("cng-prod-east", opeServiceB, "cng.acme.example.com"),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttr(opeResource, "hostname", "cng.acme.example.com"),
					tfresource.TestCheckResourceAttr(opeResource, "is_custom_hostname", "true"),
				),
			},
			{
				// The same config plans nothing.
				Config: outboundPrivateEndpointConfig("cng-prod-east", opeServiceB, "cng.acme.example.com"),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				// Import reads the endpoint back from its IDs.
				ResourceName:      opeResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					id := s.RootModule().Resources[opeResource].Primary.ID
					return fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s", id, opeClusterID, opeProjectID, opeOrgID), nil
				},
			},
			{
				// Removing the hostname replaces the endpoint, which then gets an assigned one.
				Config: outboundPrivateEndpointConfig("cng-prod-east", opeServiceB, ""),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: tfresource.TestCheckResourceAttr(opeResource, "is_custom_hostname", "false"),
			},
			{
				// An endpoint deleted outside Terraform drops out of state and is created again.
				PreConfig: fake.removeAll,
				Config:    outboundPrivateEndpointConfig("cng-prod-east", opeServiceB, ""),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionCreate),
					},
				},
				Check: tfresource.TestCheckResourceAttr(opeResource, "status", outboundpe.StatusPendingAcceptance),
			},
		},
	})

	assert.Equal(t, 1, fake.renames)
}

// TestOutboundPrivateEndpointApplyTimeHostnameReplaces sets hostname from a value only known at
// apply. Update can only rename, so a hostname that may change has to plan a replacement up front.
func TestOutboundPrivateEndpointApplyTimeHostnameReplaces(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)
	fake := newFakeOutboundPrivateEndpoints()

	// terraform_data's output reads as unknown in the plan whenever its input changes.
	config := func(hostname string) string {
		return fmt.Sprintf(`
resource "terraform_data" "hostname" {
  input = %q
}

resource "couchbase-capella_outbound_private_endpoint" "test" {
  organization_id = %q
  project_id      = %q
  cluster_id      = %q
  name            = "cng-prod"
  service_name    = %q
  hostname        = terraform_data.hostname.output
}
`, hostname, opeOrgID, opeProjectID, opeClusterID, opeServiceA)
	}

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: outboundPrivateEndpointFactories(t, fake),
		Steps: []tfresource.TestStep{
			{
				Config: config("cng.acme.example.com"),
				Check:  tfresource.TestCheckResourceAttr(opeResource, "hostname", "cng.acme.example.com"),
			},
			{
				Config: config("cng2.acme.example.com"),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionDestroyBeforeCreate),
						plancheck.ExpectUnknownValue(opeResource, tfjsonpath.New("hostname")),
					},
				},
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttr(opeResource, "hostname", "cng2.acme.example.com"),
					tfresource.TestCheckResourceAttr(opeResource, "is_custom_hostname", "true"),
				),
			},
		},
	})

	assert.Equal(t, 0, fake.renames)
}

func TestOutboundPrivateEndpointFailedCreateStaysInState(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)
	fake := newFakeOutboundPrivateEndpoints()
	fake.createStatus = outboundpe.StatusFailed
	fake.createFailureReason = "serviceZoneUnavailable"

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: outboundPrivateEndpointFactories(t, fake),
		CheckDestroy: func(*terraform.State) error {
			if n := fake.count(); n != 0 {
				return fmt.Errorf("%d endpoints left after destroy", n)
			}
			return nil
		},
		Steps: []tfresource.TestStep{
			{
				Config:      outboundPrivateEndpointConfig("cng-prod", opeServiceA, ""),
				ExpectError: regexp.MustCompile(`is failed\s+\(serviceZoneUnavailable\)`),
			},
			{
				// The failed endpoint was kept in state and tainted, so the next apply replaces it.
				Config: outboundPrivateEndpointConfig("cng-prod", opeServiceA, ""),
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(opeResource, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				ExpectError: regexp.MustCompile(`is failed\s+\(serviceZoneUnavailable\)`),
			},
		},
	})

	// The second apply replaced the tainted endpoint, and the test's destroy removed the second one.
	assert.Equal(t, 2, fake.creates)
	assert.Equal(t, 2, fake.deletes)
}

func TestOutboundPrivateEndpointCreateSendsHostnameOnlyWhenSet(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)
	fake := newFakeOutboundPrivateEndpoints()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: outboundPrivateEndpointFactories(t, fake),
		Steps: []tfresource.TestStep{
			{
				Config: outboundPrivateEndpointConfig("cng-prod", opeServiceA, "cng.acme.example.com"),
				Check: func(*terraform.State) error {
					want := outboundpe.CreateOutboundPrivateEndpointRequest{
						Name:         "cng-prod",
						EndpointType: "interface",
						ServiceName:  opeServiceA,
						Hostname:     "cng.acme.example.com",
					}
					if fake.lastCreate != want {
						return fmt.Errorf("create sent %+v, want %+v", fake.lastCreate, want)
					}
					return nil
				},
			},
		},
	})
}

func TestOutboundPrivateEndpointExecuteRetryingConflicts(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)

	tests := []struct {
		name      string
		answers   []int // HTTP status per attempt; 409 carries code
		code      int
		wantCalls int
		wantErr   bool
	}{
		{name: "busy then accepted", answers: []int{http.StatusConflict, http.StatusConflict, http.StatusAccepted}, code: outboundpe.ErrCodeOperationInProgress, wantCalls: 3},
		{name: "endpoint busy then accepted", answers: []int{http.StatusConflict, http.StatusAccepted}, code: outboundpe.ErrCodeBusy, wantCalls: 2},
		{name: "already connected is not retried", answers: []int{http.StatusConflict, http.StatusAccepted}, code: 16008, wantCalls: 1, wantErr: true},
		{name: "busy on every attempt gives up", answers: []int{http.StatusConflict}, code: outboundpe.ErrCodeOperationInProgress, wantCalls: 5, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			handler := func(w http.ResponseWriter, _ *http.Request) {
				status := tt.answers[min(calls, len(tt.answers)-1)]
				calls++
				if status == http.StatusConflict {
					writeAPIError(w, status, tt.code, "conflict")
					return
				}
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"id":"x"}`))
			}
			r := &OutboundPrivateEndpoint{Data: newTestProviderData(t, handler)}

			cfg := api.EndpointCfg{Url: r.HostURL + "/x", Method: http.MethodPost, SuccessStatus: http.StatusAccepted}
			_, err := r.executeRetryingConflicts(t.Context(), cfg, outboundpe.CreateOutboundPrivateEndpointRequest{Name: "n"})

			assert.Equal(t, tt.wantCalls, calls)
			if tt.wantErr {
				var apiErr *api.Error
				require.ErrorAs(t, err, &apiErr)
				assert.Equal(t, tt.code, apiErr.Code)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestOutboundPrivateEndpointWaitUntilDeleted(t *testing.T) {
	fastOutboundPrivateEndpointPolling(t)

	failed := func(reason string) string {
		return fmt.Sprintf(`{"id":"e","status":"failed","failureReason":%q}`, reason)
	}
	tests := []struct {
		name    string
		bodies  []string // "" answers 404
		wantErr string
	}{
		{name: "deleting then gone", bodies: []string{`{"id":"e","status":"deleting"}`, ""}},
		{name: "gone at once", bodies: []string{""}},
		{name: "earlier failure then gone", bodies: []string{failed("createFailed"), `{"id":"e","status":"deleting"}`, ""}},
		{name: "teardown failed", bodies: []string{`{"id":"e","status":"deleting"}`, failed("deleteFailed")}, wantErr: "the teardown failed (deleteFailed), delete it again to retry"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls int
			handler := func(w http.ResponseWriter, _ *http.Request) {
				body := tt.bodies[min(calls, len(tt.bodies)-1)]
				calls++
				if body == "" {
					writeAPIError(w, http.StatusNotFound, 16004, "not found")
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(body))
			}
			r := &OutboundPrivateEndpoint{Data: newTestProviderData(t, handler)}

			err := r.waitUntilDeleted(t.Context(), opeOrgID, opeProjectID, opeClusterID, "e")
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, len(tt.bodies), calls)
		})
	}
}

func TestNewOutboundPrivateEndpointState(t *testing.T) {
	audit := types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())

	t.Run("gcp has no provider config and no hostname yet", func(t *testing.T) {
		got, err := providerschema.NewOutboundPrivateEndpoint(t.Context(), &outboundpe.GetOutboundPrivateEndpointResponse{
			Id:            "e",
			Name:          "n",
			CloudProvider: outboundpe.CloudProvider{Type: "gcp", Region: "us-east1"},
			EndpointType:  "interface",
			ServiceName:   "projects/p/regions/us-east1/serviceAttachments/s",
			Status:        outboundpe.StatusPending,
		}, opeOrgID, opeProjectID, opeClusterID, audit)
		require.NoError(t, err)

		assert.True(t, got.ProviderConfig.IsNull())
		assert.True(t, got.Hostname.IsNull())
		assert.True(t, got.ProviderEndpointId.IsNull())
		assert.True(t, got.FailureReason.IsNull())
		assert.Equal(t, types.ListValueMust(types.StringType, nil).String(), got.IpAddresses.String())
		assert.Equal(t, opeClusterID, got.ClusterId.ValueString())
	})

	t.Run("aws without an address yet keeps empty lists", func(t *testing.T) {
		got, err := providerschema.NewOutboundPrivateEndpoint(t.Context(), &outboundpe.GetOutboundPrivateEndpointResponse{
			Id:             "e",
			CloudProvider:  outboundpe.CloudProvider{Type: "aws", Region: "us-east-1"},
			Status:         outboundpe.StatusPending,
			ProviderConfig: &outboundpe.ProviderConfig{AWS: &outboundpe.AWSConfig{}},
		}, opeOrgID, opeProjectID, opeClusterID, audit)
		require.NoError(t, err)

		var config providerschema.OutboundPrivateEndpointProviderConfig
		require.False(t, got.ProviderConfig.As(t.Context(), &config, basetypes.ObjectAsOptions{}).HasError())
		var aws providerschema.OutboundPrivateEndpointAWSConfig
		require.False(t, config.AWS.As(t.Context(), &aws, basetypes.ObjectAsOptions{}).HasError())

		assert.True(t, aws.RegionalDnsName.IsNull())
		assert.Empty(t, aws.IpAddresses.Elements())
		assert.False(t, aws.IpAddresses.IsNull())
		assert.Empty(t, aws.AvailabilityZones.Elements())
		assert.False(t, aws.AvailabilityZones.IsNull())
	})
}
