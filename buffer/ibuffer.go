package buffer

type IBuffer interface {
	// Reset removes all entries from buffer
	Reset()
	// Push Returns index or error if maximum size buffer limit is reached
 	Push(data []byte) (int, error)
	// Len returns number of entries kept in buffer
	Len() int
	// Capacity returns number of allocated bytes for buffer
	Capacity() int
	// Get reads entry from index
	Get(index int) ([]byte, error)
	// Remove delete entries by index
	Remove(index int) error
	// GetPlaceholderCount buffer store entry count
	GetPlaceholderCount() int
	// GetPlaceholderIndex get all index
	GetPlaceholderIndex() []int
}