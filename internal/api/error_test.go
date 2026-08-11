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

func Test_IsTransientInternalServerError(t *testing.T) {
	transient500 := Error{
		HttpStatusCode: http.StatusInternalServerError,
		Code:           transientInternalServerErrorCode,
		Message:        "An internal server error occurred.",
	}

	type test struct {
		err     error
		name    string
		expBool bool
	}

	tests := []test{
		{
			name:    "Transient 500 received from Capella V4 Api",
			err:     &transient500,
			expBool: true,
		},
		{
			name:    "Wrapped transient 500 received from Capella V4 Api",
			err:     fmt.Errorf("received error: %w", &transient500),
			expBool: true,
		},
		{
			name: "500 with a different error code is not transient",
			err: &Error{
				HttpStatusCode: http.StatusInternalServerError,
				Code:           9999,
			},
			expBool: false,
		},
		{
			name: "Forbidden error carrying the transient code is not transient",
			err: &Error{
				HttpStatusCode: http.StatusForbidden,
				Code:           transientInternalServerErrorCode,
			},
			expBool: false,
		},
		{
			name:    "Error other than received from Capella V4 Api",
			err:     internalerrors.ErrAllowListIdCannotBeEmpty,
			expBool: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expBool, IsTransientInternalServerError(test.err))
		})
	}
}
