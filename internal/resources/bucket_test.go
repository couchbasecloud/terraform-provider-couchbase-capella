package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testBucketID = "YnVja2V0"

// bucketErrorHandler answers every request with an API error that has the given status and message.
func bucketErrorHandler(status int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{"code": status, "httpStatusCode": status, "message": message})
	}
}

// bucketState returns a bucket state that has only the given attributes set.
func bucketState(t *testing.T, attrs map[string]string) tfsdk.State {
	t.Helper()
	schema := BucketSchema()
	state := tfsdk.State{Schema: schema, Raw: tftypes.NewValue(schema.Type().TerraformType(t.Context()), nil)}
	for name, value := range attrs {
		assertNoDiags(t, state.SetAttribute(t.Context(), path.Root(name), value))
	}
	return state
}

func TestBucketRead(t *testing.T) {
	const turnedOffMessage = "Temporarily unavailable while the Cluster is in the Turned Off state."

	refreshState := map[string]string{
		"id":              testBucketID,
		"organization_id": testOrgID,
		"project_id":      testProjectID,
		"cluster_id":      testClusterID,
	}
	importState := map[string]string{
		"id": fmt.Sprintf("id=%s,cluster_id=%s,project_id=%s,organization_id=%s",
			testBucketID, testClusterID, testProjectID, testOrgID),
	}

	tests := []struct {
		name        string
		state       map[string]string
		status      int
		message     string
		wantWarning string
		wantError   string
		wantRemoved bool
	}{
		{
			name:        "a turned off cluster is a warning and keeps the state",
			state:       refreshState,
			status:      http.StatusUnprocessableEntity,
			message:     turnedOffMessage,
			wantWarning: "Bucket not refreshed",
		},
		{
			name:      "an import while the cluster is turned off is an error",
			state:     importState,
			status:    http.StatusUnprocessableEntity,
			message:   turnedOffMessage,
			wantError: "Error reading bucket",
		},
		{
			name:      "a 422 that is not about the cluster state is an error",
			state:     refreshState,
			status:    http.StatusUnprocessableEntity,
			message:   "Request validation error.",
			wantError: "Error reading bucket",
		},
		{
			name:        "a deleted bucket is removed from state",
			state:       refreshState,
			status:      http.StatusNotFound,
			message:     "The requested bucket does not exist.",
			wantRemoved: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bucket{Data: newTestProviderData(t, bucketErrorHandler(tt.status, tt.message))}

			prior := bucketState(t, tt.state)
			resp := resource.ReadResponse{State: prior}
			b.Read(t.Context(), resource.ReadRequest{State: prior}, &resp)

			switch {
			case tt.wantError != "":
				require.True(t, resp.Diagnostics.HasError())
				assert.Equal(t, tt.wantError, resp.Diagnostics.Errors()[0].Summary())
			case tt.wantWarning != "":
				assert.False(t, resp.Diagnostics.HasError())
				require.Len(t, resp.Diagnostics.Warnings(), 1)
				assert.Equal(t, tt.wantWarning, resp.Diagnostics.Warnings()[0].Summary())
				assert.True(t, resp.State.Raw.Equal(prior.Raw), "the state should not change")
			case tt.wantRemoved:
				assertNoDiags(t, resp.Diagnostics)
				assert.True(t, resp.State.Raw.IsNull(), "the bucket should be removed from state")
			}
		})
	}
}
