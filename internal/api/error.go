package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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

// HumanReadableError returns the API error as a plain-text, user-friendly string
// that directly surfaces the message and hint from the Capella API response
// instead of dumping raw JSON.
func (e *Error) HumanReadableError() string {
	var sb strings.Builder
	sb.WriteString(e.Message)

	if e.Hint != "" {
		sb.WriteString("\nHint: ")
		sb.WriteString(e.Hint)
	}

	sb.WriteString(fmt.Sprintf(" (code: %d, HTTP status: %d)", e.Code, e.HttpStatusCode))

	return sb.String()
}

// ParseReadableError extracts a human-readable error string from an error.
// If the error is an api.Error, it returns the message and hint as plain text.
// Otherwise it falls back to err.Error().
func ParseReadableError(err error) string {
	var apiError *Error
	if errors.As(err, &apiError) {
		return apiError.HumanReadableError()
	}
	return err.Error()
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
