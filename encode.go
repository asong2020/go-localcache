package localcache

import (
	"encoding/binary"
	"unsafe"
)

const (
	timestampSizeInBytes = 8                                                       // Number of bytes used for timestamp
	hashSizeInBytes      = 8                                                       // Number of bytes used for hash
	keySizeInBytes       = 2                                                       // Number of bytes used for size of entry key
	headersSizeInBytes   = timestampSizeInBytes + hashSizeInBytes + keySizeInBytes                // Number of bytes used for all headers
)


func wrapEntry(timestamp uint64, key string, hash uint64, entry []byte) []byte {
	keyLength := len(key)
	blobLength := len(entry) + keyLength + headersSizeInBytes
	blob := make([]byte, blobLength)

	binary.LittleEndian.PutUint64(blob, timestamp)
	binary.LittleEndian.PutUint64(blob[timestampSizeInBytes:], hash)
	binary.LittleEndian.PutUint16(blob[timestampSizeInBytes+hashSizeInBytes:], uint16(keyLength))
	copy(blob[headersSizeInBytes:], key)
	copy(blob[headersSizeInBytes+keyLength:], entry)

	return blob[:blobLength]
}

func readKeyFromEntry(data []byte) string {
	length := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, length)
	copy(dst, data[headersSizeInBytes:headersSizeInBytes+length])
	return bytesToString(dst)
}

func readEntry(data []byte) []byte {
	length := binary.LittleEndian.Uint16(data[timestampSizeInBytes+hashSizeInBytes:])

	dst := make([]byte, len(data) - int(length + headersSizeInBytes))
	copy(dst, data[headersSizeInBytes+length:])

	return dst
}

func readExpireAtFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data)
}

func readHashFromEntry(data []byte) uint64 {
	return binary.LittleEndian.Uint64(data[timestampSizeInBytes:])
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}