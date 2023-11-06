package api

import (
	"terraform-provider-capella/internal/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseError(t *testing.T) {
	var (
		errMessage = "Error received from Capella V4 Api"
	)

	type test struct {
		name      string
		err       error
		expString string
	}

	tests := []test{
		{
			name:      "Error received from Capella V4 Api",
			err:       Error{Message: errMessage},
			expString: errMessage,
		},
		{
			name:      "Error other than received from Capella V4 Api",
			err:       errors.ErrAllowListIdCannotBeEmpty,
			expString: errors.ErrAllowListIdCannotBeEmpty.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			errString := ParseError(test.err)

			assert.Equal(t, errString, test.expString)
		})
	}
}

func Test_CheckResourceNotFound(t *testing.T) {
	var (
		errMessage = "example error message"
	)

	type test struct {
		name      string
		err       error
		expString string
		expBool   bool
	}

	tests := []test{
		{
			name: "Error received from Capella V4 Api",
			err: Error{
				HttpStatusCode: 500,
				Message:        errMessage,
			},
			expString: errMessage,
			expBool:   false,
		},
		{
			name: "Error received from Capella V4 Api - Resource Not Found",
			err: Error{
				HttpStatusCode: 404,
				Message:        errMessage,
			},
			expString: errMessage,
			expBool:   true,
		},
		{
			name:      "Error other than received from Capella V4 Api",
			err:       errors.ErrAllowListIdCannotBeEmpty,
			expString: errors.ErrAllowListIdCannotBeEmpty.Error(),
			expBool:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resourceNotFound, errString := CheckResourceNotFoundError(test.err)

			assert.Equal(t, errString, test.expString)
			assert.Equal(t, resourceNotFound, test.expBool)
		})
	}
}
