package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestItem is a simple struct for testing pagination
type TestItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// createMockPaginatedResponse creates a mock API response with pagination
func createMockPaginatedResponse(data []TestItem, currentPage, nextPage, lastPage, totalItems int, clusterStats map[string]int64) []byte {
	response := map[string]interface{}{
		"data": data,
		"cursor": map[string]interface{}{
			"pages": map[string]interface{}{
				"page":       currentPage,
				"next":       nextPage,
				"last":       lastPage,
				"perPage":    25,
				"totalItems": totalItems,
			},
		},
	}

	// Add clusterStats if provided (for testing metadata extraction)
	if clusterStats != nil {
		response["clusterStats"] = clusterStats
	}

	jsonBytes, _ := json.Marshal(response)
	return jsonBytes
}

func TestGetPaginated_SinglePage(t *testing.T) {
	// Setup mock server that returns a single page
	items := []TestItem{
		{Id: "1", Name: "item1"},
		{Id: "2", Name: "item2"},
		{Id: "3", Name: "item3"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request has pagination params
		assert.Contains(t, r.URL.RawQuery, "page=1")
		assert.Contains(t, r.URL.RawQuery, "perPage=25")

		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse(items, 1, 0, 1, 3, nil))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "1", result[0].Id)
	assert.Equal(t, "item1", result[0].Name)
	assert.Equal(t, "2", result[1].Id)
	assert.Equal(t, "3", result[2].Id)
}

func TestGetPaginated_MultiplePages(t *testing.T) {
	// Setup mock server that returns multiple pages
	page1Items := make([]TestItem, 25)
	for i := 0; i < 25; i++ {
		page1Items[i] = TestItem{Id: string(rune('a' + i)), Name: "item" + string(rune('a'+i))}
	}

	page2Items := []TestItem{
		{Id: "z", Name: "itemZ"},
		{Id: "y", Name: "itemY"},
		{Id: "x", Name: "itemX"},
		{Id: "w", Name: "itemW"},
		{Id: "v", Name: "itemV"},
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)

		if requestCount == 1 {
			// First page - has next page
			w.Write(createMockPaginatedResponse(page1Items, 1, 2, 2, 30, nil))
		} else {
			// Second page - no next page
			w.Write(createMockPaginatedResponse(page2Items, 2, 0, 2, 30, nil))
		}
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result, 30) // 25 from page 1 + 5 from page 2
	assert.Equal(t, 2, requestCount, "Should have made 2 requests for 2 pages")
}

func TestGetPaginatedWithMeta_SinglePage(t *testing.T) {
	items := []TestItem{
		{Id: "1", Name: "item1"},
		{Id: "2", Name: "item2"},
	}

	clusterStats := map[string]int64{
		"freeMemoryInMb":  640,
		"maxReplicas":     2,
		"totalMemoryInMb": 1040,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse(items, 1, 0, 1, 2, clusterStats))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginatedWithMeta[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.NotNil(t, result.RawFirstPage)

	// Verify we can extract clusterStats from RawFirstPage
	var meta struct {
		ClusterStats struct {
			FreeMemoryInMb  int64 `json:"freeMemoryInMb"`
			MaxReplicas     int64 `json:"maxReplicas"`
			TotalMemoryInMb int64 `json:"totalMemoryInMb"`
		} `json:"clusterStats"`
	}
	err = json.Unmarshal(result.RawFirstPage, &meta)
	require.NoError(t, err)
	assert.Equal(t, int64(640), meta.ClusterStats.FreeMemoryInMb)
	assert.Equal(t, int64(2), meta.ClusterStats.MaxReplicas)
	assert.Equal(t, int64(1040), meta.ClusterStats.TotalMemoryInMb)
}

func TestGetPaginatedWithMeta_MultiplePages_MetadataFromFirstPage(t *testing.T) {
	// This test verifies that RawFirstPage contains metadata from the first page,
	// even when there are multiple pages

	page1Items := []TestItem{{Id: "1", Name: "item1"}}
	page2Items := []TestItem{{Id: "2", Name: "item2"}}

	// Only first page has clusterStats
	page1ClusterStats := map[string]int64{
		"freeMemoryInMb":  100,
		"maxReplicas":     3,
		"totalMemoryInMb": 500,
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)

		if requestCount == 1 {
			w.Write(createMockPaginatedResponse(page1Items, 1, 2, 2, 2, page1ClusterStats))
		} else {
			// Second page doesn't have clusterStats (simulating real API behavior)
			w.Write(createMockPaginatedResponse(page2Items, 2, 0, 2, 2, nil))
		}
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginatedWithMeta[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result.Data, 2) // Combined from both pages

	// Verify RawFirstPage contains metadata from first page
	var meta struct {
		ClusterStats struct {
			FreeMemoryInMb  int64 `json:"freeMemoryInMb"`
			MaxReplicas     int64 `json:"maxReplicas"`
			TotalMemoryInMb int64 `json:"totalMemoryInMb"`
		} `json:"clusterStats"`
	}
	err = json.Unmarshal(result.RawFirstPage, &meta)
	require.NoError(t, err)
	assert.Equal(t, int64(100), meta.ClusterStats.FreeMemoryInMb)
	assert.Equal(t, int64(3), meta.ClusterStats.MaxReplicas)
	assert.Equal(t, int64(500), meta.ClusterStats.TotalMemoryInMb)
}

func TestGetPaginated_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse([]TestItem{}, 1, 0, 1, 0, nil))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result, 0)
}

