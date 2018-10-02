package hash

import (
	"hash"
	"hash/fnv"
)

var fnvHash32 hash.Hash32 = fnv.New32a()

func FNV32a(s string) uint32 {
	fnvHash32.Write([]byte(s))
	defer fnvHash32.Reset()

	return fnvHash32.Sum32()
}
