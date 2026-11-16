package resources

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/appservice"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
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

func TestAppServicePublicAccess_changePublicAccess(t *testing.T) {
	origTimeout, origInterval := appServicePublicAccessTimeout, appServicePublicAccessPollInterval
	appServicePublicAccessTimeout, appServicePublicAccessPollInterval = time.Second, time.Millisecond
	t.Cleanup(func() {
		appServicePublicAccessTimeout, appServicePublicAccessPollInterval = origTimeout, origInterval
	})

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
