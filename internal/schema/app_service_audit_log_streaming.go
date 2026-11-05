package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppServiceAuditLogStreaming defines the Terraform state for the app service audit log streaming resource.
type AppServiceAuditLogStreaming struct {
	// OrganizationId is the ID of the organization to which the Capella cluster belongs.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the ID of the project to which the Capella cluster belongs.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the ID of the cluster for which the app service is deployed.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the ID of the app service for which audit log streaming is configured.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// StreamingEnabled determines whether audit log streaming is enabled.
	StreamingEnabled types.Bool `tfsdk:"streaming_enabled"`

	// OutputType is the log collector type.
	OutputType types.String `tfsdk:"output_type"`

	// Credentials contains the credentials for the configured log collector. The API does not return them.
	Credentials types.Object `tfsdk:"credentials"`

	// DisabledAppEndpoints is the list of App Endpoints excluded from audit log streaming.
	DisabledAppEndpoints types.Set `tfsdk:"disabled_app_endpoints"`

	// LogStreamingState is the current state of audit log streaming.
	LogStreamingState types.String `tfsdk:"log_streaming_state"`
}

// Validate validates the state and returns parsed IDs. It handles both normal reads and terraform import.
func (a *AppServiceAuditLogStreaming) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId: a.OrganizationId,
		ProjectId:      a.ProjectId,
		ClusterId:      a.ClusterId,
		AppServiceId:   a.AppServiceId,
	}

	IDs, err := validateSchemaState(state, AppServiceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

// AsLogStreamingCredentials converts the Credentials types.Object into a LogStreamingCredentials struct.
// Returns nil if the object is null or unknown.
func (a *AppServiceAuditLogStreaming) AsLogStreamingCredentials(ctx context.Context) (*LogStreamingCredentials, diag.Diagnostics) {
	if a.Credentials.IsNull() || a.Credentials.IsUnknown() {
		return nil, nil
	}
	var creds LogStreamingCredentials
	diags := a.Credentials.As(ctx, &creds, basetypes.ObjectAsOptions{})
	return &creds, diags
}
