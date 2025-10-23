package snapshot_backup

type Progress struct {
	Status State  `json:"status"`
	Time   string `json:"time"`
}

type CMEK struct {
	ID         string `json:"id"`
	ProviderID string `json:"providerId"`
}

type CMEKs []CMEK

type Server struct {
	Version string `json:"version"`
}

type CrossRegionCopy struct {
	RegionCode string `json:"regionCode"`
	Status     State  `json:"status"`
	Time       string `json:"time"`
}

type SnapshotBackup struct {
	ClusterID         string            `json:"clusterId"`
	CreatedAt         string            `json:"createdAt"`
	Expiration        string            `json:"expiration"`
	ID                string            `json:"id"`
	Progress          Progress          `json:"progress"`
	ProjectID         string            `json:"projectId"`
	Retention         int64             `json:"retention"`
	CrossRegionCopies []CrossRegionCopy `json:"crossRegionCopies"`
	CMEK              CMEKs             `json:"cmek"`
	Server            Server            `json:"server"`
	Size              int               `json:"size"`
	OrganizationID    string            `json:"tenantId"`
	Type              string            `json:"type"`
}

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

type CreateSnapshotBackupRequest struct {
	Retention     int64    `json:"retention"`
	RegionsToCopy []string `json:"regionsToCopy"`
}

type CreateSnapshotBackupResponse struct {
	ID string `json:"backupId"`
}

type ListSnapshotBackupsResponse struct {
	Data []SnapshotBackup `json:"data"`
}

type EditBackupRetentionRequest struct {
	Retention int64 `json:"retention"`
}

type ListSnapshotRestoresResponse struct {
	Data []SnapshotRestore `json:"data"`
}
