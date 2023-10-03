package bucket

const (
	SeqNo ConflictResolution = "seqno"

	LWW ConflictResolution = "lww"
)

type ConflictResolution string
