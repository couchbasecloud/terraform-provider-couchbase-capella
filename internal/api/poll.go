package api

import "context"

// PollError returns timeoutErr when ctx has ended, and err otherwise.
//
// Requests are cancelled with their context, so a deadline that passes while a
// poller's request is in flight fails the request with the context's error. A
// poller that returns request errors should pass them through PollError, so a
// deadline reports the poller's timeout error whether it passes during a request
// or between requests.
func PollError(ctx context.Context, err, timeoutErr error) error {
	if ctx.Err() != nil {
		return timeoutErr
	}
	return err
}
