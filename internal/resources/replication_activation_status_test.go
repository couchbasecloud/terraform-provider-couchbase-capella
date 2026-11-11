package resources

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	replicationapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/replication"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const testActivationStatusPath = testReplicationPath + "/activationStatus"

// activationStatusBackend is an in-memory stand-in for the replication GET and the pause and
// resume endpoints. The HTTP codes and the replication statuses come from scripted lists; the last
// one repeats. A status is used only by a GET that returns 200.
type activationStatusBackend struct {
	t           *testing.T
	mu          sync.Mutex
	statuses    []string
	getCodes    []int
	changeCodes []int
	okGets      int
	calls       map[string]int
}

func newActivationStatusBackend(t *testing.T, statuses ...string) *activationStatusBackend {
	return &activationStatusBackend{
		t:           t,
		statuses:    statuses,
		getCodes:    []int{http.StatusOK},
		changeCodes: []int{http.StatusNoContent},
		calls:       map[string]int{},
	}
}

func scripted(values []int, call int) int {
	return values[min(call-1, len(values)-1)]
}

// writeActivationStatusError writes the 422 that the API returns while a cluster is not healthy,
// or a plain error for other codes.
func writeActivationStatusError(w http.ResponseWriter, code int) {
	message := ""
	if code == http.StatusUnprocessableEntity {
		message = testClusterNotHealthyMessage
	}
	writeReplicationError(w, code, message)
}

func (b *activationStatusBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	key := r.Method + " " + r.URL.Path
	b.calls[key]++

	switch key {
	case http.MethodGet + " " + testReplicationPath:
		if code := scripted(b.getCodes, b.calls[key]); code != http.StatusOK {
			writeActivationStatusError(w, code)
			return
		}
		b.okGets++
		status := b.statuses[min(b.okGets-1, len(b.statuses)-1)]
		_, _ = fmt.Fprintf(w, `{"id": %q, "status": %q}`, testReplicationID, status)

	case http.MethodPost + " " + testActivationStatusPath, http.MethodDelete + " " + testActivationStatusPath:
		changes := b.calls[http.MethodPost+" "+testActivationStatusPath] + b.calls[http.MethodDelete+" "+testActivationStatusPath]
		if code := scripted(b.changeCodes, changes); code != http.StatusNoContent {
			writeActivationStatusError(w, code)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		b.t.Errorf("unexpected request: %s", key)
		w.WriteHeader(http.StatusTeapot)
	}
}

func (b *activationStatusBackend) callCount(method string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls[method+" "+testActivationStatusPath]
}

// fastActivationStatusPolling shrinks the poll timing for the duration of a test.
func fastActivationStatusPolling(t *testing.T, timeout time.Duration) {
	t.Helper()
	origInterval, origTimeout := replicationActivationStatusPollInterval, replicationActivationStatusTimeout
	replicationActivationStatusPollInterval = time.Millisecond
	replicationActivationStatusTimeout = timeout
	t.Cleanup(func() {
		replicationActivationStatusPollInterval = origInterval
		replicationActivationStatusTimeout = origTimeout
	})
}

func testActivationStatusModel(state string) providerschema.ReplicationActivationStatus {
	return *providerschema.NewReplicationActivationStatus(testOrgID, testProjectID, testClusterID, testReplicationID, state)
}

func activationStatusPlanFrom(t *testing.T, ctx context.Context, model providerschema.ReplicationActivationStatus) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: ReplicationActivationStatusSchema()}
	assertNoDiags(t, plan.Set(ctx, model))
	return plan
}

func activationStatusStateFrom(t *testing.T, ctx context.Context, model providerschema.ReplicationActivationStatus) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: ReplicationActivationStatusSchema()}
	assertNoDiags(t, state.Set(ctx, model))
	return state
}

func emptyActivationStatusState(ctx context.Context) tfsdk.State {
	s := ReplicationActivationStatusSchema()
	return tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
}

