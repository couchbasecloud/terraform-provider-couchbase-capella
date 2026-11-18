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

func TestAppServiceSelfServiceUpgradesDataSourceRead(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		expectError bool
		expected    func(*testing.T, []providerschema.AppServiceMaintenanceData)
	}{
		{
			name:   "maps type to kind and leaves completed_at null until it finishes",
			status: http.StatusOK,
			body: `{"data":[` +
				`{"id":"u1","type":"upgrade","initiatedBy":"selfService","currentState":"pending","queueAt":"2026-08-24T02:30:00Z"},` +
				`{"id":"u2","type":"upgrade","initiatedBy":"selfService","currentState":"completed","versionId":"3.1.8",` +
				`"queueAt":"2026-08-01T00:00:00Z","completedAt":"2026-08-01T01:00:00Z"}` +
				`],"cursor":{"pages":{"page":1,"next":0,"last":1,"perPage":25,"totalItems":2}}}`,
			expected: func(t *testing.T, data []providerschema.AppServiceMaintenanceData) {
				require.Len(t, data, 2)
				assert.Equal(t, "u1", data[0].Id.ValueString())
				assert.Equal(t, "upgrade", data[0].Kind.ValueString())
				assert.True(t, data[0].CompletedAt.IsNull())
				assert.True(t, data[0].VersionId.IsNull())
				assert.Equal(t, "3.1.8", data[1].VersionId.ValueString())
				assert.False(t, data[1].CompletedAt.IsNull())
			},
		},
		{
			name:   "no upgrades is an empty list",
			status: http.StatusOK,
			body:   `{"data":[],"cursor":{"pages":{"page":1,"next":0,"last":1,"perPage":25,"totalItems":0}}}`,
			expected: func(t *testing.T, data []providerschema.AppServiceMaintenanceData) {
				assert.Empty(t, data)
			},
		},
		{
			name:        "API error is a diagnostic",
			status:      http.StatusForbidden,
			body:        `{"code":2,"message":"This feature is disabled for your organization.","httpStatusCode":403}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.True(t, strings.HasSuffix(r.URL.Path, "/appservices/appservice-1/selfServiceUpgrades"), r.URL.Path)
				assert.Empty(t, r.URL.Query().Get("sortBy"), "the list endpoint takes no sortBy")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)
			d := &AppServiceSelfServiceUpgrades{Data: &providerschema.Data{
				ClientV1: &api.Client{Client: srv.Client()},
				HostURL:  srv.URL,
			}}

			s := AppServiceSelfServiceUpgradesSchema()
			config := tfsdk.Config{Schema: s}
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, providerschema.AppServiceSelfServiceUpgrades{
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
			var got providerschema.AppServiceSelfServiceUpgrades
			require.False(t, resp.State.Get(ctx, &got).HasError())
			tt.expected(t, got.Data)
		})
	}
}
