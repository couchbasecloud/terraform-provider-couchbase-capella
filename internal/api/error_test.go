package api

import (
	"fmt"
	"net/http"
	"testing"

	internalerrors "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"

	"github.com/stretchr/testify/assert"
)

func Test_ParseError(t *testing.T) {
	var (
		errMessage = "Error received from Capella V4 Api"
		apiError   = Error{Message: errMessage}
	)

	type test struct {
		err       error
		expOutput interface{}
		name      string
	}

	tests := []test{
		{
			name:      "Error received from Capella V4 Api",
			err:       &Error{Message: errMessage},
			expOutput: Error{Message: errMessage}.CompleteError(),
		},
		{
			name:      "Wrapped error received from Capella V4 Api",
			err:       fmt.Errorf("received error: %w", &apiError),
			expOutput: Error{Message: errMessage}.CompleteError(),
		},
		{
			name:      "Error other than received from Capella V4 Api",
			err:       internalerrors.ErrAllowListIdCannotBeEmpty,
			expOutput: internalerrors.ErrAllowListIdCannotBeEmpty.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errString := ParseError(test.err)

			assert.Equal(t, test.expOutput, errString)
		})
	}
}

func Test_CheckResourceNotFound(t *testing.T) {
	var (
		errMessage = "example error message"
		err500     = Error{
			HttpStatusCode: http.StatusInternalServerError,
			Message:        errMessage,
		}
		err404 = Error{
			HttpStatusCode: http.StatusNotFound,
			Message:        errMessage,
		}
	)

	type test struct {
		err       error
		expOutput interface{}
		name      string
		expBool   bool
	}

	tests := []test{
		{
			name:      "Error received from Capella V4 Api",
			err:       &err500,
			expOutput: err500.CompleteError(),
			expBool:   false,
		},
		{
			name:      "Error received from Capella V4 Api - Resource Not Found",
			err:       &err404,
			expOutput: err404.CompleteError(),
			expBool:   true,
		},
		{
			name:      "Wrapped 500 error received from Capella V4 Api",
			err:       fmt.Errorf("received error: %w", &err500),
			expOutput: err500.CompleteError(),
			expBool:   false,
		},
		{
			name:      "Wrapped 404 error received from Capella V4 Api",
			err:       fmt.Errorf("received error: %w", &err404),
			expOutput: err404.CompleteError(),
			expBool:   true,
		},
		{
			name:      "Error other than received from Capella V4 Api",
			err:       internalerrors.ErrAllowListIdCannotBeEmpty,
			expOutput: internalerrors.ErrAllowListIdCannotBeEmpty.Error(),
			expBool:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourceNotFound, errString := CheckResourceNotFoundError(test.err)

			assert.Equal(t, test.expOutput, errString)
			assert.Equal(t, test.expBool, resourceNotFound)
		})
	}
}

func Test_HumanReadableError(t *testing.T) {
	type test struct {
		apiError       Error
		name           string
		expectContains []string
	}

	tests := []test{
		{
			name: "Message and hint are shown as plain text",
			apiError: Error{
				Code:           422,
				HttpStatusCode: http.StatusUnprocessableEntity,
				Message:        "invalid tenantID. You can locate your tenantID by following the instructions listed here: https://learn.microsoft.com/partner-center/find-ids-and-domain-names",
				Hint:           "Check your Azure Active Directory tenant ID",
			},
			expectContains: []string{
				"invalid tenantID",
				"https://learn.microsoft.com/partner-center/find-ids-and-domain-names",
				"Hint: Check your Azure Active Directory tenant ID",
				"(code: 422, HTTP status: 422)",
			},
		},
		{
			name: "Message without hint omits hint line",
			apiError: Error{
				Code:           500,
				HttpStatusCode: http.StatusInternalServerError,
				Message:        "internal server error",
			},
			expectContains: []string{
				"internal server error",
				"(code: 500, HTTP status: 500)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.apiError.HumanReadableError()

			for _, s := range test.expectContains {
				assert.Contains(t, result, s)
			}
			// Should never contain JSON braces
			assert.NotContains(t, result, `{"code"`)
		})
	}
}

func Test_ParseReadableError(t *testing.T) {
	type test struct {
		err            error
		name           string
		expectContains []string
	}

	tests := []test{
		{
			name: "API error is returned as human-readable text",
			err: &Error{
				Code:           422,
				HttpStatusCode: http.StatusUnprocessableEntity,
				Message:        "invalid tenantID. You can locate your tenantID by following the instructions listed here: https://learn.microsoft.com/partner-center/find-ids-and-domain-names",
			},
			expectContains: []string{
				"invalid tenantID",
				"https://learn.microsoft.com/partner-center/find-ids-and-domain-names",
				"(code: 422, HTTP status: 422)",
			},
		},
		{
			name: "API error with hint surfaces both message and hint",
			err: &Error{
				Code:           422,
				HttpStatusCode: http.StatusUnprocessableEntity,
				Message:        "invalid subscriptionId provided",
				Hint:           "Please check your Azure subscription ID",
			},
			expectContains: []string{
				"invalid subscriptionId provided",
				"Hint: Please check your Azure subscription ID",
			},
		},
		{
			name: "Wrapped API error is still surfaced as readable text",
			err: fmt.Errorf("received error: %w", &Error{
				Code:           422,
				HttpStatusCode: http.StatusUnprocessableEntity,
				Message:        "invalid tenantID",
			}),
			expectContains: []string{
				"invalid tenantID",
				"(code: 422, HTTP status: 422)",
			},
		},
		{
			name: "Non-API error falls back to err.Error()",
			err:  internalerrors.ErrClusterIdCannotBeEmpty,
			expectContains: []string{
				internalerrors.ErrClusterIdCannotBeEmpty.Error(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ParseReadableError(test.err)

			for _, s := range test.expectContains {
				assert.Contains(t, result, s)
			}
		})
	}
}
