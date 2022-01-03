package localcache

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrShardCount = errors.New("shard count must be power of two")
	ErrBytes = errors.New("maxBytes must be greater than 0")
)

const (
	defaultBucketCount = 256
	defaultMaxBytes = 512 * 1024 * 1024 // 512M
	defaultCleanTIme = time.Minute * 10
	defaultStatsEnabled = false
	defaultCleanupEnabled = false
)

type cache struct {
	// hashFunc represents used hash func
	hashFunc HashFunc
	// bucketCount represents the number of segments within a cache instance. value must be a power of two.
	bucketCount uint64
	// bucketMask is bitwise AND applied to the hashVal to find the segment id.
	bucketMask uint64
	// segment is shard
	segments []*segment
	// segment lock
	locks    []sync.RWMutex
	// close cache
	close chan struct{}
}


// NewCache constructor cache instance
func NewCache(opts ...Opt) (ICache, error) {
	options := &options{
		hashFunc: NewDefaultHashFunc(),
		bucketCount: defaultBucketCount,
		maxBytes: defaultMaxBytes,
		cleanTime: defaultCleanTIme,
		statsEnabled: defaultStatsEnabled,
		cleanupEnabled: defaultCleanupEnabled,
	}
	for _, each := range opts{
		each(options)
	}

	if !isPowerOfTwo(options.bucketCount){
		return nil, ErrShardCount
	}

	if options.maxBytes <= 0 {
		return nil, ErrBytes
	}

	segments := make([]*segment, options.bucketCount)
	locks := make([]sync.RWMutex, options.bucketCount)

	maxSegmentBytes := (options.maxBytes + options.bucketCount - 1) / options.bucketCount
	for index := range segments{
		segments[index] = newSegment(maxSegmentBytes, options.statsEnabled)
	}

	c := &cache{
		hashFunc: options.hashFunc,
		bucketCount: options.bucketCount,
		bucketMask: options.bucketCount - 1,
		segments: segments,
		locks: locks,
		close: make(chan struct{}),
	}
    if options.cleanupEnabled {
		go c.cleanup(options.cleanTime)
	}

	return c, nil
}

func (c *cache) Set(key string, value []byte) error  {
	hashKey := c.hashFunc.Sum64(key)
	bucketIndex := hashKey&c.bucketMask
	c.locks[bucketIndex].Lock()
	defer c.locks[bucketIndex].Unlock()
	err := c.segments[bucketIndex].set(key, hashKey, value, defaultExpireTime)
	return err
}

func (c *cache) Get(key string) ([]byte, error)  {
	hashKey := c.hashFunc.Sum64(key)
	bucketIndex := hashKey&c.bucketMask
	c.locks[bucketIndex].RLock()
	defer c.locks[hashKey&c.bucketMask].RUnlock()
	entry, err := c.segments[bucketIndex].get(key, hashKey)
	if err != nil{
		return nil, err
	}
	return entry,nil
}

func (c *cache) SetWithTime(key string, value []byte, expired time.Duration) error{
	hashKey := c.hashFunc.Sum64(key)
	bucketIndex := hashKey&c.bucketMask
	c.locks[bucketIndex].Lock()
	defer c.locks[bucketIndex].Unlock()
	err := c.segments[bucketIndex].set(key, hashKey, value, expired)
	return err
}

func (c *cache) Delete(key string) error{
	hashKey := c.hashFunc.Sum64(key)
	bucketIndex := hashKey&c.bucketMask
	c.locks[bucketIndex].Lock()
	defer c.locks[bucketIndex].Unlock()
	err := c.segments[bucketIndex].delete(hashKey)
	return err
}

func (c *cache) Len() int {
	length := 0
	for index :=0; index < int(c.bucketCount); index++{
		c.locks[index].RLock()
		length += c.segments[index].len()
		c.locks[index].RUnlock()
	}
	return length
}

func (c *cache) Capacity() int {
	capacity := 0
	for index := 0; index < int(c.bucketCount); index++{
		c.locks[index].RLock()
		capacity += c.segments[index].capacity()
		c.locks[index].RUnlock()
	}
	return capacity
}

func (c *cache) Close() error {
	close(c.close)
	return nil
}

func (c *cache) cleanup(cleanTime time.Duration)  {
	ticker := time.NewTicker(cleanTime)
	defer ticker.Stop()
	for {
		select {
		case t := <- ticker.C:
			for index := 0; index < int(c.bucketCount); index++{
				c.locks[index].Lock()
				c.segments[index].cleanup(t.Unix())
				c.locks[index].Unlock()
			}
		case <- c.close:
			return
		}
	}
}

func (c *cache) Stats() Stats {
	s := Stats{}
	for _, shard := range c.segments {
		tmp := shard.getStats()
		s.Hits += tmp.Hits
		s.Misses += tmp.Misses
		s.DelHits += tmp.DelHits
		s.DelMisses += tmp.DelMisses
		s.Collisions += tmp.Collisions
	}
	return s
}

func (c *cache) GetKeyHit(key string) int64 {
	hashKey := c.hashFunc.Sum64(key)
	bucketIndex := hashKey&c.bucketMask
	c.locks[bucketIndex].Lock()
	defer c.locks[bucketIndex].Unlock()
	hit := c.segments[bucketIndex].getKeyHit(key)
	return hit
}