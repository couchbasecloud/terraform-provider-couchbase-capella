package snapshot_backup

type SnapshotRestore struct {
	ClusterID      string `json:"clusterId"`
	CreatedAt      string `json:"createdAt"`
	ID             string `json:"id"`
	ProjectID      string `json:"projectId"`
	RestoreTo      string `json:"restoreTo"`
	Snapshot       string `json:"snapshot"`
	Status         State  `json:"status"`
	OrganizationID string `json:"tenantId"`
}

type ListSnapshotRestoresResponse struct {
	Data []SnapshotRestore `json:"data"`
}
