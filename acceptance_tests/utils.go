package acceptance_tests

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

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

// generateRandomCIDR generates a random CIDR block across all three RFC 1918
// private ranges to avoid conflicts with existing clusters in the shared test
// organization. It randomly picks one of:
//
//	10.A.B.0/23    (A: 0-255, B: even 0-254)  → 32,768 blocks
//	172.C.D.0/21   (C: 16-31, D: multiple of 8) → 512 blocks
//	192.168.E.0/24 (E: 0-255)                  → 256 blocks
func generateRandomCIDR() string {
	// Use crypto/rand for cryptographically secure randomness
	buf := make([]byte, 3)
	if _, err := cryptorand.Read(buf); err != nil {
		// If crypto/rand fails, this indicates a serious system issue.
		// For test utilities, it's appropriate to panic rather than continue with weak randomness.
		panic(fmt.Sprintf("failed to generate random CIDR: crypto/rand.Read failed: %v", err))
	}

	// Pick one of three RFC 1918 ranges weighted roughly by address space.
	switch int(buf[0]) % 3 {
	case 0: // 10.A.B.0/23 — most common, largest pool
		secondOctet := int(buf[1])
		thirdOctet := int(buf[2]) & 0xFE // even for /23
		return fmt.Sprintf("10.%d.%d.0/23", secondOctet, thirdOctet)
	case 1: // 172.C.D.0/21
		secondOctet := 16 + int(buf[1])%16    // 16..31
		thirdOctet := (int(buf[2]) % 32) * 8  // 0, 8, 16, ..., 248
		return fmt.Sprintf("172.%d.%d.0/21", secondOctet, thirdOctet)
	default: // 192.168.E.0/24
		thirdOctet := int(buf[2])
		return fmt.Sprintf("192.168.%d.0/24", thirdOctet)
	}
}

func newTestClient(t *testing.T) *providerschema.Data {
	t.Helper()

	retryingHTTP := apigen.NewRetryHTTPClient(context.Background(), timeout, false)
	clientV2, err := apigen.NewClientWithResponses(globalHost, apigen.WithHTTPClient(retryingHTTP), apigen.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+globalToken)
		return nil
	}))
	if err != nil {
		t.Fatalf("failed to create testing V2 API client: %v", err)
	}

	providerData := &providerschema.Data{
		HostURL:  globalHost,
		Token:    globalToken,
		ClientV1: api.NewClient(timeout),
		ClientV2: clientV2,
	}
	return providerData
}
