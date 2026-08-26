package resources

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"gotest.tools/assert"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testCredentialID   = "e8d94b7a-d03e-4f28-b14a-65e0487c764a" //nolint:gosec // G101: a resource UUID, not a credential
	testCredentialName = "AdvancedCredential"
	testUserRole       = "developer"
)

// databaseCredentialGetHandler serves the V4 GET database credential endpoint, carrying
// userRoles only when withUserRoles is set - omitting them is the response shape that used
// to flip credential_type to basic (AV-139985).
//
// The body is a literal rather than a marshalled api.GetDatabaseCredentialResponse so the
// test pins the real wire format: marshalling the provider's own struct would round-trip
// successfully even if its json tags did not match what the API sends.
func databaseCredentialGetHandler(t *testing.T, withUserRoles bool, gotPath *string) http.HandlerFunc {
	t.Helper()

	var userRoles string
	if withUserRoles {
		userRoles = fmt.Sprintf(`"userRoles": [%q],`, testUserRole)
	}
	body := fmt.Sprintf(`{
		"id": %q,
		"name": %q,
		"organizationId": %q,
		"projectId": %q,
		"clusterId": %q,
		"access": [],
		%s
		"audit": {
			"createdAt": "2026-08-01T00:00:00Z", "createdBy": "tf-acc",
			"modifiedAt": "2026-08-01T00:00:00Z", "modifiedBy": "tf-acc", "version": 1
		}
	}`, testCredentialID, testCredentialName, testOrgID, testProjectID, testClusterID, userRoles)

	return func(w http.ResponseWriter, r *http.Request) {
		*gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("writing stub response: %v", err)
		}
	}
}

// managedCredentialState is the prior state Create or Update would have written for a
// Terraform-managed credential: every ID resolved and credential_type concrete.
func managedCredentialState(credentialType string, userRoles []types.String, access []providerschema.Access) providerschema.DatabaseCredential {
	return providerschema.DatabaseCredential{
		Id:             types.StringValue(testCredentialID),
		Name:           types.StringValue(testCredentialName),
		Password:       types.StringValue("Secret12$#"),
		OrganizationId: types.StringValue(testOrgID),
		ProjectId:      types.StringValue(testProjectID),
		ClusterId:      types.StringValue(testClusterID),
		Audit:          types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes()),
		CredentialType: types.StringValue(credentialType),
		UserRoles:      userRoles,
		Access:         access,
	}
}

// importedCredentialState is the state ImportState leaves behind: the comma separated
// import string in id and every other attribute null, credential_type included. This is
// the only case where the credential type has to be derived from the GET response.
func importedCredentialState() providerschema.DatabaseCredential {
	importID := fmt.Sprintf("id=%s,organization_id=%s,project_id=%s,cluster_id=%s",
		testCredentialID, testOrgID, testProjectID, testClusterID)

	return providerschema.DatabaseCredential{
		Id:             types.StringValue(importID),
		Name:           types.StringNull(),
		Password:       types.StringNull(),
		OrganizationId: types.StringNull(),
		ProjectId:      types.StringNull(),
		ClusterId:      types.StringNull(),
		Audit:          types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes()),
		CredentialType: types.StringNull(),
	}
}

