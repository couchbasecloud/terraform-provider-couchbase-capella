package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

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

// Response struct is used to encapsulate the response details
type Response struct {
	Response *http.Response
	Body     []byte
}

// RequestCfg is used to encapsulate request details to endpoints
type EndpointCfg struct {
	// Url is url of the endpoint to be contacted
	Url string

	// Method is the HTTP method to be requested.
	Method string

	// SuccessStatus represents the HTTP status code associated
	// with a successful response from the endpoint.
	SuccessStatus int
}

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
		requestBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errors.ErrMarshallingPayload, err)
		}
	}

	req, err := http.NewRequest(endpointCfg.Method, endpointCfg.Url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrConstructingRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
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
