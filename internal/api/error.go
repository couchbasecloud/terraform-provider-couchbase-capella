package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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

// RetriesExhaustedError is returned when all retry attempts have been exhausted
// for a retryable error (rate limit, service unavailable, or gateway timeout).
// It wraps the original sentinel error so callers can still inspect the cause
// with errors.Is.
type RetriesExhaustedError struct {
	Original error
	Attempts int
	Elapsed  time.Duration
}

func (e *RetriesExhaustedError) Error() string {
	return fmt.Sprintf("request failed after %d attempts over %s: %s",
		e.Attempts, e.Elapsed.Truncate(time.Millisecond), e.Original.Error())
}

func (e *RetriesExhaustedError) Unwrap() error {
	return e.Original
}
