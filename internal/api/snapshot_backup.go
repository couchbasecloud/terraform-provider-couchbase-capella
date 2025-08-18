package api

type Progress struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}

type CMEK struct {
	ID         string `json:"id"`
	ProviderID string `json:"providerId"`
}

type Server struct {
	Version string `json:"version"`
}

type ProjectLevelSnapshotBackup struct {
	ClusterID    string   `json:"clusterId"`
	CreatedAt    string   `json:"createdAt"`
	Expiration   string   `json:"expiration"`
	ID           string   `json:"id"`
	Progress     Progress `json:"progress"`
	ProjectID    string   `json:"projectId"`
	AppService   string   `json:"appService"`
	CMEK         []CMEK   `json:"cmek"`
	Retention    string   `json:"retention"`
	Server       Server   `json:"server"`
	DatabaseSize string   `json:"databaseSize"`
	TenantID     string   `json:"tenantId"`
	Type         string   `json:"type"`
}

type SnapshotBackup struct {
	ClusterID  string   `json:"clusterId"`
	CreatedAt  string   `json:"createdAt"`
	Expiration string   `json:"expiration"`
	BackupID   string   `json:"id"`
	Progress   Progress `json:"progress"`
	ProjectID  string   `json:"projectId"`
	Retention  string   `json:"retention"`
	CMEK       []CMEK   `json:"cmek"`
	Server     string   `json:"server"`
	Size       string   `json:"size"`
	TenantID   string   `json:"tenantId"`
	Type       string   `json:"type"`
}

type SnapshotRestore struct {
	ClusterID  string   `json:"clusterId"`
	CreatedAt  string   `json:"createdAt"`
	Expiration string   `json:"expiration"`
	BackupID   string   `json:"id"`
	Progress   Progress `json:"progress"`
	ProjectID  string   `json:"projectId"`
	Retention  string   `json:"retention"`
	Server     string   `json:"server"`
	Size       string   `json:"size"`
	TenantID   string   `json:"tenantId"`
	Type       string   `json:"type"`
}

type ProjectLevelSnapshotBackups struct {
	ClusterID          string                     `json:"clusterId"`
	ClusterName        string                     `json:"clusterName"`
	CreationDateTime   string                     `json:"creationDateTime"`
	CreatedBy          string                     `json:"createdBy"`
	CurrentStatus      string                     `json:"currentStatus"`
	CloudProvider      string                     `json:"cloudProvider"`
	Region             string                     `json:"region"`
	MostRecentSnapshot ProjectLevelSnapshotBackup `json:"mostRecentSnapshot"`
	OldestSnapshot     ProjectLevelSnapshotBackup `json:"oldestSnapshot"`
}

type CloudProvider struct {
	Type   string `json:"type"`
	Region string `json:"region"`
	CIDR   string `json:"cidr"`
}

type Availability struct {
	Type string `json:"type"`
}

type Support struct {
	Plan     string `json:"plan"`
	Timezone string `json:"timezone"`
}

type CreateSnapshotBackupRequest struct {
	Retention int `json:"retention"`
}

type CreateSnapshotBackupResponse struct {
	BackupID string `json:"backupId"`
}

type ListSnapshotBackupsResponse struct {
	Data []SnapshotBackup `json:"data"`
}

type ListSnapshotRestoresResponse struct {
	Data []SnapshotRestore `json:"data"`
}

type EditBackupRetentionRequest struct {
	Retention int `json:"retention"`
}

type RestoreSnapshotBackupResponse struct {
	RestoreID string `json:"restoreId"`
}

type ListProjectLevelSnapshotBackupsResponse struct {
	Data []ProjectLevelSnapshotBackups `json:"data"`
}

type CloneClusterBackupRequest struct {
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	CloudProvider CloudProvider `json:"cloudProvider"`
	Availability  Availability  `json:"availability"`
	Zones         []string      `json:"zones"`
	Support       Support       `json:"support"`
}

type CloneClusterBackupResponse struct {
	RestoreID string `json:"restoreId"`
	ClusterID string `json:"clusterId"`
}
