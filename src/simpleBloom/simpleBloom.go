package simpleBloom

import (
	"github.com/twmb/murmur3"
	"hash"
)

func newBloom() []hash.Hash32 {
	// Initialize the bloom filter
	seeds := []int{1, 2, 3, 4, 5}
	hash := make([]hash.Hash32, len(seeds))
	for _, seed := range seeds {
		hash = append(hash, murmur3.SeedNew32(uint32(seed)))
	}
	return hash
}

// Simple bloom filter implementation using hash32 functions
func addElement(bitArray []bool, element string, hash []hash.Hash32) {
	bitSize := len(bitArray)
	for _, h := range hash {
		h.Write([]byte(element))
		bitArray[h.Sum32()%uint32(bitSize)] = true
		h.Reset()
	}
}

func checkElement(bitArray []bool, element string, hash []hash.Hash32) bool {
	bitSize := len(bitArray)
	for _, h := range hash {
		_, err := h.Write([]byte(element))
		if err != nil {
			return false
		}
		if !bitArray[h.Sum32()%uint32(bitSize)] {
			return false
		}
		h.Reset()
	}
	return true
}
