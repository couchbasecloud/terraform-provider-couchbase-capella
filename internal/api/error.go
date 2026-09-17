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

// RetryExhaustedError is returned when a retryable 503 or 504 response persisted
// past the retry budget. It reports the attempt count and the elapsed time so an
// operator can tell a genuinely unavailable service from a merely slow one.
type RetryExhaustedError struct {
	// LastErr is the error produced by the final attempt.
	LastErr error

	// Attempts is the number of requests issued, including the first.
	Attempts int

	// Elapsed is the wall-clock time spent across all attempts.
	Elapsed time.Duration
}

func (e *RetryExhaustedError) Error() string {
	return fmt.Sprintf("%v after %d attempts over %s",
		e.LastErr, e.Attempts, e.Elapsed.Round(time.Millisecond))
}

// Unwrap exposes the final attempt's error so that callers can still match the
// underlying cause with errors.Is.
func (e *RetryExhaustedError) Unwrap() error { return e.LastErr }
