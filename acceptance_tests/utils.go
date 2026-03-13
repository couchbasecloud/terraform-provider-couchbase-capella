package acceptance_tests

import (
	"context"
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
