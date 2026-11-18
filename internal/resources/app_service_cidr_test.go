package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const allowedCIDRTestId = "cidr-id"

// allowedCIDRBackend creates allowed CIDRs with the ID allowedCIDRTestId and lists entries with one entry on each
// page, so that an entry after the first one is only found through pagination.
type allowedCIDRBackend struct {
	mu          sync.Mutex
	entries     []api.AppServiceAllowedCIDRResponse
	createBody  map[string]any
	accessModes []string
}

func (b *allowedCIDRBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if r.Method == http.MethodPost {
		_ = json.NewDecoder(r.Body).Decode(&b.createBody)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(api.AppServiceAllowedCIDRResponse{Id: allowedCIDRTestId})
		return
	}

	b.accessModes = r.URL.Query()["accessMode"]
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	var (
		data []api.AppServiceAllowedCIDRResponse
		next int
	)
	if page <= len(b.entries) {
		data = b.entries[page-1 : page]
	}
	if page < len(b.entries) {
		next = page + 1
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"data":   data,
		"cursor": api.Cursor{Pages: api.Pages{Page: page, Next: next}},
	})
}

func allowedCIDREntry(id, accessMode string) api.AppServiceAllowedCIDRResponse {
	return api.AppServiceAllowedCIDRResponse{Id: id, Cidr: "203.0.113.0/24", AccessMode: accessMode}
}

func allowedCIDRModel(id, accessMode string) providerschema.AppServiceCIDR {
	model := providerschema.AppServiceCIDR{
		Id:             types.StringValue(id),
		OrganizationId: types.StringValue("org"),
		ProjectId:      types.StringValue("project"),
		ClusterId:      types.StringValue("cluster"),
		AppServiceId:   types.StringValue("appservice"),
		Cidr:           types.StringValue("203.0.113.0/24"),
		AccessMode:     types.StringValue(accessMode),
		Audit:          types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes()),
	}
	if accessMode == "" {
		model.AccessMode = types.StringNull()
	}
	return model
}

// importedAllowedCIDRModel returns the state that an import gives to Read: the import ID only.
func importedAllowedCIDRModel() providerschema.AppServiceCIDR {
	return providerschema.AppServiceCIDR{
		Id:         types.StringValue("id=" + allowedCIDRTestId + ",app_service_id=appservice,cluster_id=cluster,project_id=project,organization_id=org"),
		AccessMode: types.StringNull(),
		Audit:      types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes()),
	}
}

func getAllowedCIDRModel(t *testing.T, state tfsdk.State) providerschema.AppServiceCIDR {
	t.Helper()
	var got providerschema.AppServiceCIDR
	diags := state.Get(context.Background(), &got)
	require.False(t, diags.HasError(), diags)
	return got
}

func TestAppServiceCidr_Create(t *testing.T) {
	tests := []struct {
		name              string
		accessMode        string
		entryAccessMode   string
		expectCreateBody  map[string]any
		expectAccessModes []string
	}{
		{
			name:             "omitted access mode sends no accessMode and stays null in the state",
			entryAccessMode:  api.AllowedCIDRAccessModeAdmin,
			expectCreateBody: map[string]any{"cidr": "203.0.113.0/24"},
		},
		{
			name:              "admin sends accessMode to the create and the list",
			accessMode:        api.AllowedCIDRAccessModeAdmin,
			entryAccessMode:   api.AllowedCIDRAccessModeAdmin,
			expectCreateBody:  map[string]any{"cidr": "203.0.113.0/24", "accessMode": "admin"},
			expectAccessModes: []string{api.AllowedCIDRAccessModeAdmin},
		},
		{
			name:              "public sends accessMode to the create and the list",
			accessMode:        api.AllowedCIDRAccessModePublic,
			entryAccessMode:   api.AllowedCIDRAccessModePublic,
			expectCreateBody:  map[string]any{"cidr": "203.0.113.0/24", "accessMode": "public"},
			expectAccessModes: []string{api.AllowedCIDRAccessModePublic},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			backend := &allowedCIDRBackend{
				entries: []api.AppServiceAllowedCIDRResponse{allowedCIDREntry(allowedCIDRTestId, test.entryAccessMode)},
			}
			a := &AppServiceCidr{Data: newTestProviderData(t, backend.handler)}

			plan := tfsdk.Plan{Schema: AllowedCIDRsSchema()}
			diags := plan.Set(ctx, allowedCIDRModel("", test.accessMode))
			require.False(t, diags.HasError(), diags)
			resp := resource.CreateResponse{State: tfsdk.State{Schema: AllowedCIDRsSchema()}}

			a.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			assert.Equal(t, test.expectCreateBody, backend.createBody)
			assert.Equal(t, test.expectAccessModes, backend.accessModes)
			got := getAllowedCIDRModel(t, resp.State)
			assert.Equal(t, allowedCIDRTestId, got.Id.ValueString())
			assert.Equal(t, allowedCIDRModel("", test.accessMode).AccessMode, got.AccessMode)
		})
	}
}

