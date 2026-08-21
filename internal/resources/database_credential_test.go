package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testCredentialID   = "e8d94b7a-d03e-4f28-b14a-65e0487c764a"
	testCredentialName = "AdvancedCredential"
	testUserRole       = "developer"
)

// databaseCredentialGetHandler serves the V4 GET database credential endpoint with a
// fixed credential, carrying userRoles only when withUserRoles is set. Omitting them
// is the response shape that used to flip credential_type to basic (AV-139985); the
// omitempty tag on the field drops it from the body entirely.
func databaseCredentialGetHandler(t *testing.T, withUserRoles bool, gotPath *string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		*gotPath = r.URL.Path

		body := api.GetDatabaseCredentialResponse{
			Id:             uuid.MustParse(testCredentialID),
			Name:           testCredentialName,
			OrganizationId: testOrgID,
			ProjectId:      testProjectID,
			ClusterId:      testClusterID,
		}
		if withUserRoles {
			body.UserRoles = []string{testUserRole}
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Errorf("encoding stub response: %v", err)
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
