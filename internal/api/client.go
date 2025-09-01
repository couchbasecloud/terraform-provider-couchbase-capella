package api

import (
	"bytes"
	"context"
	"encoding/json"
	goer "errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/version"
)

const clientName = "terraform-provider-couchbase-capella"

var userAgent = fmt.Sprintf("%s/%s", clientName, version.ProviderVersion)

// Client is responsible for constructing and executing HTTP requests.
type Client struct {
	*http.Client
}

// NewClient instantiates a new Client with the provided timeout.
func NewClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
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
}

// defaultWaitAttempt re-attempt http request after 2 seconds.
const defaultWaitAttempt = time.Second * 2

// Execute is used to construct and execute a HTTP request.
// It then returns the response.
func (c *Client) Execute(
	endpointCfg EndpointCfg,
	payload any,
	authToken string,
	headers map[string]string,
) (response *Response, err error) {
	var requestBody []byte
	if payload != nil {
		if content, ok := headers["Content-Type"]; ok && content == "application/javascript" {
			requestBody = []byte(payload.(string))
		} else {
			requestBody, err = json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", errors.ErrMarshallingPayload, err)
			}
		}
	}

	req, err := http.NewRequest(endpointCfg.Method, endpointCfg.Url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConstructingRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("User-Agent", userAgent)
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	apiRes, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
	}
	defer apiRes.Body.Close()

	responseBody, err := io.ReadAll(apiRes.Body)
	if err != nil {
		return
	}

	if apiRes.StatusCode != endpointCfg.SuccessStatus {
		var apiError Error
		if err := json.Unmarshal(responseBody, &apiError); err != nil {
			return nil, fmt.Errorf(
				"unexpected code: %d, expected: %d, body: %s",
				apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)
		}
		return nil, &apiError
	}

	return &Response{
		Response: apiRes,
		Body:     responseBody,
	}, nil
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
	var dur time.Duration
	if payload != nil {
		if content, ok := headers["Content-Type"]; ok && content == "application/javascript" {
			requestBody = []byte(payload.(string))
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
			return nil, dur, fmt.Errorf("%s: %w", errors.ErrConstructingRequest, err)
		}

		req.Header.Set("Authorization", "Bearer "+authToken)
		req.Header.Set("User-Agent", userAgent)
		for header, value := range headers {
			req.Header.Set(header, value)
		}
		apiRes, err := c.Do(req)
		if err != nil {
			return nil, dur, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
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
				return nil, dur, fmt.Errorf("error parsing Retry-After value from response header")
			}
			dur = time.Second * time.Duration(retryAfter)
			return nil, dur, errors.ErrRatelimit
		case http.StatusGatewayTimeout:
			var apiError Error
			if err := json.Unmarshal(responseBody, &apiError); err != nil {
				return nil, dur, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)
			}

			if apiError.Code == 7001 {
				return nil, 0, errors.ErrGatewayTimeoutForIndexDDL
			}

			return nil, dur, errors.ErrGatewayTimeout
		default:
			var apiError Error
			if err := json.Unmarshal(responseBody, &apiError); err != nil {
				return nil, dur, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)
			}
			if apiError.Code == 0 {
				return nil, dur, fmt.Errorf(
					"unexpected code: %d, expected: %d, body: %s",
					apiRes.StatusCode, endpointCfg.SuccessStatus, responseBody)

			}
			return nil, dur, &apiError
		}

		return &Response{
			Response: apiRes,
			Body:     responseBody,
		}, dur, nil
	}

	return exec(ctx, fn, defaultWaitAttempt)
}

func exec(
	ctx context.Context, fn func() (response *Response, dur time.Duration, err error), waitOnReattempt time.Duration,
) (*Response, error) {
	timer := time.NewTimer(time.Millisecond)

	var (
		err      error
		backOff  time.Duration
		response *Response
	)

	const timeout = time.Minute * 10

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out executing request against api: %w", err)
		case <-timer.C:
			response, backOff, err = fn()
			switch {
			case err == nil:
				return response, nil
			case goer.Is(err, errors.ErrRatelimit):
			case !goer.Is(err, errors.ErrGatewayTimeout):
				return response, err
			}

			if backOff > 0 {
				timer.Reset(backOff)
			} else {
				timer.Reset(waitOnReattempt)
			}
		}
	}
}
