package resources

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func testMaintenanceWindowModel() providerschema.AppServiceMaintenanceWindow {
	return providerschema.AppServiceMaintenanceWindow{
		OrganizationId:  types.StringValue(testOrgID),
		ProjectId:       types.StringValue(testProjectID),
		ClusterId:       types.StringValue(testClusterID),
		AppServiceId:    types.StringValue(testAppServiceID),
		PreferredDay:    types.StringValue("sunday"),
		PreferredHour:   types.Int64Value(2),
		PreferredMinute: types.Int64Value(30),
	}
}

// newTestMaintenanceWindow wires the resource to a server that answers every request with
// status and body, and counts the requests it receives.
func newTestMaintenanceWindow(t *testing.T, status int, body string) (*AppServiceMaintenanceWindow, *int32) {
	t.Helper()
	var requests int32
	data := newTestProviderData(t, func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	})
	return &AppServiceMaintenanceWindow{Data: data}, &requests
}

func TestAppServiceMaintenanceWindowSchemaValidators(t *testing.T) {
	attrs := AppServiceMaintenanceWindowSchema().Attributes

	validateString := func(name, value string) bool {
		attr, ok := attrs[name].(*schema.StringAttribute)
		require.True(t, ok, name)
		resp := &validator.StringResponse{}
		for _, v := range attr.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue(value)}, resp)
		}
		return !resp.Diagnostics.HasError()
	}
	validateInt64 := func(name string, value int64) bool {
		attr, ok := attrs[name].(*schema.Int64Attribute)
		require.True(t, ok, name)
		resp := &validator.Int64Response{}
		for _, v := range attr.Validators {
			v.ValidateInt64(context.Background(), validator.Int64Request{ConfigValue: types.Int64Value(value)}, resp)
		}
		return !resp.Diagnostics.HasError()
	}

	t.Run("accepts every lowercase weekday", func(t *testing.T) {
		for _, day := range []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"} {
			assert.True(t, validateString("preferred_day", day), day)
		}
	})
	t.Run("rejects capitalised or unknown days", func(t *testing.T) {
		for _, day := range []string{"Sunday", "SUNDAY", "funday", ""} {
			assert.False(t, validateString("preferred_day", day), day)
		}
	})
	t.Run("accepts hour and minute bounds", func(t *testing.T) {
		assert.True(t, validateInt64("preferred_hour", 0))
		assert.True(t, validateInt64("preferred_hour", 23))
		assert.True(t, validateInt64("preferred_minute", 0))
		assert.True(t, validateInt64("preferred_minute", 59))
	})
	t.Run("rejects out of range hour and minute", func(t *testing.T) {
		assert.False(t, validateInt64("preferred_hour", 24))
		assert.False(t, validateInt64("preferred_hour", -1))
		assert.False(t, validateInt64("preferred_minute", 60))
		assert.False(t, validateInt64("preferred_minute", -1))
	})
}

func TestAppServiceMaintenanceWindowCreateConflict(t *testing.T) {
	ctx := context.Background()
	r, requests := newTestMaintenanceWindow(t, http.StatusConflict,
		`{"code":11074,"message":"exists","httpStatusCode":409}`)

	s := AppServiceMaintenanceWindowSchema()
	plan := tfsdk.Plan{Schema: s}
	require.False(t, plan.Set(ctx, testMaintenanceWindowModel()).HasError())
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s}}

	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	require.True(t, resp.Diagnostics.HasError())
	assert.Equal(t, "App Service maintenance window already exists", resp.Diagnostics.Errors()[0].Summary())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "terraform import couchbase-capella_app_service_maintenance_window")
	assert.Equal(t, int32(1), atomic.LoadInt32(requests), "a 409 must not be followed by any other request")
	assert.True(t, resp.State.Raw.IsNull(), "nothing may be written to state")
}

func TestAppServiceMaintenanceWindowCreate(t *testing.T) {
	ctx := context.Background()
	var methods []string
	data := newTestProviderData(t, func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"preferredDay":"sunday","preferredHour":2,"preferredMinute":30}`))
	})
	r := &AppServiceMaintenanceWindow{Data: data}

	s := AppServiceMaintenanceWindowSchema()
	plan := tfsdk.Plan{Schema: s}
	require.False(t, plan.Set(ctx, testMaintenanceWindowModel()).HasError())
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s}}

	r.Create(ctx, resource.CreateRequest{Plan: plan}, resp)

	require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
	require.Empty(t, resp.Diagnostics.Warnings(), resp.Diagnostics)
	assert.Equal(t, []string{http.MethodPost, http.MethodGet}, methods)
	var got providerschema.AppServiceMaintenanceWindow
	require.False(t, resp.State.Get(ctx, &got).HasError())
	assert.Equal(t, testMaintenanceWindowModel(), got)
}

func TestAppServiceMaintenanceWindowRead(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		body          string
		expectError   bool
		expectRemoved bool
	}{
		{
			name:          "404 removes the resource from state",
			status:        http.StatusNotFound,
			body:          `{"code":11075,"message":"not found","httpStatusCode":404}`,
			expectRemoved: true,
		},
		{
			name:        "500 keeps the state",
			status:      http.StatusInternalServerError,
			body:        `{"code":10000,"message":"boom","httpStatusCode":500}`,
			expectError: true,
		},
		{
			name:        "403 keeps the state",
			status:      http.StatusForbidden,
			body:        `{"code":1002,"message":"not enabled","httpStatusCode":403}`,
			expectError: true,
		},
		{
			name:   "200 refreshes the state",
			status: http.StatusOK,
			body:   `{"preferredDay":"wednesday","preferredHour":23,"preferredMinute":59}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r, _ := newTestMaintenanceWindow(t, tt.status, tt.body)

			s := AppServiceMaintenanceWindowSchema()
			state := tfsdk.State{Schema: s}
			require.False(t, state.Set(ctx, testMaintenanceWindowModel()).HasError())
			resp := &resource.ReadResponse{State: state}

			r.Read(ctx, resource.ReadRequest{State: state}, resp)

			assert.Equal(t, tt.expectError, resp.Diagnostics.HasError())
			assert.Equal(t, tt.expectRemoved, resp.State.Raw.IsNull())
			if tt.expectRemoved {
				return
			}

			var got providerschema.AppServiceMaintenanceWindow
			require.False(t, resp.State.Get(ctx, &got).HasError())
			want := testMaintenanceWindowModel()
			if tt.status == http.StatusOK {
				want.PreferredDay = types.StringValue("wednesday")
				want.PreferredHour = types.Int64Value(23)
				want.PreferredMinute = types.Int64Value(59)
			}
			assert.Equal(t, want, got)
		})
	}
}

func TestAppServiceMaintenanceWindowDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	r, _ := newTestMaintenanceWindow(t, http.StatusNotFound,
		`{"code":11075,"message":"not found","httpStatusCode":404}`)

	state := tfsdk.State{Schema: AppServiceMaintenanceWindowSchema()}
	require.False(t, state.Set(ctx, testMaintenanceWindowModel()).HasError())
	resp := &resource.DeleteResponse{State: state}

	r.Delete(ctx, resource.DeleteRequest{State: state}, resp)

	assert.False(t, resp.Diagnostics.HasError(), "a window that is already gone is not an error")
}
