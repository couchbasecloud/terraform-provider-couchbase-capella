package replication

import "fmt"

// ClusterReplicationsURL is the URL to create a replication on a source cluster,
// or to list the replications of a cluster.
func ClusterReplicationsURL(host, organizationId, projectId, clusterId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/replications",
		host, organizationId, projectId, clusterId,
	)
}

// ReplicationURL is the URL to get, update or delete one replication.
func ReplicationURL(host, organizationId, projectId, clusterId, replicationId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/replications/%s",
		host, organizationId, projectId, clusterId, replicationId,
	)
}

// ReplicationJobURL is the URL to get the status of an async replication create job.
func ReplicationJobURL(host, organizationId, projectId, clusterId, jobId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/replications/jobs/%s",
		host, organizationId, projectId, clusterId, jobId,
	)
}

// ActivationStatusURL is the URL to pause (DELETE) or resume (POST) one replication.
func ActivationStatusURL(host, organizationId, projectId, clusterId, replicationId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/clusters/%s/replications/%s/activationStatus",
		host, organizationId, projectId, clusterId, replicationId,
	)
}

// ProjectReplicationsURL is the URL to list the replications of a project.
func ProjectReplicationsURL(host, organizationId, projectId string) string {
	return fmt.Sprintf(
		"%s/v4/organizations/%s/projects/%s/replications",
		host, organizationId, projectId,
	)
}
