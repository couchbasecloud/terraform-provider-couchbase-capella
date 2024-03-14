package errors

import "errors"

var (
	// ErrIdMissing is returned when an expected Id was not found after an import.
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

	// ErrUnmarshallingResponse is returned when a HTTP response failrf to unmarshal.
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

	// ErrUnableToReadCapellaUser is returned when the provider failed to read a requested Capella user.
	ErrUnableToReadCapellaUser = errors.New("could not read Capella user, please contact Couchbase Capella Support")

	// ErrApiKeyIdCannotBeEmpty is returned when an ApiKeyId was required for a request but was not included.
	ErrApiKeyIdCannotBeEmpty = errors.New("api key ID cannot be empty, please contact Couchbase Capella Support")

	// ErrApiKeyIdMissing is returned when an expected ErrApiKeyIdMissing was not found after an import.
	ErrApiKeyIdMissing = errors.New("api key ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrBucketIdCannotBeEmpty is returned when an BucketId was required for a request but was not included.
	ErrBucketIdCannotBeEmpty = errors.New("bucket ID cannot be empty, please contact Couchbase Capella Support")

	// ErrInvalidImport is returned when when the IDs supplied to terraform import did not match those expected for the resource.
	ErrInvalidImport = errors.New("terraform import parameters did not match those expected for the resource, please check provider documentation for syntax")

	// ErrValidatingResource is returned when the validation of an existing resource state failed.
	ErrValidatingResource = errors.New("could not validate resource state, please contact Couchbase Capella Support")

	// ErrNotFound is returned when a resource is requested but not found.
	ErrNotFound = errors.New("A resource was requested but could not be found")

	// ErrAppServiceIdCannotBeEmpty is returned when an AppServiceId was required for a request but was not included.
	ErrAppServiceIdCannotBeEmpty = errors.New("App Service ID cannot be empty, please contact Couchbase Capella Support")

	// ErrUnableToUpdateAppServiceName is returned when an app service name was updated.
	ErrUnableToUpdateAppServiceName = errors.New("app service name cannot be updated")

	// ErrRefreshingState is returned when a resource state failed to refresh.
	ErrRefreshingState = errors.New("failed to refresh the state of a resource, please contact Couchbase Capella Support")

	// ErrConvertingServiceGroups is returned when terraform fails to convert a clusters service groups.
	ErrConvertingServiceGroups = errors.New("failed to convert cluster service groups, please contact Couchbase Capella Support")

	// ErrGcpIopsCannotBeSet is returned when iops is set for GCP cluster.
	ErrGcpIopsCannotBeSet = errors.New("iops for gcp cluster cannot be set")

	// ErrConvertingCidr is returned when terraform fails to convert a CIDR.
	ErrConvertingCidr = errors.New("failed to convert CIDR, please contact Couchbase Capella Support")

	// ErrReadingAWSDisk is returned when an AWS disk read fails.
	ErrReadingAWSDisk = errors.New("failed to read AWS disk, please contact Couchbase Capella Support")

	// ErrReadingAzureDisk is returned when an Azure disk read fails.
	ErrReadingAzureDisk = errors.New("failed to read Azure disk, please contact Couchbase Capella Support")

	// ErrReadingGCPDisk is returned when a GCP disk read fails.
	ErrReadingGCPDisk = errors.New("failed to read GCP disk, please contact Couchbase Capella Support")
	// ErrBucketIdMissing is returned when an expected Bucket Id was not found after an import.
	ErrBucketIdMissing = errors.New("bucket ID is missing or was passed incorrectly, please check provider documentation for syntax")

	ErrRestoreTimesMustNotBeSetWhileCreateBackup = errors.New("restore times must not be set while create backup")

	// ErrTFVarHostIsNotSet is returned when TF_VAR_host is not set.
	ErrTFVarHostIsNotSet = errors.New("TF_VAR_host is not set")

	// ErrTFVARAuthTokenIsNotSet is returned when TF_VAR_auth_token is not set.
	ErrTFVARAuthTokenIsNotSet = errors.New("TF_VAR_auth_token is not set")

	// ErrTFVAROrganizationIdIsNotSet is returned when TF_VAR_organization_id is not set.
	ErrTFVAROrganizationIdIsNotSet = errors.New("TF_VAR_organization_id is not set")

	// ErrClusterCreationTimeoutAfterInitiation is returned when cluster creation
	// is timeout after initiation.
	ErrClusterCreationTimeoutAfterInitiation = errors.New("cluster creation status transition timed out after initiation")

	// ErrGatewayTimeout is returned when a gateway operation times out.
	ErrGatewayTimeout = errors.New("gateway timeout")

	// ErrNotTrimmed is returned when any attribute has leading or trailing spaces.
	ErrNotTrimmed = errors.New("attribute has leading or trailing spaces")

	// ErrIfMatchCannotBeSetWhileCreate is returned when if_match is set during create operation.
	ErrIfMatchCannotBeSetWhileCreate = errors.New("if_match attribute cannot be set during create operation")

	// ErrIfMatchCannotBeSetWhileCreate is returned when if_match is set during create operation.
	ErrInvalidSampleBucketName = errors.New("sample bucket name can only be travel-sample, beer-sample, gamesim-sample")
)
