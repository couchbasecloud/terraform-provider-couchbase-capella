package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"
	"gotest.tools/assert/cmp"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	oidcTestOrgID       = "aaaaaaaa-0000-0000-0000-000000000001"
	oidcTestProjectID   = "aaaaaaaa-0000-0000-0000-000000000002"
	oidcTestClusterID   = "aaaaaaaa-0000-0000-0000-000000000003"
	oidcTestAppSvcID    = "aaaaaaaa-0000-0000-0000-000000000004"
	oidcTestEndpoint    = "test-endpoint"
	oidcTestProviderID  = "prov-1"
	oidcTestIssuer      = "https://accounts.google.com"
	oidcTestClientID    = "example-client-id"
	oidcTestImportIDFmt = "organization_id=%s,project_id=%s,cluster_id=%s,app_service_id=%s,app_endpoint_name=%s,provider_id=%s"
)

// TestAppEndpointOidcProviderReadRestoresIdentityAfterImport is the AV-145496 guard.
//
// ImportStatePassthroughID stores the whole composite ID in app_endpoint_name and leaves
// the other identity attributes null, so Read is the only thing that can restore them.
// mapResponseToState does not, and Validate returns a parsed map without mutating the
// model, so the imported resource keeps null IDs and a composite endpoint name. All five
// are requiresReplace, so the first plan after a successful import destroys and recreates
// the provider. app_endpoint_cors.go refreshes into a fresh state built from the parsed
// IDs and does not have this problem.
func TestAppEndpointOidcProviderReadRestoresIdentityAfterImport(t *testing.T) {
	ctx := context.Background()

	r := &AppEndpointOidcProvider{Data: newTestProviderData(t, oidcProviderGetHandler(t))}

	// The state the framework hands Read straight after terraform import.
	imported := providerschema.AppEndpointOidcProvider{
		AppEndpointName: types.StringValue(fmt.Sprintf(
			oidcTestImportIDFmt,
			oidcTestOrgID, oidcTestProjectID, oidcTestClusterID,
			oidcTestAppSvcID, oidcTestEndpoint, oidcTestProviderID,
		)),
	}

	priorState := tfsdk.State{Schema: AppEndpointOidcProviderSchema()}
	assertNoDiags(t, priorState.Set(ctx, imported))

	resp := resource.ReadResponse{State: priorState}
	r.Read(ctx, resource.ReadRequest{State: priorState}, &resp)
	assertNoDiags(t, resp.Diagnostics)

	var got providerschema.AppEndpointOidcProvider
	assertNoDiags(t, resp.State.Get(ctx, &got))

	// assert.Check rather than assert.Equal so a failure reports every attribute that
	// was lost, not just the first.
	assert.Check(t, cmp.Equal(got.OrganizationId.ValueString(), oidcTestOrgID),
		"organization_id must be restored from the import ID; it is requiresReplace")
	assert.Check(t, cmp.Equal(got.ProjectId.ValueString(), oidcTestProjectID),
		"project_id must be restored from the import ID; it is requiresReplace")
	assert.Check(t, cmp.Equal(got.ClusterId.ValueString(), oidcTestClusterID),
		"cluster_id must be restored from the import ID; it is requiresReplace")
	assert.Check(t, cmp.Equal(got.AppServiceId.ValueString(), oidcTestAppSvcID),
		"app_service_id must be restored from the import ID; it is requiresReplace")
	assert.Check(t, cmp.Equal(got.AppEndpointName.ValueString(), oidcTestEndpoint),
		"app_endpoint_name must be normalised to the endpoint name, not left as the composite import ID")
}

// TestAppEndpointOidcProviderCreateRecordsProviderOnBadResponseBody is the AV-145497 guard.
//
// The POST has already created the provider by the time the body is decoded, so returning
// on an unmarshal failure without writing state loses a resource that exists in Capella:
// the next apply creates a second one. Either state records enough to address it, or the
// diagnostic has to tell the operator what was created and must be removed by hand.
func TestAppEndpointOidcProviderCreateRecordsProviderOnBadResponseBody(t *testing.T) {
	ctx := context.Background()

	r := &AppEndpointOidcProvider{Data: newTestProviderData(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// 201: the provider now exists remotely. The body is not valid JSON.
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("not json"))
	})}

	plan := tfsdk.Plan{Schema: AppEndpointOidcProviderSchema()}
	assertNoDiags(t, plan.Set(ctx, providerschema.AppEndpointOidcProvider{
		OrganizationId:  types.StringValue(oidcTestOrgID),
		ProjectId:       types.StringValue(oidcTestProjectID),
		ClusterId:       types.StringValue(oidcTestClusterID),
		AppServiceId:    types.StringValue(oidcTestAppSvcID),
		AppEndpointName: types.StringValue(oidcTestEndpoint),
		Issuer:          types.StringValue(oidcTestIssuer),
		ClientId:        types.StringValue(oidcTestClientID),
		Register:        types.BoolNull(),
		DiscoveryUrl:    types.StringNull(),
		UsernameClaim:   types.StringNull(),
		RolesClaim:      types.StringNull(),
		UserPrefix:      types.StringNull(),
		ProviderId:      types.StringNull(),
		IsDefault:       types.BoolNull(),
	}))

	resp := resource.CreateResponse{State: tfsdk.State{Schema: AppEndpointOidcProviderSchema()}}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected an error diagnostic for an undecodable create response")
	}

	// Terraform drops a null state, so a provider that was really created is orphaned
	// unless the operator is told about it.
	if resp.State.Raw.IsNull() {
		summaries := ""
		for _, d := range resp.Diagnostics.Errors() {
			summaries += fmt.Sprintf("\n  %s: %s", d.Summary(), d.Detail())
		}
		t.Errorf("Create wrote no state after a 201, so the OIDC provider is orphaned in "+
			"Capella and the next apply creates a second one; the diagnostics do not name "+
			"the created provider either:%s", summaries)
	}
}

// oidcProviderGetHandler answers the single-provider GET that Read issues.
func oidcProviderGetHandler(t *testing.T) http.HandlerFunc {
	t.Helper()
	want := fmt.Sprintf(
		"/v4/organizations/%s/projects/%s/clusters/%s/appservices/%s/appEndpoints/%s/oidcProviders/%s",
		oidcTestOrgID, oidcTestProjectID, oidcTestClusterID, oidcTestAppSvcID,
		oidcTestEndpoint, oidcTestProviderID,
	)
	return func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != want {
			t.Errorf("unexpected request path\n got: %s\nwant: %s", req.URL.Path, want)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(api.AppEndpointOIDCProviderResponse{
			ProviderID: oidcTestProviderID,
			Issuer:     oidcTestIssuer,
			ClientID:   oidcTestClientID,
		})
	}
}
