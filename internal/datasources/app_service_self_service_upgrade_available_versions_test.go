package datasources

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAppServiceSelfServiceUpgradeAvailableVersionsDataSourceRead(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		expectError bool
		expected    func(*testing.T, []providerschema.AppServiceSelfServiceUpgradeAvailableVersion)
	}{
		{
			name:   "decodes priority as a string and release notes as a list",
			status: http.StatusOK,
			body: `{"data":[{"id":"3.1.8","title":"App Service 3.1.8","description":"Fixes.",` +
				`"priority":"high","releaseNotes":["n1"]}]}`,
			expected: func(t *testing.T, data []providerschema.AppServiceSelfServiceUpgradeAvailableVersion) {
				require.Len(t, data, 1)
				assert.Equal(t, "3.1.8", data[0].Id.ValueString())
				assert.Equal(t, "high", data[0].Priority.ValueString())
				var notes []string
				require.False(t, data[0].ReleaseNotes.ElementsAs(context.Background(), &notes, false).HasError())
				assert.Equal(t, []string{"n1"}, notes)
			},
		},
		{
			name:   "no versions is an empty list",
			status: http.StatusOK,
			body:   `{"data":[]}`,
			expected: func(t *testing.T, data []providerschema.AppServiceSelfServiceUpgradeAvailableVersion) {
				assert.Empty(t, data)
			},
		},
		{
			name:        "API error is a diagnostic",
			status:      http.StatusNotFound,
			body:        `{"code":404,"message":"app service not found","httpStatusCode":404}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.True(t, strings.HasSuffix(r.URL.Path, "/appservices/appservice-1/selfServiceUpgrades/availableVersions"), r.URL.Path)
				assert.Empty(t, r.URL.RawQuery, "the list is not paginated")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)
			d := &AppServiceSelfServiceUpgradeAvailableVersions{Data: &providerschema.Data{
				ClientV1: &api.Client{Client: srv.Client()},
				HostURL:  srv.URL,
			}}

			s := AppServiceSelfServiceUpgradeAvailableVersionsSchema()
			config := tfsdk.Config{Schema: s}
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, providerschema.AppServiceSelfServiceUpgradeAvailableVersions{
				OrganizationId: types.StringValue("org-1"),
				ProjectId:      types.StringValue("proj-1"),
				ClusterId:      types.StringValue("cluster-1"),
				AppServiceId:   types.StringValue("appservice-1"),
			}).HasError())
			config.Raw = state.Raw
			resp := &datasource.ReadResponse{State: tfsdk.State{Schema: s}}

			d.Read(ctx, datasource.ReadRequest{Config: config}, resp)

			require.Equal(t, tt.expectError, resp.Diagnostics.HasError(), resp.Diagnostics)
			if tt.expectError {
				return
			}
			var got providerschema.AppServiceSelfServiceUpgradeAvailableVersions
			require.False(t, resp.State.Get(ctx, &got).HasError())
			tt.expected(t, got.Data)
		})
	}
}
