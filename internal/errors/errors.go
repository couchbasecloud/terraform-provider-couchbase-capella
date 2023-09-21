package errors

import "errors"

var (
	ErrIdMissing                   = errors.New("some ID is missing or was passed incorrectly, please check provider documentation for syntax")
	ErrProjectIdCannotBeEmpty      = errors.New("project ID cannot be empty, please contact Couchbase Capella Support")
	ErrProjectIdMissing            = errors.New("project ID is missing or was passed incorrectly, please check provider documentation for syntax")
	ErrOrganizationIdCannotBeEmpty = errors.New("organization ID cannot be empty, please contact Couchbase Capella Support")
	ErrOrganizationIdMissing       = errors.New("organization ID is missing or was passed incorrectly, please check provider documentation for syntax")
)
