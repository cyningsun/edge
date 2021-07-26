package cache

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var val = make([]byte, 1024)

func BenchmarkLRU_Write(b *testing.B) {
	for _, concurrency := range []int{1, 16, 32} {
		for _, capacity := range []int{2 ^ 13, 2 ^ 17, 2 ^ 21, 2 ^ 24} {
			b.Run(fmt.Sprintf("%d-concurrency-%d-capacity", concurrency, capacity), func(b *testing.B) {
				writeToLRU(b, concurrency, capacity)
			})
		}
	}
}

func BenchmarkLRU_Read(b *testing.B) {
	for _, concurrency := range []int{1, 16, 32} {
		for _, capacity := range []int{2 ^ 13, 2 ^ 17, 2 ^ 21, 2 ^ 24} {
			b.Run(fmt.Sprintf("%d-concurrency-%d-capacity", concurrency, capacity), func(b *testing.B) {
				readFromLRU(b, concurrency, capacity)
			})
		}
	}
}

func BenchmarkLRU_ReadNotExists(b *testing.B) {
	for _, concurrency := range []int{1, 16, 32} {
		for _, capacity := range []int{2 ^ 13, 2 ^ 17, 2 ^ 21, 2 ^ 24} {
			b.Run(fmt.Sprintf("%d-concurrency-%d-capacity", concurrency, capacity), func(b *testing.B) {
				readFromLRUNotExists(b, concurrency, capacity)
			})
		}
	}
}

func writeToLRU(b *testing.B, concurrency, capacity int) {
	cache, _ := NewLRU(WithConcurrency(concurrency), WithCapacity(capacity))
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0

		b.ReportAllocs()
		for pb.Next() {
			cache.Set(fmt.Sprintf("key-%d-%d", id, counter), val)
			counter = counter + 1
		}
	})
}

func readFromLRU(b *testing.B, concurrency, capacity int) {
	cache, _ := NewLRU(WithConcurrency(concurrency), WithCapacity(capacity))
	for i := 0; i < b.N; i++ {
		cache.Set(strconv.Itoa(i), val)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()

		for pb.Next() {
			cache.Get(strconv.Itoa(rand.Intn(b.N)))
		}
	})
}

func readFromLRUNotExists(b *testing.B, concurrency, capacity int) {
	cache, _ := NewLRU(WithConcurrency(concurrency), WithCapacity(capacity))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		b.ReportAllocs()

		for pb.Next() {
			cache.Get(strconv.Itoa(rand.Intn(b.N)))
		}
	})
}
