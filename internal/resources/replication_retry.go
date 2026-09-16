package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

// XDCR write operations serialise behind a lock the control plane holds per cluster pair. The
// budgets below are how long each operation waits for that lock before giving up, and are
// overridden in tests to avoid waiting on real wall-clock time.
//
// Note the lock is held for roughly 400s by a create and 120s by a delete, so an operation queued
// behind someone else's create can still exhaust the shorter budget while the lock legitimately has
// time left to run.
var (
	// xdcrCreateLockRetryBudget applies to replication creation.
	xdcrCreateLockRetryBudget = 420 * time.Second

	// xdcrShortLockRetryBudget applies to update, delete, pause and resume.
	xdcrShortLockRetryBudget = 140 * time.Second

	xdcrLockRetryInterval = 15 * time.Second
)

// retryOnXDCRLock re-issues fn while the control plane reports that another XDCR operation holds
// the cluster pair lock, giving up once budget is exhausted. It returns the response body of the
// first call that was not rejected by the lock.
//
// fn returns the raw response and body rather than a typed result because each generated
// replication operation has its own response type; every one of them exposes HTTPResponse and Body,
// so callers close over their own typed call and surface those two.
//
// Detection is by HTTP 423 alone. The status is absent from the OpenAPI spec, so the generated
// client has no typed 423 field and the shared V2 retry policy never retries it; the accompanying
// error code and wording are undocumented too, which is why neither is matched on. On these
// endpoints a 423 can only mean the cluster pair lock.
func retryOnXDCRLock(
	ctx context.Context,
	budget time.Duration,
	fn func() (*http.Response, []byte, error),
) ([]byte, error) {
	deadline := time.Now().Add(budget)

	var lastLockMessage string
	for {
		response, body, err := fn()
		if err != nil {
			return nil, err
		}

		if response == nil || response.StatusCode != http.StatusLocked {
			return body, nil
		}

		lastLockMessage = parseXDCRLockMessage(body)
		tflog.Debug(ctx, "XDCR operation in progress for this cluster pair, retrying", map[string]any{
			"message": lastLockMessage,
		})

		// Clamp the wait so the helper never overruns its budget.
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, xdcrLockError(budget, lastLockMessage, nil)
		}
		wait := xdcrLockRetryInterval
		if wait > remaining {
			wait = remaining
		}

		select {
		case <-ctx.Done():
			return nil, xdcrLockError(budget, lastLockMessage, ctx.Err())
		case <-time.After(wait):
		}
	}
}

func xdcrLockError(budget time.Duration, lastLockMessage string, cause error) error {
	if cause != nil {
		return fmt.Errorf(
			"%w: gave up after %v, last message: %s: %w",
			errors.ErrXDCROperationInProgress, budget, lastLockMessage, cause,
		)
	}
	return fmt.Errorf(
		"%w: gave up after %v, last message: %s",
		errors.ErrXDCROperationInProgress, budget, lastLockMessage,
	)
}

// parseXDCRLockMessage extracts the human-readable part of a lock response. The body is only used
// for logging, so an unparsable one falls back to the raw bytes rather than failing the retry.
func parseXDCRLockMessage(body []byte) string {
	var apiErr apigen.Error
	if err := json.Unmarshal(body, &apiErr); err != nil || apiErr.Message == "" {
		return string(body)
	}
	if apiErr.Hint == "" {
		return apiErr.Message
	}
	return apiErr.Message + " (" + apiErr.Hint + ")"
}
