package resources

import (
	"fmt"
	"net/http"
	"terraform-provider-capella/internal/api"
)

// CheckResourceNotFoundError is used to check if an error is of
// type api.Error whether the error is resource not found.
func CheckResourceNotFoundError(err error) (bool, error) {
	switch err := err.(type) {
	case nil:
		return false, nil
	case api.Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, fmt.Errorf(err.CompleteError())
		}
		return true, fmt.Errorf(err.CompleteError())
	default:
		return false, err
	}
}
