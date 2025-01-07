package acceptance_tests

import (
	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
	"math/rand"
	"time"
)

const (
	// charSetAlpha are lower case alphabet letters.
	charSetAlpha = "abcdefghijklmnopqrstuvwxyz"

	// Length of the resource name we wish to generate.
	resourceNameLength = 10

	apiRequestTimeout = 60 * time.Second
)

func RandomString() string {
	result := make([]byte, resourceNameLength)
	for i := 0; i < resourceNameLength; i++ {
		result[i] = charSetAlpha[rand.Intn(len(charSetAlpha))]
	}
	return string(result)
}

func RandomStringWithPrefix(prefix string) string {
	return prefix + RandomString()
}

func NewTestClient() *providerschema.Data {
	providerData := &providerschema.Data{
		HostURL: Host,
		Token:   Token,
		Client:  api.NewClient(apiRequestTimeout),
	}
	return providerData
}
