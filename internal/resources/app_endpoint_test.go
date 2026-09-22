package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func oidcProvider(clientID string) providerschema.AppEndpointOidc {
	return providerschema.AppEndpointOidc{
		Issuer:        types.StringValue("https://accounts.google.com"),
		ClientId:      types.StringValue(clientID),
		Register:      types.BoolValue(true),
		UserPrefix:    types.StringNull(),
		DiscoveryUrl:  types.StringNull(),
		UsernameClaim: types.StringNull(),
		RolesClaim:    types.StringNull(),
		ProviderId:    types.StringNull(),
		IsDefault:     types.BoolNull(),
	}
}

func Test_preserveDisabledOidc(t *testing.T) {
	tests := []struct {
		name string
		// configOidc is the oidc value from the plan, or from the prior state during Read.
		configOidc []providerschema.AppEndpointOidc
		// stateOidc is the oidc value produced by refreshAppEndpoint.
		stateOidc []providerschema.AppEndpointOidc
		// aliasPlanAsState mirrors the Create path where the post-create refresh fails
		// and the plan, already rewritten by setAppEndpointComputedAttributesToNull, is
		// handed back as the state. config and state are then the same pointer.
		aliasPlanAsState bool
		expected         []providerschema.AppEndpointOidc
		// expectNil distinguishes a nil slice (null list) from an empty one (empty list).
		expectNil bool
	}{
		{
			name:       "null config with no remote providers stays null",
			configOidc: nil,
			stateOidc:  nil,
			expected:   nil,
			expectNil:  true,
		},
		{
			name:       "null config is not widened by a non-nil empty refresh",
			configOidc: nil,
			stateOidc:  []providerschema.AppEndpointOidc{},
			expected:   nil,
			expectNil:  true,
		},
		{
			name:       "empty list config with no remote providers stays an empty list",
			configOidc: []providerschema.AppEndpointOidc{},
			stateOidc:  nil,
			expected:   []providerschema.AppEndpointOidc{},
			expectNil:  false,
		},
		{
			name:       "remote providers are preserved for a null config",
			configOidc: nil,
			stateOidc:  []providerschema.AppEndpointOidc{oidcProvider("remote-client-id")},
			expected:   []providerschema.AppEndpointOidc{oidcProvider("remote-client-id")},
			expectNil:  false,
		},
		{
			name:       "remote providers are preserved for a configured provider",
			configOidc: []providerschema.AppEndpointOidc{oidcProvider("configured-client-id")},
			stateOidc:  []providerschema.AppEndpointOidc{oidcProvider("remote-client-id")},
			expected:   []providerschema.AppEndpointOidc{oidcProvider("remote-client-id")},
			expectNil:  false,
		},
		{
			name:       "a configured provider that the API dropped is left as drift, not an empty list",
			configOidc: []providerschema.AppEndpointOidc{oidcProvider("configured-client-id")},
			stateOidc:  nil,
			expected:   nil,
			expectNil:  true,
		},
		{
			name:             "null config stays null when the failed create hands the plan back as state",
			configOidc:       nil,
			aliasPlanAsState: true,
			expected:         nil,
			expectNil:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &providerschema.AppEndpoint{Oidc: tt.configOidc}
			state := &providerschema.AppEndpoint{Oidc: tt.stateOidc}
			if tt.aliasPlanAsState {
				requireNoDiags(t, setAppEndpointComputedAttributesToNull(context.Background(), config))
				state = config
			}

			preserveDisabledOidc(config, state)

			assert.Equal(t, tt.expected, state.Oidc)
			assert.Equal(t, tt.expectNil, state.Oidc == nil, "nil-ness decides null vs empty list in state")
		})
	}
}

