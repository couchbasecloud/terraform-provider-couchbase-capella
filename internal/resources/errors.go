package resources

import (
	"errors"
	"net/http"
	"terraform-provider-capella/internal/api"
)

// ParseApiError is used to check if an error is of type
// api.Error error and return it as a string.
func ParseApiError(err error) (bool, string) {
	var apiErr *api.Error

	if errors.As(err, &apiErr) {
		return true, apiErr.CompleteError()
	}

	return false, err.Error()
}

// CheckResourceNotFoundError is used to check if an error is of
// type api.Error whether the error is resource not found.
func CheckResourceNotFoundError(err error) (bool, string) {
	var apiErr *api.Error

	if errors.As(err, &apiErr) {
		if apiErr.HttpStatusCode != http.StatusNotFound {
			return false, apiErr.CompleteError()
		}
		return true, apiErr.CompleteError()
	}

	return false, err.Error()
}
