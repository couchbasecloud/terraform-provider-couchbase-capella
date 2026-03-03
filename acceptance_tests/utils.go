package acceptance_tests

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api"
	apigen "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/generated/api"
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
	// Use crypto/rand for cryptographically secure randomness
	buf := make([]byte, 2)
	if _, err := cryptorand.Read(buf); err != nil {
		// If crypto/rand fails, this indicates a serious system issue.
		// For test utilities, it's appropriate to panic rather than continue with weak randomness.
		panic(fmt.Sprintf("failed to generate random CIDR: crypto/rand.Read failed: %v", err))
	}

	// Second octet: 0-255
	secondOctet := int(buf[0])

	// Third octet: must be even for /23 CIDR (0, 2, 4, ..., 254)
	thirdOctet := int(buf[1]) & 0xFE // Clear the last bit to make it even

	return fmt.Sprintf("10.%d.%d.0/23", secondOctet, thirdOctet)
}

func newTestClient() *providerschema.Data {
	retryingHTTP := apigen.NewRetryHTTPClient(context.Background(), timeout, false)
	clientV2, err := apigen.NewClientWithResponses(globalHost, apigen.WithHTTPClient(retryingHTTP), apigen.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+globalToken)
		return nil
	}))
	if err != nil {
		// Can't proceed with tests if we fail to create a client, so we panic here.
		panic(fmt.Sprintf("failed to create V2 API client: %v", err))
	}

	providerData := &providerschema.Data{
		HostURL:  globalHost,
		Token:    globalToken,
		ClientV1: api.NewClient(timeout),
		ClientV2: clientV2,
	}
	return providerData
}
