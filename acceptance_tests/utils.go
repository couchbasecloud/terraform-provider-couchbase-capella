package acceptance_tests

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	providerschema "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/schema"
)

const (
	// charSetAlpha are lower case alphabet letters.
	charSetAlpha = "abcdefghijklmnopqrstuvwxyz"

	// resourceNameLength is the length of the resource name we wish to generate.
	resourceNameLength = 10
)

func randomString() string {
	result := make([]byte, resourceNameLength)
	for i := 0; i < resourceNameLength; i++ {
		result[i] = charSetAlpha[rand.Intn(len(charSetAlpha))] // #nosec G404
	}
	return string(result)
}

func randomStringWithPrefix(prefix string) string {
	return prefix + randomString()
}

// generateRandomCIDR generates a random /23 CIDR block in the 10.0.0.0/8 private range
// to avoid conflicts with existing clusters in the organization.
// Format: 10.X.Y.0/23 where X is 0-255 and Y is an even number (0, 2, 4, ..., 254)
func generateRandomCIDR() string {
	// Use crypto/rand for better randomness
	buf := make([]byte, 2)
	cryptorand.Read(buf)

	// Second octet: 0-255
	secondOctet := int(buf[0])

	// Third octet: must be even for /23 CIDR (0, 2, 4, ..., 254)
	thirdOctet := int(buf[1]) & 0xFE // Clear the last bit to make it even

	return fmt.Sprintf("10.%d.%d.0/23", secondOctet, thirdOctet)
}

func newTestClient() *providerschema.Data {
	providerData := &providerschema.Data{
		HostURL:  globalHost,
		Token:    globalToken,
		ClientV1: api.NewClient(timeout),
	}
	return providerData
}
