package api

import (
	"bytes"
	"context"
	"encoding/json"
	goer "errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/version"
)

const clientName = "terraform-provider-couchbase-capella"

var userAgent = fmt.Sprintf("%s/%s", clientName, version.ProviderVersion)

// Client is responsible for constructing and executing HTTP requests.
type Client struct {
	*http.Client

	retry retryPolicy
}

// retryPolicy bounds how a Client retries retryable 5xx responses.
type retryPolicy struct {
	// backoff returns how long to wait before the nth retry, 1-based.
	backoff func(attempt int) time.Duration

	// maxServiceErrors is how many 503/504 responses are tolerated before the
	// request is abandoned. Zero means no attempt cap, leaving the overall request
	// deadline as the only bound. It never bounds 429 retries.
	maxServiceErrors int
}

// RetryOption customises a Client's retry behaviour.
type RetryOption func(*Client)

// WithFastBackoff shrinks retry delays to milliseconds so that tests can exercise
// the retry loop without waiting for the production budget.
func WithFastBackoff() RetryOption {
	return func(c *Client) {
		c.retry.backoff = func(attempt int) time.Duration {
			if attempt < 1 {
				attempt = 1
			}

			return time.Duration(attempt) * time.Millisecond
		}
	}
}

// NewClient instantiates a new Client with the provided timeout.
func NewClient(timeout time.Duration, opts ...RetryOption) *Client {
	client := &Client{
		Client: &http.Client{
			Timeout: timeout,
		},
		retry: retryPolicy{
			backoff: exponentialJitterBackoff,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Response struct is used to encapsulate the response details.
type Response struct {
	Response *http.Response
	Body     []byte
}

// EndpointCfg is used to encapsulate request details to endpoints.
type EndpointCfg struct {
	// Url is url of the endpoint to be contacted
	Url string

	// Method is the HTTP method to be requested.
	Method string

	// SuccessStatus represents the HTTP status code associated
	// with a successful response from the endpoint.
	SuccessStatus int

	// MaxServiceErrorRetries bounds how many 503/504 responses are retried for
	// this request. Zero, the default, applies no attempt cap: retries continue
	// until the overall request deadline, which is what every endpoint relies on
	// to wait out long transient conditions such as an in-flight bucket delete.
	// Set it to FastFailServiceErrorRetries on endpoints that have no such
	// condition and should surface an unavailable service quickly instead.
	MaxServiceErrorRetries int
}

// FastFailServiceErrorRetries is a short retry budget for endpoints with no known
// long-running transient condition, where an unavailable service should surface
// quickly rather than consume the full request deadline. It is roughly 23 seconds
// with the default backoff. Opt in per request via
// EndpointCfg.MaxServiceErrorRetries; it is deliberately not the default, because
// shortening the budget globally would stop calls such as bucket create from
// waiting out an in-flight bucket delete.
const FastFailServiceErrorRetries = 5

// Bounds for the backoff applied between retries of a response that carries no
// Retry-After header.
const (
	retryWaitMin = time.Second
	retryWaitMax = 30 * time.Second
)

// exponentialJitterBackoff returns the wait before retry attempt n, 1-based:
// retryWaitMin doubled per attempt and capped at retryWaitMax, then jittered into
// the upper half of that interval. The jitter stops resources from retrying a
// struggling backend in lockstep; halving rather than zeroing keeps a minimum
// interval so the retries never become a busy loop.
func exponentialJitterBackoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}

	wait := retryWaitMax
	if attempt < 32 {
		if shifted := retryWaitMin << (attempt - 1); shifted > 0 && shifted < retryWaitMax {
			wait = shifted
		}
	}

	// The jitter only has to decorrelate retries between parallel resources, so a
	// pseudo-random source is sufficient.
	//nolint:gosec // G404: jitter does not require cryptographic randomness
	return wait/2 + time.Duration(rand.Int64N(int64(wait/2)+1))
}

// ExecuteWithRetry is used to construct and execute a HTTP request with retry.
// It then returns the response.
func (c *Client) ExecuteWithRetry(
	ctx context.Context,
	endpointCfg EndpointCfg,
	payload any,
	authToken string,
	headers map[string]string,
) (response *Response, err error) {
	var requestBody []byte
	if payload != nil {
		if content, ok := headers["Content-Type"]; ok && content == "application/javascript" {
			// json.Marshal will add escape characters to the string payload which makes it invalid javascript, this is a workaround
			js, ok := payload.(string)
			if !ok {
				return nil, fmt.Errorf("%s: %w", errors.ErrNotAString, fmt.Errorf("expected string payload for javascript content type"))
			}
			requestBody = []byte(js)
		} else {
			requestBody, err = json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrMarshallingPayload, err)
			}
		}
	}

	var fn = func() (response *Response, backoff time.Duration, err error) {
		req, err := http.NewRequest(endpointCfg.Method, endpointCfg.Url, bytes.NewReader(requestBody))
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %w", errors.ErrConstructingRequest, err)
		}

		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("User-Agent", userAgent)
		for header, value := range headers {
			req.Header.Set(header, value)
		}
		apiRes, err := c.Do(req)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %w", errors.ErrExecutingRequest, err)
		}
		defer apiRes.Body.Close()

		responseBody, err := io.ReadAll(apiRes.Body)
		if err != nil {
			return
		}

		switch apiRes.StatusCode {
		case endpointCfg.SuccessStatus:
			// success case
		case http.StatusTooManyRequests:
			header := apiRes.Header.Get("Retry-After")
			retryAfter, err := strconv.Atoi(header)
			if err != nil {
				return nil, 0, fmt.Errorf("error parsing Retry-After value from response header")
			}
			retryAfterDur := time.Second * time.Duration(retryAfter)
			tflog.Debug(ctx, "API rate limited", map[string]interface{}{
				"method":      endpointCfg.Method,
				"url":         endpointCfg.Url,
				"retry_after": retryAfterDur.Seconds(),
			})
			return nil, retryAfterDur, errors.ErrRatelimit
		case http.StatusServiceUnavailable:
			var retryAfter time.Duration
			if header := apiRes.Header.Get("Retry-After"); header != "" {
				if secs, parseErr := strconv.Atoi(header); parseErr == nil {
					retryAfter = time.Second * time.Duration(secs)
				}
			}
			tflog.Debug(ctx, "API returned 503 Service Unavailable, retrying", map[string]interface{}{
				"method": endpointCfg.Method,
				"url":    endpointCfg.Url,
			})
			return nil, retryAfter, errors.ErrServiceUnavailable
		case http.StatusGatewayTimeout:
			var apiError Error
			if err := json.Unmarshal(responseBody, &apiError); err != nil {
				return nil, 0, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)
			}

			if apiError.Code == 7001 {
				return nil, 0, errors.ErrGatewayTimeoutForIndexDDL
			}

			return nil, 0, errors.ErrGatewayTimeout
		default:
			var apiError Error
			if err := json.Unmarshal(responseBody, &apiError); err != nil {
				return nil, 0, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)
			}
			if apiError.Code == 0 {
				return nil, 0, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)

			}
			return nil, 0, &apiError
		}

		return &Response{
			Response: apiRes,
			Body:     responseBody,
		}, 0, nil
	}

	policy := c.retry
	policy.maxServiceErrors = endpointCfg.MaxServiceErrorRetries

	return exec(ctx, fn, policy)
}

