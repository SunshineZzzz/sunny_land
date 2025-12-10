package utils

import (
	"encoding/binary"
	"unsafe"
)

// float32切片转换为字节切片
func Float32ToBytes(data []float32) []byte {
	bytes := make([]byte, len(data)*4)
	for i, f := range data {
		binary.LittleEndian.PutUint32(bytes[i*4:], *(*uint32)(unsafe.Pointer(&f)))
	}
	return bytes
}
