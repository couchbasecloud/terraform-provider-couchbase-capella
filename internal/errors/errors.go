package errors

import "errors"

var (
	ErrProjectIdCannotBeEmpty      = errors.New("project ID cannot be empty, please contact Couchbase Capella Support")
	ErrOrganizationIdCannotBeEmpty = errors.New("organization ID cannot be empty, please contact Couchbase Capella Support")
)
