package api

import (
	"encoding/json"
	"errors"
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
func (c *Client) Execute(url string, method string, payload io.Reader, apiToken string) (response *Response, err error) {
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	apiRes, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer apiRes.Body.Close()

	body, err := io.ReadAll(apiRes.Body)
	if err != nil {
		return
	}

	if apiRes.StatusCode >= http.StatusBadRequest {
		var error Error
		if err := json.Unmarshal(body, &error); err != nil {
			return nil, err
		}

		return nil, errors.New("received unexpected status code")
	}

	return &Response{
		Response: apiRes,
		Body:     body,
	}, nil
}
