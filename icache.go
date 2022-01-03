package localcache

import "time"

// ICache abstract interface
type ICache interface {
	// Set value use default expire time. default does not expire.
	Set(key string, value []byte) error
	// Get value if find it. if value already expire will delete.
	Get(key string) ([]byte, error)
	// SetWithTime set value with expire time
	SetWithTime(key string, value []byte, expired time.Duration) error
	// Delete manual removes the key
	Delete(key string) error
	// Len computes number of entries in cache
	Len() int
	// Capacity returns amount of bytes store in the cache.
	Capacity() int
	// Close is used to signal a shutdown of the cache when you are done with it.
	// This allows the cleaning goroutines to exit and ensures references are not
	// kept to the cache preventing GC of the entire cache.
	Close() error
	// Stats returns cache's statistics
	Stats() Stats
	// GetKeyHit returns key hit
	GetKeyHit(key string) int64
}