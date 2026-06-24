package resources

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

func Test_validateURLBindingAuth(t *testing.T) {
	authPath := path.Root("bindings").AtName("urls").AtListIndex(0).AtName("authentication")

	tests := []struct {
		name           string
		auth           providerschema.EventingFunctionURLBindingAuthentication
		wantErrorCount int
		wantSubstr     string
	}{
		{
			name: "none with no credentials is valid",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthNone),
				Username:    types.StringNull(),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
		},
		{
			name: "none with a username is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthNone),
				Username:    types.StringValue("user"),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 1,
			wantSubstr:     "must not set username, password or bearer_token",
		},
		{
			name: "none with an unknown bearer token is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthNone),
				Username:    types.StringNull(),
				Password:    types.StringNull(),
				BearerToken: types.StringUnknown(),
			},
			wantErrorCount: 1,
		},
		{
			name: "basic with username and password is valid",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("user"),
				Password:    types.StringValue("pass"),
				BearerToken: types.StringNull(),
			},
		},
		{
			name: "basic with username but no password is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("userOnly"),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 1,
			wantSubstr:     `type "basic" requires username and password`,
		},
		{
			name: "basic with an empty password is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("user"),
				Password:    types.StringValue(""),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 1,
		},
		{
			name: "basic with an unknown password defers validation",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("user"),
				Password:    types.StringUnknown(),
				BearerToken: types.StringNull(),
			},
		},
		{
			name: "basic with a bearer token is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("user"),
				Password:    types.StringValue("pass"),
				BearerToken: types.StringValue("token"),
			},
			wantErrorCount: 1,
			wantSubstr:     `type "basic" must not set bearer_token`,
		},
		{
			name: "basic missing password and setting a bearer token yields two errors",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBasic),
				Username:    types.StringValue("user"),
				Password:    types.StringNull(),
				BearerToken: types.StringValue("token"),
			},
			wantErrorCount: 2,
		},
		{
			name: "digest with username and password is valid",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthDigest),
				Username:    types.StringValue("user"),
				Password:    types.StringValue("pass"),
				BearerToken: types.StringNull(),
			},
		},
		{
			name: "digest with no password is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthDigest),
				Username:    types.StringValue("user"),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 1,
			wantSubstr:     `type "digest" requires username and password`,
		},
		{
			name: "bearer with a token is valid",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBearer),
				Username:    types.StringNull(),
				Password:    types.StringNull(),
				BearerToken: types.StringValue("token"),
			},
		},
		{
			name: "bearer with no token is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBearer),
				Username:    types.StringNull(),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 1,
			wantSubstr:     `type "bearer" requires bearer_token`,
		},
		{
			name: "bearer with a username is rejected",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBearer),
				Username:    types.StringValue("user"),
				Password:    types.StringNull(),
				BearerToken: types.StringValue("token"),
			},
			wantErrorCount: 1,
			wantSubstr:     `type "bearer" must not set username or password`,
		},
		{
			name: "bearer missing token and setting a username yields two errors",
			auth: providerschema.EventingFunctionURLBindingAuthentication{
				Type:        types.StringValue(eventingURLAuthBearer),
				Username:    types.StringValue("user"),
				Password:    types.StringNull(),
				BearerToken: types.StringNull(),
			},
			wantErrorCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &resource.ValidateConfigResponse{}
			validateURLBindingAuth(tt.auth, authPath, resp)

			assert.Len(t, resp.Diagnostics, tt.wantErrorCount)
			if tt.wantSubstr != "" {
				var matched bool
				for _, d := range resp.Diagnostics {
					if strings.Contains(d.Detail(), tt.wantSubstr) {
						matched = true
						break
					}
				}
				assert.Truef(t, matched, "expected a diagnostic containing %q, got %v", tt.wantSubstr, resp.Diagnostics)
			}
		})
	}
}
