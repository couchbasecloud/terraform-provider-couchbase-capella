package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// publicAccessBackend answers the public access change with the scripted status codes and the App Service GET with
// the scripted public access values, or with getStatus when it is set. The last entry of each script repeats.
type publicAccessBackend struct {
	mu            sync.Mutex
	changeCodes   []int
	publicAccess  []appservice.PublicAccess
	getStatus     int
	changeMethods []string
	gets          int
}

func (b *publicAccessBackend) handler(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if r.Method == http.MethodGet && b.getStatus != 0 {
		w.WriteHeader(b.getStatus)
		_ = json.NewEncoder(w).Encode(map[string]any{"code": b.getStatus, "httpStatusCode": b.getStatus, "message": "scripted"})
		return
	}

	if r.Method == http.MethodGet {
		publicAccess := b.publicAccess[min(b.gets, len(b.publicAccess)-1)]
		b.gets++
		_ = json.NewEncoder(w).Encode(appservice.GetAppServiceResponse{PublicAccess: &publicAccess})
		return
	}

	code := b.changeCodes[min(len(b.changeMethods), len(b.changeCodes)-1)]
	b.changeMethods = append(b.changeMethods, r.Method)
	w.WriteHeader(code)
	if code != http.StatusAccepted {
		_ = json.NewEncoder(w).Encode(map[string]any{"code": code, "httpStatusCode": code, "message": "scripted"})
	}
}

func fastPublicAccessPolling(t *testing.T) {
	t.Helper()
	origTimeout, origInterval := appServicePublicAccessTimeout, appServicePublicAccessPollInterval
	appServicePublicAccessTimeout, appServicePublicAccessPollInterval = time.Second, time.Millisecond
	t.Cleanup(func() {
		appServicePublicAccessTimeout, appServicePublicAccessPollInterval = origTimeout, origInterval
	})
}

