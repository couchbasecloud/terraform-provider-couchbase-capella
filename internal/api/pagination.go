package api

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Cursor represents pagination metadata for navigating through large data sets.
type Cursor struct {
	Hrefs HRefs `json:"hrefs"`
	Pages Pages `json:"pages"`
}

// Pages represents the pagination details of the data set.
type Pages struct {
	// Page is the current page of results, starting from page 1.
	Page int `json:"page"`

	// Next is the number of the next page of results. Not set on the last page.
	Next int `json:"next"`

	// Previous is the of the previous page of results. Not set on the first page.
	Previous int `json:"previous"`

	// Last is the number of the last page of results.
	Last int `json:"last"`

	// PerPage is the number of items displayed in each page.
	PerPage int `json:"perPage"`

	// TotalItems is the total items found by the given query.
	TotalItems int `json:"totalItems"`
}

// Hrefs contains the hyperlinks for navigation through the paginated data set.
type HRefs struct {
	// First is the base URL, endpoint, and path parameters required to fetch the first page of results.
	First string `json:"first"`

	// Last is the the base URL, endpoint, and path parameters required to fetch the last page of results.
	Last string `json:"last"`

	// Previous is the base URL, endpoint, and path parameters required to fetch the previous page of results. Empty if there is no previous page.
	Previous string `json:"pages"`

	// Next is the base URL, endpoint, and path parameters required to fetch the next page of results. Empty if there is no next page.
	Next string `json:"next"`
}

type sortParameter string

const (
	SortById   = "id"
	SortByName = "name"
)

// overlay is a generic struct used to store data and cursor
// information from paginated responses.
type overlay[DataSchema any] struct {
	Data   DataSchema `json:"data"`
	Cursor Cursor     `json:"cursor"`
}

// GetPaginated is a generic function used to handle pagination. It executes a get request
// according to the supplied url parameter. It then iterates through remaining pages to
// flatten paginated responses into a single slice of responses.
func GetPaginated[DataSchema ~[]T, T any](
	ctx context.Context,
	client *Client,
	token string,
	cfg EndpointCfg,
	sortBy sortParameter,
) (DataSchema, error) {
	result, err := getPaginatedInternal[DataSchema, T](ctx, client, token, cfg, sortBy)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// PaginatedResponse contains the paginated data along with the raw JSON response
// from the first page, which can be used to extract additional metadata fields
// that are not part of the standard pagination structure.
type PaginatedResponse[DataSchema ~[]T, T any] struct {
	// Data contains all items from all pages, flattened into a single slice.
	Data DataSchema

	// RawFirstPage contains the raw JSON response from the first page.
	// This can be used to unmarshal additional top-level fields like clusterStats.
	RawFirstPage []byte
}

// GetPaginatedWithMeta is similar to GetPaginated but also returns the raw response
// from the first page, allowing callers to extract additional metadata fields
// that exist alongside the standard "data" and "cursor" fields.
//
// Example usage for extracting clusterStats:
//
//	result, err := api.GetPaginatedWithMeta[[]bucket.GetBucketResponse](ctx, client, token, cfg, api.SortById)
//	if err != nil { return err }
//
//	var meta struct {
//	    ClusterStats *bucket.ClusterStats `json:"clusterStats"`
//	}
//	json.Unmarshal(result.RawFirstPage, &meta)
func GetPaginatedWithMeta[DataSchema ~[]T, T any](
	ctx context.Context,
	client *Client,
	token string,
	cfg EndpointCfg,
	sortBy sortParameter,
) (*PaginatedResponse[DataSchema, T], error) {
	return getPaginatedInternal[DataSchema, T](ctx, client, token, cfg, sortBy)
}

// getPaginatedInternal is the common implementation for pagination logic.
// It handles fetching all pages and returns both the combined data and the raw first page response.
func getPaginatedInternal[DataSchema ~[]T, T any](
	ctx context.Context,
	client *Client,
	token string,
	cfg EndpointCfg,
	sortBy sortParameter,
) (*PaginatedResponse[DataSchema, T], error) {
	var (
		responses    DataSchema
		rawFirstPage []byte
		page         = 1
		perPage      = 25
		baseUrl      = cfg.Url
	)

	for {
		cfg.Url = baseUrl + fmt.Sprintf("?page=%d&perPage=%d", page, perPage)
		if string(sortBy) != "" {
			cfg.Url += fmt.Sprintf("&sortBy=%s", string(sortBy))
		}
		cfg.Method = http.MethodGet

		response, err := client.ExecuteWithRetry(
			ctx,
			cfg,
			nil,
			token,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errors.ErrExecutingRequest, err)
		}

		// Store the raw response from the first page for metadata extraction
		if page == 1 {
			rawFirstPage = response.Body
		}

		var decoded overlay[DataSchema]
		err = json.Unmarshal(response.Body, &decoded)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errors.ErrUnmarshallingResponse, err)
		}

		responses = append(responses, decoded.Data...)

		cursor := decoded.Cursor

		if cursor.Pages.Next == 0 {
			break
		}

		page = cursor.Pages.Next
	}

	return &PaginatedResponse[DataSchema, T]{
		Data:         responses,
		RawFirstPage: rawFirstPage,
	}, nil
}