func TestActivationStatus(t *testing.T) {
	tests := []struct {
		status replicationapi.Status
		want   string
	}{
		{replicationapi.StatusPending, "running"},
		{replicationapi.StatusRunning, "running"},
		{replicationapi.StatusPausing, "paused"},
		{replicationapi.StatusPaused, "paused"},
		{replicationapi.StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, activationStatus(tt.status))
		})
	}
}

func TestReplicationActivationStatusCreate(t *testing.T) {
	ctx := context.Background()
	fastActivationStatusPolling(t, 100*time.Millisecond)

	tests := []struct {
		name        string
		desired     string
		statuses    []string
		getCodes    []int
		changeCodes []int
		wantPauses  int
		wantResumes int
		wantErr     string
	}{
		{
			name:       "pauses a running replication and waits until it is pausing",
			desired:    "paused",
			statuses:   []string{"running", "pausing"},
			wantPauses: 1,
		},
		{
			name:        "resumes a paused replication and waits until it is pending",
			desired:     "running",
			statuses:    []string{"paused", "pending"},
			wantResumes: 1,
		},
		{
			name:     "sends no request when the replication already has the requested status",
			desired:  "paused",
			statuses: []string{"pausing"},
		},
		{
			name:        "a 409 is a success when a GET shows the requested status",
			desired:     "paused",
			statuses:    []string{"running", "paused"},
			changeCodes: []int{http.StatusConflict},
			wantPauses:  1,
		},
		{
			name:        "a 409 is an error when a GET does not show the requested status",
			desired:     "paused",
			statuses:    []string{"running"},
			changeCodes: []int{http.StatusConflict},
			wantPauses:  1,
			wantErr:     "Conflict",
		},
		{
			name:       "retries the GET while a cluster is not healthy",
			desired:    "paused",
			statuses:   []string{"running", "pausing"},
			getCodes:   []int{http.StatusUnprocessableEntity, http.StatusOK},
			wantPauses: 1,
		},
		{
			name:        "retries the pause while a cluster is not healthy",
			desired:     "paused",
			statuses:    []string{"running", "running", "pausing"},
			changeCodes: []int{http.StatusUnprocessableEntity, http.StatusNoContent},
			wantPauses:  2,
		},
		{
			name:        "a failed status after the request is an error",
			desired:     "running",
			statuses:    []string{"paused", "failed"},
			wantResumes: 1,
			wantErr:     errors.ErrReplicationActivationStatusFailed.Error(),
		},
		{
			name:       "a status that does not change is a timeout",
			desired:    "paused",
			statuses:   []string{"running"},
			wantPauses: 1,
			wantErr:    errors.ErrReplicationActivationStatusTimeout.Error(),
		},
		{
			name:     "a missing replication is an error that explains which cluster to use",
			desired:  "paused",
			getCodes: []int{http.StatusNotFound},
			wantErr:  "reverse_replication_id and the target cluster ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newActivationStatusBackend(t, tt.statuses...)
			if tt.getCodes != nil {
				backend.getCodes = tt.getCodes
			}
			if tt.changeCodes != nil {
				backend.changeCodes = tt.changeCodes
			}
			r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

			resp := resource.CreateResponse{State: emptyActivationStatusState(ctx)}
			r.Create(ctx, resource.CreateRequest{Plan: activationStatusPlanFrom(t, ctx, testActivationStatusModel(tt.desired))}, &resp)

			assert.Equal(t, tt.wantPauses, backend.callCount(http.MethodDelete))
			assert.Equal(t, tt.wantResumes, backend.callCount(http.MethodPost))

			if tt.wantErr != "" {
				require.True(t, resp.Diagnostics.HasError())
				assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.wantErr)
				assert.True(t, resp.State.Raw.IsNull(), "no state is saved after an error")
				return
			}

			assertNoDiags(t, resp.Diagnostics)
			var got providerschema.ReplicationActivationStatus
			assertNoDiags(t, resp.State.Get(ctx, &got))
			assert.Equal(t, tt.desired, got.State.ValueString())
		})
	}
}

