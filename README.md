# lru-cache

A thread-safe LRU (Least Recently Used) cache library written in Go, featuring two concurrency backends: **Mutex** (`sync.RWMutex`) and **Actor** (goroutine + channels). Ideal for learning concurrency patterns in Go or using as a production‑ready cache.

## Features

- ✅ **LRU eviction** – O(1) get/set with doubly linked list + map.
- ✅ **Two concurrency backends** – choose between `Mutex` (fast, simple) and `Actor` (channel‑based, no locks).
- ✅ **Eviction callback** – hook when an entry is removed.
- ✅ **Fully tested** – passes `go test -race` and includes benchmarks.

## Backend Comparison

| Backend | Pros | Cons
|:---------|:------:|:------|
|Mutex|	Lower overhead, higher throughput for read‑heavy workloads|	Uses locks, subject to contention|
|Actor|	No locks, pure channel communication, easier to reason about|	Slightly higher latency, one background goroutine|

Switch backend by passing lru.ActorBackend to New().

## Running Tests and Benchmarks
```bash
# Run race detector tests
go test -race -v

# Run benchmarks (memory allocation included)
go test -bench="^Benchmark" -benchmem
```
Example output:

```text
BenchmarkMutexSet-12    3529717     359.4 ns/op     102 B/op        4 allocs/op
BenchmarkActorSet-12    2213468     530.2 ns/op     40 B/op         3 allocs/op
```

## TODO:

- [ ] Atomic GetOrSet – compute and store only once even under high concurrency.
- [ ] TTL support – per‑item expiration with automatic cleanup (janitor goroutine).
- [ ] Statistics – track hits, misses, and evictions.
