package errors

import "errors"

var (
	ErrIdMissing = errors.New("some ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrAllowListIdCannotBeEmpty = errors.New("allowlist ID cannot be empty, please contact Couchbase Capella Support")
	ErrAllowListIdMissing       = errors.New("allowList ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrClusterIdCannotBeEmpty = errors.New("cluster ID cannot be empty, please contact Couchbase Capella Support")
	ErrClusterIdMissing       = errors.New("cluster ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrProjectIdCannotBeEmpty = errors.New("project ID cannot be empty, please contact Couchbase Capella Support")
	ErrProjectIdMissing       = errors.New("project ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrOrganizationIdCannotBeEmpty = errors.New("organization ID cannot be empty, please contact Couchbase Capella Support")
	ErrOrganizationIdMissing       = errors.New("organization ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrDatabaseCredentialIdCannotBeEmpty = errors.New("database credential ID cannot be empty, please contact Couchbase Capella Support")
	ErrDatabaseCredentialIdMissing       = errors.New("database credential ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrUserIdCannotBeEmpty = errors.New("user ID cannot be empty, please contact Couchbase Capella Support")
	ErrUserIdMissing       = errors.New("user ID is missing or was passed incorrectly, please check provider documentation for syntax")
)
