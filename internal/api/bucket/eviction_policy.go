package bucket

const (
	// EvictionValue is the default policy for Couchbase and indicates that
	// no members of the cache will be evicted. This improves performance
	// but restricts the size of the bucket to be in-memory.
	EvictionValue BucketEviction = "valueOnly"

	// EvictionFull allows for cache member eviction and will allow the size
	// of the bucket to be unbounded.
	EvictionFull BucketEviction = "fullEviction"

	// EvictionNone only applies to ephemeral buckets and should not be used
	// in Couchbase Capella.
	EvictionNone BucketEviction = "noEviction"

	// EvictionNRU only applies to ephemeral buckets and should not be used
	// in Couchbase Capella.
	EvictionNRU BucketEviction = "nruEviction"
)

type BucketEviction string
