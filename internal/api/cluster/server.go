package cluster

// Service is a string for each individual service in Couchbase Server.
// This field can accept enums: kv, index, n1ql, fts, analytics and eventing.
type Service string

// ServiceGroup The set of nodes that share the same disk,
// number of nodes and services.
type ServiceGroup struct {
	Node *Node `json:"node,omitempty"`

	// NumOfNodes is the number of nodes. The minimum number of
	// nodes for the cluster can be 3 and maximum can be 27 nodes.
	// Additional service groups can have 2 nodes minimum and 24 nodes maximum.
	NumOfNodes *int `json:"numOfNodes,omitempty"`

	// Services is the couchbase service to run on the node.
	Services *[]Service `json:"services,omitempty"`
}

// CouchbaseServer defines version for the Couchbase Server to be launched during the creation of the Capella cluster.
type CouchbaseServer struct {
	// Version is version of the Couchbase Server to be installed
	// in the cluster. Refer to documentation
	// [here](https://docs.couchbase.com/cloud/clusters/upgrade-database.html#server-version-maintenance-support)
	// for list of supported versions. The latest Couchbase Server version
	// will be deployed by default.
	Version *string `json:"version,omitempty"`
}

// AreEqual returns true if the two arrays contain the same elements, without any extra values, False otherwise.
func AreEqual[T comparable](array1 []T, array2 []T) bool {
	if len(array1) != len(array2) {
		return false
	}
	set1 := make(map[T]bool)
	for _, element := range array1 {
		set1[element] = true
	}

	for _, element := range array2 {
		if !set1[element] {
			return false
		}
	}

	return len(set1) == len(array1)
}
