package resources

import (
	"encoding/json"
	"testing"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCreateDatabaseCredentialUpdateRequestPassword(t *testing.T) {
	tests := []struct {
		name             string
		plannedPassword  types.String
		previousPassword types.String
		wantPassword     string
		wantPasswordKey  bool
	}{
		{
			name:             "unchanged password is omitted",
			plannedPassword:  types.StringValue("existing-password"),
			previousPassword: types.StringValue("existing-password"),
		},
		{
			name:             "changed password is included",
			plannedPassword:  types.StringValue("new-password"),
			previousPassword: types.StringValue("existing-password"),
			wantPassword:     "new-password",
			wantPasswordKey:  true,
		},
		{
			name:             "unknown password is omitted",
			plannedPassword:  types.StringUnknown(),
			previousPassword: types.StringValue("existing-password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createDatabaseCredentialUpdateRequest(
				providerschema.DatabaseCredential{Password: tt.plannedPassword},
				providerschema.DatabaseCredential{Password: tt.previousPassword},
			)

			require.Equal(t, tt.wantPassword, request.Password)

			payload, err := json.Marshal(request)
			require.NoError(t, err)

			var fields map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(payload, &fields))
			_, hasPassword := fields["password"]
			require.Equal(t, tt.wantPasswordKey, hasPassword)
		})
	}
}