// TestDatabaseCredentialReadPreservesCredentialType covers the refresh path for
// credential_type and user_roles. The V4 GET response reports neither reliably -
// credential_type is absent from it altogether - so Read carries both forward from prior
// state and derives them only on import. Before AV-139985 Read re-derived credential_type
// from the presence of userRoles, so a response omitting them rewrote an advanced
// credential as basic and the RequiresReplace on credential_type destroyed it.
func TestDatabaseCredentialReadPreservesCredentialType(t *testing.T) {
	ctx := context.Background()

	advancedRoles := []types.String{types.StringValue(testUserRole)}
	basicAccess := []providerschema.Access{{Privileges: []types.String{types.StringValue("read")}}}

	tests := []struct {
		name          string
		prior         providerschema.DatabaseCredential
		withUserRoles bool
		wantType      string
		wantRoles     []string
	}{
		{
			name:          "advanced credential, response carries roles",
			prior:         managedCredentialState(credentialTypeAdvanced, advancedRoles, nil),
			withUserRoles: true,
			wantType:      credentialTypeAdvanced,
			wantRoles:     []string{testUserRole},
		},
		{
			// The regression. Identical to the case above but for the missing roles.
			name:          "advanced credential, response omits roles",
			prior:         managedCredentialState(credentialTypeAdvanced, advancedRoles, nil),
			withUserRoles: false,
			wantType:      credentialTypeAdvanced,
			wantRoles:     []string{testUserRole},
		},
		{
			name:          "basic credential stays basic with no roles",
			prior:         managedCredentialState(credentialTypeBasic, nil, basicAccess),
			withUserRoles: false,
			wantType:      credentialTypeBasic,
			wantRoles:     nil,
		},
		{
			// A basic credential must not be promoted if the response ever reports roles.
			name:          "basic credential is not promoted by roles in the response",
			prior:         managedCredentialState(credentialTypeBasic, nil, basicAccess),
			withUserRoles: true,
			wantType:      credentialTypeBasic,
			wantRoles:     nil,
		},
		{
			name:          "import derives advanced from roles in the response",
			prior:         importedCredentialState(),
			withUserRoles: true,
			wantType:      credentialTypeAdvanced,
			wantRoles:     []string{testUserRole},
		},
		{
			name:          "import derives basic when the response has no roles",
			prior:         importedCredentialState(),
			withUserRoles: false,
			wantType:      credentialTypeBasic,
			wantRoles:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotPath string
			r := &DatabaseCredential{
				Data: newTestProviderData(t, databaseCredentialGetHandler(t, tc.withUserRoles, &gotPath)),
			}

			priorState := tfsdk.State{Schema: DatabaseCredentialSchema()}
			assertNoDiags(t, priorState.Set(ctx, tc.prior))

			// The framework seeds the response state with prior state before calling Read.
			resp := resource.ReadResponse{State: priorState}
			r.Read(ctx, resource.ReadRequest{State: priorState}, &resp)
			assertNoDiags(t, resp.Diagnostics)

			var got providerschema.DatabaseCredential
			assertNoDiags(t, resp.State.Get(ctx, &got))

			assert.Equal(t, got.CredentialType.ValueString(), tc.wantType,
				"credential_type must never change on a refresh")
			assert.DeepEqual(t, roleNames(got.UserRoles), tc.wantRoles)
			assert.Equal(t, got.Id.ValueString(), testCredentialID,
				"the credential id must be resolved to the id from the response")
			assert.Equal(t, gotPath, fmt.Sprintf("/v4/organizations/%s/projects/%s/clusters/%s/users/%s",
				testOrgID, testProjectID, testClusterID, testCredentialID))
		})
	}
}

// roleNames unwraps user role names for comparison, preserving a nil slice so that a
// null user_roles attribute is distinguishable from an empty one.
func roleNames(roles []types.String) []string {
	if roles == nil {
		return nil
	}
	names := make([]string, len(roles))
	for i, role := range roles {
		names[i] = role.ValueString()
	}
	return names
}

func assertNoDiags(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags.Errors())
	}
}

