package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default Capella URL
const HostURL string = "https://example.com"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// AuthResponse -
type AuthResponse struct {
	Token string `json:"token"`
}

type Response struct {
	Body         []byte
	HTTPResponse *http.Response
}

type Error struct {
	Code           int    `json:"code"`
	Hint           string `json:"hint"`
	HttpStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) CompleteError() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return e.Message
	}
	return string(jsonData)
}

// NewClient -
func NewClient(host, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Capella URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if token != nil {
		c.Token = *token
	}

	return &c, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*Response, error) {
	req.WithContext(ctx)

	token := c.Token
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusAccepted {
		var error Error
		if err := json.Unmarshal(body, &error); err != nil {
			return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
		}
		return nil, error
	}

	resp := Response{
		Body:         body,
		HTTPResponse: res,
	}

	return &resp, nil
}
