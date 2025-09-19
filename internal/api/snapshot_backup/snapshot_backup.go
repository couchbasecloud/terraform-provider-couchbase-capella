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

type SnapshotBackup struct {
	AppService     string   `json:"appService"`
	ClusterID      string   `json:"clusterId"`
	CreatedAt      string   `json:"createdAt"`
	Expiration     string   `json:"expiration"`
	ID             string   `json:"id"`
	Progress       Progress `json:"progress"`
	ProjectID      string   `json:"projectId"`
	Retention      int      `json:"retention"`
	CMEK           CMEKs    `json:"cmek"`
	Server         Server   `json:"server"`
	Size           int      `json:"size"`
	OrganizationID string   `json:"tenantId"`
	Type           string   `json:"type"`
}

type CreateSnapshotBackupRequest struct {
	Retention int `json:"retention"`
}

type CreateSnapshotBackupResponse struct {
	ID string `json:"backupId"`
}

type ListSnapshotBackupsResponse struct {
	Data []SnapshotBackup `json:"data"`
}

type EditBackupRetentionRequest struct {
	Retention int `json:"retention"`
}
