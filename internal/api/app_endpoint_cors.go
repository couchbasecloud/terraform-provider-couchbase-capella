package api

// CorsRequest represents the request body for creating/updating CORS configuration.
type CorsRequest struct {
	// Origin is a list of allowed origins for CORS.
	Origin []string `json:"origin"`

	// LoginOrigin is a list of allowed login origins for CORS.
	LoginOrigin []string `json:"loginOrigin"`

	// Headers is a list of allowed headers for CORS.
	Headers []string `json:"headers"`

	// MaxAge specifies the duration (in seconds) for which the results of a preflight request can be cached.
	MaxAge int64 `json:"maxAge"`

	// Disabled indicates whether CORS is disabled.
	Disabled bool `json:"disabled"`
}

// CorsResponse represents the response body for CORS configuration.
type CorsResponse struct {
	// Origin is a list of allowed origins for CORS.
	Origin []string `json:"origin"`

	// LoginOrigin is a list of allowed login origins for CORS.
	LoginOrigin []string `json:"loginOrigin"`

	// Headers is a list of allowed headers for CORS.
	Headers []string `json:"headers"`

	// MaxAge specifies the duration (in seconds) for which the results of a preflight request can be cached.
	MaxAge int64 `json:"maxAge"`

	// Disabled indicates whether CORS is disabled.
	Disabled bool `json:"disabled"`
}
