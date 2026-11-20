package datasources

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
)

func TestMapAccessFromAPIMapsAllResources(t *testing.T) {
	got := mapAccessFromAPI([]api.Access{
		{Privileges: []string{"dataRead"}, Resources: &api.AccessibleResources{Buckets: []api.Bucket{{Name: "b"}}}},
		{Privileges: []string{"externalCatalogRead"}, Resources: &api.AccessibleResources{Catalogs: []string{"*"}}},
		{Privileges: []string{"credentialConsumer"}, Resources: &api.AccessibleResources{Credentials: []string{"cred-1"}}},
		{Privileges: []string{"analyticsAdmin"}},
	})

	assert.Len(t, got[0].Resources.Buckets, 1)
	assert.Nil(t, got[0].Resources.Catalogs)
	assert.Equal(t, []types.String{types.StringValue("*")}, got[1].Resources.Catalogs)
	assert.Equal(t, []types.String{types.StringValue("cred-1")}, got[2].Resources.Credentials)
	assert.Nil(t, got[3].Resources)
}