func TestGetPaginated_WithSortParameter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify sortBy parameter is included
		assert.Contains(t, r.URL.RawQuery, "sortBy=name")

		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse([]TestItem{{Id: "1", Name: "test"}}, 1, 0, 1, 1, nil))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, SortByName)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetPaginated_WithoutSortParameter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify sortBy parameter is NOT included when empty
		assert.NotContains(t, r.URL.RawQuery, "sortBy")

		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse([]TestItem{{Id: "1", Name: "test"}}, 1, 0, 1, 1, nil))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, "")

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGetPaginated_ThreePages(t *testing.T) {
	// Test with exactly 3 pages to verify pagination loop works correctly
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)

		switch requestCount {
		case 1:
			w.Write(createMockPaginatedResponse(
				[]TestItem{{Id: "1", Name: "page1"}},
				1, 2, 3, 3, nil,
			))
		case 2:
			w.Write(createMockPaginatedResponse(
				[]TestItem{{Id: "2", Name: "page2"}},
				2, 3, 3, 3, nil,
			))
		case 3:
			w.Write(createMockPaginatedResponse(
				[]TestItem{{Id: "3", Name: "page3"}},
				3, 0, 3, 3, nil, // next=0 means last page
			))
		}
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	result, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg, SortById)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, 3, requestCount, "Should have made 3 requests for 3 pages")
	assert.Equal(t, "page1", result[0].Name)
	assert.Equal(t, "page2", result[1].Name)
	assert.Equal(t, "page3", result[2].Name)
}

func TestGetPaginatedWithMeta_BackwardCompatibility(t *testing.T) {
	// Verify that GetPaginatedWithMeta returns the same data as GetPaginated
	items := []TestItem{
		{Id: "1", Name: "item1"},
		{Id: "2", Name: "item2"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse(items, 1, 0, 1, 2, nil))
	}))
	defer server.Close()

	client := NewClient(10 * time.Second)
	cfg := EndpointCfg{
		Url:           server.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	// Get results from both functions
	resultWithMeta, err := GetPaginatedWithMeta[[]TestItem](context.Background(), client, "test-token", cfg, SortById)
	require.NoError(t, err)

	// Reset server for second call
	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(createMockPaginatedResponse(items, 1, 0, 1, 2, nil))
	}))
	defer server2.Close()

	cfg2 := EndpointCfg{
		Url:           server2.URL,
		Method:        http.MethodGet,
		SuccessStatus: http.StatusOK,
	}

	resultSimple, err := GetPaginated[[]TestItem](context.Background(), client, "test-token", cfg2, SortById)
	require.NoError(t, err)

	// Verify both return the same data
	assert.Equal(t, resultSimple, resultWithMeta.Data)
}
