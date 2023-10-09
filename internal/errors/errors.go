package errors

import "errors"

var (
	ErrIdMissing = errors.New("some ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrUserIdCannotBeEmpty = errors.New("user ID cannot be empty, please contact Couchbase Capella Support")
	ErrUserIdMissing       = errors.New("user ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrAllowListIdCannotBeEmpty = errors.New("allowlist ID cannot be empty, please contact Couchbase Capella Support")
	ErrAllowListIdMissing       = errors.New("allowList ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrClusterIdCannotBeEmpty = errors.New("cluster ID cannot be empty, please contact Couchbase Capella Support")
	ErrClusterIdMissing       = errors.New("cluster ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrProjectIdCannotBeEmpty  = errors.New("project ID cannot be empty, please contact Couchbase Capella Support")
	ErrProjectIdMissing        = errors.New("project ID is missing or was passed incorrectly, please check provider documentation for syntax")
	ErrUnableToUpdateProjectId = errors.New("unable to update projectId")

	ErrOrganizationIdCannotBeEmpty = errors.New("organization ID cannot be empty, please contact Couchbase Capella Support")
	ErrOrganizationIdMissing       = errors.New("organization ID is missing or was passed incorrectly, please check provider documentation for syntax")
	ErrUnableToUpdateOrgId         = errors.New("unable to update organizationId")

	ErrDatabaseCredentialIdCannotBeEmpty = errors.New("database credential ID cannot be empty, please contact Couchbase Capella Support")
	ErrDatabaseCredentialIdMissing       = errors.New("database credential ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrEmailCannotBeEmpty = errors.New("email cannot be empty, please contact Couchbase Capella Support")

	ErrOrganizationRolesCannotBeEmpty = errors.New("organization roles cannot be empty, please contact Couchbase Capella Support")

	ErrUnableToUpdateServerVersion = errors.New("unable to update server version")

	ErrUnableToUpdateAvailabilityType = errors.New("unable to update availability type")

	ErrUnableToUpdateCloudProvider = errors.New("unable to update cloud provider")

	// ErrMarshallingPayload is returned when a payload has failed to
	// marshal into a request body.
	ErrMarshallingPayload = errors.New("failed to marshal payload")

	ErrUnmarshallingResponse = errors.New("failed to unmarshal response")

	// ErrConstructingRequest is returned when a HTTP.NewRequest has failed.
	ErrConstructingRequest = errors.New("failed to construct request")

	// ErrExecutingRequest is returned when a HTTP request has failed to execute.
	ErrExecutingRequest = errors.New("failed to execute request")

	ErrUnableToConvertAuditData = errors.New("failed to convert audit data")

	ErrUnableToImportResource = errors.New("failed to import resource")

	ErrUnsupportedCloudProvider = errors.New("cloud provider is not supported")

	ErrUnableToReadCapellaUser = errors.New("Could not read Capella user, please contact Couchbase Capella Support")
)
