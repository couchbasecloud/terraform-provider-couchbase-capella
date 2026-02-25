package schema

import (
	"fmt"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// AppServiceLogStreamingBase contains the common fields shared by both the resource
// and datasource models for app service log streaming.
type AppServiceLogStreamingBase struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the app service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the app service for which log streaming is configured.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// OutputType is the log collector type (datadog, dynatrace, elastic, generic_http, loki, splunk, sumologic).
	OutputType types.String `tfsdk:"output_type"`

	// ConfigState is the current configuration state of log streaming (enabled, enabling, disabled, disabling, paused, pausing, errored).
	ConfigState types.String `tfsdk:"config_state"`

	// StreamingState indicates if logs are being successfully streamed from the App Service nodes (degraded, healthy, unhealthy, unknown, unsupported).
	StreamingState types.String `tfsdk:"streaming_state"`
}

// AppServiceLogStreaming defines the Terraform state for the app service log streaming resource.
type AppServiceLogStreaming struct {
	AppServiceLogStreamingBase

	// Credentials contains the credentials for the configured log collector.
	Credentials *LogStreamingCredentials `tfsdk:"credentials"`
}

// LogStreamingCredentials contains the credential configuration for log streaming.
// Only one of the nested credential types should be set based on the output_type.
type LogStreamingCredentials struct {
	// Datadog credentials for Datadog log collector.
	Datadog *DatadogCredentials `tfsdk:"datadog"`

	// Dynatrace credentials for Dynatrace log collector.
	Dynatrace *DynatraceCredentials `tfsdk:"dynatrace"`

	// Elastic credentials for Elasticsearch log collector.
	Elastic *ElasticCredentials `tfsdk:"elastic"`

	// GenericHttp credentials for generic HTTP log collector.
	GenericHttp *GenericHttpCredentials `tfsdk:"generic_http"`

	// Loki credentials for Grafana Loki log collector.
	Loki *LokiCredentials `tfsdk:"loki"`

	// Splunk credentials for Splunk log collector.
	Splunk *SplunkCredentials `tfsdk:"splunk"`

	// Sumologic credentials for SumoLogic log collector.
	Sumologic *SumologicCredentials `tfsdk:"sumologic"`
}

// DatadogCredentials contains credentials for Datadog log collector.
type DatadogCredentials struct {
	// ApiKey is the API key for authentication.
	ApiKey types.String `tfsdk:"api_key"`

	// Url is the DataDog log ingestion URL.
	Url types.String `tfsdk:"url"`
}

// DynatraceCredentials contains credentials for Dynatrace log collector.
type DynatraceCredentials struct {
	// ApiToken is the token for the Dynatrace log collector.
	ApiToken types.String `tfsdk:"api_token"`

	// Url is the URL for the Dynatrace log collector.
	Url types.String `tfsdk:"url"`
}

// ElasticCredentials contains credentials for Elasticsearch log collector.
type ElasticCredentials struct {
	// User is the username for the Elasticsearch log collector.
	User types.String `tfsdk:"user"`

	// Password is the password for the Elasticsearch log collector.
	Password types.String `tfsdk:"password"`

	// Url is the URL for the Elasticsearch log collector.
	Url types.String `tfsdk:"url"`
}

// GenericHttpCredentials contains credentials for generic HTTP log collector.
type GenericHttpCredentials struct {
	// User is the username for HTTP authentication (optional).
	User types.String `tfsdk:"user"`

	// Password is the password for HTTP authentication (optional).
	Password types.String `tfsdk:"password"`

	// Url is the URL for the generic HTTP log collector.
	Url types.String `tfsdk:"url"`
}

// LokiCredentials contains credentials for Grafana Loki log collector.
type LokiCredentials struct {
	// User is the username for the Grafana Loki log collector.
	User types.String `tfsdk:"user"`

	// Password is the password for the Grafana Loki log collector.
	Password types.String `tfsdk:"password"`

	// Url is the URL for the Grafana Loki log collector.
	Url types.String `tfsdk:"url"`
}

// SplunkCredentials contains credentials for Splunk log collector.
type SplunkCredentials struct {
	// SplunkToken is the token for the Splunk log collector.
	SplunkToken types.String `tfsdk:"splunk_token"`

	// Url is the URL for the Splunk log collector.
	Url types.String `tfsdk:"url"`
}

// SumologicCredentials contains credentials for SumoLogic log collector.
type SumologicCredentials struct {
	// Url is the SumoLogic signed URL for log ingestion.
	Url types.String `tfsdk:"url"`
}