// Test_appEndpointStateOidc is the regression test for CBSE-23701: a null oidc must
// reach Terraform as a null list, never cty.ListValEmpty.
func Test_appEndpointStateOidc(t *testing.T) {
	configuredProvider := oidcProvider("example-client-id")
	configuredProvider.ProviderId = types.StringValue("stale-provider-id")
	configuredProvider.IsDefault = types.BoolValue(true)

	tests := []struct {
		name string
		oidc []providerschema.AppEndpointOidc
		// expectNull true: Terraform value should be null list, not empty list.
		expectNull bool
		// expectComputedNulled is plan.Oidc after setAppEndpointComputedAttributesToNull.
		expectComputedNulled []providerschema.AppEndpointOidc
	}{
		{
			name:                 "nil slice is a null list",
			oidc:                 nil,
			expectNull:           true,
			expectComputedNulled: nil,
		},
		{
			name:                 "empty slice is an empty list",
			oidc:                 []providerschema.AppEndpointOidc{},
			expectNull:           false,
			expectComputedNulled: []providerschema.AppEndpointOidc{},
		},
		{
			name:       "configured provider keeps required fields and nulls computed ones",
			oidc:       []providerschema.AppEndpointOidc{configuredProvider},
			expectNull: false,
			// oidcProvider is the same provider with provider_id and is_default unset.
			expectComputedNulled: []providerschema.AppEndpointOidc{oidcProvider("example-client-id")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			endpoint := appEndpointWithoutOptionals()
			endpoint.Oidc = tt.oidc
			assert.Equal(t, tt.expectNull, setStateAndReadAttribute(ctx, t, endpoint, "oidc").IsNull(),
				"framework encoding of the slice")

			plan := appEndpointWithoutOptionals()
			plan.Oidc = tt.oidc
			requireNoDiags(t, setAppEndpointComputedAttributesToNull(ctx, plan))

			assert.Equal(t, tt.expectComputedNulled, plan.Oidc)
			assert.Equal(t, tt.expectNull, plan.Oidc == nil,
				"nil-ness after setAppEndpointComputedAttributesToNull")
			assert.Equal(t, tt.expectNull, setStateAndReadAttribute(ctx, t, plan, "oidc").IsNull(),
				"encoding after setAppEndpointComputedAttributesToNull")
		})
	}
}

// appEndpointWithoutOptionals returns an App Endpoint with every attribute that the
// schema requires populated, so that only the attribute under test varies.
func appEndpointWithoutOptionals() *providerschema.AppEndpoint {
	return &providerschema.AppEndpoint{
		OrganizationId:   types.StringValue("8a5d3a1c-7b0e-4c9d-9f2a-1e6b4d8c3f70"),
		ProjectId:        types.StringValue("3f9c2d5e-6a1b-4e8f-8c7d-2b5a9e4f1c60"),
		ClusterId:        types.StringValue("9d75ba14-52a6-42dc-988e-0d911693b591"),
		AppServiceId:     types.StringValue("65c09dea-919b-4748-82f3-c6c921f65928"),
		Bucket:           types.StringValue("travel-sample"),
		Name:             types.StringValue("pickandpack"),
		UserXattrKey:     types.StringNull(),
		DeltaSyncEnabled: types.BoolNull(),
		Scopes: types.MapNull(types.ObjectType{
			AttrTypes: providerschema.AppEndpointScope{}.AttributeTypes(),
		}),
		RequireResync: types.MapNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{"items": types.SetType{ElemType: types.StringType}},
		}),
		State:      types.StringNull(),
		AdminURL:   types.StringNull(),
		MetricsURL: types.StringNull(),
		PublicURL:  types.StringNull(),
	}
}

// setStateAndReadAttribute writes endpoint into a state backed by the real resource
// schema and returns the raw value of name, which is where null and empty lists
// become distinguishable.
func setStateAndReadAttribute(
	ctx context.Context, t *testing.T, endpoint *providerschema.AppEndpoint, name string,
) tftypes.Value {
	t.Helper()

	appEndpointSchema := AppEndpointSchema()
	state := tfsdk.State{
		Schema: appEndpointSchema,
		Raw:    tftypes.NewValue(appEndpointSchema.Type().TerraformType(ctx), nil),
	}
	requireNoDiags(t, state.Set(ctx, endpoint))

	raw, err := state.Raw.ApplyTerraform5AttributePathStep(tftypes.AttributeName(name))
	require.NoError(t, err)

	value, ok := raw.(tftypes.Value)
	require.True(t, ok, "attribute %q is not a tftypes.Value", name)

	return value
}

// requireNoDiags fails the test when diags carries an error, reporting the diagnostics.
func requireNoDiags(t *testing.T, diags diag.Diagnostics) {
	t.Helper()

	require.Falsef(t, diags.HasError(), "unexpected diagnostics: %v", diags.Errors())
}
