package bucket

const (
	Couchstore StorageBackend = "couchstore"
	Magma      StorageBackend = "magma"
)

type StorageBackend string
