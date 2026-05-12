package lru

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentAccess(t *testing.T) {
	caches := []Cache{
		NewMutexCache(3),
		NewActorCache(3),
	}

	for _, cache := range caches {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				key := fmt.Sprintf("key-%d", i%5)
				cache.Set(key, i)
				cache.Get(key)
			}(i)
		}
		wg.Wait()
	}
}

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