func TestReplicationActivationStatusRead(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		prior     string
		status    string
		wantState string
	}{
		{name: "a pausing replication is shown as paused", prior: "paused", status: "pausing", wantState: "paused"},
		{name: "a pending replication is shown as running", prior: "running", status: "pending", wantState: "running"},
		{name: "a failed replication is shown as failed, so the plan shows a change", prior: "running", status: "failed", wantState: "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := newActivationStatusBackend(t, tt.status)
			r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

			prior := activationStatusStateFrom(t, ctx, testActivationStatusModel(tt.prior))
			resp := resource.ReadResponse{State: prior}
			r.Read(ctx, resource.ReadRequest{State: prior}, &resp)
			assertNoDiags(t, resp.Diagnostics)

			var got providerschema.ReplicationActivationStatus
			assertNoDiags(t, resp.State.Get(ctx, &got))
			assert.Equal(t, tt.wantState, got.State.ValueString())
		})
	}

	t.Run("a deleted replication is removed from state", func(t *testing.T) {
		backend := newActivationStatusBackend(t)
		backend.getCodes = []int{http.StatusNotFound}
		r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

		prior := activationStatusStateFrom(t, ctx, testActivationStatusModel("running"))
		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

		assertNoDiags(t, resp.Diagnostics)
		assert.True(t, resp.State.Raw.IsNull())
	})

	t.Run("an unhealthy source cluster is a warning and keeps the state", func(t *testing.T) {
		backend := newActivationStatusBackend(t)
		backend.getCodes = []int{http.StatusUnprocessableEntity}
		r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

		prior := activationStatusStateFrom(t, ctx, testActivationStatusModel("running"))
		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)

		assert.False(t, resp.Diagnostics.HasError())
		require.Len(t, resp.Diagnostics.Warnings(), 1)
		assert.Equal(t, "Replication activation status not refreshed", resp.Diagnostics.Warnings()[0].Summary())
		assert.True(t, resp.State.Raw.Equal(prior.Raw))
	})

	t.Run("import resolves the IDs from the import string", func(t *testing.T) {
		backend := newActivationStatusBackend(t, "paused")
		r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

		prior := activationStatusStateFrom(t, ctx, providerschema.ReplicationActivationStatus{
			OrganizationId: types.StringNull(),
			ProjectId:      types.StringNull(),
			ClusterId:      types.StringNull(),
			ReplicationId: types.StringValue(fmt.Sprintf("replication_id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
				testReplicationID, testClusterID, testProjectID, testOrgID)),
			State: types.StringNull(),
		})
		resp := resource.ReadResponse{State: prior}
		r.Read(ctx, resource.ReadRequest{State: prior}, &resp)
		assertNoDiags(t, resp.Diagnostics)

		var got providerschema.ReplicationActivationStatus
		assertNoDiags(t, resp.State.Get(ctx, &got))
		assert.Equal(t, testActivationStatusModel("paused"), got)
	})
}

func TestReplicationActivationStatusUpdate(t *testing.T) {
	ctx := context.Background()
	fastActivationStatusPolling(t, 100*time.Millisecond)

	backend := newActivationStatusBackend(t, "paused", "running")
	r := &ReplicationActivationStatus{Data: newTestProviderData(t, backend.handler)}

	prior := activationStatusStateFrom(t, ctx, testActivationStatusModel("paused"))
	resp := resource.UpdateResponse{State: prior}
	r.Update(ctx, resource.UpdateRequest{
		Plan:  activationStatusPlanFrom(t, ctx, testActivationStatusModel("running")),
		State: prior,
	}, &resp)
	assertNoDiags(t, resp.Diagnostics)

	assert.Equal(t, 1, backend.callCount(http.MethodPost))
	var got providerschema.ReplicationActivationStatus
	assertNoDiags(t, resp.State.Get(ctx, &got))
	assert.Equal(t, "running", got.State.ValueString())
}
