package localcache


func isPowerOfTwo(number uint64) bool {
	return (number != 0) && (number&(number-1)) == 0
}