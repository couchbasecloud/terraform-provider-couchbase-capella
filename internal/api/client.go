package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is responsible for constructing and executing HTTP requests.
type Client struct {
	*http.Client
}

// NewClient instantiates a new Client.
func NewClient(timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Reponse stuct is used to encapsulate the response details
type Response struct {
	Response *http.Response
	Body     []byte
}

// Execute is used to construct and execute a HTTP request.
// It then returns the response.
func (c *Client) Execute(url string, method string, payload any, authToken string, headers map[string]string) (response *Response, err error) {
	var requestBody []byte
	if payload != nil {
		requestBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	for header, value := range headers {
		req.Header.Set(header, value)
	}

	apiRes, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()

	responseBody, err := io.ReadAll(apiRes.Body)
	if err != nil {
		return
	}

	if apiRes.StatusCode >= http.StatusBadRequest {
		var apiError Error
		if err := json.Unmarshal(responseBody, &apiError); err != nil {
			return nil, err
		}

		return nil, errors.New("received unexpected status code")
	}

	return &Response{
		Response: apiRes,
		Body:     responseBody,
	}, nil
}
