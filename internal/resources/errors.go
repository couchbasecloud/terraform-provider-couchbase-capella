package resources

import (
	"errors"
	"net/http"
	"terraform-provider-capella/internal/api"
)

// ParseError is used to check if an error is of type
// api.Error error and return it as a string.
func ParseError(err error) string {
	var apiErr *api.Error

	if errors.As(err, &apiErr) {
		return apiErr.CompleteError()
	}

	return err.Error()
}

// CheckResourceNotFoundError is used to check if an error is of
// type api.Error and whether the error is resource not found.
//
// Note: If the error is other than not found, the error string
// will be returned along with a bool value of false.
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
