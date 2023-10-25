package api

// Cursor represents pagination metadata for navigating through large data sets.
type Cursor struct {
	// Pages represents the pagination details of the data set.
	Pages Pages `json:"pages"`

	// Hrefs contains the hyperlinks for navigation through the paginated data set.
	Hrefs HRefs `json:"hrefs"`
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
