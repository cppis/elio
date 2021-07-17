package elio

import (
	"fmt"
	"unsafe"
)

// Fnv32 fnv32 type definition
type Fnv32 uint32

// FromUint32 from uint32
func (f *Fnv32) FromUint32(v uint32) {
	*f = Fnv32(v)
}

// ToUint32 to uint32
func (f *Fnv32) ToUint32() uint32 {
	return uint32(*f)
}

// FromString from string
func (f *Fnv32) FromString(k string) {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(k); i++ {
		hash *= prime32
		hash ^= uint32(k[i])
	}

	f.FromUint32(hash)
}

// // Fnv32 fnv32
// func Fnv32(key string) uint32 {
// 	hash := uint32(2166136261)
// 	const prime32 = uint32(16777619)
// 	for i := 0; i < len(key); i++ {
// 		hash *= prime32
// 		hash ^= uint32(key[i])
// 	}
// 	return hash
// }

// Fnv64 fnv64 type definition
type Fnv64 uint64

// FromUint64 from uint64
func (f *Fnv64) FromUint64(v uint64) {
	*f = Fnv64(v)
}

// ToUint64 to uint64
func (f *Fnv64) ToUint64() uint64 {
	return uint64(*f)
}

// FromString from string
func (f *Fnv64) FromString(k string) {
	hash := uint64(2166136261)
	const prime64 = uint64(16777619)
	for i := 0; i < len(k); i++ {
		hash *= prime64
		hash ^= uint64(k[i])
	}

	f.FromUint64(hash)
}

// FromPointer from pointer
func (f *Fnv64) FromPointer(p unsafe.Pointer) {
	// TODO: pointer 가 짝수만 나와 문자열로 수정
	//f.FromUint64(uint64(uintptr(p)))
	f.FromString(fmt.Sprintf("%p", p))
}
