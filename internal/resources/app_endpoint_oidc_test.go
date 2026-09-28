package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAppEndpointOidcProviderStateIsKnownAfterFailedRefresh is the AV-145498 regression
// guard. Create and Update refresh the resource with a GET whose failure is only a
// warning, so the interim state they wrote is the state Terraform keeps. Terraform
// rejects any state still holding an unknown after apply ("still indicated an unknown
// value") and taints the resource, so every Computed and Optional+Computed attribute
// must be resolved to a concrete value or null before that interim write.
//
// The plan below marks every computed attribute unknown, which is what the framework
// hands Create when the config sets only the required attributes. Adding a Computed or
// Optional+Computed attribute to this resource means adding it to
// nullifyUnsetOidcProviderFields, or this test fails.
func TestAppEndpointOidcProviderStateIsKnownAfterFailedRefresh(t *testing.T) {
	s := AppEndpointOidcProviderSchema()

	// Guard the guard: if the schema grows a computed attribute the model below does
	// not mark unknown, the test would pass without exercising it.
	computed := map[string]bool{}
	for name, attr := range s.Attributes {
		if attr.IsComputed() {
			computed[name] = true
		}
	}
	for _, name := range []string{
		"register", "discovery_url", "username_claim", "roles_claim", "user_prefix",
		"provider_id", "is_default",
	} {
		if !computed[name] {
			t.Errorf("%s is no longer computed; update this test", name)
		}
		delete(computed, name)
	}
	for name := range computed {
		t.Errorf("%s is computed but not covered by this test; add it here and to nullifyUnsetOidcProviderFields", name)
	}

	minimalConfigPlan := func() providerschema.AppEndpointOidcProvider {
		return providerschema.AppEndpointOidcProvider{
			OrganizationId:  types.StringValue("org"),
			ProjectId:       types.StringValue("proj"),
			ClusterId:       types.StringValue("clus"),
			AppServiceId:    types.StringValue("appsvc"),
			AppEndpointName: types.StringValue("ep"),
			Issuer:          types.StringValue("https://accounts.google.com"),
			ClientId:        types.StringValue("example-client-id"),
			Register:        types.BoolUnknown(),
			DiscoveryUrl:    types.StringUnknown(),
			UsernameClaim:   types.StringUnknown(),
			RolesClaim:      types.StringUnknown(),
			UserPrefix:      types.StringUnknown(),
			ProviderId:      types.StringUnknown(),
			IsDefault:       types.BoolUnknown(),
		}
	}

	t.Run("create", func(t *testing.T) {
		plan := minimalConfigPlan()

		// Mirrors Create: provider_id comes from the POST response, is_default is not
		// in it, and the remaining unknowns are nullified.
		plan.ProviderId = types.StringValue("prov-123")
		plan.IsDefault = types.BoolNull()
		nullifyUnsetOidcProviderFields(&plan)

		assertStateFullyKnown(t, s, plan)

		if plan.ProviderId.ValueString() != "prov-123" {
			t.Errorf("provider_id must survive a failed refresh, or Read, Update and Delete "+
				"cannot address the remote provider; got %v", plan.ProviderId)
		}
	})

	t.Run("update", func(t *testing.T) {
		plan := minimalConfigPlan()

		// A create whose refresh failed leaves is_default null in state, which the
		// framework re-marks unknown in the next update plan.
		state := minimalConfigPlan()
		state.ProviderId = types.StringValue("prov-123")
		state.IsDefault = types.BoolNull()

		// Mirrors Update.
		plan.ProviderId = state.ProviderId
		plan.IsDefault = state.IsDefault
		nullifyUnsetOidcProviderFields(&plan)

		assertStateFullyKnown(t, s, plan)
	})
}

// assertStateFullyKnown writes the model into a state of the given schema and reports
// every attribute still unknown, which is what Terraform rejects after apply.
func assertStateFullyKnown(t *testing.T, s schema.Schema, model any) {
	t.Helper()

	state := tfsdk.State{Schema: s}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("State.Set: %v", diags)
	}
	if state.Raw.IsFullyKnown() {
		return
	}

	var obj map[string]tftypes.Value
	if err := state.Raw.As(&obj); err != nil {
		t.Fatalf("state is not fully known and could not be decoded: %v", err)
	}
	for name, v := range obj {
		if !v.IsFullyKnown() {
			t.Errorf("%s is still unknown after apply; nullifyUnsetOidcProviderFields must cover it", name)
		}
	}
}