func TestAppServicePublicAccess_changePublicAccess(t *testing.T) {
	fastPublicAccessPolling(t)

	tests := []struct {
		name                 string
		requested            appservice.PublicAccess
		changeCodes          []int
		publicAccess         []appservice.PublicAccess
		getStatus            int
		expectChangeMethods  []string
		expectPublicAccess   string
		expectStatus         string
		expectErr            error
		expectErrContains    string
		expectNoRefreshState bool
	}{
		{
			name:                "restrict polls until restricted",
			requested:           appservice.PublicAccessRestricted,
			changeCodes:         []int{http.StatusAccepted},
			publicAccess:        []appservice.PublicAccess{appservice.PublicAccessRestricting, appservice.PublicAccessRestricted},
			expectChangeMethods: []string{http.MethodPost},
			expectPublicAccess:  "restricted",
			expectStatus:        "restricted",
		},
		{
			name:                "unrestrict sends delete and polls until unrestricted",
			requested:           appservice.PublicAccessUnrestricted,
			changeCodes:         []int{http.StatusAccepted},
			publicAccess:        []appservice.PublicAccess{appservice.PublicAccessUnrestricting, appservice.PublicAccessUnrestricted},
			expectChangeMethods: []string{http.MethodDelete},
			expectPublicAccess:  "unrestricted",
			expectStatus:        "unrestricted",
		},
		{
			name:                "409 is retried until the other change completes",
			requested:           appservice.PublicAccessRestricted,
			changeCodes:         []int{http.StatusConflict, http.StatusConflict, http.StatusAccepted},
			publicAccess:        []appservice.PublicAccess{appservice.PublicAccessRestricted},
			expectChangeMethods: []string{http.MethodPost, http.MethodPost, http.MethodPost},
			expectPublicAccess:  "restricted",
			expectStatus:        "restricted",
		},
		{
			name:                "failed change stops polling and returns the refreshed state",
			requested:           appservice.PublicAccessRestricted,
			changeCodes:         []int{http.StatusAccepted},
			publicAccess:        []appservice.PublicAccess{appservice.PublicAccessRestricting, appservice.PublicAccessRestrictFailed},
			expectChangeMethods: []string{http.MethodPost},
			expectPublicAccess:  "restricted",
			expectStatus:        "restrictFailed",
			expectErr:           errors.ErrAppServicePublicAccessChangeFailed,
		},
		{
			name:                 "422 is not retried",
			requested:            appservice.PublicAccessUnrestricted,
			changeCodes:          []int{http.StatusUnprocessableEntity},
			publicAccess:         []appservice.PublicAccess{appservice.PublicAccessUnrestricted},
			expectChangeMethods:  []string{http.MethodDelete},
			expectErrContains:    "scripted",
			expectNoRefreshState: true,
		},
		{
			name:                 "status GET error is returned without retry",
			requested:            appservice.PublicAccessRestricted,
			changeCodes:          []int{http.StatusAccepted},
			getStatus:            http.StatusNotFound,
			expectChangeMethods:  []string{http.MethodPost},
			expectErrContains:    "scripted",
			expectNoRefreshState: true,
		},
		{
			name:                 "change that never completes times out",
			requested:            appservice.PublicAccessRestricted,
			changeCodes:          []int{http.StatusAccepted},
			publicAccess:         []appservice.PublicAccess{appservice.PublicAccessRestricting},
			expectChangeMethods:  []string{http.MethodPost},
			expectErr:            context.DeadlineExceeded,
			expectNoRefreshState: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			backend := &publicAccessBackend{
				changeCodes:  test.changeCodes,
				publicAccess: test.publicAccess,
				getStatus:    test.getStatus,
			}
			a := &AppServicePublicAccess{Data: newTestProviderData(t, backend.handler)}

			state, err := a.changePublicAccess(context.Background(), "org", "project", "cluster", "appservice", test.requested)

			assert.Equal(t, test.expectChangeMethods, backend.changeMethods)
			switch {
			case test.expectErr != nil:
				require.ErrorIs(t, err, test.expectErr)
			case test.expectErrContains != "":
				require.ErrorContains(t, err, test.expectErrContains)
			default:
				require.NoError(t, err)
			}

			if test.expectNoRefreshState {
				assert.Nil(t, state)
				return
			}
			require.NotNil(t, state)
			assert.Equal(t, "appservice", state.AppServiceId.ValueString())
			assert.Equal(t, test.expectPublicAccess, state.PublicAccess.ValueString())
			assert.Equal(t, test.expectStatus, state.PublicAccessStatus.ValueString())
		})
	}
}

func TestAppServicePublicAccess_failedChange(t *testing.T) {
	fastPublicAccessPolling(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		requested  appservice.PublicAccess
		prior      appservice.PublicAccess
		inProgress appservice.PublicAccess
		failed     appservice.PublicAccess
	}{
		{
			name:       "restrict",
			requested:  appservice.PublicAccessRestricted,
			prior:      appservice.PublicAccessUnrestricted,
			inProgress: appservice.PublicAccessRestricting,
			failed:     appservice.PublicAccessRestrictFailed,
		},
		{
			name:       "unrestrict",
			requested:  appservice.PublicAccessUnrestricted,
			prior:      appservice.PublicAccessRestricted,
			inProgress: appservice.PublicAccessUnrestricting,
			failed:     appservice.PublicAccessUnrestrictFailed,
		},
	}

	for _, test := range tests {
		t.Run(test.name+" create saves the state with a warning", func(t *testing.T) {
			a := newPublicAccessResource(t, test.inProgress, test.failed)
			resp := resource.CreateResponse{State: tfsdk.State{Schema: AppServicePublicAccessSchema()}}

			a.Create(ctx, resource.CreateRequest{Plan: publicAccessPlan(t, test.requested)}, &resp)

			assertFailedChangeWarning(t, resp.Diagnostics)
			assertPublicAccessState(t, resp.State, test.requested, test.failed)
		})

		t.Run(test.name+" update errors and keeps the prior state", func(t *testing.T) {
			a := newPublicAccessResource(t, test.inProgress, test.failed)
			prior := publicAccessState(t, test.prior)
			resp := resource.UpdateResponse{State: prior}

			a.Update(ctx, resource.UpdateRequest{Plan: publicAccessPlan(t, test.requested), State: prior}, &resp)

			require.True(t, resp.Diagnostics.HasError())
			assert.Equal(t, "Error changing App Service public access", resp.Diagnostics.Errors()[0].Summary())
			assertPublicAccessState(t, resp.State, test.prior, test.prior)
		})

		t.Run(test.name+" read warns until the change recovers", func(t *testing.T) {
			a := newPublicAccessResource(t, test.failed, test.requested)

			resp := resource.ReadResponse{State: publicAccessState(t, test.failed)}
			a.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)

			assertFailedChangeWarning(t, resp.Diagnostics)
			assertPublicAccessState(t, resp.State, test.requested, test.failed)

			resp = resource.ReadResponse{State: resp.State}
			a.Read(ctx, resource.ReadRequest{State: resp.State}, &resp)

			assert.Empty(t, resp.Diagnostics)
			assertPublicAccessState(t, resp.State, test.requested, test.requested)
		})
	}
}

