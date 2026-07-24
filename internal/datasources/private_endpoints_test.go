package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

// TestPrivateEndpointsSchemaMatchesStructNonEmpty guards against a struct/schema
// mismatch on the couchbase-capella_private_endpoints data source. The nested
// "data" attributes must line up exactly with the PrivateEndpointData struct, or
// terraform-plugin-framework fails state conversion the moment the endpoint list
// is non-empty (the empty-list path never exercises the mismatch, which is why
// the acceptance test that only checks data.# = 0 did not catch it).
func TestPrivateEndpointsSchemaMatchesStructNonEmpty(t *testing.T) {
	state := &tfsdk.State{Schema: PrivateEndpointsSchema()}

	pe := providerschema.PrivateEndpoints{
		OrganizationId:     types.StringValue("00000000-0000-0000-0000-000000000000"),
		ProjectId:          types.StringValue("11111111-1111-1111-1111-111111111111"),
		ClusterId:          types.StringValue("22222222-2222-2222-2222-222222222222"),
		PrivateEndpointDNS: types.StringValue("cb.private.example.com"),
		Data: []providerschema.PrivateEndpointData{
			{
				Id:          types.StringValue("vpce-1"),
				Status:      types.StringValue("pendingAcceptance"),
				ServiceName: types.StringValue("com.amazonaws.vpce.us-east-1.vpce-svc-1"),
			},
			{
				Id:          types.StringValue("vpce-2"),
				Status:      types.StringValue("linked"),
				ServiceName: types.StringValue("com.amazonaws.vpce.us-east-1.vpce-svc-2"),
			},
		},
	}

	diags := state.Set(context.Background(), &pe)
	if diags.HasError() {
		t.Fatalf("setting a non-empty private endpoints list must not error, got: %v", diags)
	}
}
