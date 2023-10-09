package api

import "errors"

var (
	// ErrMarshallingPayload is returned when a payload has failed to
	// marshal into a request body.
	ErrMarshallingPayload = errors.New("failed to marshal payload")

	// ErrConstructingRequest is returned when a HTTP.NewRequest has failed.
	ErrConstructingRequest = errors.New("failed to construct request")

	// ErrExecutingRequest is returned when a HTTP request has failed to execute.
	ErrExecutingRequest = errors.New("failed to execute request")
)
