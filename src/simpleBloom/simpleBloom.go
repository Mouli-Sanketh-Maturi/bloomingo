package simpleBloom

import (
	"github.com/spaolacci/murmur3"
	"hash"
)
import "math/rand"

type bloom struct {
	bitArray []bool
	hash     []hash.Hash32
}

const defaultHashFunctionCount = 2

func NewBloom(size int, seeds []int) *bloom {
	// Initialize the bloom filter
	bloomFilter := &bloom{bitArray: make([]bool, size)}

	if seeds == nil || len(seeds) == 0 {
		seeds = make([]int, defaultHashFunctionCount)
		for i := 0; i < defaultHashFunctionCount; i++ {
			seeds[i] = rand.Int()
		}
	}

	hash := make([]hash.Hash32, 0)
	for _, seed := range seeds {
		hash = append(hash, murmur3.New32WithSeed(uint32(seed)))
	}

	bloomFilter.hash = hash

	return bloomFilter
}

// Simple bloom filter implementation using hash32 functions

func AddElement(bloomFilter *bloom, element string) {
	bitSize := len(bloomFilter.bitArray)
	for _, h := range bloomFilter.hash {
		h.Write([]byte(element))
		bloomFilter.bitArray[h.Sum32()%uint32(bitSize)] = true
		h.Reset()
	}
}

func CheckElement(bloomFilter *bloom, element string) bool {
	bitSize := len(bloomFilter.bitArray)
	for _, h := range bloomFilter.hash {
		_, err := h.Write([]byte(element))
		if err != nil {
			return false
		}
		if !bloomFilter.bitArray[h.Sum32()%uint32(bitSize)] {
			return false
		}
		h.Reset()
	}
	return true
}
