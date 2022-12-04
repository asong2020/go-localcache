package localcache

import "time"

type clock interface {
	Epoch(ttl time.Duration) int64
	TimeStamp() int64
}

type systemClock struct {
}

func (c systemClock) Epoch(ttl time.Duration) int64 {
	return time.Now().Add(ttl).Unix()
}

func (c systemClock) TimeStamp() int64 {
	return time.Now().Unix()
}