func exec(
	ctx context.Context, fn func() (response *Response, dur time.Duration, err error), policy retryPolicy,
) (*Response, error) {
	timer := time.NewTimer(time.Millisecond)
	defer timer.Stop()

	var (
		err      error
		backOff  time.Duration
		response *Response
	)

	const timeout = time.Minute * 10

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	// attempts counts every request issued, whereas serviceErrors counts only the
	// retryable 5xx responses that the attempt cap applies to.
	var attempts, serviceErrors int
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			// A persistent 503/504 normally ends here rather than at an attempt cap,
			// so report the attempts and elapsed time instead of a bare deadline.
			if serviceErrors > 0 {
				return nil, &RetryExhaustedError{
					LastErr:  err,
					Attempts: attempts,
					Elapsed:  time.Since(start),
				}
			}

			return nil, fmt.Errorf("timed out executing request against api: %w", ctx.Err())
		case <-timer.C:
			attempts++
			response, backOff, err = fn()
			switch {
			case err == nil:
				return response, nil
			case goer.Is(err, errors.ErrRatelimit):
			case goer.Is(err, errors.ErrServiceUnavailable):
				serviceErrors++
			case !goer.Is(err, errors.ErrGatewayTimeout):
				return response, err
			default:
				serviceErrors++
			}

			if policy.maxServiceErrors > 0 && serviceErrors > policy.maxServiceErrors {
				return nil, &RetryExhaustedError{
					LastErr:  err,
					Attempts: attempts,
					Elapsed:  time.Since(start),
				}
			}

			if backOff > 0 {
				timer.Reset(backOff)
			} else {
				timer.Reset(policy.backoff(serviceErrors))
			}
		}
	}
}
