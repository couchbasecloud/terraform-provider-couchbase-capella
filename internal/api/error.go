package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Error tracks the error structure received from Capella V4 APIs.
type Error struct {
	Hint           string `json:"hint"`
	Message        string `json:"message"`
	Code           int    `json:"code"`
	HttpStatusCode int    `json:"httpStatusCode"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(`{"code":%d,"hint":"%s","httpStatusCode":%d,"message":"%s"}`,
		e.Code, e.Hint, e.HttpStatusCode, e.Message,
	)
}

func (e *Error) Is(target error) bool {
	var t *Error
	return errors.As(target, &t)
}

func (e Error) CompleteError() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"code":%d,"hint":"%s","httpStatusCode":%d,"message":"%s"}`,
			e.Code, e.Hint, e.HttpStatusCode, e.Message,
		)
	}
	return string(jsonData)
}

// ParseError is used to check if an error is of type
// api.Error error and return it as a string.
func ParseError(err error) string {
	var apiError *Error
	switch {
	case errors.As(err, &apiError):
		return apiError.CompleteError()
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
	var apiError *Error
	switch {
	case errors.As(err, &apiError):
		if apiError.HttpStatusCode != http.StatusNotFound {
			return false, apiError.CompleteError()
		}
		return true, apiError.CompleteError()
	default:
		return false, err.Error()
	}
}

// IsForbiddenError checks whether the given error is an api.Error with HTTP status 403.
func IsForbiddenError(err error) bool {
	var apiError *Error
	return errors.As(err, &apiError) && apiError.HttpStatusCode == http.StatusForbidden
}

// transientInternalServerErrorCode is the Capella error code carried by the transient
// backend 500 ("Something went wrong on our end. We are actively investigating the issue.").
const transientInternalServerErrorCode = 10000

// IsTransientInternalServerError checks whether the given error is an api.Error with
// HTTP status 500 and the transient backend error code. Capella returns this while it
// finishes propagating state after a write, and it clears on its own within minutes, so
// callers may retry it. Other 500s indicate real failures and should surface immediately.
func IsTransientInternalServerError(err error) bool {
	var apiError *Error
	return errors.As(err, &apiError) &&
		apiError.HttpStatusCode == http.StatusInternalServerError &&
		apiError.Code == transientInternalServerErrorCode
}
