package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	testProjectReplicationsPath = "/v4/organizations/org-1/projects/proj-1/replications"
	testClusterReplicationsPath = "/v4/organizations/org-1/projects/proj-1/clusters/cluster-1/replications"
)

// replicationsBackend is a stand-in for the list endpoints.
type replicationsBackend struct {
	t     *testing.T
	path  string
	ids   []string
	mu    sync.Mutex
	calls int
}

func (b *replicationsBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.calls++

	if r.Method != http.MethodGet || r.URL.Path != b.path {
		b.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))
	last := (len(b.ids) + perPage - 1) / perPage
	next := 0
	if page < last {
		next = page + 1
	}

	data := make([]map[string]any, 0, len(b.ids))
	for _, id := range b.ids {
		data = append(data, map[string]any{
			"id":            id,
			"sourceCluster": "source",
			"targetCluster": "target",
			"status":        "running",
			"direction":     "oneWay",
			"audit":         map[string]any{"createdBy": "user-1", "createdAt": "2026-10-01T00:00:00Z"},
		})
	}
	require.NoError(b.t, json.NewEncoder(w).Encode(map[string]any{
		"data":   data,
		"cursor": map[string]any{"pages": map[string]any{"page": page, "next": next, "last": last, "perPage": perPage}},
	}))
}

func readReplications(t *testing.T, backend *replicationsBackend, clusterId types.String) *datasource.ReadResponse {
	t.Helper()
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(backend.handler))
	t.Cleanup(srv.Close)
	data := &providerschema.Data{ClientV1: &api.Client{Client: srv.Client()}, HostURL: srv.URL}

	s := ReplicationsSchema()
	config := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	diags := config.Set(ctx, providerschema.Replications{
		OrganizationId: types.StringValue("org-1"),
		ProjectId:      types.StringValue("proj-1"),
		ClusterId:      clusterId,
	})
	require.False(t, diags.HasError(), "%v", diags)

	resp := &datasource.ReadResponse{State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}}
	(&Replications{Data: data}).Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: s, Raw: config.Raw}}, resp)
	require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics.Errors())
	return resp
}

func replicationIds(t *testing.T, resp *datasource.ReadResponse) []string {
	t.Helper()
	var state providerschema.Replications
	require.False(t, resp.State.Get(context.Background(), &state).HasError())
	require.NotNil(t, state.Data, "data must be an empty list, not null")

	ids := make([]string, 0, len(state.Data))
	for _, entry := range state.Data {
		ids = append(ids, entry.Id.ValueString())
	}
	return ids
}

func TestReplicationsRead(t *testing.T) {
	t.Run("removes the duplicates that the API returns on every page, and sorts by ID", func(t *testing.T) {
		ids := make([]string, 0, 30)
		want := make([]string, 0, 30)
		for i := 30; i >= 1; i-- {
			ids = append(ids, fmt.Sprintf("repl-%02d", i))
		}
		for i := 1; i <= 30; i++ {
			want = append(want, fmt.Sprintf("repl-%02d", i))
		}
		backend := &replicationsBackend{t: t, path: testProjectReplicationsPath, ids: ids}

		resp := readReplications(t, backend, types.StringNull())

		assert.Equal(t, 2, backend.calls, "30 entries with 25 for each page are 2 pages")
		assert.Equal(t, want, replicationIds(t, resp))
	})

	t.Run("cluster_id selects the cluster endpoint", func(t *testing.T) {
		backend := &replicationsBackend{t: t, path: testClusterReplicationsPath, ids: []string{"repl-1"}}

		resp := readReplications(t, backend, types.StringValue("cluster-1"))

		assert.Equal(t, []string{"repl-1"}, replicationIds(t, resp))
	})

	t.Run("no replications is an empty list", func(t *testing.T) {
		backend := &replicationsBackend{t: t, path: testProjectReplicationsPath}

		resp := readReplications(t, backend, types.StringNull())

		assert.Empty(t, replicationIds(t, resp))
	})
}
