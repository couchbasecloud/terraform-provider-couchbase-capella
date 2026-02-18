package datasources

import (
	"encoding/json"
	"testing"

	clusterapi "github.com/couchbasecloud/terraform-provider-couchbase-capella/internal/api/cluster"
	"github.com/stretchr/testify/assert"
)

func TestClusterStats_Unmarshall(t *testing.T) {
	jsonResponse := `{
    "freeMemoryInMb": 1871,
    "maxReplicas": 2,
    "totalMemoryInMb": 2071
}`

	var stats clusterapi.GetClusterStatsResponse
	err := json.Unmarshal([]byte(jsonResponse), &stats)
	assert.NoError(t, err)

	assert.Equal(t, int64(1871), stats.FreeMemoryInMb)
	assert.Equal(t, int64(2), stats.MaxReplicas)
	assert.Equal(t, int64(2071), stats.TotalMemoryInMb)
}
