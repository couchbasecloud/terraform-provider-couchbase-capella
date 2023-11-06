package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error tracks the error structure received from Capella V4 APIs
type Error struct {
	// Code is the HTTP Status code sent by Capella V4 APIs.
	// Common code include: 200, 201, 202, 400, 403, 404, 409, 412, 422, 500
	Code int `json:"code"`
	// Hint tells us why this error occurred and if there is a way to fix it easily.
	Hint string `json:"hint"`
	// Code is the HTTP Status code sent by Capella V4 APIs.
	// Common code include: 200, 201, 202, 400, 403, 404, 409, 412, 422, 500
	HttpStatusCode int `json:"httpStatusCode"`
	// Message is the exact error message sent by the Capella V4 API
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e Error) CompleteError() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return e.Message
	}
	return string(jsonData)
}

// ParseError is used to check if an error is of type
// api.Error error and return it as a string.
func ParseError(err error) string {
	switch err := err.(type) {
	case Error:
		return err.CompleteError()
	default:
		return err.Error()
	}
}

// CheckResourceNotFoundError is used to check if an error is of
// type api.Error and whether the error is resource not found.
//
// Note: If the error is other than not found, the error string
// will be returned along with a bool value of false.
func CheckResourceNotFoundError(err error) (bool, string) {
	switch err := err.(type) {
	case Error:
		if err.HttpStatusCode != http.StatusNotFound {
			return false, err.CompleteError()
		}
		return true, err.CompleteError()
	default:
		return false, err.Error()
	}
}
