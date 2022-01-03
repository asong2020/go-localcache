package localcache

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/asong2020/go-localcache/buffer"
	"time"
)

var (
	// ErrEntryNotFound is an error type struct which is returned when entry was not found for provided key
	ErrEntryNotFound = errors.New("Entry not found")
	ErrExpireTimeInvalid = errors.New("Entry expire time invalid")
)

const (
	segmentSizeBits = 40
	maxSegmentSize uint64 = 1 << segmentSizeBits
	segmentSize = 32 * 1024 // 32kb
	defaultExpireTime = 10 * time.Minute
)

type segment struct {
	hashmap map[uint64]uint32
	entries buffer.IBuffer
	clock   clock
	evictList  *list.List
	stats IStats
}

func newSegment(bytes uint64, statsEnabled bool) *segment {
	if bytes == 0 {
		panic(fmt.Errorf("bytes cannot be zero"))
	}
	if bytes >= maxSegmentSize{
		panic(fmt.Errorf("too big bytes=%d; should be smaller than %d", bytes, maxSegmentSize))
	}
	capacity := (bytes + segmentSize - 1) / segmentSize
	entries := buffer.NewBuffer(int(capacity))
	entries.Reset()
	return &segment{
		entries: entries,
		hashmap: make(map[uint64]uint32),
		clock:   &systemClock{},
		evictList: list.New(),
		stats: newStats(statsEnabled),
	}
}

func (s *segment) set(key string, hashKey uint64, value []byte, expireTime time.Duration) error {
	if expireTime <= 0{
		return ErrExpireTimeInvalid
	}
	expireAt := uint64(s.clock.Epoch(expireTime))

	if previousIndex, ok := s.hashmap[hashKey]; ok {
		if err := s.entries.Remove(int(previousIndex)); err != nil{
			return err
		}
		delete(s.hashmap, hashKey)
	}

	entry := wrapEntry(expireAt, key, hashKey, value)
	for {
		index, err := s.entries.Push(entry)
		if err == nil {
			s.hashmap[hashKey] = uint32(index)
			s.evictList.PushFront(index)
			return nil
		}
		ele := s.evictList.Back()
		if err := s.entries.Remove(ele.Value.(int)); err != nil{
			return err
		}
		s.evictList.Remove(ele)
	}
}

func (s *segment) getWarpEntry(key string, hashKey uint64) ([]byte,error) {
	index, ok := s.hashmap[hashKey]
	if !ok {
		s.stats.miss()
		return nil, ErrEntryNotFound
	}
	entry, err := s.entries.Get(int(index))
	if err != nil{
		s.stats.miss()
		return nil, err
	}
	if entry == nil{
		s.stats.miss()
		return nil, ErrEntryNotFound
	}

	if entryKey := readKeyFromEntry(entry); key != entryKey {
		s.stats.collision()
		return nil, ErrEntryNotFound
	}
	return entry, nil
}

func (s *segment) get(key string, hashKey uint64) ([]byte, error) {
	currentTimestamp := s.clock.TimeStamp()
	entry, err := s.getWarpEntry(key, hashKey)
	if err != nil{
		return nil, err
	}
	res := readEntry(entry)

	expireAt := int64(readExpireAtFromEntry(entry))
	if currentTimestamp - expireAt >= 0{
		_ = s.entries.Remove(int(s.hashmap[hashKey]))
		delete(s.hashmap, hashKey)
		return nil, ErrEntryNotFound
	}
	s.stats.hit(key)

	return res, nil
}

func (s *segment) len() int {
	res := len(s.hashmap)
	return res
}

func (s *segment) capacity() int {
	res := s.entries.Capacity()
	return res
}

func (s *segment) delete(hashKey uint64) error {
	index,ok := s.hashmap[hashKey]
	if !ok {
		s.stats.delMiss()
		return ErrEntryNotFound
	}

	if err := s.entries.Remove(int(index)); err != nil{
		return err
	}

	delete(s.hashmap, hashKey)
	s.stats.delHit()
	return nil
}

func (s *segment) cleanup(currentTimestamp int64) {
	indexs := s.entries.GetPlaceholderIndex()
	for index := range indexs {
		entry, err := s.entries.Get(index)
		if err != nil || entry == nil{
			continue
		}
		expireAt := int64(readExpireAtFromEntry(entry))
		if currentTimestamp - expireAt >= 0{
			hash := readHashFromEntry(entry)
			delete(s.hashmap, hash)
			_ = s.entries.Remove(index)
			continue
		}
	}
}

func (s *segment) getStats() Stats {
	res := Stats{
		Hits:       s.stats.getHits(),
		Misses:     s.stats.getMisses(),
		DelHits:    s.stats.getDelHits(),
		DelMisses:  s.stats.getDelMisses(),
		Collisions: s.stats.getCollisions(),
	}
	return res
}

func (s *segment) getKeyHit(key string) int64 {
	return s.stats.getKeyHits(key)
}