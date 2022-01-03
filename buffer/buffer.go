package buffer

const (
	defaultIndex = 0
	defaultCount = 0
)

var (
	ErrIndexOutOFBounds = &bufferError{"Index out of range."}
	ErrInvalidIndex = &bufferError{"Index must be greater than or equal zero. Invalid index."}
	ErrBufferFull   = &bufferError{"Buffer full."}
)

type Buffer struct {
	array [][]byte
	capacity int
	index int
	// maxCount = capacity - 1
	count int
	// availableSpace If any objects are removed after the buffer is full, the idle index is logged.
	// Avoid array "wormhole"
	availableSpace map[int]struct{}
	// placeholder record the index that buffer has stored.
	placeholder map[int]struct{}
}

type bufferError struct {
	message string
}

// Error returns error message
func (e *bufferError) Error() string {
	return e.message
}

func NewBuffer(capacity int) IBuffer {
	return &Buffer{
		array: make([][]byte, capacity),
		capacity: capacity,
		index: defaultIndex,
		availableSpace: make(map[int]struct{}, capacity),
		placeholder: make(map[int]struct{}, capacity),
	}
}

func (b *Buffer) Push(data []byte) (int, error)  {
	availableSpaceLen := len(b.availableSpace)
	if b.index >= b.capacity && availableSpaceLen == 0 {
		return 0, ErrBufferFull
	}
	index := 0
	if b.index >= b.capacity && availableSpaceLen != 0 {
		index = b.rangeGetAvailableSpace()
	} else {
		index = b.index
	}

	dataLen := len(data)
	b.array[index] = make([]byte, dataLen)

	copy(b.array[index][0:], data[:dataLen])
	b.count++
	if b.index <= b.capacity {
		b.index++
	}
	b.placeholder[index] = struct{}{}

	return index, nil
}

func (b *Buffer) rangeGetAvailableSpace() int {
	for key :=  range b.availableSpace{
		delete(b.availableSpace, key)
		return key
	}
	return 0
}

func (b *Buffer) Reset()  {
	b.index = defaultIndex
	b.count = defaultCount
	b.array =  make([][]byte, b.capacity)
	b.availableSpace = make(map[int]struct{}, b.capacity)
	b.placeholder = make(map[int]struct{}, b.capacity)
}

func (b *Buffer) Len() int{
	return b.count
}

func (b *Buffer) Capacity() int{
	return b.capacity
}

func (b *Buffer) GetPlaceholderCount() int {
	return len(b.placeholder)
}

func (b *Buffer) GetPlaceholderIndex() []int {
	res := make([]int, 0, len(b.placeholder))
	for index := range b.placeholder{
		res = append(res, index)
	}
	return res
}

func (b *Buffer) Get(index int) ([]byte, error){
	if index > b.capacity {
		return nil, ErrIndexOutOFBounds
	}
	if index < 0 {
		return nil, ErrInvalidIndex
	}
	item := b.array[index]
	return item, nil
}

func (b *Buffer) Remove(index int) error {
	if index > b.capacity{
		return ErrIndexOutOFBounds
	}
	if index < 0 {
		return ErrInvalidIndex
	}
	b.array[index] = nil
	b.count--
	b.availableSpace[index] = struct{}{}
	delete(b.placeholder, index)
	return nil
}