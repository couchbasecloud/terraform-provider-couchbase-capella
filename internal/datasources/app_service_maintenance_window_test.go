package datasources

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func TestAppServiceMaintenanceWindowDataSourceRead(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		expectError bool
		wantDay     types.String
		wantHour    types.Int64
		wantMinute  types.Int64
	}{
		{
			name:       "window set fills the attributes",
			status:     http.StatusOK,
			body:       `{"preferredDay":"sunday","preferredHour":2,"preferredMinute":30}`,
			wantDay:    types.StringValue("sunday"),
			wantHour:   types.Int64Value(2),
			wantMinute: types.Int64Value(30),
		},
		{
			name:       "no window set leaves the attributes null",
			status:     http.StatusNotFound,
			body:       fmt.Sprintf(`{"code":%d,"message":"no window","httpStatusCode":404}`, errorCodeForMaintenanceWindowNotFound),
			wantDay:    types.StringNull(),
			wantHour:   types.Int64Null(),
			wantMinute: types.Int64Null(),
		},
		{
			name:        "unknown app service is an error",
			status:      http.StatusNotFound,
			body:        `{"code":404,"message":"app service not found","httpStatusCode":404}`,
			expectError: true,
		},
		{
			name:        "unknown cluster is an error",
			status:      http.StatusNotFound,
			body:        `{"code":4025,"message":"cluster not found","httpStatusCode":404}`,
			expectError: true,
		},
		{
			name:        "feature not enabled is an error",
			status:      http.StatusForbidden,
			body:        `{"code":1002,"message":"not enabled","httpStatusCode":403}`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)
			d := &AppServiceMaintenanceWindow{Data: &providerschema.Data{
				ClientV1: &api.Client{Client: srv.Client()},
				HostURL:  srv.URL,
			}}

			s := AppServiceMaintenanceWindowSchema()
			config := tfsdk.Config{Schema: s}
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, providerschema.AppServiceMaintenanceWindow{
				OrganizationId:  types.StringValue("org-1"),
				ProjectId:       types.StringValue("proj-1"),
				ClusterId:       types.StringValue("cluster-1"),
				AppServiceId:    types.StringValue("appservice-1"),
				PreferredDay:    types.StringNull(),
				PreferredHour:   types.Int64Null(),
				PreferredMinute: types.Int64Null(),
			}).HasError())
			config.Raw = state.Raw
			resp := &datasource.ReadResponse{State: tfsdk.State{Schema: s}}

			d.Read(ctx, datasource.ReadRequest{Config: config}, resp)

			require.Equal(t, tt.expectError, resp.Diagnostics.HasError(), resp.Diagnostics)
			if tt.expectError {
				return
			}
			var got providerschema.AppServiceMaintenanceWindow
			require.False(t, resp.State.Get(ctx, &got).HasError())
			assert.Equal(t, "appservice-1", got.AppServiceId.ValueString())
			assert.Equal(t, tt.wantDay, got.PreferredDay)
			assert.Equal(t, tt.wantHour, got.PreferredHour)
			assert.Equal(t, tt.wantMinute, got.PreferredMinute)
		})
	}
}
