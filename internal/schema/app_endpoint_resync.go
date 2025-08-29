package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/errors"
)

// AppEndpointResync defines the Terraform schema for an App Endpoint Resync resource.
// It includes both the request fields and additional metadata fields required for resource management.
type AppEndpointResync struct {
	// OrganizationId is the Capella tenant id associated with the App Service.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the id of the Capella project associated with the App Service.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the id of the Capella cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the id of the Capella App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpoint is the endpoint name for the App Service.
	AppEndpoint types.String `tfsdk:"app_endpoint"`

	// Scopes is a map where each key represents a scope name and the value
	// is a list of collection names within that scope to be resynced.
	// This field maps to the JSON request body structure.
	Scopes types.Map `tfsdk:"scopes"`

	// CollectionsProcessing contains a map of collections currently being processed,
	// organized by scope. This field is populated in the response.
	CollectionsProcessing types.Map `tfsdk:"collections_processing"`

	// DocsChanged represents the number of documents that have been changed
	// during the resync operation.
	DocsChanged types.Int64 `tfsdk:"docs_changed"`

	// DocsProcessed represents the total number of documents that have been
	// processed during the resync operation.
	DocsProcessed types.Int64 `tfsdk:"docs_processed"`

	// LastError contains the last error message encountered during the resync
	// operation, if any.
	LastError types.String `tfsdk:"last_error"`

	// StartTime represents the timestamp when the resync operation was initiated.
	StartTime types.String `tfsdk:"start_time"`

	// State indicates the current state of the resync operation
	// (e.g., "running", "completed", "error", etc.).
	State types.String `tfsdk:"state"`
}

func (a AppEndpointResync) Validate() (map[Attr]string, error) {
	state := map[Attr]basetypes.StringValue{
		OrganizationId:  a.OrganizationId,
		ProjectId:       a.ProjectId,
		ClusterId:       a.ClusterId,
		AppServiceId:    a.AppServiceId,
		AppEndpointName: a.AppEndpoint,
	}

	IDs, err := validateSchemaState(state, AppEndpointName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errors.ErrValidatingResource, err)
	}

	return IDs, nil
}

type AppEndpointResyncData struct {
	// OrganizationId is the Capella tenant id associated with the App Service.
	OrganizationId types.String `tfsdk:"organization_id"`

	// ProjectId is the id of the Capella project associated with the App Service.
	ProjectId types.String `tfsdk:"project_id"`

	// ClusterId is the id of the Capella cluster associated with the App Service.
	ClusterId types.String `tfsdk:"cluster_id"`

	// AppServiceId is the id of the Capella App Service.
	AppServiceId types.String `tfsdk:"app_service_id"`

	// AppEndpoint is the endpoint name for the App Service.
	AppEndpoint types.String `tfsdk:"app_endpoint"`

	// CollectionsProcessing contains a map of collections currently being processed,
	// organized by scope. This field is populated in the response.
	CollectionsProcessing types.Map `tfsdk:"collections_processing"`

	// DocsChanged represents the number of documents that have been changed
	// during the resync operation.
	DocsChanged types.Int64 `tfsdk:"docs_changed"`

	// DocsProcessed represents the total number of documents that have been
	// processed during the resync operation.
	DocsProcessed types.Int64 `tfsdk:"docs_processed"`

	// LastError contains the last error message encountered during the resync
	// operation, if any.
	LastError types.String `tfsdk:"last_error"`

	// StartTime represents the timestamp when the resync operation was initiated.
	StartTime types.String `tfsdk:"start_time"`

	// State indicates the current state of the resync operation
	// (e.g., "running", "completed", "error", etc.).
	State types.String `tfsdk:"state"`
}
