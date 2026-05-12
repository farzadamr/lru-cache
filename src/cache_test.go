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
