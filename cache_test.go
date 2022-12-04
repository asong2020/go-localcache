package localcache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type cacheTestSuite struct {
	suite.Suite
}

func TestCacheTestSuite(t *testing.T) {
	suite.Run(t, new(cacheTestSuite))
}

func (h *cacheTestSuite) SetupSuite() {}

func (h *cacheTestSuite) TestSetAndGet() {
	cache, err := NewCache()
	assert.Equal(h.T(), nil, err)
	key := "asong"
	value := []byte("公众号：Golang梦工厂")

	err = cache.Set(key, value)
	assert.Equal(h.T(), nil, err)

	res, err := cache.Get(key)
	assert.Equal(h.T(), nil, err)
	assert.Equal(h.T(), value, res)
	h.T().Logf("get value is %s", string(res))
}

func (h *cacheTestSuite) TestSetOverLimit()  {
	cache, err := NewCache()
	assert.Equal(h.T(), nil, err)
	value := []byte("公众号：Golang梦工厂")
	for index := 0; index < 1000000; index++{
		key := fmt.Sprintf("asong%08d", index)
		err = cache.Set(key, value)
		assert.Equal(h.T(), nil, err)
	}
}

func (h *cacheTestSuite) TestLen() {
	cache, err := NewCache()
	assert.Equal(h.T(), nil, err)

	value := []byte("公众号：Golang梦工厂")
	for index := 0; index < 1000; index++{
		key := fmt.Sprintf("asong%03d", index)
		err = cache.Set(key, value)
		assert.Equal(h.T(), nil, err)
	}

	length := cache.Len()
	assert.Equal(h.T(), 1000, length)
	h.T().Logf("length == %d", length)
}

func (h *cacheTestSuite) TestCapacity() {
	cache, err := NewCache()
	assert.Equal(h.T(), nil, err)
	res := cache.Capacity()
	assert.Equal(h.T(), 16384, res)
	h.T().Logf("capacity == %d", res)
}

func (h *cacheTestSuite) TestDelete()  {
	cache, err := NewCache()
	assert.Equal(h.T(), nil, err)

	key := "asong"
	value := []byte("公众号：Golang梦工厂")
	err = cache.Set(key, value)
	assert.Equal(h.T(), nil, err)

	res, err := cache.Get(key)
	assert.Equal(h.T(), nil, err)
	assert.Equal(h.T(), value, res)

	err = cache.Delete(key)
	assert.Equal(h.T(), nil, err)
}

func (h *cacheTestSuite) TestCleanup()  {
	cache, err := NewCache(SetCleanTime(15 * time.Second))
	assert.Equal(h.T(), nil, err)

	value := []byte("公众号：Golang梦工厂")
	for index := 0; index < 1000; index++{
		key := fmt.Sprintf("asong%03d", index)
		err = cache.SetWithTime(key, value, time.Second * 10)
		assert.Equal(h.T(), nil, err)
	}

	time.Sleep(20 * time.Second)
	for index := 0; index < 1000; index++{
		key := fmt.Sprintf("asong%03d", index)
		_, err = cache.Get(key)
		assert.Equal(h.T(), ErrEntryNotFound, err)
	}
}

func (h *cacheTestSuite) TestStats() {
	h.T().Parallel()
	cache, err := NewCache(SetStatsEnabled(true))
	assert.Equal(h.T(), nil, err)

	value := []byte("公众号：Golang梦工厂")
	for index := 0; index < 100; index++{
		key := fmt.Sprintf("asong%03d", index)
		err = cache.Set(key, value)
		assert.Equal(h.T(), nil, err)
	}

	for i := 0; i < 10; i++ {
		entry, err := cache.Get(fmt.Sprintf("asong%03d", i))
		assert.Equal(h.T(), nil, err)
		assert.Equal(h.T(), value, entry)
	}

	for i := 100; i < 110; i++ {
		entry, err := cache.Get(fmt.Sprintf("asong%03d", i))
		assert.Equal(h.T(), ErrEntryNotFound, err)
		assert.Equal(h.T(), []byte(nil), entry)
	}

	for i := 10; i < 20; i++ {
		err := cache.Delete(fmt.Sprintf("asong%03d", i))
		assert.Equal(h.T(), nil, err)
	}

	for i := 110; i < 120; i++ {
		err := cache.Delete(fmt.Sprintf("asong%03d", i))
		assert.Equal(h.T(), ErrEntryNotFound, err)
	}

	stats := cache.Stats()
	assert.Equal(h.T(), int64(10), stats.Hits)
	assert.Equal(h.T(), int64(10), stats.Misses)
	assert.Equal(h.T(), int64(10), stats.DelHits)
	assert.Equal(h.T(), int64(10), stats.DelMisses)
}