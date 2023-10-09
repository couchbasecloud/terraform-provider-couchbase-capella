package errors

import "errors"

var (
	//ErIdMissing is returned when an expected Id was not found after an import.
	ErrIdMissing = errors.New("some ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrUserIdCannotBeEmpty is returned when a User Id was required for a request but was not included.
	ErrUserIdCannotBeEmpty = errors.New("user ID cannot be empty, please contact Couchbase Capella Support")
	// ErrUserIdMissing is returned when an expected User Id was not found after an import.
	ErrUserIdMissing = errors.New("user ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrAllowListIdCannotBeEmpty is returned when an AllowList Id was required for a request but was not included.
	ErrAllowListIdCannotBeEmpty = errors.New("allowlist ID cannot be empty, please contact Couchbase Capella Support")
	// ErrAllowListIdMissing is returned when an expected AllowList Id was not found after an import.
	ErrAllowListIdMissing = errors.New("allowList ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrClusterIdCannotBeEmpty is returned when a Cluster Id was required for a request but was not included.
	ErrClusterIdCannotBeEmpty = errors.New("cluster ID cannot be empty, please contact Couchbase Capella Support")
	// ErrClusterIdMissing is returned when an expected Cluster Id was not found after an import.
	ErrClusterIdMissing = errors.New("cluster ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrProjectIdCannotBeEmpty is returned when a Project Id was required for a request but was not included.
	ErrProjectIdCannotBeEmpty = errors.New("project ID cannot be empty, please contact Couchbase Capella Support")
	// ErrProjectIdMissing is returned when an expected Project Id was not found after an import.
	ErrProjectIdMissing = errors.New("project ID is missing or was passed incorrectly, please check provider documentation for syntax")
	// ErrUnableToUpdateProjectId is returned when an update to a projectId was unsuccessful.
	ErrUnableToUpdateProjectId = errors.New("unable to update projectId")

	// ErrOrganizationIdCannotBeEmpty is returned when an Organization Id was required for a request but was not included.
	ErrOrganizationIdCannotBeEmpty = errors.New("organization ID cannot be empty, please contact Couchbase Capella Support")
	// ErrOrganizationIdMissing is returned when an expected Organization Id was not found after an import.
	ErrOrganizationIdMissing = errors.New("organization ID is missing or was passed incorrectly, please check provider documentation for syntax")
	// ErrUnableToUpdateOrganizationId is returned when an update to a projectId was unsuccessful.
	ErrUnableToUpdateOrganizationId = errors.New("unable to update organizationId")

	// ErrDatabaseCredentialIdCannotBeEmpty is returned when a Database Credential Id was required for a request but was not included.
	ErrDatabaseCredentialIdCannotBeEmpty = errors.New("database credential ID cannot be empty, please contact Couchbase Capella Support")
	// ErrDatabaseCredentialIdMissing is returned when an expected DatabaseCredential Id was not found after an import.
	ErrDatabaseCredentialIdMissing = errors.New("database credential ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrEmailCannotBeEmpty is returned when an email address was required for a request but was not included.
	ErrEmailCannotBeEmpty = errors.New("email cannot be empty, please contact Couchbase Capella Support")

	// ErrOrganizationRolesCannotBeEmpty is returned when organization roles were required for a request but were not included.
	ErrOrganizationRolesCannotBeEmpty = errors.New("organization roles cannot be empty, please contact Couchbase Capella Support")

	// ErrUnableToUpdateServerVersion is returned when it is not possible to update the couchbase server version.
	ErrUnableToUpdateServerVersion = errors.New("unable to update couchbase server version")

	// ErrUnableToUpdateAvailabilityType is returned when it is not possible to update the availability type.
	ErrUnableToUpdateAvailabilityType = errors.New("unable to update availability type")

	// ErrUnableToUpdateCloudProvider is returned when when it is not possible to update the cloud provider.
	ErrUnableToUpdateCloudProvider = errors.New("unable to update cloud provider")

	// ErrMarshallingPayload is returned when a payload has failed to marshal into a request body.
	ErrMarshallingPayload = errors.New("failed to marshal payload")

	// ErrUnmarshallingPayload is returned when a HTTP response failrf to unmarshal
	ErrUnmarshallingResponse = errors.New("failed to unmarshal response")

	// ErrConstructingRequest is returned when a HTTP.NewRequest has failed.
	ErrConstructingRequest = errors.New("failed to construct request")

	// ErrExecutingRequest is returned when a HTTP request has failed to execute.
	ErrExecutingRequest = errors.New("failed to execute request")

	// ErrUnableToConvertAuditData is returned when an attempt to convert audit data from
	// terraform types.String to types string has failed.
	ErrUnableToConvertAuditData = errors.New("failed to convert audit data")

	// ErrUnableToImportResource is returned when a resource failed to be imported.
	ErrUnableToImportResource = errors.New("failed to import resource")

	// ErrUnsupportedCloudProvider is returned when an invalid cloud provider was requested.
	ErrUnsupportedCloudProvider = errors.New("cloud provider is not supported")

	// ErrUnableToReadCapellaUser is returned when the provider failed to read a requested Capella user
	ErrUnableToReadCapellaUser = errors.New("could not read Capella user, please contact Couchbase Capella Support")
)
