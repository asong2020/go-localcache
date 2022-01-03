package buffer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type bufferTestSuite struct {
	suite.Suite
}

func TestBufferTestSuite(t *testing.T) {
	suite.Run(t, new(bufferTestSuite))
}

func (b *bufferTestSuite) SetupSuite() {}

func (b *bufferTestSuite) TestPush() {
	b.T().Parallel()

	buffer := NewBuffer(1)
	entry := []byte("Hello")

	index, err := buffer.Push(entry)

	assert.Equal(b.T(), 0, index)
	assert.Equal(b.T(), nil, err)
}

func (b *bufferTestSuite) TestPushAndGet() {
	b.T().Parallel()

	buffer := NewBuffer(1)
	entry := []byte("Hello")

	index, err := buffer.Push(entry)

	assert.Equal(b.T(), 0, index)
	assert.Equal(b.T(), nil, err)

	res,err := buffer.Get(index)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), entry, res)
	b.T().Logf("res is %s\n", string(res))
}

func (b *bufferTestSuite) TestLen() {
	b.T().Parallel()

	buffer := NewBuffer(1)
	entry := []byte("Hello")

	index, err := buffer.Push(entry)

	assert.Equal(b.T(), 0, index)
	assert.Equal(b.T(), nil, err)

	assert.Equal(b.T(), 1, buffer.Len())
}

func (b *bufferTestSuite) TestCapacity() {
	b.T().Parallel()

	buffer := NewBuffer(1)
	entry := []byte("Hello")

	index, err := buffer.Push(entry)

	assert.Equal(b.T(), 0, index)
	assert.Equal(b.T(), nil, err)

	assert.Equal(b.T(), 1, buffer.Capacity())
}

func (b *bufferTestSuite) TestPushAndRemove() {
	b.T().Parallel()

	buffer := NewBuffer(1)
	entry := []byte("Hello")

	index, err := buffer.Push(entry)

	assert.Equal(b.T(), 0, index)
	assert.Equal(b.T(), nil, err)

	res,err := buffer.Get(index)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), entry, res)

	err = buffer.Remove(index)
	assert.Equal(b.T(), nil, err)

	res,err = buffer.Get(index)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), []byte(nil), res)
}

func (b *bufferTestSuite) TestRest()  {
	b.T().Parallel()

	buffer := NewBuffer(10)

	index1, err := buffer.Push([]byte("asong"))
	assert.Equal(b.T(), 0, index1)
	assert.Equal(b.T(), nil, err)

	index2, err := buffer.Push([]byte("公众号：Golang梦工厂"))
	assert.Equal(b.T(), 1, index2)
	assert.Equal(b.T(), nil, err)

	res1, err := buffer.Get(index1)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), []byte("asong"), res1)

	res2, err := buffer.Get(index2)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), []byte("公众号：Golang梦工厂"), res2)

	buffer.Reset()
	assert.Equal(b.T(), 0, buffer.Len())
	res1, err = buffer.Get(index1)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), []byte(nil), res1)

	res2, err = buffer.Get(index2)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), []byte(nil), res2)
}

func (b *bufferTestSuite) TestAvailableSpace()  {
	buffer := NewBuffer(3)

	prefix := "asong"
	for i := 0; i < 3; i++{
		key := fmt.Sprintf(prefix+"%02d",i)
		index, err := buffer.Push([]byte(key))
		assert.Equal(b.T(), i, index)
		assert.Equal(b.T(), nil, err)
	}

	entry := []byte("公众号：Golang梦工厂")
	_,err := buffer.Push(entry)
	assert.Equal(b.T(), ErrBufferFull, err)

	err = buffer.Remove(1)
	assert.Equal(b.T(), nil, err)

	index,err := buffer.Push(entry)
	assert.Equal(b.T(), nil, err)
	assert.Equal(b.T(), 1, index)

	_, err = buffer.Push(entry)
	assert.Equal(b.T(), ErrBufferFull, err)
}

func (b *bufferTestSuite) TestGetPlaceholderCount() {
	b.T().Parallel()
	buffer := NewBuffer(3)
	prefix := "asong"
	for i := 0; i < 3; i++{
		key := fmt.Sprintf(prefix+"%02d",i)
		index, err := buffer.Push([]byte(key))
		assert.Equal(b.T(), i, index)
		assert.Equal(b.T(), nil, err)
	}

	count := buffer.GetPlaceholderCount()
	assert.Equal(b.T(), 3, count)
}

func (b *bufferTestSuite) TestGetPlaceholderIndex() {
	b.T().Parallel()
	buffer := NewBuffer(3)
	prefix := "asong"
	for i := 0; i < 3; i++{
		key := fmt.Sprintf(prefix+"%02d",i)
		index, err := buffer.Push([]byte(key))
		assert.Equal(b.T(), i, index)
		assert.Equal(b.T(), nil, err)
	}

	expected := []int{0, 1, 2}
	indexs := buffer.GetPlaceholderIndex()
	assert.Equal(b.T(), expected, indexs)
}