func TestAppServiceCidr_Read(t *testing.T) {
	tests := []struct {
		name              string
		state             providerschema.AppServiceCIDR
		entries           []api.AppServiceAllowedCIDRResponse
		expectAccessModes []string
		expectAccessMode  types.String
		expectRemoved     bool
	}{
		{
			name:              "access mode in the state lists that access mode",
			state:             allowedCIDRModel(allowedCIDRTestId, api.AllowedCIDRAccessModePublic),
			entries:           []api.AppServiceAllowedCIDRResponse{allowedCIDREntry(allowedCIDRTestId, "public")},
			expectAccessModes: []string{api.AllowedCIDRAccessModePublic},
			expectAccessMode:  types.StringValue(api.AllowedCIDRAccessModePublic),
		},
		{
			name:              "import lists both access modes and sets public",
			state:             importedAllowedCIDRModel(),
			entries:           []api.AppServiceAllowedCIDRResponse{allowedCIDREntry(allowedCIDRTestId, "public")},
			expectAccessModes: []string{api.AllowedCIDRAccessModeAdmin, api.AllowedCIDRAccessModePublic},
			expectAccessMode:  types.StringValue(api.AllowedCIDRAccessModePublic),
		},
		{
			name:              "import lists both access modes and keeps admin null",
			state:             importedAllowedCIDRModel(),
			entries:           []api.AppServiceAllowedCIDRResponse{allowedCIDREntry(allowedCIDRTestId, "admin")},
			expectAccessModes: []string{api.AllowedCIDRAccessModeAdmin, api.AllowedCIDRAccessModePublic},
			expectAccessMode:  types.StringNull(),
		},
		{
			name:             "state without access mode sends no accessMode and keeps admin null",
			state:            allowedCIDRModel(allowedCIDRTestId, ""),
			entries:          []api.AppServiceAllowedCIDRResponse{allowedCIDREntry(allowedCIDRTestId, "admin")},
			expectAccessMode: types.StringNull(),
		},
		{
			name:  "admin state sends no accessMode and finds the entry on a later page",
			state: allowedCIDRModel(allowedCIDRTestId, api.AllowedCIDRAccessModeAdmin),
			entries: []api.AppServiceAllowedCIDRResponse{
				allowedCIDREntry("other-id", "admin"),
				allowedCIDREntry(allowedCIDRTestId, "admin"),
			},
			expectAccessMode: types.StringValue(api.AllowedCIDRAccessModeAdmin),
		},
		{
			name:              "missing entry removes the resource from the state",
			state:             allowedCIDRModel(allowedCIDRTestId, api.AllowedCIDRAccessModePublic),
			entries:           []api.AppServiceAllowedCIDRResponse{allowedCIDREntry("other-id", "public")},
			expectAccessModes: []string{api.AllowedCIDRAccessModePublic},
			expectRemoved:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			backend := &allowedCIDRBackend{entries: test.entries}
			a := &AppServiceCidr{Data: newTestProviderData(t, backend.handler)}

			state := tfsdk.State{Schema: AllowedCIDRsSchema()}
			diags := state.Set(ctx, test.state)
			require.False(t, diags.HasError(), diags)
			resp := resource.ReadResponse{State: state}

			a.Read(ctx, resource.ReadRequest{State: state}, &resp)

			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			assert.Equal(t, test.expectAccessModes, backend.accessModes)
			if test.expectRemoved {
				assert.True(t, resp.State.Raw.IsNull())
				return
			}
			got := getAllowedCIDRModel(t, resp.State)
			assert.Equal(t, allowedCIDRTestId, got.Id.ValueString())
			assert.Equal(t, "org", got.OrganizationId.ValueString())
			assert.Equal(t, "appservice", got.AppServiceId.ValueString())
			assert.Equal(t, "203.0.113.0/24", got.Cidr.ValueString())
			assert.Equal(t, test.expectAccessMode, got.AccessMode)
		})
	}
}

