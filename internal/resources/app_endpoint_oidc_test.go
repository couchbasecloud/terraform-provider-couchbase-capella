package resources

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestAppEndpointOidcProviderCreateUnparseableResponse covers a 201 whose body cannot be
// decoded: the provider exists in Capella but its ID is unknown, so state stays empty and
// the diagnostic must tell the operator to remove it manually (AV-145497).
func TestAppEndpointOidcProviderCreateUnparseableResponse(t *testing.T) {
	ctx := context.Background()

	var posts int
	r := &AppEndpointOidcProvider{Data: newTestProviderData(t, func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method, "only the create POST is expected")
		posts++
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("not json"))
	})}

	oidcSchema := AppEndpointOidcProviderSchema()
	plan := tfsdk.Plan{Schema: oidcSchema}
	requireNoDiags(t, plan.Set(ctx, providerschema.AppEndpointOidcProvider{
		OrganizationId:  types.StringValue("org"),
		ProjectId:       types.StringValue("project"),
		ClusterId:       types.StringValue("cluster"),
		AppServiceId:    types.StringValue("app-service"),
		AppEndpointName: types.StringValue("my-endpoint"),
		Issuer:          types.StringValue("https://accounts.google.com"),
		ClientId:        types.StringValue("client"),
		DiscoveryUrl:    types.StringNull(),
		Register:        types.BoolNull(),
		RolesClaim:      types.StringNull(),
		UserPrefix:      types.StringNull(),
		UsernameClaim:   types.StringNull(),
		ProviderId:      types.StringUnknown(),
		IsDefault:       types.BoolUnknown(),
	}))

	resp := resource.CreateResponse{
		State: tfsdk.State{
			Schema: oidcSchema,
			Raw:    tftypes.NewValue(oidcSchema.Type().TerraformType(ctx), nil),
		},
	}
	r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)

	assert.Equal(t, 1, posts, "create should POST exactly once")
	require.True(t, resp.Diagnostics.HasError(), "an unparseable create response must be an error")
	require.Len(t, resp.Diagnostics.Errors(), 1)

	detail := resp.Diagnostics.Errors()[0].Detail()
	assert.Contains(t, detail, "my-endpoint")
	assert.Contains(t, detail, "not saved to Terraform state")
	assert.Contains(t, detail, "remove the OIDC provider manually")

	assert.True(t, resp.State.Raw.IsNull(), "state must not be set without a provider ID")
}
