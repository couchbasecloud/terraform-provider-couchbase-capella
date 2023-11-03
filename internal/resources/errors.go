package resources

import (
	"errors"
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
)

// CheckApiError is used to check if an error is of type
// api.Error error and return it as a string.
func CheckApiError(err error) string {
	var apiErr *api.Error

	if errors.As(err, &apiErr) {
		return apiErr.CompleteError()
	}

	return err.Error()
}

// CheckResourceNotFoundError is used to check if an error is of
// type api.Error whether the error is resource not found.
func CheckResourceNotFoundError(err error) (bool, error) {
	var apiErr *api.Error

	if errors.As(err, &apiErr) {
		if apiErr.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(apiErr.CompleteError())
		}
		return true, fmt.Errorf(apiErr.CompleteError())
	}

	return false, err
}
