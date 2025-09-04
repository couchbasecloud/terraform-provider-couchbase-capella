package resources

import (
	"context"
	"fmt"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ImportFilter{}
	_ resource.ResourceWithConfigure   = &ImportFilter{}
	_ resource.ResourceWithImportState = &ImportFilter{}
)

// ImportFilter is the resource implementation for managing App Endpoint Import Filters.
type ImportFilter struct {
	*providerschema.Data
}

// NewImportFilter creates a new ImportFilter resource for provider initialization.
func NewImportFilter() resource.Resource {
	return &ImportFilter{}
}

// Metadata returns the resource type name.
func (r *ImportFilter) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_import_filter"
}

// Schema defines the Terraform schema for this resource.
func (r *ImportFilter) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ImportFilterSchema()
}

// Configure sets provider-defined data, clients, etc.
func (r *ImportFilter) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.Data = data
}

// Create upserts the import filter.
func (r *ImportFilter) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ImportFilter
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104546) - Implement API call to putImportFilter and set state.
}

// Read fetches the current import filter and refreshes state.
func (r *ImportFilter) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ImportFilter
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104546) - Implement API call to getImportFilter and update state.
}

// Update updates the import filter.
func (r *ImportFilter) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.ImportFilter
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104546) - Implement API call to putImportFilter and update state.
}

// Delete deletes the import filter (resets/removes).
func (r *ImportFilter) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.ImportFilter
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104546) - Implement API call to deleteImportFilter.
}

// ImportState imports a resource into Terraform state.
func (r *ImportFilter) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO (AV-104546) - Implement import logic similar to other resources (comma-separated IDs).
}
