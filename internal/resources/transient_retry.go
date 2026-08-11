package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

// These are vars rather than consts so unit tests can shorten them; production
// behavior is unchanged.
var (
	// transientRetryWindow bounds how long a request is retried while Capella keeps
	// returning the transient backend 500.
	transientRetryWindow = 5 * time.Minute

	// transientRetryInterval is the delay between attempts within the retry window.
	transientRetryInterval = 15 * time.Second
)

// executeWithTransientRetry executes an API request, retrying for up to
// transientRetryWindow while Capella returns the transient backend 500
// (api.IsTransientInternalServerError). That 500 appears for minutes at a time while
// the backend propagates state after a write, and api.Client.ExecuteWithRetry does not
// cover it — that helper retries only 429, 503 and 504. Any other error, including
// other 500s, is returned immediately rather than burning the retry window.
//
// operation names the request in the retry-exhausted error, e.g. "schedule create".
func executeWithTransientRetry(
	ctx context.Context,
	client *api.Client,
	token, operation string,
	cfg api.EndpointCfg,
	payload any,
) (*api.Response, error) {
	deadline := time.Now().Add(transientRetryWindow)

	var lastErr error
	exhausted := func(cause error) error {
		return fmt.Errorf("retry window (%v) exhausted for %s: %w", transientRetryWindow, operation, cause)
	}

	for {
		response, err := client.ExecuteWithRetry(ctx, cfg, payload, token, nil)
		if err == nil {
			return response, nil
		}
		if ctx.Err() != nil {
			cause := lastErr
			if cause == nil {
				cause = err
			}
			return nil, exhausted(cause)
		}
		if !api.IsTransientInternalServerError(err) {
			return response, err
		}

		lastErr = err
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, exhausted(lastErr)
		}

		tflog.Debug(ctx, "request returned a transient 500; retrying", map[string]any{
			"operation": operation,
			"err":       err.Error(),
		})

		sleep := transientRetryInterval
		if sleep > remaining {
			sleep = remaining
		}
		select {
		case <-ctx.Done():
			return nil, exhausted(lastErr)
		case <-time.After(sleep):
		}
	}
}
