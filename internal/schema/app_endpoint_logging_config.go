package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
)

type LoggingConfig struct {
	OrganizationId  types.String   `tfsdk:"organization_id"`
	ProjectId       types.String   `tfsdk:"project_id"`
	ClusterId       types.String   `tfsdk:"cluster_id"`
	AppServiceId    types.String   `tfsdk:"app_service_id"`
	AppEndpointName types.String   `tfsdk:"app_endpoint_name"`
	LogLevel        types.String   `tfsdk:"log_level"`
	LogKeys         []types.String `tfsdk:"log_keys"`
}

func (l LoggingConfig) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"organization_id":   types.StringType,
		"project_id":        types.StringType,
		"cluster_id":        types.StringType,
		"app_service_id":    types.StringType,
		"app_endpoint_name": types.StringType,
		"log_level":         types.StringType,
		"log_keys":          types.SetType{ElemType: types.StringType},
	}
}

func NewLoggingConfig(loggingConfig api.ConsoleLoggingConfig, organizationId, projectId, clusterId, appServiceId, appEndpointName string) *LoggingConfig {
	return &LoggingConfig{
		OrganizationId:  types.StringValue(organizationId),
		ProjectId:       types.StringValue(projectId),
		ClusterId:       types.StringValue(clusterId),
		AppServiceId:    types.StringValue(appServiceId),
		AppEndpointName: types.StringValue(appEndpointName),
		LogLevel:        types.StringValue(*loggingConfig.LogLevel),
		LogKeys:         StringsToBaseStrings(*loggingConfig.LogKeys),
	}
}

// Validate is used to verify that IDs have been properly imported.
func (l LoggingConfig) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  l.OrganizationId,
		ProjectId:       l.ProjectId,
		ClusterId:       l.ClusterId,
		AppServiceId:    l.AppServiceId,
		AppEndpointName: l.AppEndpointName,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}
	return IDs, nil
}
