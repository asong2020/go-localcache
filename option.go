package localcache

import "time"

// options.go
type options struct {
	hashFunc HashFunc
	bucketCount uint64
	maxBytes uint64
	cleanTime time.Duration
	statsEnabled bool
	cleanupEnabled bool
}

type Opt func(options *options)

func SetHashFunc(hashFunc HashFunc) Opt {
	return func(opt *options) {
		opt.hashFunc = hashFunc
	}
}

func SetShardCount(count uint64) Opt {
	return func(opt *options) {
		opt.bucketCount = count
	}
}

func SetMaxBytes(maxBytes uint64) Opt {
	return func(opt *options) {
		opt.maxBytes = maxBytes
	}
}

func SetCleanTime(time time.Duration) Opt {
	return func(opt *options) {
		opt.cleanTime = time
	}
}

func SetStatsEnabled(enabled bool) Opt {
	return func(opt *options) {
		opt.statsEnabled = enabled
	}
}

func SetCleanupEnabled(enabled bool) Opt {
	return func(opt *options) {
		opt.cleanupEnabled = enabled
	}
}