// TestDatabaseCredentialTypeReplacement covers when a credential_type diff is allowed to
// destroy and recreate the credential. The V4 API cannot convert a credential between
// basic and advanced, so a practitioner-driven type change must still replace. An
// attribute merely absent from an older state file must not: credential_type was added by
// AV-124932, so every state a released provider wrote decodes it as null and the Default
// then plans "basic" against that null. Under a plain RequiresReplace this destroyed and
// recreated every existing credential on the first plan after the upgrade, regenerating
// the auto-generated passwords applications were holding (AV-139981).
//
// A null prior state is not on its own proof of a backfill, so the legacy cases below pin
// both readings of one: the practitioner may convert the credential in the same plan that
// upgrades the provider, and may equally pin the default explicitly. Only the planned value
// separates them.
func TestDatabaseCredentialTypeReplacement(t *testing.T) {
	ctx := context.Background()
	credentialSchema := DatabaseCredentialSchema()

	attribute, ok := credentialSchema.Attributes["credential_type"].(*schema.StringAttribute)
	if !ok {
		t.Fatalf("credential_type is %T, want *schema.StringAttribute", credentialSchema.Attributes["credential_type"])
	}

	basicAccess := []providerschema.Access{{Privileges: []types.String{types.StringValue("read")}}}
	advancedRoles := []types.String{types.StringValue(testUserRole)}

	basic := managedCredentialState(credentialTypeBasic, nil, basicAccess)
	advanced := managedCredentialState(credentialTypeAdvanced, advancedRoles, nil)

	// The regression fixture: a basic credential as a released provider stored it, with no
	// credential_type key at all. The planned value is the Default backfilling it.
	upgraded := managedCredentialState(credentialTypeBasic, nil, basicAccess)
	upgraded.CredentialType = types.StringNull()

	tests := []struct {
		name string
		// A nil state is a create and a nil plan is a destroy, both of which the framework
		// signals with a null raw value rather than a null attribute.
		state       *providerschema.DatabaseCredential
		plan        *providerschema.DatabaseCredential
		configValue types.String
		wantReplace bool
	}{
		{
			// The regression. Nothing in the configuration changed.
			name:        "upgrade backfills a state written before the attribute existed",
			state:       &upgraded,
			plan:        &basic,
			configValue: types.StringNull(),
			wantReplace: false,
		},
		{
			// Converting in the same plan that upgrades the provider. The prior state is
			// null exactly as in the backfill above, but the V4 API cannot promote a basic
			// credential, so this must still replace.
			name:        "legacy state converted to advanced replaces",
			state:       &upgraded,
			plan:        &advanced,
			configValue: types.StringValue(credentialTypeAdvanced),
			wantReplace: true,
		},
		{
			// The mirror image: pinning the default explicitly while upgrading is not a
			// change, so keying off "was it configured" instead of the planned value would
			// destroy the credential here.
			name:        "legacy state pinned to the default does not replace",
			state:       &upgraded,
			plan:        &basic,
			configValue: types.StringValue(credentialTypeBasic),
			wantReplace: false,
		},
		{
			name:        "an unchanged credential type is left alone",
			state:       &basic,
			plan:        &basic,
			configValue: types.StringNull(),
			wantReplace: false,
		},
		{
			name:        "basic to advanced replaces",
			state:       &basic,
			plan:        &advanced,
			configValue: types.StringValue(credentialTypeAdvanced),
			wantReplace: true,
		},
		{
			// Dropping credential_type from an advanced credential's configuration falls
			// back to the Default, which is still a real type change.
			name:        "advanced to basic replaces",
			state:       &advanced,
			plan:        &basic,
			configValue: types.StringNull(),
			wantReplace: true,
		},
		{
			name:        "create does not replace",
			state:       nil,
			plan:        &basic,
			configValue: types.StringNull(),
			wantReplace: false,
		},
		{
			name:        "destroy does not replace",
			state:       &basic,
			plan:        nil,
			configValue: types.StringNull(),
			wantReplace: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := planmodifier.StringRequest{
				Path:        path.Root("credential_type"),
				State:       tfsdk.State{Schema: credentialSchema, Raw: credentialRawValue(t, ctx, credentialSchema, tc.state)},
				StateValue:  credentialTypeOf(tc.state),
				Plan:        tfsdk.Plan{Schema: credentialSchema, Raw: credentialRawValue(t, ctx, credentialSchema, tc.plan)},
				PlanValue:   credentialTypeOf(tc.plan),
				ConfigValue: tc.configValue,
			}

			// Each modifier sees the plan value the previous one produced, matching
			// fwserver.AttributeModifyPlan.
			var gotReplace bool
			for _, modifier := range attribute.PlanModifiers {
				resp := planmodifier.StringResponse{PlanValue: req.PlanValue}
				modifier.PlanModifyString(ctx, req, &resp)
				req.PlanValue = resp.PlanValue
				gotReplace = gotReplace || resp.RequiresReplace
			}

			assert.Equal(t, gotReplace, tc.wantReplace,
				"credential_type replacement decision for %q", tc.name)
		})
	}
}

// credentialRawValue renders v as the whole-resource value the framework puts in a plan
// modifier request. RequiresReplaceIf consults req.State.Raw and req.Plan.Raw before it
// looks at any attribute value, so a request carrying only StateValue and PlanValue would
// take the create branch and pass no matter what the modifier does. A nil v is that null
// raw value: no prior state on create, no planned state on destroy.
func credentialRawValue(t *testing.T, ctx context.Context, s schema.Schema, v *providerschema.DatabaseCredential) tftypes.Value {
	t.Helper()

	if v == nil {
		return tftypes.NewValue(s.Type().TerraformType(ctx), nil)
	}

	state := tfsdk.State{Schema: s}
	assertNoDiags(t, state.Set(ctx, *v))
	return state.Raw
}

func credentialTypeOf(v *providerschema.DatabaseCredential) types.String {
	if v == nil {
		return types.StringNull()
	}
	return v.CredentialType
}
