package localcache

import (
	"github.com/stretchr/testify/suite"
	"math/rand"
	"testing"
	"time"
)


type hashFuncTestSuite struct {
	suite.Suite
	randString [100000]string
	letterRunes []rune
	fnv HashFunc
	djb HashFunc
}

func TestHashFuncTestSuite(t *testing.T) {
	suite.Run(t, new(hashFuncTestSuite))
}

func (h *hashFuncTestSuite) SetupSuite() {
	rand.Seed(time.Now().UnixNano())
	h.letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 100000; i++{
		h.randString[i] = h.RandStringRunes(rand.Intn(10))
	}
	h.fnv = NewDefaultHashFunc()
	h.djb = NewHashWithDjb()
}

func (h *hashFuncTestSuite) RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = h.letterRunes[rand.Intn(len(h.letterRunes))]
	}
	return string(b)
}

func (h *hashFuncTestSuite) TestHashFuncTestSuite_FNVConflict()  {
	m := make(map[uint64]string)
	count := 0
	for index := 0; index < len(h.randString); index++ {
		if h.randString[index] == ""{
			continue
		}
		res := h.fnv.Sum64(h.randString[index])
		if value, ok := m[res]; ok &&  value != h.randString[index]{
			count++
		}else {
			m[res] = h.randString[index]
		}
	}
	h.T().Logf("%d length fnv hash conflict count is %d", len(h.randString), count)
}

func (h *hashFuncTestSuite) TestHashFuncTestSuite_DJBConflict() {
	m := make(map[uint64]string)
	count := 0
	for index := 0; index < len(h.randString); index++ {
		if h.randString[index] == ""{
			continue
		}
		res := h.djb.Sum64(h.randString[index])
		if value, ok := m[res]; ok && value != h.randString[index]{
			count++
		}else {
			m[res] = h.randString[index]
		}
	}
	h.T().Logf("%d length djb hash conflict count is %d", len(h.randString), count)
}