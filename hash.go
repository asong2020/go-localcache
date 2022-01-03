package localcache

// HashFunc is responsible for generating unsigned 64-bit hash of provided string
type HashFunc interface {
	Sum64(string) uint64
}

func NewDefaultHashFunc() HashFunc {
	return fnv64a{}
}

func NewHashWithDjb() HashFunc {
	return newTimes33()
}