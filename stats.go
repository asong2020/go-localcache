package localcache

import (
	"sync/atomic"
)

// Stats stores cache statistics
type Stats struct {
	// Hits is a number of successfully found keys
	Hits int64 `json:"hits"`
	// Misses is a number of not found keys
	Misses int64 `json:"misses"`
	// DelHits is a number of successfully deleted keys
	DelHits int64 `json:"delete_hits"`
	// DelMisses is a number of not deleted keys
	DelMisses int64 `json:"delete_misses"`
	// Collisions is a number of happened key-collisions
	Collisions int64 `json:"collisions"`
	// hashmapStats record key hit
	hashmapStats map[string]int64
	statsEnabled bool
}

func newStats(enabled bool) *Stats {
	s := &Stats{
		statsEnabled: enabled,
	}
	if enabled {
		s.hashmapStats = make(map[string]int64)
	}
	return s
}

func (s *Stats) miss() {
	if !s.statsEnabled {
		return
	}
	atomic.AddInt64(&s.Misses, 1)
}

func (s *Stats) delHit() {
	if !s.statsEnabled {
		return
	}
	atomic.AddInt64(&s.DelHits, 1)
}

func (s *Stats) delMiss() {
	if !s.statsEnabled {
		return
	}
	atomic.AddInt64(&s.DelMisses, 1)
}

func (s *Stats) collision() {
	if !s.statsEnabled {
		return
	}
	atomic.AddInt64(&s.Collisions, 1)
}

func (s *Stats) hit(key string)  {
	if !s.statsEnabled {
		return
	}
	atomic.AddInt64(&s.Hits, 1)
	s.hashmapStats[key]++
}

func (s *Stats) getMisses() int64 {
	return atomic.LoadInt64(&s.Misses)
}

func (s *Stats) getDelHits() int64 {
	return atomic.LoadInt64(&s.DelHits)
}

func (s *Stats) getDelMisses() int64 {
	return atomic.LoadInt64(&s.DelMisses)
}

func (s *Stats) getCollisions() int64 {
	return atomic.LoadInt64(&s.Collisions)
}

func (s *Stats) getHits() int64 {
	return atomic.LoadInt64(&s.Hits)
}

func (s *Stats)getKeyHits(key string) int64{
	c := s.hashmapStats[key]
	return c
}