package datasources

import (
	"context"
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

func TestAppServiceMaintenancesDataSourceRead(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		expectError bool
		wantIDs     []string
	}{
		{
			name:   "lists every record",
			status: http.StatusOK,
			body: `{"data":[` +
				`{"id":"req1","type":"upgrade","initiatedBy":"couchbase","currentState":"pending","queueAt":"2026-09-05T02:30:00Z"},` +
				`{"id":"req2","type":"upgrade","initiatedBy":"selfService","currentState":"completed","versionId":"v1",` +
				`"queueAt":"2026-08-01T00:00:00Z","completedAt":"2026-08-01T01:00:00Z"}` +
				`],"cursor":{"pages":{"page":1,"next":0,"last":1,"perPage":25,"totalItems":2}}}`,
			wantIDs: []string{"req1", "req2"},
		},
		{
			name:    "no records is an empty list",
			status:  http.StatusOK,
			body:    `{"data":[],"cursor":{"pages":{"page":1,"next":0,"last":1,"perPage":25,"totalItems":0}}}`,
			wantIDs: []string{},
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
				assert.Empty(t, r.URL.Query().Get("sortBy"), "the list endpoint takes no sortBy")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)
			d := &AppServiceMaintenances{Data: &providerschema.Data{
				ClientV1: &api.Client{Client: srv.Client()},
				HostURL:  srv.URL,
			}}

			s := AppServiceMaintenancesSchema()
			config := tfsdk.Config{Schema: s}
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, providerschema.AppServiceMaintenances{
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
			var got providerschema.AppServiceMaintenances
			require.False(t, resp.State.Get(ctx, &got).HasError())
			ids := make([]string, 0, len(got.Data))
			for _, m := range got.Data {
				ids = append(ids, m.Id.ValueString())
			}
			assert.Equal(t, tt.wantIDs, ids)
		})
	}
}
