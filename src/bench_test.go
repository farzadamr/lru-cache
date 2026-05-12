package lru

import (
	"fmt"
	"testing"
)

func BenchmarkMutexSet(b *testing.B) {
	cache := NewMutexCache(1000)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(fmt.Sprintf("key%d", i), i)
			i++
		}
	})
}

func BenchmarkActorSet(b *testing.B) {
	cache := NewActorCache(1000)
	defer cache.(*actorCache).Close()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(fmt.Sprintf("key%d", i), i)
			i++
		}
	})
}
