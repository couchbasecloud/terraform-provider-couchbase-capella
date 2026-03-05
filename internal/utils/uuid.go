package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUID parses a string into a UUID, returning a descriptive error if parsing fails.
// The fieldName parameter is used to construct a human-readable error message indicating
// which field contained the invalid UUID value.
func ParseUUID(fieldName, value string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return parsed, nil
}

// ParseHierarchyUUIDs parses the four standard resource-hierarchy ID strings
// (organization, project, cluster, appService) into UUID values in a single call.
// It returns the first parse error encountered, with the field name included in
// the error message, making it a drop-in replacement for the repeated
// parseUUIDs / mapIDsToUUIDs helper methods that were duplicated across
// several resources and datasources.
func ParseHierarchyUUIDs(organizationId, projectId, clusterId, appServiceId string) (orgUUID, projectUUID, clusterUUID, appServiceUUID uuid.UUID, err error) {
	if orgUUID, err = ParseUUID("organization_id", organizationId); err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}
	if projectUUID, err = ParseUUID("project_id", projectId); err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}
	if clusterUUID, err = ParseUUID("cluster_id", clusterId); err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}
	if appServiceUUID, err = ParseUUID("app_service_id", appServiceId); err != nil {
		return uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, uuid.UUID{}, err
	}
	return orgUUID, projectUUID, clusterUUID, appServiceUUID, nil
}