// Validate validates the app service log streaming state and returns parsed IDs.
// It handles both normal reads and terraform import scenarios.
func (a *AppServiceLogStreamingBase) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// setFromAPIResponse populates the computed fields from an API response.
func (b *AppServiceLogStreamingBase) setFromAPIResponse(apiResponse *apigen.GetLogStreamingResponse) {
	b.OutputType = types.StringNull()
	if apiResponse.OutputType != nil {
		b.OutputType = types.StringValue(string(*apiResponse.OutputType))
	}

	b.ConfigState = types.StringNull()
	if apiResponse.ConfigState != nil {
		b.ConfigState = types.StringValue(string(*apiResponse.ConfigState))
	}

	b.StreamingState = types.StringNull()
	if apiResponse.StreamingState != nil {
		b.StreamingState = types.StringValue(string(*apiResponse.StreamingState))
	}
}

// NewAppServiceLogStreaming creates a new AppServiceLogStreaming from API response data.
func NewAppServiceLogStreaming(
	organizationId, projectId, clusterId, appServiceId string,
	apiResponse *apigen.GetLogStreamingResponse,
	existingCredentials *LogStreamingCredentials,
) *AppServiceLogStreaming {
	result := &AppServiceLogStreaming{
		AppServiceLogStreamingBase: AppServiceLogStreamingBase{
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			AppServiceId:   types.StringValue(appServiceId),
		},
		Credentials: existingCredentials, // Preserve credentials from plan/state since API doesn't return them
	}

	result.setFromAPIResponse(apiResponse)

	return result
}

// AppServiceLogStreamingData defines the Terraform state for the app service log streaming datasource.
// This is separate from AppServiceLogStreaming (the resource model) because the datasource
// does not include credentials.
type AppServiceLogStreamingData struct {
	AppServiceLogStreamingBase
}

// NewAppServiceLogStreamingData creates a new AppServiceLogStreamingData from API response data.
// This is the datasource constructor — it omits credentials since the API does not return them.
func NewAppServiceLogStreamingData(
	organizationId, projectId, clusterId, appServiceId string,
	apiResponse *apigen.GetLogStreamingResponse,
) *AppServiceLogStreamingData {
	result := &AppServiceLogStreamingData{
		AppServiceLogStreamingBase: AppServiceLogStreamingBase{
			OrganizationId: types.StringValue(organizationId),
			ProjectId:      types.StringValue(projectId),
			ClusterId:      types.StringValue(clusterId),
			AppServiceId:   types.StringValue(appServiceId),
		},
	}

	result.setFromAPIResponse(apiResponse)

	return result
}

// AttributeTypes returns the attribute types for LogStreamingCredentials.
func (c LogStreamingCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"datadog":      types.ObjectType{AttrTypes: DatadogCredentials{}.AttributeTypes()},
		"dynatrace":    types.ObjectType{AttrTypes: DynatraceCredentials{}.AttributeTypes()},
		"elastic":      types.ObjectType{AttrTypes: ElasticCredentials{}.AttributeTypes()},
		"generic_http": types.ObjectType{AttrTypes: GenericHttpCredentials{}.AttributeTypes()},
		"loki":         types.ObjectType{AttrTypes: LokiCredentials{}.AttributeTypes()},
		"splunk":       types.ObjectType{AttrTypes: SplunkCredentials{}.AttributeTypes()},
		"sumologic":    types.ObjectType{AttrTypes: SumologicCredentials{}.AttributeTypes()},
	}
}

// AttributeTypes returns the attribute types for DatadogCredentials.
func (c DatadogCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"api_key": types.StringType,
		"url":     types.StringType,
	}
}

// AttributeTypes returns the attribute types for DynatraceCredentials.
func (c DynatraceCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"api_token": types.StringType,
		"url":       types.StringType,
	}
}

// AttributeTypes returns the attribute types for ElasticCredentials.
func (c ElasticCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":     types.StringType,
		"password": types.StringType,
		"url":      types.StringType,
	}
}

// AttributeTypes returns the attribute types for GenericHttpCredentials.
func (c GenericHttpCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":     types.StringType,
		"password": types.StringType,
		"url":      types.StringType,
	}
}

// AttributeTypes returns the attribute types for LokiCredentials.
func (c LokiCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":     types.StringType,
		"password": types.StringType,
		"url":      types.StringType,
	}
}

// AttributeTypes returns the attribute types for SplunkCredentials.
func (c SplunkCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"splunk_token": types.StringType,
		"url":          types.StringType,
	}
}

// AttributeTypes returns the attribute types for SumologicCredentials.
func (c SumologicCredentials) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url": types.StringType,
	}
}
