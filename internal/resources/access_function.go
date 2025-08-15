package resources

import (
	"context"
	"fmt"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AccessFunction{}
	_ resource.ResourceWithConfigure   = &AccessFunction{}
	_ resource.ResourceWithImportState = &AccessFunction{}
)

const errorMessageAfterAccessFunctionCreation = "Access function creation is successful, but encountered an error while checking the current" +
	" state of the access function. Please run `terraform plan` after 1-2 minutes to know the" +
	" current access function state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileAccessFunctionCreation = "There is an error during access function creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// AccessFunction is the AccessFunction resource implementation.
type AccessFunction struct {
	*providerschema.Data
}

func NewAccessFunction() resource.Resource {
	return &AccessFunction{}
}

// Metadata returns the access function resource type name.
func (r *AccessFunction) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_function"
}

// Schema defines the schema for the access function resource.
func (r *AccessFunction) Schema(ctx context.Context, rsc resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = AccessFunctionSchema()
}

// Configure sets provider-defined data, clients, etc. that is passed to data sources or resources in the provider.
func (r *AccessFunction) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new access function.
func (r *AccessFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.AccessFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104559) - Implement access function creation logic.
}

// Read refreshes the Terraform state with the latest data.
func (r *AccessFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.AccessFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104559) - Implement access function read logic.
}

// Update updates the access function.
func (r *AccessFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan providerschema.AccessFunction
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104559) - Implement access function update logic.
}

// Delete deletes the access function.
func (r *AccessFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state providerschema.AccessFunction
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO (AV-104559) - Implement access function deletion logic.
}

// ImportState imports a resource into Terraform state.
func (r *AccessFunction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO (AV-104559) - Implement access function import logic.
}