func TestAllowedCIDRsSchema_AccessModeReplace(t *testing.T) {
	ctx := context.Background()
	cidrSchema := AllowedCIDRsSchema()
	attribute, ok := cidrSchema.Attributes["access_mode"].(*schema.StringAttribute)
	require.True(t, ok, "access_mode is %T", cidrSchema.Attributes["access_mode"])

	tests := []struct {
		name        string
		state       string
		plan        string
		wantReplace bool
	}{
		{name: "omitted to admin does not replace", plan: api.AllowedCIDRAccessModeAdmin},
		{name: "admin to omitted does not replace", state: api.AllowedCIDRAccessModeAdmin},
		{name: "unchanged public does not replace", state: "public", plan: "public"},
		{name: "omitted to public replaces", plan: "public", wantReplace: true},
		{name: "public to omitted replaces", state: "public", wantReplace: true},
		{name: "admin to public replaces", state: api.AllowedCIDRAccessModeAdmin, plan: "public", wantReplace: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := tfsdk.State{Schema: cidrSchema}
			diags := state.Set(ctx, allowedCIDRModel(allowedCIDRTestId, test.state))
			require.False(t, diags.HasError(), diags)
			plan := tfsdk.Plan{Schema: cidrSchema}
			diags = plan.Set(ctx, allowedCIDRModel(allowedCIDRTestId, test.plan))
			require.False(t, diags.HasError(), diags)

			req := planmodifier.StringRequest{
				Path:       path.Root("access_mode"),
				State:      state,
				StateValue: allowedCIDRModel("", test.state).AccessMode,
				Plan:       plan,
				PlanValue:  allowedCIDRModel("", test.plan).AccessMode,
			}
			var gotReplace bool
			for _, modifier := range attribute.PlanModifiers {
				resp := planmodifier.StringResponse{PlanValue: req.PlanValue}
				modifier.PlanModifyString(ctx, req, &resp)
				gotReplace = gotReplace || resp.RequiresReplace
			}

			assert.Equal(t, test.wantReplace, gotReplace)
		})
	}
}

func TestAppServiceCidr_Update(t *testing.T) {
	tests := []struct {
		name       string
		state      string
		accessMode string
	}{
		{name: "admin added to an omitted access mode", accessMode: api.AllowedCIDRAccessModeAdmin},
		{name: "admin removed from the config", state: api.AllowedCIDRAccessModeAdmin},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			a := &AppServiceCidr{}

			state := tfsdk.State{Schema: AllowedCIDRsSchema()}
			diags := state.Set(ctx, allowedCIDRModel(allowedCIDRTestId, test.state))
			require.False(t, diags.HasError(), diags)
			plan := tfsdk.Plan{Schema: AllowedCIDRsSchema()}
			diags = plan.Set(ctx, allowedCIDRModel(allowedCIDRTestId, test.accessMode))
			require.False(t, diags.HasError(), diags)
			resp := resource.UpdateResponse{State: state}

			a.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &resp)

			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			assert.Equal(t, allowedCIDRModel(allowedCIDRTestId, test.accessMode), getAllowedCIDRModel(t, resp.State))
		})
	}
}