func TestAppServicePublicAccessSchema_PublicAccessValidator(t *testing.T) {
	attribute, ok := AppServicePublicAccessSchema().Attributes["public_access"].(*schema.StringAttribute)
	require.True(t, ok, "public_access is %T", AppServicePublicAccessSchema().Attributes["public_access"])

	tests := []struct {
		name      string
		value     string
		expectErr bool
	}{
		{name: "restricted is accepted", value: "restricted"},
		{name: "unrestricted is accepted", value: "unrestricted"},
		{name: "in progress status is rejected", value: "restricting", expectErr: true},
		{name: "unknown value is rejected", value: "foo", expectErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp := &validator.StringResponse{}
			for _, v := range attribute.Validators {
				v.ValidateString(context.Background(), validator.StringRequest{ConfigValue: types.StringValue(test.value)}, resp)
			}

			assert.Equal(t, test.expectErr, resp.Diagnostics.HasError(), resp.Diagnostics)
		})
	}
}

// newPublicAccessResource returns a resource whose App Service GETs return publicAccess in order and whose public
// access changes are accepted.
func newPublicAccessResource(t *testing.T, publicAccess ...appservice.PublicAccess) *AppServicePublicAccess {
	t.Helper()
	backend := &publicAccessBackend{changeCodes: []int{http.StatusAccepted}, publicAccess: publicAccess}
	return &AppServicePublicAccess{Data: newTestProviderData(t, backend.handler)}
}

func publicAccessModel(publicAccess appservice.PublicAccess) *providerschema.AppServicePublicAccess {
	return providerschema.NewAppServicePublicAccess("org", "project", "cluster", "appservice", publicAccess)
}

func publicAccessPlan(t *testing.T, publicAccess appservice.PublicAccess) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: AppServicePublicAccessSchema()}
	diags := plan.Set(context.Background(), publicAccessModel(publicAccess))
	require.False(t, diags.HasError(), diags)
	return plan
}

func publicAccessState(t *testing.T, publicAccess appservice.PublicAccess) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: AppServicePublicAccessSchema()}
	diags := state.Set(context.Background(), publicAccessModel(publicAccess))
	require.False(t, diags.HasError(), diags)
	return state
}

func assertFailedChangeWarning(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	require.False(t, diags.HasError(), diags)
	require.Len(t, diags.Warnings(), 1)
	assert.Equal(t, publicAccessChangeFailedSummary, diags.Warnings()[0].Summary())
}

func assertPublicAccessState(t *testing.T, state tfsdk.State, publicAccess, status appservice.PublicAccess) {
	t.Helper()
	var got providerschema.AppServicePublicAccess
	diags := state.Get(context.Background(), &got)
	require.False(t, diags.HasError(), diags)
	assert.Equal(t, string(publicAccess), got.PublicAccess.ValueString())
	assert.Equal(t, string(status), got.PublicAccessStatus.ValueString())
}
