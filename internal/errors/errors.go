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

	// ErrPeerIdMissing is returned when an expected Peer Id was not found after an import.
	ErrPeerIdMissing = errors.New("peer ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrAzureTenantIdMissing is returned when an expected Azure Tenant Id was not found after an import.
	ErrAzureTenantIdMissing = errors.New("azure Tenant ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrSubscriptionIdMissing is returned when an expected Azure Subscription Id was not found after an import.
	ErrSubscriptionIdMissing = errors.New("azure Subscription ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrVNetIdMissing is returned when an expected Azure Vnet Id was not found after an import.
	ErrVNetIdMissing = errors.New("azure Vnet name is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrResourceGroup is returned when an expected Azure Vnet Resource group was not found after an import.
	ErrResourceGroup = errors.New("azure Vnet Resource group is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrVnetPeeringServicePrincipal is returned when an expected Azure Vnet Service Principal or object id was not found after an import.
	ErrVnetPeeringServicePrincipal = errors.New("azure object id or the Azure Vnet Service Principal is missing or was passed incorrectly, please check provider documentation for syntax")

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

	// ErrUnableToUpdateCloudProvider is returned when it is not possible to update the cloud provider.
	ErrUnableToUpdateCloudProvider = errors.New("unable to update cloud provider")

	// ErrMarshallingPayload is returned when a payload has failed to marshal into a request body.
	ErrMarshallingPayload = errors.New("failed to marshal payload")

	// ErrUnmarshallingResponse is returned when a HTTP response failed to unmarshal.
	ErrUnmarshallingResponse = errors.New("failed to unmarshal response")

	// ErrUnmarshallingAWSConfigResponse is returned when a HTTP response failed to unmarshal.
	ErrUnmarshallingAWSConfigResponse = errors.New("failed to unmarshal aws config response")

	// ErrUnmarshallingGCPConfigResponse is returned when a HTTP response failed to unmarshal.
	ErrUnmarshallingGCPConfigResponse = errors.New("failed to unmarshal GCP config response")

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

	// ErrConvertingZone is returned when terraform fails to convert a Zone.
	ErrConvertingZone = errors.New("failed to convert Zones, please contact Couchbase Capella Support")

	// ErrReadingAWSDisk is returned when an AWS disk read fails.
	ErrReadingAWSDisk = errors.New("failed to read AWS disk, please contact Couchbase Capella Support")

	// ErrReadingAzureDisk is returned when an Azure disk read fails.
	ErrReadingAzureDisk = errors.New("failed to read Azure disk, please contact Couchbase Capella Support")

	// ErrReadingGCPDisk is returned when a GCP disk read fails.
	ErrReadingGCPDisk = errors.New("failed to read GCP disk, please contact Couchbase Capella Support")

	// ErrReadingAWSConfig is returned when an AWS disk read fails.
	ErrReadingAWSConfig = errors.New("failed to read AWS config, please contact Couchbase Capella Support")

	// ErrReadingGCPConfig is returned when a GCP disk read fails.
	ErrReadingGCPConfig = errors.New("failed to read GCP config, please contact Couchbase Capella Support")

	// ErrReadingAzureConfig is returned when a GCP disk read fails.
	ErrReadingAzureConfig = errors.New("failed to read Azure config, please contact Couchbase Capella Support")

	// ErrReadingProviderConfig is returned when one or more of the fields in the provider config of the csp is missing.
	ErrReadingProviderConfig = errors.New("failed to read the provider config as one or more of the fields in the config is missing, please contact Couchbase Capella Support")

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

	// ErrRatelimit is returned when the Capella API reaches a rate limit for the same API key.
	ErrRatelimit = errors.New("api key reached the ratelimit")

	// ErrScopeNameMissing is returned when an expected ScopeName was not found after an import.
	ErrScopeNameMissing = errors.New("scope Name is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrInvalidSampleBucketName is returned when sample bucket name is not valid.
	ErrInvalidSampleBucketName = errors.New("sample bucket name can only be travel-sample, beer-sample, gamesim-sample")

	// ErrOnoffStateCannotBeEmpty is returned when cluster on/off state is required for a request but was not included.
	ErrOnoffStateCannotBeEmpty = errors.New("on/off state cannot be empty, please mention the state in which you want your cluster to be")

	// ErrEndpointIdMissing is returned when an expected endpoint ID was not found after an import.
	ErrEndpointIdMissing = errors.New("endpoint ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrVPCIDMissing is returned when an expected AWS VPC ID was not found after an import.
	ErrVPCIDMissing = errors.New("AWS VPC ID is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrVirtualNetworkMissing is returned when an expected Azure virtual network was not found after an import.
	ErrVirtualNetworkMissing = errors.New("Azure virtual network is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrResourceGroupName is returned when an expected Azure resource group was not found after an import.
	ErrResourceGroupName = errors.New("Azure resource group is missing or was passed incorrectly, please check provider documentation for syntax")

	// ErrConvertingProviderConfig is returned when terraform fails to convert a network peer provider config.
	ErrConvertingProviderConfig = errors.New("failed to convert network peer provider config, please contact Couchbase Capella Support")

	// ErrProviderConfigCannotBeEmpty is returned when the provider_config was required for a request but was empty.
	ErrProviderConfigCannotBeEmpty = errors.New("provider_config cannot be empty, it should be populated with one of- aws_config, gcp_config or azure_config. Please contact Couchbase Capella Support")

	ErrPrivateEndpointServiceTimeout = errors.New("changing private endpoint service status timed out after initiation")

	ErrBucketCreationStatusTimeout = errors.New("bucket backup creation status transition timed out after initiation")

	ErrAppServiceCreationStatusTimeout = errors.New("app service creation status transition timed out after initiation")

	ErrMonitorTimeout = errors.New("timed out while watching indexes")

	ErrConcurrentIndexCreation = errors.New("another index create request is in progress")

	ErrorMessageWhileFreeTierBucketCreation = errors.New("There is an error during free tier bucket creation. Please check in Capella to see if any hanging resources")

	ErrorMessageAfterFreeTierClusterCreationInitiation = errors.New("Cluster creation is initiated, but encountered an error while checking the current" +
		" state of the cluster. Please run `terraform plan` after 4-5 minutes to know the" +
		" current status of the cluster. Additionally, run `terraform apply --refresh-only` to update" +
		" the state from remote, unexpected error: ")

	ErrorMessageWhileFreeTierClusterCreation = errors.New("There is an error during cluster creation. Please check in Capella to see if any hanging resources" +
		" have been created, unexpected error: ")
	ErrUnableToConvertAppServiceCompute = errors.New("failed to convert app service compute")

	ErrFreeTierCreateAppServiceError   = errors.New("There is an error during app service creation. Please check in Capella to see if any hanging resources" + " have been created, unexpected error: ")
	ErrFreeTierAppServiceAfterCreation = errors.New("App Service creation is initiated, but encountered an error while checking the current" +
		" state of the app service. Please run `terraform plan` after 4-5 minutes to know the" +
		" current status of the app service. Additionally, run `terraform apply --refresh-only` to update" +
		" the state from remote, unexpected error: ")
)
